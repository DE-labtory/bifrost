/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	Mux                 *mux.Mux
	info                HostInfo
	server              *grpc.Server
	onConnectionHandler OnConnectionHandler
	Signer              Signer
	Formatter           Formatter
}

func New(myConnInfo HostInfo, mux *mux.Mux, onConnectionHandler OnConnectionHandler, signer Signer, formatter Formatter) *BifrostHost {

	host := &BifrostHost{
		Mux:                 mux,
		info:                myConnInfo,
		onConnectionHandler: onConnectionHandler,
		Signer:              signer,
		Formatter:           formatter,
	}

	return host
}

func (bih BifrostHost) ConnectToPeer(address Address) (conn.Connection, error) {

	endPointAddress := stream.Address{IP: address.Ip}
	grpcConn, err := stream.NewClientConn(endPointAddress, false, nil)

	streamWrapper, err := stream.Connect(grpcConn)

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

		info := bih.getPublicInfo()

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

			connectedConnInfo, err := bih.pubConnInfoToConnInfo(envelope.Payload)

			if err != nil {
				return nil, err
			}

			conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.Mux)

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

		//log.Printf("Received payload [%s]", envelope.Payload)
		log.Printf("Received payload [%s]", m.Payload)

		info := bih.getPublicInfo()
		envelope, err := bih.createSignedEnvelope(CONNECTION_ESTABLISH, info)

		if err = streamServer.Send(envelope); err != nil {
			return err
		}

		connectedConnInfo, err := bih.pubConnInfoToConnInfo(envelope.Payload)

		//validate connectedInfo
		if err != nil {
			return err
		}

		_, cf := context.WithCancel(context.Background())
		streamWrapper := stream.NewServerStreamWrapper(streamServer, cf)

		conn, err := conn.NewConnection(*connectedConnInfo, streamWrapper, bih.Mux)
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

func (bih BifrostHost) pubConnInfoToConnInfo(payload []byte) (*conn.ConnInfo, error) {

	pubConnInfo := &conn.PublicConnInfo{}
	err := json.Unmarshal(payload, pubConnInfo)

	if err != nil {
		return nil, err
	}

	pubKey := bih.Formatter.FromByte(pubConnInfo.Pubkey, pubConnInfo.CurveOpt)

	return &conn.ConnInfo{
		Id:      conn.ID(pubConnInfo.Id),
		Address: pubConnInfo.Address,
		PubKey:  pubKey,
	}, nil
}

func (bih BifrostHost) createSignedEnvelope(protocol string, data interface{}) (*pb.Envelope, error) {

	payload, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	pub := bih.Formatter.ToByte(bih.info.PubKey)
	sig, err := bih.Signer.Sign()

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

func (bih BifrostHost) getPublicInfo() *conn.PublicConnInfo {
	publicConnInfo := &conn.PublicConnInfo{}
	publicConnInfo.Id = bih.info.Id.ToString()
	publicConnInfo.Address = bih.info.Address

	bytePubKey := bih.Formatter.ToByte(bih.info.PubKey)

	publicConnInfo.Pubkey = bytePubKey
	publicConnInfo.CurveOpt = bih.Formatter.GetCurveOpt(bih.info.PubKey)

	return publicConnInfo
}

func (bih BifrostHost) handleError(err error) {

}
