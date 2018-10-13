package mocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net"

	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/iLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type MockStreamServer struct {
	countRecv int32
	countSend int32
	peerInfo  bifrost.PeerInfo
}

func NewMockStreamServer(peerInfo bifrost.PeerInfo) *MockStreamServer {
	return &MockStreamServer{
		countRecv: 0,
		countSend: 0,
		peerInfo:  peerInfo,
	}
}

func (s *MockStreamServer) Send(envelope *pb.Envelope) error {
	iLogger.Info(nil, "[Bifrost] Mock send func called")

	s.countSend = s.countSend + 1

	if s.countSend == 1 {
		if envelope.Type == pb.Envelope_REQUEST_PEERINFO {
			return nil
		}
		return errors.New("invalid protocol")
	}

	mockServer := NewMockServer()

	if s.countSend == 2 {
		valid, _, _ := mockServer.ValidateResponsePeerInfo(envelope)

		if valid {
			return nil
		} else {
			return errors.New("invaild peerinfo")
		}
	}
	return nil
}

func (s *MockStreamServer) Recv() (*pb.Envelope, error) {

	s.countRecv = s.countRecv + 1

	if s.countRecv == 1 {
		payload, _ := json.Marshal(s.peerInfo)

		envelope := &pb.Envelope{}
		envelope.Type = pb.Envelope_RESPONSE_PEERINFO
		envelope.Payload = payload
		return envelope, nil
	}

	return nil, nil
}

func (MockStreamServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (MockStreamServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (MockStreamServer) SetTrailer(metadata.MD) {
	panic("implement me")
}

type MockvalueCtx struct {
	context.Context
	key, val interface{}
}

func (c *MockvalueCtx) String() string {
	return fmt.Sprintf("%v.WithValue(%#v, %#v)", c.Context, c.key, c.val)
}

func (c *MockvalueCtx) Value(key interface{}) interface{} {

	return "127.0.0.1:7777"
}

func (MockStreamServer) Context() context.Context {
	return &MockvalueCtx{}
}

func (MockStreamServer) SendMsg(m interface{}) error {
	panic("implement me")
}

func (MockStreamServer) RecvMsg(m interface{}) error {
	panic("implement me")
}

type MockConnectionHandler func(stream pb.StreamService_BifrostStreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	Rh  MockRecvHandler
	Ch  MockConnectionHandler
	Clh MockCloseHandler
}

type MockHandler struct{}

func (h MockHandler) ServeRequest(message bifrost.Message) {

}

func (h MockHandler) ServeError(conn bifrost.Connection, err error) {

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
		iLogger.Errorf(nil, "[Bifrost] Failed to listen: %v", err.Error())
	}

	iLogger.Infof(nil, "listen on [%s]...", ipAddress)

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, mockServer)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			iLogger.Fatalf(nil, "failed to serve: %v", err.Error())
			s.Stop()
			lis.Close()
		}
	}()

	return s, lis
}

func NewMockKeyOpts() bifrost.KeyOpts {
	pri, pub, err := NewMockKeyPair()
	if err != nil {
		iLogger.Fatalf(nil, err.Error())
	}

	return bifrost.KeyOpts{
		PubKey: pub,
		PriKey: pri,
	}
}

func NewMockServer() *server.Server {
	keyOpt := NewMockKeyOpts()
	mockCrypto := NewMockCrypto()
	mockCrypto.Signer.(*MockECDSASigner).KeyID = keyOpt.PubKey.ID()

	s := server.New(keyOpt, mockCrypto, nil)

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
