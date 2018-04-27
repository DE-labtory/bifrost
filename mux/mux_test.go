package mux

import (
	"testing"

	"github.com/it-chain/bifrost/grpc"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	//when
	mux := NewMux()

	//then
	mux.Handle(Protocol("test1"), func(message conn.OutterMessage) {

	})

	err := mux.Handle(Protocol("test1"), func(message conn.OutterMessage) {

	})

	mux.Handle(Protocol("test3"), func(message conn.OutterMessage) {

	})

	//result
	assert.Error(t, err, "Asd")
	assert.Equal(t, len(mux.registerHandled), 2)
}

func TestMux_Handle(t *testing.T) {
	//when
	mux := NewMux()

	mux.Handle(Protocol("exist"), func(message conn.OutterMessage) {

	})

	hf := mux.match(Protocol("exist"))
	hf2 := mux.match(Protocol("do not exist"))

	assert.NotNil(t, hf)
	assert.Nil(t, hf2)
}

func TestMux_ServeRequest(t *testing.T) {

	//when
	mux := NewMux()

	mux.Handle(Protocol("exist"), func(message conn.OutterMessage) {
		assert.Equal(t, message.Data, []byte("hello"))
	})

	message := conn.OutterMessage{}
	message.Data = []byte("hello")
	message.Envelope = &pb.Envelope{Protocol: "exist"}

	mux.ServeRequest(message)
}
