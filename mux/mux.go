/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mux

import (
	"errors"
	"sync"

	"github.com/it-chain/bifrost/conn"
)

type Protocol string

type HandlerFunc func(message conn.OutterMessage)

type ErrorFunc func(conn conn.Connection, err error)

type Mux struct {
	sync.RWMutex
	RegisterHandled map[Protocol]*Handle
	errorFunc       ErrorFunc
}

type Handle struct {
	handlerFunc HandlerFunc
}

func NewMux() *Mux {
	return &Mux{
		RegisterHandled: make(map[Protocol]*Handle),
	}
}

func (mux *Mux) Handle(protocol Protocol, handler HandlerFunc) error {

	mux.Lock()
	defer mux.Unlock()

	_, ok := mux.RegisterHandled[protocol]

	if ok {
		return errors.New("already exist protocol")
	}

	mux.RegisterHandled[protocol] = &Handle{handler}
	return nil
}

func (mux *Mux) Match(protocol Protocol) HandlerFunc {

	mux.Lock()
	defer mux.Unlock()

	handle, ok := mux.RegisterHandled[protocol]

	if ok {
		return handle.handlerFunc
	}

	return nil
}

func (mux *Mux) ServeRequest(msg conn.OutterMessage) {

	protocol := msg.Envelope.Protocol

	handleFunc := mux.Match(Protocol(protocol))

	if handleFunc != nil {
		handleFunc(msg)
	}
}

func (mux *Mux) ServeError(conn conn.Connection, err error) {

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
