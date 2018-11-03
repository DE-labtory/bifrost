package server

import (
	"context"
	"errors"
	"net"

	"encoding/json"

	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/iLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	onConnectionHandler OnConnectionHandler
	onErrorHandler      OnErrorHandler
	pubKey              bifrost.Key
	ip                  string
	lis                 net.Listener
	metaData            map[string]string
	bifrost.Crypto
}

func (s Server) BifrostStream(streamServer pb.StreamService_BifrostStreamServer) error {
	//1. RquestPeer를 통해 나에게 Stream연결을 보낸 ConnInfo의정보를 확인
	//2. ConnInfo의정보를정보를 기반으로 Connection을 생성
	//3. 생성완료후 OnConnectionHandler를 통해 처리한다.

	ip := extractRemoteAddress(streamServer)

	_, cf := context.WithCancel(context.Background())
	streamWrapper := bifrost.NewServerStreamWrapper(streamServer, cf)

	peerKey, metaData, err := s.handShake(streamWrapper)

	if err != nil {
		return err
	}

	conn, err := bifrost.NewConnection(ip, metaData, peerKey, streamWrapper, s.Crypto)

	if s.onConnectionHandler != nil {
		s.onConnectionHandler(conn)
	}

	return nil
}

func (s Server) handShake(streamWrapper bifrost.StreamWrapper) (bifrost.Key, map[string]string, error) {

	err := requestInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, nil, err
	}

	peerKey, metaData, err := s.getClientInfo(streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, nil, err
	}

	err = s.sendInfo(streamWrapper, metaData)

	if err != nil {
		streamWrapper.Close()
		return nil, nil, err
	}

	iLogger.Info(nil, "[Bifrost] Handshake success")

	return peerKey, metaData, nil
}

func requestInfo(streamWrapper bifrost.StreamWrapper) error {
	envelope := &pb.Envelope{Type: pb.Envelope_REQUEST_PEERINFO}

	if err := streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func (s Server) sendInfo(streamWrapper bifrost.StreamWrapper, metaData map[string]string) error {

	envelope, err := bifrost.BuildResponsePeerInfo(s.ip, s.pubKey, metaData)

	if err != nil {
		return errors.New("fail to build info")
	}

	if err = streamWrapper.Send(envelope); err != nil {
		return err
	}

	return nil
}

func (s Server) getClientInfo(streamWrapper bifrost.StreamWrapper) (bifrost.Key, map[string]string, error) {

	env, err := bifrost.RecvWithTimeout(3*time.Second, streamWrapper)

	if err != nil {
		return nil, nil, err
	}

	if env.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		iLogger.Info(nil, "[Bifrost] Invaild message type")
		return nil, nil, errors.New("invalid message type")
	}

	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := s.Crypto.RecoverKeyFromByte(peerInfo.PubKeyBytes, peerInfo.IsPrivate)
	if err != nil {
		return nil, nil, err
	}

	return pubKey, peerInfo.MetaData, nil
}

func (s Server) validateRequestPeerInfo(envelope *pb.Envelope) (bool, string, bifrost.Key) {

	if envelope.GetType() != pb.Envelope_REQUEST_PEERINFO {
		iLogger.Info(nil, "[Bifrost] Invaild message type")
		return false, "", nil
	}
	return s.ValidatePeerInfo(envelope)
}

func (s Server) ValidateResponsePeerInfo(envelope *pb.Envelope) (bool, string, bifrost.Key) {

	if envelope.GetType() != pb.Envelope_RESPONSE_PEERINFO {
		iLogger.Info(nil, "[Bifrost] Invaild message type")
		return false, "", nil
	}
	return s.ValidatePeerInfo(envelope)
}

func (s Server) ValidatePeerInfo(envelope *pb.Envelope) (bool, string, bifrost.Key) {

	iLogger.Infof(nil, "[Bifrost] Received payload [%s]", envelope.Payload)

	peerInfo := &bifrost.PeerInfo{}

	err := json.Unmarshal(envelope.Payload, peerInfo)

	if err != nil {
		iLogger.Errorf(nil, "[Bifrost] Fail to unmarshal message [%s]", err.Error())
		return false, "", nil
	}

	pubKey, err := s.Crypto.RecoverKeyFromByte(peerInfo.PubKeyBytes, peerInfo.IsPrivate)

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

func New(key bifrost.KeyOpts, crypto bifrost.Crypto, metaData map[string]string) *Server {
	return &Server{
		pubKey:   key.PubKey,
		Crypto:   crypto,
		metaData: metaData,
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
		iLogger.Infof(nil, "[Bifrost] Listen error: %s", err.Error())
	}
	defer lis.Close()

	g := grpc.NewServer()

	defer g.Stop()
	pb.RegisterStreamServiceServer(g, s)
	reflection.Register(g)

	s.lis = lis

	iLogger.Infof(nil, "[Bifrost] Listen... on: [%s]", ip)
	if err := g.Serve(lis); err != nil {
		iLogger.Infof(nil, "[Bifrost] Listen... on: [%s]", ip)
		g.Stop()
		lis.Close()
	}
}

func (s *Server) Stop() {

	if s.lis != nil {
		s.lis.Close()
	}
}
