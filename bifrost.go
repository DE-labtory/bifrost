package bifrost

import (
	"context"
	"encoding/json"

	"time"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
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

type Bifrost struct {
	mux        *mux.Mux
	myConnInfo conn.MyConnectionInfo
	connectionInfo
}

func NewHost(server *grpc.Server) *Bifrost {

	mux := mux.NewMux()

	host := &Bifrost{
		mux: mux,
	}

	//set default handler
	//mux.Handle(REQUEST_IDENTITY, host.handleRequestIdentity)
	//mux.Handle(CONNECTION_ESTABLISH, host.handleConnectionEstablish)
	mux.HandleError(host.handleError)

	return host
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

//func (bih BifrostHost) handleConnectionEstablish(message stream.OutterMessage) {
//	//peer추가
//	//todo verify 추가
//
//	connectedPeer := peer.ConnenctionInfo{}
//	err := json.Unmarshal(message.Data, &connectedPeer)
//
//	if err != nil {
//		return
//	}
//
//	streamHandler := message.Stream
//}
//
//func (bih BifrostHost) handleRequestIdentity(message stream.OutterMessage) {
//
//	info := bih.identity.GetPublicInfo()
//
//	envelope, err := bih.createEnvelope(REQUEST_IDENTITY, info)
//
//	if err != nil {
//		return
//	}
//
//	message.Respond(envelope, nil, nil)
//}

func (bih BifrostHost) ConnectToPeer(peer peer.ConnenctionInfo) error {

	endPointAddress := stream.Address{IP: peer.Address.IP}
	grpc_conn, err := stream.NewClientConn(endPointAddress, false, nil)

	streamWrapper, err := stream.Connect(grpc_conn)

	//handshake
	// 1. wait identity request
	// 2. send identity
	// 3. connection Established

	// 1.
	envelope, err := recvWithTimeout(2, streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return err
	}

	// 2.
	if IsRequestIdentityProtocol(envelope.GetProtocol()) {
		info := bih.identity.GetPublicInfo()

		envelope, err := bih.createEnvelope(REQUEST_IDENTITY, info)

		if err != nil {
			return err
		}

		err = streamWrapper.Send(envelope)

		if err != nil {
			streamWrapper.Close()
			return err
		}

		// 3.
		envelope, err = recvWithTimeout(2, streamWrapper)

		if err != nil {
			streamWrapper.Close()
			return err
		}

		if IsConnectionIstablishProtocol(envelope.GetProtocol()) {

		}
	}

	if err != nil {
		return err
	}

	return nil
}

func recvWithTimeout(seconds int, wrapper stream.StreamWrapper) (*pb.Envelope, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	c := make(chan *pb.Envelope, 1)
	errch := make(chan error, 1)

	go func() {
		envelope, err := wrapper.Recv()
		if err != nil {
			errch <- err
		}
		c <- envelope
	}()

	select {
	case <-ctx.Done():
		//timeoutted body
		return nil, ctx.Err()
	case err := <-errch:
		return nil, err
	case ok := <-c:
		//okay body
		return ok, nil
	}
}

func IsRequestIdentityProtocol(protocol string) bool {

	if protocol == REQUEST_IDENTITY {
		return true
	}
	return false
}

func IsConnectionIstablishProtocol(protocol string) bool {

	if protocol == REQUEST_IDENTITY {
		return true
	}
	return false
}
