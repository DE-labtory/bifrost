package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"net"

	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type StreamServer struct {
	countRecv int32
	countSend int32
	peerInfo  bifrost.PeerInfo
}

func NewMockStreamServer(peerInfo bifrost.PeerInfo) *StreamServer {
	return &StreamServer{
		countRecv: 0,
		countSend: 0,
		peerInfo:  peerInfo,
	}
}

func (s *StreamServer) Send(envelope *pb.Envelope) error {
	log.Print("Mock send func called")

	s.countSend = s.countSend + 1

	if s.countSend == 1 {
		if envelope.Type == pb.Envelope_REQUEST_PEERINFO {
			return nil
		}
		return errors.New("invalid protocol")
	}

	if s.countSend == 2 {
		bool, _, _ := ValidateResponsePeerInfo(envelope)

		if bool {
			return nil
		} else {
			return errors.New("invaild peerinfo")
		}
	}
	return nil
}

func (s *StreamServer) Recv() (*pb.Envelope, error) {

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

func (StreamServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (StreamServer) SetTrailer(metadata.MD) {
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

func (StreamServer) Context() context.Context {
	return &MockvalueCtx{}
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

func GetKeyOpts(path string) bifrost.KeyOpts {

	km, err := key.NewKeyManager(path)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	return bifrost.KeyOpts{
		PubKey: pub,
		PriKey: pri,
	}
}

func GetServer(path string) *Server {

	keyOpt := GetKeyOpts(path)

	s := New(keyOpt)

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

func (MockStreamWrapper) GetStream() bifrost.Stream {
	panic("implement me")
}
