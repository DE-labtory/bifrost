package bifrost

import (
	"fmt"
	"log"
	"testing"
	"time"

	"os"

	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestNewStreamHandler(t *testing.T) {

	var connectionHandler = func(stream pb.StreamService_BifrostStreamServer) {
		//result
		fmt.Print("connected")
	}

	serverIP := "127.0.0.1:9999"
	mockServer := &MockServer{ch: connectionHandler}
	server1, listner1 := ListenMockServer(mockServer, serverIP)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithTimeout(3*time.Second))
	grpc_conn, err := grpc.Dial(serverIP, opts...)

	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, _ := context.WithCancel(context.Background())
	streamServiceClient := pb.NewStreamServiceClient(grpc_conn)
	_, err = streamServiceClient.BifrostStream(ctx)

	defer func() {
		server1.Stop()
		listner1.Close()
	}()

	time.Sleep(3 * time.Second)
}

func TestGrpcConnection_Send(t *testing.T) {

	path := "./key"
	keyOpts := getKeyOpts(path)
	defer os.RemoveAll(path)

	grpc_conn, err := dial(serverIP)
	assert.NoError(t, err)

	streamWrapper, err := NewClientStreamWrapper(grpc_conn)

	if err != nil {
		log.Fatal(err.Error())
	}

	conn, err := NewConnection(serverIP, streamWrapper)

	if err != nil {
		fmt.Errorf("error")
	}

	var success = func(interface{}) {
		fmt.Printf("success")
	}

	var fail = func(err error) {
		t.Fail()
	}

	envelope := &pb.Envelope{}
	envelope.Payload = []byte("hello")

	//then
	conn.Send(envelope, success, fail)

	time.Sleep(2 * time.Second)
}

//
//func TestStreamHandler_Send(t *testing.T) {
//
//	//when
//	var recvHandler = func(envelope *pb.Envelope) {
//		//result
//		assert.Equal(t, envelope.Payload, []byte("hello"))
//	}
//
//	serverIP := "127.0.0.1:9999"
//	mockServer := &MockServer{rh: recvHandler}
//	server1, listner1 := ListenMockServer(mockServer, serverIP)
//
//	defer func() {
//		server1.Stop()
//		listner1.Close()
//	}()
//
//	grpc_conn, err := NewClientConn(serverIP, false, nil)
//
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	streamWrapper, err := NewClientStreamWrapper(grpc_conn)
//
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	conn, err := NewConnection(ConnInfo{}, streamWrapper, Handler{})
//
//	if err != nil {
//		fmt.Errorf("error")
//	}
//
//	var success = func(interface{}) {
//		fmt.Printf("success")
//	}
//
//	var fail = func(err error) {
//		t.Fail()
//	}
//
//	envelope := &pb.Envelope{}
//	envelope.Payload = []byte("hello")
//
//	//then
//	conn.Send(envelope, success, fail)
//
//	time.Sleep(2 * time.Second)
//}
//
//func TestStreamHandler_Close(t *testing.T) {
//
//	//when
//	closed := false
//
//	var closeHandler = func() {
//		//should be closed
//		closed = true
//	}
//
//	var connectionHandler = func(stream pb.StreamService_StreamServer) {
//		//result
//		fmt.Print("connected")
//	}
//
//	serverIP := "127.0.0.1:9999"
//	mockServer := &MockServer{ch: connectionHandler, clh: closeHandler}
//	server1, listner1 := ListenMockServer(mockServer, serverIP)
//
//	//address := Address{Ip: serverIP}
//	grpc_conn, err := NewClientConn(serverIP, false, nil)
//
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	defer func() {
//		server1.Stop()
//		listner1.Close()
//	}()
//
//	//need to wait to connect
//	time.Sleep(1 * time.Second)
//	streamWrapper, err := NewClientStreamWrapper(grpc_conn)
//
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//
//	time.Sleep(1 * time.Second)
//	conn, err := NewConnection(ConnInfo{}, streamWrapper, Handler{})
//
//	//then
//	conn.Close()
//	//result
//	time.Sleep(3 * time.Second)
//	assert.Equal(t, closed, true)
//}
