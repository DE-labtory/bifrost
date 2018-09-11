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

package bifrost_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"crypto/ecdsa"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	mux2 "github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockGenerator struct {
}

func (generator MockGenerator) GenerateKey() (*ecdsa.PrivateKey, error) {
	return new(ecdsa.PrivateKey), nil
}

type MockSigner struct {
}

func (signer *MockSigner) Sign() ([]byte, error) {
	return []byte("signature"), nil
}

type MockFormatter struct {
}

func (formatter *MockFormatter) ToByte(*ecdsa.PublicKey) []byte {
	return []byte("byte format of ecdsa public key")
}

func (formatter *MockFormatter) FromByte([]byte, int) *ecdsa.PublicKey {
	return new(ecdsa.PublicKey)
}

func (formatter *MockFormatter) GetCurveOpt(pubKey *ecdsa.PublicKey) int {
	return *new(int)
}

type MockIdGetter struct {
}

func (idGetter *MockIdGetter) GetID(key *ecdsa.PublicKey) bifrost.ID {
	return *new(bifrost.ID)
}

type MockServer struct {
}

func (ms MockServer) Stream(stream pb.StreamService_StreamServer) error {
	mockGenerator := MockGenerator{}
	mockFormatter := MockFormatter{}

	pri, err := mockGenerator.GenerateKey()
	pub := &pri.PublicKey

	envelope := &pb.Envelope{}
	envelope.Protocol = bifrost.REQUEST_CONNINFO
	err = stream.Send(envelope)

	if err != nil {
		log.Fatalf(err.Error())
	}

	connectionInfo, err := stream.Recv()

	log.Printf("Received Connection Info is [%s]", connectionInfo)

	if err != nil {
		log.Fatalf(err.Error())
	}

	b := mockFormatter.ToByte(pub)

	pci := conn.PublicConnInfo{}
	pci.Id = "test1"
	pci.Address = conn.Address{IP: "127.0.0.1"}
	pci.Pubkey = b
	pci.CurveOpt = mockFormatter.GetCurveOpt(pub)

	envelope2 := &pb.Envelope{}
	envelope2.Protocol = bifrost.CONNECTION_ESTABLISH
	payload, err := json.Marshal(pci)
	if err != nil {
		log.Fatalf(err.Error())
	}
	envelope2.Payload = payload

	err = stream.Send(envelope2)

	if err != nil {
		log.Fatalf(err.Error())
	}

	testEnvelope, err := stream.Recv()

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Recevied Test envelop is [%s]", testEnvelope)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}

func ListenMockServer(mockServer pb.StreamServiceServer, ipAddress string) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ipAddress)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, mockServer)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()

	return s, lis
}

func TestBifrostHost_ConnectToPeer(t *testing.T) {

	serverIP := "127.0.0.1:9999"
	mockServer := &MockServer{}
	server1, listner1 := ListenMockServer(mockServer, serverIP)

	mockGenerator := MockGenerator{}

	priv, err := mockGenerator.GenerateKey()
	assert.Nil(t, err)

	idGetter := MockIdGetter{}
	address, err := conn.ToAddress("127.0.0.1:8888")
	assert.NoError(t, err)
	myconnectionInfo := bifrost.NewHostInfo(address, priv, &idGetter)
	mux := mux2.NewMux()
	mockSigner := &MockSigner{}
	mockFormatter := &MockFormatter{}

	host := bifrost.New(myconnectionInfo, mux, nil, mockSigner, mockFormatter)

	connection, err := host.ConnectToPeer(bifrost.Address{Ip: "127.0.0.1:9999"})
	assert.Nil(t, err)
	log.Printf("Sending data...")
	connection.Send(&pb.Envelope{Payload: []byte("test1")}, nil, nil)

	assert.Equal(t, "test1", connection.GetConnInfo().Id.ToString())

	time.Sleep(2 * time.Second)
	server1.Stop()
	listner1.Close()
}

func TestBifrostHost_Stream(t *testing.T) {

	mockGenerator := MockGenerator{}

	priv, err := mockGenerator.GenerateKey()

	idGetter := MockIdGetter{}
	address, err := conn.ToAddress("127.0.0.1:8888")
	assert.NoError(t, err)
	myconnectionInfo := bifrost.NewHostInfo(address, priv, &idGetter)
	mux := mux2.NewMux()
	mockSigner := &MockSigner{}
	mockFormatter := &MockFormatter{}

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
		assert.Equal(t, connection.GetConnInfo().Address.IP, "127.0.0.1:8888")
	}

	serverHost := bifrost.New(myconnectionInfo, mux, OnConnectionHandler, mockSigner, mockFormatter)
	serverIP := "127.0.0.1:8888"
	server1, listner1 := ListenMockServer(serverHost, serverIP)

	defer func() {
		server1.Stop()
		listner1.Close()
	}()

	clientHost := bifrost.New(myconnectionInfo, mux, nil, mockSigner, mockFormatter)

	connection, err := clientHost.ConnectToPeer(bifrost.Address{Ip: serverIP})

	fmt.Println(connection)

	if err != nil {
		fmt.Printf("error is [%s]", err.Error())
	}

	//fmt.Println(connection)

	time.Sleep(2 * time.Second)
}
