package server

import (
	"context"
	"errors"
	"net"

	"encoding/json"

	"time"

	"crypto/ecdsa"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/logger"
	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	onConnectionHandler OnConnectionHandler
	onErrorHandler      OnErrorHandler
	priKey              *ecdsa.PrivateKey
	pubKey              *ecdsa.PublicKey
	ip                  string
	lis                 net.Listener
	idGetter            bifrost.IDGetter
	formatter           bifrost.Formatter
	signer              bifrost.Signer
	verifier            bifrost.Verifier
}

func (s Server) BifrostStream(streamServer pb.StreamService_BifrostStreamServer) error {
	//1. RquestPeer를 통해 나에게 Stream연결을 보낸 ConnInfo의정보를 확인
	//2. ConnInfo의정보를정보를 기반으로 Connection을 생성
	//3. 생성완료후 OnConnectionHandler를 통해 처리한다.

	ip := extractRemoteAddress(streamServer)

	_, cf := context.WithCancel(context.Background())
	streamWrapper := bifrost.NewServerStreamWrapper(streamServer, cf)

	pub, err := s.handShake(streamWrapper, s.pubKey)

	if err != nil {
		return err
	}

	conn, err := bifrost.NewConnection(ip, s.priKey, pub, streamWrapper, s.idGetter, s.formatter, s.signer, s.verifier)

	if s.onConnectionHandler != nil {
		s.onConnectionHandler(conn)
	}

	return nil
}

func (s Server) handShake(streamWrapper bifrost.StreamWrapper, pubKey *ecdsa.PublicKey) (*ecdsa.PublicKey, error) {

	err := requestInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	pub, err := s.getClientInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	err = s.sendInfo(streamWrapper, pubKey)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	logger.Infof(nil, "handshake success")

	return pub, nil
}

func requestInfo(streamWrapper bifrost.StreamWrapper) error {
	envelope := &pb.Envelope{Type: pb.Envelope_REQUEST_PEERINFO}

	if err := streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func (s Server) sendInfo(streamWrapper bifrost.StreamWrapper, pubKey *ecdsa.PublicKey) error {

	envelope, err := bifrost.BuildResponsePeerInfo(pubKey, s.formatter)

	if err != nil {
		return errors.New("fail to build info")
	}

	if err = streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func (s Server) getClientInfo(streamWrapper bifrost.StreamWrapper) (*ecdsa.PublicKey, error) {

	env, err := bifrost.RecvWithTimeout(3*time.Second, streamWrapper)

	if err != nil {
		return nil, err
	}

	if env.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		logger.Infof(nil, "Invaild message type")
		return nil, errors.New("invalid message type")
	}

	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)

	if err != nil {
		return nil, err
	}

	pubKey := s.formatter.FromByte(peerInfo.Pubkey, peerInfo.CurveOpt)

	return pubKey, nil
}

func (s Server) validateRequestPeerInfo(envelope *pb.Envelope) (bool, string, *ecdsa.PublicKey) {

	if envelope.GetType() != pb.Envelope_REQUEST_PEERINFO {
		logger.Infof(nil, "Invaild message type")
		return false, "", nil
	}
	return s.ValidatePeerInfo(envelope)
}

func (s Server) ValidateResponsePeerInfo(envelope *pb.Envelope) (bool, string, *ecdsa.PublicKey) {

	if envelope.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		logger.Infof(nil, "Invaild message type")
		return false, "", nil
	}
	return s.ValidatePeerInfo(envelope)
}

func (s Server) ValidatePeerInfo(envelope *pb.Envelope) (bool, string, *ecdsa.PublicKey) {

	logger.Infof(nil, "Received payload [%s]", envelope.Payload)

	peerInfo := &bifrost.PeerInfo{}

	err := json.Unmarshal(envelope.Payload, peerInfo)

	if err != nil {
		logger.Infof(nil, "fail to unmarshal message [%s]", err.Error())
		return false, "", nil
	}

	pubKey := s.formatter.FromByte(peerInfo.Pubkey, peerInfo.CurveOpt)

	return true, peerInfo.IP, pubKey
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

func New(key bifrost.KeyOpts, idGetter bifrost.IDGetter, formatter bifrost.Formatter, signer bifrost.Signer, verifier bifrost.Verifier) *Server {
	return &Server{
		priKey:    key.PriKey,
		pubKey:    key.PubKey,
		idGetter:  idGetter,
		formatter: formatter,
		signer:    signer,
		verifier:  verifier,
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
		logger.Fatal(nil, err.Error())
	}

	defer lis.Close()

	g := grpc.NewServer()

	defer g.Stop()
	pb.RegisterStreamServiceServer(g, s)
	reflection.Register(g)

	s.lis = lis
	logger.Infof(nil, "Listen... on: [%s]", ip)
	if err := g.Serve(lis); err != nil {
		logger.Infof(nil, "failed to serve: %v", err)
		g.Stop()
		lis.Close()
	}
}

func (s *Server) Stop() {

	if s.lis != nil {
		s.lis.Close()
	}
}
