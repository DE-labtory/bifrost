package bifrost

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/it-chain/bifrost/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockConnectionHandler func(stream pb.StreamService_BifrostStreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	rh  MockRecvHandler
	ch  MockConnectionHandler
	clh MockCloseHandler
}

type Handler struct{}

func (h Handler) ServeRequest(message Message) {

}

func (h Handler) ServeError(conn Connection, err error) {

}

func (ms MockServer) BifrostStream(stream pb.StreamService_BifrostStreamServer) error {

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