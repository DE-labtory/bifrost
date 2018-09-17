package client_test

import (
	"testing"

	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/bifrost/server"
	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	// given
	clientIP := "127.0.0.1:12345"
	pri, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assert.NoError(t, err)

	clientOpt := client.ClientOpts{
		Ip:     clientIP,
		PubKey: &pri.PublicKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	serverIP := "127.0.0.1:43213"
	s := server.GetServer()
	s.OnConnection(func(connection bifrost.Connection) {
		defer connection.Close()

		if err := connection.Start(); err != nil {
			connection.Close()
		}
	})
	go s.Listen(serverIP)
	time.Sleep(3 * time.Second)

	// when
	testConn, err := client.Dial(serverIP, clientOpt, grpcOpt, bifrost.GetMockCrypto())
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
