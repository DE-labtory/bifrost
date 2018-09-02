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

package stream

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockConnectionHandler func(stream pb.StreamService_StreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	rh  MockRecvHandler
	ch  MockConnectionHandler
	clh MockCloseHandler
}

func (ms MockServer) Stream(stream pb.StreamService_StreamServer) error {

	if ms.ch != nil {
		ms.ch(stream)
	}

	for {
		envelope, err := stream.Recv()

		//fmt.Printf(err.Error())

		if err == io.EOF {
			return nil
		}

		if err != nil {
			if ms.clh != nil {
				ms.clh()
			}
			return err
		}

		if ms.rh != nil {
			ms.rh(envelope)
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

func TestConnect(t *testing.T) {

	//when
	connectionFlag := false
	var connectionHandler = func(stream pb.StreamService_StreamServer) {
		//result
		connectionFlag = true
	}

	var recvHandler = func(envelope *pb.Envelope) {
		//result
		assert.Equal(t, envelope.Payload, []byte("hello"))
	}

	serverIP := "127.0.0.1:9999"
	mockServer := &MockServer{ch: connectionHandler, rh: recvHandler}
	server1, listner1 := ListenMockServer(mockServer, serverIP)

	defer func() {
		server1.Stop()
		listner1.Close()
	}()

	address := Address{IP: serverIP}
	grpc_conn, _ := NewClientConn(address, false, nil)

	//then
	_, err := Connect(grpc_conn)

	if err != nil {

	}

	time.Sleep(1 * time.Second)

	assert.Equal(t, true, connectionFlag)
}
