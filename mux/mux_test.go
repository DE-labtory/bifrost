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
	"testing"

	"github.com/it-chain/bifrost/conn"
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
