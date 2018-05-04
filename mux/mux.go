package mux

import (
	"errors"
	"sync"

	"github.com/it-chain/bifrost"
)

type Protocol string

type HandlerFunc func(message bifrost.Message)

type ErrorFunc func(conn bifrost.Connection, err error)

type Mux struct {
	sync.RWMutex
	registerHandled map[Protocol]*Handle
	errorFunc       ErrorFunc
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

	mux.Lock()
	defer mux.Unlock()

	handle, ok := mux.registerHandled[protocol]

	if ok {
		return handle.handlerFunc
	}

	return nil
}

func (mux *Mux) ServeRequest(msg bifrost.Message) {

	protocol := msg.Envelope.Protocol

	handleFunc := mux.match(Protocol(protocol))

	if handleFunc != nil {
		handleFunc(msg)
	}
}

func (mux *Mux) ServeError(conn bifrost.Connection, err error) {

	mux.Lock()
	defer mux.Unlock()

	if mux.errorFunc != nil {
		mux.errorFunc(conn, err)
	}
}

func (mux *Mux) HandleError(errorfunc ErrorFunc) {

	mux.Lock()
	defer mux.Unlock()

	mux.errorFunc = errorfunc
}
