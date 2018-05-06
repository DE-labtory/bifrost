package bifrost

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockConnectionHandler func(stream pb.StreamService_BifrostStreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	Rh  MockRecvHandler
	Ch  MockConnectionHandler
	Clh MockCloseHandler
}

type MockHandler struct{}

func (h MockHandler) ServeRequest(message Message) {

}

func (h MockHandler) ServeError(conn Connection, err error) {

}

func (ms MockServer) BifrostStream(stream pb.StreamService_BifrostStreamServer) error {

	if ms.Ch != nil {
		ms.Ch(stream)
	}

	for {
		envelope, err := stream.Recv()

		//fmt.Printf(err.Error())

		if err == io.EOF {
			return nil
		}

		if err != nil {
			if ms.Clh != nil {
				ms.Clh()
			}
			return err
		}

		if ms.Rh != nil {
			ms.Rh(envelope)
		}
	}
}

func ListenMockServer(mockServer pb.StreamServiceServer, ipAddress string) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ipAddress)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, mockServer)
	reflection.Register(s)

	fmt.Printf("listen..")

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()

	return s, lis
}

func GetKeyOpts(path string) KeyOpts {

	km, err := key.NewKeyManager(path)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	return KeyOpts{
		PubKey: pub,
		PriKey: pri,
	}
}

type SendCallBack func(envelope *pb.Envelope)
type CloseCallBack func()

type MockStreamWrapper struct {
	sendCallBack  SendCallBack
	closeCallBack CloseCallBack
}

func (msw MockStreamWrapper) Send(envelope *pb.Envelope) error {
	msw.sendCallBack(envelope)
	return nil
}

func (MockStreamWrapper) Recv() (*pb.Envelope, error) {
	panic("implement me")
}

func (msw MockStreamWrapper) Close() {
	msw.closeCallBack()
}

func (MockStreamWrapper) GetStream() Stream {
	panic("implement me")
}
