package server

import (
	"context"
	"errors"
	"log"
	"net"

	"encoding/json"

	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	onConnectionHandler OnConnectionHandler
	onErrorHandler      OnErrorHandler
	priKey              key.PriKey
	pubKey              key.PubKey
	ip                  string
	lis                 net.Listener
}

func (s Server) BifrostStream(streamServer pb.StreamService_BifrostStreamServer) error {
	//1. RquestPeer를 통해 나에게 Stream연결을 보낸 ConnInfo의정보를 확인
	//2. ConnInfo의정보를정보를 기반으로 Connection을 생성
	//3. 생성완료후 OnConnectionHandler를 통해 처리한다.

	ip := extractRemoteAddress(streamServer)

	_, cf := context.WithCancel(context.Background())
	streamWrapper := bifrost.NewServerStreamWrapper(streamServer, cf)

	pub, err := handShake(streamWrapper, s.pubKey)

	if err != nil {
		return err
	}

	conn, err := bifrost.NewConnection(ip, s.priKey, pub, streamWrapper)

	if s.onConnectionHandler != nil {
		s.onConnectionHandler(conn)
	}

	return nil
}

func handShake(streamWrapper bifrost.StreamWrapper, pubKey key.PubKey) (key.PubKey, error) {

	err := requestInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	pub, err := getClientInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	err = sendInfo(streamWrapper, pubKey)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	log.Printf("handshake success")

	return pub, nil
}

func requestInfo(streamWrapper bifrost.StreamWrapper) error {
	envelope := &pb.Envelope{Type: pb.Envelope_REQUEST_PEERINFO}

	if err := streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func sendInfo(streamWrapper bifrost.StreamWrapper, pubKey key.PubKey) error {

	envelope, err := bifrost.BuildResponsePeerInfo(pubKey)

	if err != nil {
		return errors.New("fail to build info")
	}

	if err = streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func getClientInfo(streamWrapper bifrost.StreamWrapper) (key.PubKey, error) {

	env, err := bifrost.RecvWithTimeout(3*time.Second, streamWrapper)

	if err != nil {
		return nil, err
	}

	if env.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		log.Printf("Invaild message type")
		return nil, errors.New("Invalid Message Type")
	}

	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)

	if err != nil {
		return nil, err
	}

	pubKey, err := bifrost.ByteToPubKey(peerInfo.Pubkey, peerInfo.KeyGenOpt)

	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func validateRequestPeerInfo(envelope *pb.Envelope) (bool, string, key.PubKey) {

	if envelope.GetType() != pb.Envelope_REQUEST_PEERINFO {
		log.Printf("Invaild message type")
		return false, "", nil
	}
	return ValidatePeerInfo(envelope)
}

func ValidateResponsePeerInfo(envelope *pb.Envelope) (bool, string, key.PubKey) {

	if envelope.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		log.Printf("Invaild message type")
		return false, "", nil
	}
	return ValidatePeerInfo(envelope)
}

func ValidatePeerInfo(envelope *pb.Envelope) (bool, string, key.PubKey) {

	log.Printf("Received payload [%s]", envelope.Payload)

	peerInfo := &bifrost.PeerInfo{}

	err := json.Unmarshal(envelope.Payload, peerInfo)

	if err != nil {
		log.Printf("fail to unmarshal message [%s]", err.Error())
		return false, "", nil
	}

	pubKey, err := bifrost.ByteToPubKey(peerInfo.Pubkey, peerInfo.KeyGenOpt)

	if err != nil {
		log.Printf("Invaild Pubkey [%s]", err.Error())
		return false, "", nil
	}

	return true, peerInfo.Ip, pubKey
}

func extractRemoteAddress(stream pb.StreamService_BifrostStreamServer) string {
	var remoteAddress string
	p, ok := peer.FromContext(stream.Context())
	if ok {
		if address := p.Addr; address != nil {
			remoteAddress = address.String()
		}
	}
	return remoteAddress
}

type OnConnectionHandler func(connection bifrost.Connection)
type OnErrorHandler func(err error)

func New(key bifrost.KeyOpts) *Server {
	return &Server{
		priKey: key.PriKey,
		pubKey: key.PubKey,
	}
}

func (s *Server) OnConnection(handler OnConnectionHandler) {

	if handler == nil {
		return
	}

	s.onConnectionHandler = handler
}

func (s *Server) OnError(handler OnErrorHandler) {

	if handler == nil {
		return
	}

	s.onErrorHandler = handler
}

func (s *Server) Listen(ip string) {

	lis, err := net.Listen("tcp", ip)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer lis.Close()

	g := grpc.NewServer()

	defer g.Stop()
	pb.RegisterStreamServiceServer(g, s)
	reflection.Register(g)

	s.lis = lis
	log.Printf("Listen... on: [%s]", ip)
	if err := g.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
		g.Stop()
		lis.Close()
	}
}

func (s *Server) Stop() {

	if s.lis != nil {
		s.lis.Close()
	}
}
