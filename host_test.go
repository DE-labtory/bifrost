package bifrost

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/it-chain/bifrost/conn"
	mux2 "github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/it-chain/heimdall"
)

type MockServer struct {
}

func (ms MockServer) Stream(stream pb.StreamService_StreamServer) error {

	pri, err := heimdall.GenerateKey(heimdall.SECP384R1)
	pub := &pri.PublicKey

	envelope := &pb.Envelope{}
	envelope.Protocol = REQUEST_CONNINFO
	err = stream.Send(envelope)

	if err != nil {
		log.Fatalf(err.Error())
	}

	connectionInfo, err := stream.Recv()

	log.Printf("Received Connection Info is [%s]", connectionInfo)

	if err != nil {
		log.Fatalf(err.Error())
	}

	b := heimdall.PubKeyToBytes(pub)

	pci := conn.PublicConnInfo{}
	pci.Id = "test1"
	pci.Address = conn.Address{IP: "127.0.0.1"}
	pci.Pubkey = b
	pci.CurveOpt = heimdall.CurveToCurveOpt(pub.Curve)

	envelope2 := &pb.Envelope{}
	envelope2.Protocol = CONNECTION_ESTABLISH
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

	priv, err := heimdall.GenerateKey(heimdall.SECP384R1)
	assert.Nil(t, err)
	pub := &priv.PublicKey

	myconnectionInfo := NewHostInfo(conn.Address{IP: "127.0.0.1:8888"}, pub, priv)
	mux := mux2.NewMux()

	host := New(myconnectionInfo, mux, nil)

	connection, err := host.ConnectToPeer(Address{Ip: "127.0.0.1:9999"})
	assert.Nil(t, err)
	log.Printf("Sending data...")
	connection.Send(&pb.Envelope{Payload: []byte("test1")}, nil, nil)

	assert.Equal(t, "test1", connection.GetConnInfo().Id.ToString())

	time.Sleep(2 * time.Second)
	server1.Stop()
	listner1.Close()
}

func TestBifrostHost_Stream(t *testing.T) {

	priv, err := heimdall.GenerateKey(heimdall.SECP384R1)
	pub := &priv.PublicKey

	myconnectionInfo := NewHostInfo(conn.Address{IP: "127.0.0.1:8888"}, pub, priv)
	mux := mux2.NewMux()

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
		assert.Equal(t, connection.GetConnInfo().Address.IP, "127.0.0.1:8888")
	}

	serverHost := New(myconnectionInfo, mux, OnConnectionHandler)
	serverIP := "127.0.0.1:8888"
	server1, listner1 := ListenMockServer(serverHost, serverIP)

	defer func() {
		server1.Stop()
		listner1.Close()
	}()

	clientHost := New(myconnectionInfo, mux, nil)

	connection, err := clientHost.ConnectToPeer(Address{Ip: serverIP})

	fmt.Println(connection)

	if err != nil {
		fmt.Printf("error is [%s]", err.Error())
	}

	//fmt.Println(connection)

	time.Sleep(2 * time.Second)
}
