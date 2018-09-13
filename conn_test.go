package bifrost

import (
	"fmt"
	"log"
	"testing"
	"time"

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
	mockServer := &MockServer{Ch: connectionHandler}
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
		grpc_conn.Close()
	}()

	time.Sleep(3 * time.Second)
}

func TestGrpcConnection_Send(t *testing.T) {

	keyOpts := GetKeyOpts()

	mockStreamWrapper := MockStreamWrapper{sendCallBack: func(envelope *pb.Envelope) {

	}}
	mockIDGetter := MockIdGetter{}
	mockFormatter := MockFormatter{}
	mockSigner := MockSigner{}
	mockVerifier := MockVerifier{}

	conn, err := NewConnection("127.0.0.1", keyOpts.PriKey, keyOpts.PubKey, mockStreamWrapper, &mockIDGetter, &mockFormatter, &mockSigner, &mockVerifier)

	mockStreamWrapper.sendCallBack = func(envelope *pb.Envelope) {

		//then
		assert.Equal(t, envelope.Protocol, "test1")
		assert.Equal(t, envelope.Payload, []byte("jun"))
		assert.True(t, conn.(*GrpcConnection).Verify(envelope, keyOpts.PubKey))
	}

	assert.NoError(t, err)

	go func() {
		if err := conn.Start(); err != nil {
			conn.Close()
		}
	}()

	//when
	conn.Send([]byte("jun"), "test1", nil, nil)
}

func TestGrpcConnection_GetPeerKey(t *testing.T) {

	//given
	keyOpts := GetKeyOpts()

	mockStreamWrapper := MockStreamWrapper{}
	mockIDGetter := MockIdGetter{}
	mockFormatter := MockFormatter{}
	mockSigner := MockSigner{}
	mockVerifier := MockVerifier{}

	conn, err := NewConnection("127.0.0.1", keyOpts.PriKey, keyOpts.PubKey, mockStreamWrapper, &mockIDGetter, &mockFormatter, &mockSigner, &mockVerifier)

	assert.NoError(t, err)

	go func() {
		if err := conn.Start(); err != nil {
			conn.Close()
		}
	}()

	//when
	k := conn.GetPeerKey()

	//then
	assert.Equal(t, k, keyOpts.PubKey)
}

func TestGrpcConnection_Close(t *testing.T) {

	//given
	keyOpts := GetKeyOpts()

	mockStreamWrapper := MockStreamWrapper{}
	mockStreamWrapper.closeCallBack = func() {
		assert.True(t, true)
	}
	mockIDGetter := MockIdGetter{}
	mockFormatter := MockFormatter{}
	mockSigner := MockSigner{}
	mockVerifier := MockVerifier{}

	conn, err := NewConnection("127.0.0.1", keyOpts.PriKey, keyOpts.PubKey, mockStreamWrapper, &mockIDGetter, &mockFormatter, &mockSigner, &mockVerifier)

	assert.NoError(t, err)

	go func() {
		if err := conn.Start(); err != nil {
			assert.NotNil(t, err)
		}
	}()

	conn.Close()
}
