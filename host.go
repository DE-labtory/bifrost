package bifrost

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"errors"

	"fmt"

	"log"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/stream"
	"google.golang.org/grpc"
)

const (
	REQUEST_CONNINFO     = "/requestConnInfo"
	CONNECTION_ESTABLISH = "/connectionEstablish"
)

type Host interface {
	//Register(*grpc.Server)
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

type OnConnectionHandler func(conn.Connection)

type BifrostHost struct {
	mux                 *mux.Mux
	myConnInfo          conn.MyConnectionInfo
	server              *grpc.Server
	onConnectionHandler OnConnectionHandler
}

func New(myConnInfo conn.MyConnectionInfo, mux *mux.Mux, onConnectionHandler OnConnectionHandler) *BifrostHost {

	host := &BifrostHost{
		mux:                 mux,
		myConnInfo:          myConnInfo,
		onConnectionHandler: onConnectionHandler,
	}

	return host
}

func (bih BifrostHost) createEnvelope(protocol string, data interface{}) (*pb.Envelope, error) {

	payload, err := json.Marshal(data)
	//todo signing process
	if err != nil {
		return nil, err
	}

	pub, err := bih.myConnInfo.PubKey.ToPEM()

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

func (bih BifrostHost) ConnectToPeer(address Address) (conn.Connection, error) {

	endPointAddress := stream.Address{IP: address.Ip}
	grpc_conn, err := stream.NewClientConn(endPointAddress, false, nil)

	streamWrapper, err := stream.Connect(grpc_conn)

	if err != nil {
		return nil, err
	}

	//handshake
	// 1. wait identity request
	// 2. send identity
	// 3. connection Established

	// 1.
	envelope, err := recvWithTimeout(2, streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	// 2.
	if IsRequestConnInfoProtocol(envelope.GetProtocol()) {
		info := bih.myConnInfo.GetPublicInfo()

		envelope, err := bih.createEnvelope(REQUEST_CONNINFO, info)

		if err != nil {
			return nil, err
		}

		err = streamWrapper.Send(envelope)

		if err != nil {
			streamWrapper.Close()
			return nil, err
		}

		// 3.
		envelope, err = recvWithTimeout(2, streamWrapper)

		if err != nil {
			streamWrapper.Close()
			return nil, err
		}

		if IsConnectionIstablishProtocol(envelope.GetProtocol()) {

			log.Printf("Received payload [%s]", envelope.Payload)

			connectedConnInfo := &conn.ConnenctionInfo{}
			err := json.Unmarshal(envelope.Payload, connectedConnInfo)

			if err != nil {
				return nil, err
			}

			conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.mux)

			go func() {
				if err != conn.Start() {
					conn.Close()
				}
			}()

			return conn, nil
		}
	}

	return nil, errors.New("Not a Request Identity Protocol")
}

func (bih BifrostHost) Stream(streamServer pb.StreamService_StreamServer) error {
	//1. RquestPeer를 통해 나에게 Stream연결을 보낸 ConnInfo의정보를 확인
	//2. ConnInfo의정보를정보를 기반으로 Connection을 생성
	//3. 생성완료후 OnConnectionHandler를 통해 처리한다.

	//todo siging
	info := bih.myConnInfo.GetPublicInfo()
	envelope, err := bih.createEnvelope(REQUEST_CONNINFO, info)

	err = streamServer.Send(envelope)

	if err != nil {
		return err
	}

	if m, err := recvWithTimeout(2, streamServer); err == nil {

		wg := sync.WaitGroup{}
		wg.Add(1)

		if !IsRequestConnInfoProtocol(m.GetProtocol()) {
			return errors.New(fmt.Sprintf("Not a request connInfo protocol [%s]", m.GetProtocol()))
		}

		log.Printf("Received payload [%s]", envelope.Payload)

		connectedConnInfo := &conn.ConnenctionInfo{}
		err := json.Unmarshal(m.Payload, connectedConnInfo)

		//validate connectedInfo
		if err != nil {
			return err
		}

		_, cf := context.WithCancel(context.Background())
		streamWrapper := stream.NewServerStreamWrapper(streamServer, cf)

		conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.mux)

		go func() {
			if err != conn.Start() {
				conn.Close()
			}
		}()

		bih.onConnectionHandler(conn)

		wg.Wait()
	}

	return nil
}

func recvWithTimeout(seconds int, wrapper stream.Stream) (*pb.Envelope, error) {

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

func IsRequestConnInfoProtocol(protocol string) bool {

	if protocol == REQUEST_CONNINFO {
		return true
	}
	return false
}

func IsConnectionIstablishProtocol(protocol string) bool {

	if protocol == CONNECTION_ESTABLISH {
		return true
	}
	return false
}
