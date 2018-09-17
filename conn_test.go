package bifrost_test

import (
	"testing"
	"time"

	"sync"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func TestNewStreamHandler(t *testing.T) {
	// given
	wg := sync.WaitGroup{}
	wg.Add(1)
	var connectionHandler = func(stream pb.StreamService_BifrostStreamServer) {
		//result
		t.Log("connected")
		wg.Done()
	}

	serverIP := "127.0.0.1:9999"
	mockServer := &bifrost.MockServer{Ch: connectionHandler}
	server1, listener1 := bifrost.ListenMockServer(mockServer, serverIP)
	assert.NotNil(t, server1)
	assert.NotNil(t, listener1)

	time.Sleep(3 * time.Second)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTimeout(3*time.Second))

	grpc_conn, err := grpc.Dial(serverIP, opts...)
	assert.NoError(t, err)

	ctx, _ := context.WithCancel(context.Background())
	streamServiceClient := pb.NewStreamServiceClient(grpc_conn)
	_, err = streamServiceClient.BifrostStream(ctx)
	assert.NoError(t, err)

	defer func() {
		server1.Stop()
		listener1.Close()
		grpc_conn.Close()
	}()

	wg.Wait()
}

func TestGrpcConnection_Send(t *testing.T) {

	keyOpts := bifrost.GetKeyOpts()

	mockStreamWrapper := bifrost.MockStreamWrapper{SendCallBack: func(envelope *pb.Envelope) {

	}}
	crypto := bifrost.GetMockCrypto()

	conn, err := bifrost.NewConnection("127.0.0.1:1234", keyOpts.PubKey, mockStreamWrapper, crypto)

	mockStreamWrapper.SendCallBack = func(envelope *pb.Envelope) {
		//then
		assert.Equal(t, envelope.Protocol, "test1")
		assert.Equal(t, envelope.Payload, []byte("jun"))
		assert.True(t, conn.(*bifrost.GrpcConnection).Verify(envelope))
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
	keyOpts := bifrost.GetKeyOpts()
	mockStreamWrapper := bifrost.MockStreamWrapper{}
	crypto := bifrost.GetMockCrypto()

	conn, err := bifrost.NewConnection("127.0.0.1:1234", keyOpts.PubKey, mockStreamWrapper, crypto)
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	//given
	keyOpts := bifrost.GetKeyOpts()

	mockStreamWrapper := bifrost.MockStreamWrapper{}
	mockStreamWrapper.CloseCallBack = func() {
		// then
		assert.True(t, true)
		wg.Done()
	}
	crypto := bifrost.GetMockCrypto()

	conn, err := bifrost.NewConnection("127.0.0.1:1234", keyOpts.PubKey, mockStreamWrapper, crypto)
	assert.NoError(t, err)

	go func() {
		if err := conn.Start(); err != nil {
			assert.NotNil(t, err)
		}
	}()

	// when
	conn.Close()
	wg.Wait()
}

func TestGrpcConnection_GetIP(t *testing.T) {
	keyOpts := bifrost.GetKeyOpts()

	mockStreamWrapper := bifrost.MockStreamWrapper{}
	crypto := bifrost.GetMockCrypto()

	conn, err := bifrost.NewConnection("127.0.0.1:1234", keyOpts.PubKey, mockStreamWrapper, crypto)
	assert.NoError(t, err)
	ipAddr := conn.GetIP()
	assert.Equal(t, bifrost.Address{IP: "127.0.0.1:1234"}, ipAddr)
}
