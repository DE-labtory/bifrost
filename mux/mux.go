package mux

import (
	"errors"
	"sync"

	"github.com/DE-labtory/bifrost"
)

type Protocol string

type HandlerFunc func(message bifrost.Message)

type ErrorFunc func(conn bifrost.Connection, err error)

type DefaultMux struct {
	sync.RWMutex
	registerHandled map[Protocol]*Handle
	errorFunc       ErrorFunc
}

type Handle struct {
	handlerFunc HandlerFunc
}

func New() *DefaultMux {
	return &DefaultMux{
		registerHandled: make(map[Protocol]*Handle),
	}
}

func (mux *DefaultMux) Handle(protocol Protocol, handler HandlerFunc) error {

	mux.Lock()
	defer mux.Unlock()

	_, ok := mux.registerHandled[protocol]

	if ok {
		return errors.New("already exist protocol")
	}

	mux.registerHandled[protocol] = &Handle{handler}
	return nil
}

func (mux *DefaultMux) match(protocol Protocol) HandlerFunc {

	mux.Lock()
	defer mux.Unlock()

	handle, ok := mux.registerHandled[protocol]

	if ok {
		return handle.handlerFunc
	}

	return nil
}

func (mux *DefaultMux) ServeRequest(msg bifrost.Message) {

	protocol := msg.Envelope.Protocol

	handleFunc := mux.match(Protocol(protocol))

	if handleFunc != nil {
		handleFunc(msg)
	}
}

func (mux *DefaultMux) ServeError(conn bifrost.Connection, err error) {

	mux.Lock()
	defer mux.Unlock()

	if mux.errorFunc != nil {
		mux.errorFunc(conn, err)
	}
}

func (mux *DefaultMux) HandleError(errorFunc ErrorFunc) {

	mux.Lock()
	defer mux.Unlock()

	mux.errorFunc = errorFunc
}
