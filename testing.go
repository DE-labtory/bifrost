package bifrost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"net"

	"time"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type StreamServer struct {
	countRecv int32
	countSend int32
	peerInfo  PeerInfo
}

func NewMockStreamServer(peerInfo PeerInfo) *StreamServer {
	return &StreamServer{
		countRecv: 0,
		countSend: 0,
		peerInfo:  peerInfo,
	}
}

func (s StreamServer) Send(envelope *pb.Envelope) error {
	log.Print("Mock send func called")

	s.countSend = s.countSend + 1

	if s.countSend == 1 {
		if envelope.Type == pb.Envelope_REQUEST_PEERINFO {
			return nil
		}
		return errors.New("invalid protocol")
	}

	if s.countSend == 2 {
		bool, _, _ := ValidateRequestPeerInfo(envelope)

		if bool {
			return nil
		} else {
			return errors.New("invaild peerinfo")
		}
	}
	return nil
}

func (s StreamServer) Recv() (*pb.Envelope, error) {

	s.countRecv = s.countRecv + 1

	if s.countRecv == 1 {
		payload, _ := json.Marshal(s.peerInfo)

		envelope := &pb.Envelope{}
		envelope.Type = pb.Envelope_REQUEST_PEERINFO
		envelope.Payload = payload
		return envelope, nil
	}

	return nil, nil
}

func (StreamServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (StreamServer) Context() context.Context {
	panic("implement me")
}

func (StreamServer) SendMsg(m interface{}) error {
	panic("implement me")
}

func (StreamServer) RecvMsg(m interface{}) error {
	panic("implement me")
}

type MockConnectionHandler func(stream pb.StreamService_BifrostStreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	rh  MockRecvHandler
	ch  MockConnectionHandler
	clh MockCloseHandler
}

type Handler struct{}

func (h Handler) ServeRequest(message Message) {

}

func (h Handler) ServeError(conn Connection, err error) {

}

func (ms MockServer) BifrostStream(stream pb.StreamService_BifrostStreamServer) error {

	if ms.ch != nil {
		ms.ch(stream)
	}

	for {
		envelope, err := stream.Recv()

		//fmt.Printf(err.Error())

		if err == io.EOF {
			return nil
		}

		if err != nil {
			if ms.clh != nil {
				ms.clh()
			}
			return err
		}

		if ms.rh != nil {
			ms.rh(envelope)
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

func getKeyOpts(path string) KeyOpts {

	km, err := key.NewKeyManager(path)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	return KeyOpts{
		pubKey: pub,
		priKey: pri,
	}
}

func getServer(path string) *Server {

	keyOpt := getKeyOpts(path)

	s := NewServer(keyOpt)

	return s
}

func dial(serverIP string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithTimeout(3*time.Second))
	grpc_conn, err := grpc.Dial(serverIP, opts...)

	if err != nil {
		return nil, err
	}

	return grpc_conn, nil
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
