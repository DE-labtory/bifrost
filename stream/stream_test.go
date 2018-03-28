package stream

import (
	"testing"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

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

	address := conn.Address{IP: serverIP}
	grpc_conn, _ := conn.NewConnectionWithAddress(address, false, nil)

	//then
	_, err := Connect(grpc_conn, Handler{})

	if err != nil {

	}

	assert.Equal(t, connectionFlag, true)
}
