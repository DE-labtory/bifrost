package mux_test

import (
	"testing"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMux_ServeRequest(t *testing.T) {

	// given
	testMux := mux.New()

	testMux.Handle(mux.Protocol("exist"), func(message bifrost.Message) {
		// then
		assert.Equal(t, message.Data, []byte("hello"))
	})

	message := bifrost.Message{}
	message.Data = []byte("hello")
	message.Envelope = &pb.Envelope{Protocol: "exist"}

	// when
	testMux.ServeRequest(message)
}

func TestMux_ServeError(t *testing.T) {
	// given
	targetIP := "127.0.0.1"
	testMux := mux.New()

	testMux.HandleError(func(conn bifrost.Connection, err error) {
		// then
		assert.Equal(t, err.Error(), "testError")
	})

	conn, err := bifrost.GetMockConnection(targetIP)
	assert.NoError(t, err)

	// when
	testMux.ServeError(conn, errors.New("testError"))
}
