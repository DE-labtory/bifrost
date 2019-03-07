package client_test

import (
	"testing"

	"time"

	"github.com/DE-labtory/bifrost"
	"github.com/DE-labtory/bifrost/client"
	"github.com/DE-labtory/bifrost/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	// given
	clientIP := "127.0.0.1:12345"
	keyPair := mocks.NewMockKeyOpts()

	clientOpt := client.ClientOpts{
		Ip:     clientIP,
		PubKey: keyPair.PubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	serverIP := "127.0.0.1:43213"
	s := mocks.NewMockServer()
	s.OnConnection(func(connection bifrost.Connection) {
		defer connection.Close()

		if err := connection.Start(); err != nil {
			connection.Close()
		}
	})
	go s.Listen(serverIP)
	time.Sleep(3 * time.Second)

	// when
	testConn, err := client.Dial(serverIP, nil, clientOpt, grpcOpt, mocks.NewMockCrypto())
	go func() {
		defer testConn.Close()
		if err := testConn.Start(); err != nil {
			testConn.Close()
		}
	}()

	// then
	assert.NoError(t, err)
	assert.Equal(t, testConn.GetIP(), bifrost.Address{serverIP})
}
