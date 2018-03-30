package bifrost

import (
	"log"
	"testing"

	"net"

	"encoding/json"

	"sync"

	"os"

	"github.com/it-chain/bifrost/conn"
	mux2 "github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockServer struct {
}

func (ms MockServer) Stream(stream pb.StreamService_StreamServer) error {

	envelope := &pb.Envelope{}
	envelope.Protocol = REQUEST_IDENTITY
	err := stream.Send(envelope)

	if err != nil {
		log.Fatalf(err.Error())
	}

	connectionInfo, err := stream.Recv()

	log.Printf("Received Connection Info is [%s]", connectionInfo)

	if err != nil {
		log.Fatalf(err.Error())
	}

	envelope2 := &pb.Envelope{}
	envelope2.Protocol = CONNECTION_ESTABLISH
	payload, err := json.Marshal(conn.ConnenctionInfo{Id: "123"})
	envelope2.Payload = payload

	err = stream.Send(envelope2)

	if err != nil {
		log.Fatalf(err.Error())
	}

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

func TestNew(t *testing.T) {

	serverIP := "127.0.0.1:9999"
	mockServer := &MockServer{}
	server1, listner1 := ListenMockServer(mockServer, serverIP)

	defer func() {
		server1.Stop()
		listner1.Close()
	}()

	km, err := key.NewKeyManager("~/key")

	defer os.RemoveAll("~/key")

	priv, pub, err := km.GenerateKey(key.RSA4096)

	myconnectionInfo := conn.NewMyConnectionInfo(conn.FromRsaPubKey(pub), conn.Address{IP: "127.0.0.1:8888"}, pub, priv)
	connStore := conn.NewConnectionStore()
	mux := mux2.NewMux()

	host := New(myconnectionInfo, connStore, mux, nil)

	connection, err := host.ConnectToPeer(Address{Ip: "127.0.0.1:9999"})
	//
	//fmt.Print(err)
	//fmt.Print(connection)
	assert.Nil(t, err)
	assert.Equal(t, connection.GetConnInfo().Id, conn.ID("123"))
}
