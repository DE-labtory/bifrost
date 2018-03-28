package host

import (
	"encoding/json"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/msg"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/peer"
	"github.com/it-chain/bifrost/stream"
	"google.golang.org/grpc"
)

const (
	REQUEST_IDENTITY     = "/requestIdentity"
	CONNECTION_ESTABLISH = "/connectionEstablish"
)

type Host interface {
	Register(*grpc.Server)
}

type Address struct {
	Ip string
}

func NewAddress(ipAddress string) Address {
	// validate ip pattern
	return Address{
		Ip: ipAddress,
	}
}

type BifrostHost struct {
	server   *grpc.Server
	mux      *mux.Mux
	identity peer.Identity
}

func NewHost(server *grpc.Server) *BifrostHost {

	mux := mux.NewMux()

	host := &BifrostHost{
		server: server,
		mux:    mux,
	}

	//set default handler
	mux.Handle(REQUEST_IDENTITY, host.handleRequestIdentity)
	mux.Handle(CONNECTION_ESTABLISH, host.handleConnectionEstablish)
	mux.HandleError(host.handleError)

	return host
}

func (bih BifrostHost) ConnectToPeer(peer peer.Peer) error {

	endPointAddress := conn.Address{IP: peer.Address.IP}
	grpc_conn, err := conn.NewConnectionWithAddress(endPointAddress, false, nil)

	if err != nil {
		return err
	}

	streamWrapper, err := stream.NewClientStreamWrapper(grpc_conn)

	if err != nil {
		return err
	}

	streamHandler, err := stream.NewStreamHandler(streamWrapper, bih.mux)

	if err != nil {
		return err
	}
}

func (bih BifrostHost) createEnvelope(protocol string, data interface{}) (*pb.Envelope, error) {

	payload, err := json.Marshal(data)
	//todo signing process
	if err != nil {
		return nil, err
	}

	pub, err := bih.identity.PubKey.ToPEM()

	if err != nil {
		return nil, err
	}

	envelope := &pb.Envelope{}
	envelope.Protocol = protocol
	envelope.Payload = payload
	envelope.Pubkey = pub

	return envelope, nil
}

func (bih BifrostHost) handleError(err error) {

}

func (bih BifrostHost) handleConnectionEstablish(message msg.OutterMessage) {
	//peer추가
	//
}

func (bih BifrostHost) handleRequestIdentity(message msg.OutterMessage) {

	info := bih.identity.GetPublicInfo()

	envelope, err := bih.createEnvelope(REQUEST_IDENTITY, info)

	if err != nil {
		return
	}

	message.Respond(envelope, nil, nil)
}

//func NewHost(address Address) Host {
//	lis, err := net.Listen("tcp", address.Ip)
//
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	s := grpc.NewServer()
//	s.RegisterService(pb.StreamServiceServer{}, &BifrostHost)
//	reflection.Register(s)
//
//	return Bifrost{
//		server: s,
//	}
//}
