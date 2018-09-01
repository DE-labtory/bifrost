package bifrost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/stream"
	"github.com/it-chain/heimdall"
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
	info                HostInfo
	server              *grpc.Server
	onConnectionHandler OnConnectionHandler
}

func New(myConnInfo HostInfo, mux *mux.Mux, onConnectionHandler OnConnectionHandler) *BifrostHost {

	host := &BifrostHost{
		mux:                 mux,
		info:                myConnInfo,
		onConnectionHandler: onConnectionHandler,
	}

	return host
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
	envelope, err := recvWithTimeout(10, streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	// 2.
	if IsRequestConnInfoProtocol(envelope.GetProtocol()) {
		info := bih.info.GetPublicInfo()

		envelope, err := bih.createSignedEnvelope(REQUEST_CONNINFO, info)

		if err != nil {
			return nil, err
		}

		err = streamWrapper.Send(envelope)

		if err != nil {
			streamWrapper.Close()
			return nil, err
		}

		// 3.
		envelope, err = recvWithTimeout(3, streamWrapper)

		if err != nil {
			streamWrapper.Close()
			return nil, err
		}

		if IsConnectionIstablishProtocol(envelope.GetProtocol()) {

			log.Printf("Received payload [%s]", envelope.Payload)

			connectedConnInfo, err := pubConnInfoToConnInfo(envelope.Payload)

			if err != nil {
				return nil, err
			}

			conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.mux)

			go func() {
				if err = conn.Start(); err != nil {
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

	var s struct{}
	envelope, err := bih.createSignedEnvelope(REQUEST_CONNINFO, s)

	err = streamServer.Send(envelope)

	if err != nil {
		return err
	}

	if m, err := recvWithTimeout(3, streamServer); err == nil {

		wg := sync.WaitGroup{}
		wg.Add(1)

		if !IsRequestConnInfoProtocol(m.GetProtocol()) {
			return errors.New(fmt.Sprintf("Not a request connInfo protocol [%s]", m.GetProtocol()))
		}

		log.Printf("Received payload [%s]", envelope.Payload)

		info := bih.info.GetPublicInfo()
		envelope, err := bih.createSignedEnvelope(CONNECTION_ESTABLISH, info)

		if err = streamServer.Send(envelope); err != nil {
			return err
		}

		connectedConnInfo, err := pubConnInfoToConnInfo(envelope.Payload)

		//validate connectedInfo
		if err != nil {
			return err
		}

		_, cf := context.WithCancel(context.Background())
		streamWrapper := stream.NewServerStreamWrapper(streamServer, cf)

		conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.mux)
		defer conn.Close()

		go func() {
			if err = conn.Start(); err != nil {
				conn.Close()
				wg.Done()
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

func pubConnInfoToConnInfo(payload []byte) (*conn.ConnInfo, error) {

	pubConnInfo := &conn.PublicConnInfo{}
	err := json.Unmarshal(payload, pubConnInfo)

	if err != nil {
		return nil, err
	}

	connectedConnInfo, err := conn.FromPublicConnInfo(*pubConnInfo)

	if err != nil {
		return nil, err
	}

	return connectedConnInfo, nil
}

func (bih BifrostHost) createSignedEnvelope(protocol string, data interface{}) (*pb.Envelope, error) {

	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	pub := heimdall.PubKeyToBytes(bih.info.PubKey)

	sig, err := heimdall.Sign(bih.info.PriKey, payload, nil, heimdall.SHA384)

	if err != nil {
		return nil, err
	}

	envelope := &pb.Envelope{}
	envelope.Protocol = protocol
	envelope.Payload = payload
	envelope.Pubkey = pub
	envelope.Signature = sig

	return envelope, nil
}

func (bih BifrostHost) handleError(err error) {

}
