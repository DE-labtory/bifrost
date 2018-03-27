package mux

import (
	"sync"

	"errors"

	"github.com/it-chain/bifrost/msg"
)

type Protocol string

type HandlerFunc func(message msg.OutterMessage)

type Mux struct {
	sync.RWMutex
	registerHandled map[Protocol]*Handle
}

type Handle struct {
	handlerFunc HandlerFunc
}

func NewMux() *Mux {
	return &Mux{
		registerHandled: make(map[Protocol]*Handle),
	}
}

func (mux *Mux) Handle(protocol Protocol, handler HandlerFunc) error {

	mux.Lock()
	defer mux.Unlock()

	_, ok := mux.registerHandled[protocol]

	if ok {
		return errors.New("already exist protocol")
	}

	mux.registerHandled[protocol] = &Handle{handler}
	return nil
}

func (mux *Mux) match(protocol Protocol) HandlerFunc {
	handle, ok := mux.registerHandled[protocol]

	if ok {
		return handle.handlerFunc
	}

	return nil
}

func (mux *Mux) ServeRequest(msg msg.OutterMessage) {
	protocol := msg.Envelope.Protocol

	handleFunc := mux.match(Protocol(protocol))

	if handleFunc != nil {
		handleFunc(msg)
	}
}
