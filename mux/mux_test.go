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

package mux_test

import (
	"testing"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	//when
	testMux := mux.NewMux()

	//then
	testMux.Handle(mux.Protocol("test1"), func(message conn.OutterMessage) {

	})

	err := testMux.Handle(mux.Protocol("test1"), func(message conn.OutterMessage) {

	})

	testMux.Handle(mux.Protocol("test3"), func(message conn.OutterMessage) {

	})

	//result
	assert.Error(t, err, "Asd")
	assert.Equal(t, len(testMux.RegisterHandled), 2)
}

func TestMux_Handle(t *testing.T) {
	//when
	testMux := mux.NewMux()

	testMux.Handle(mux.Protocol("exist"), func(message conn.OutterMessage) {

	})

	hf := testMux.Match(mux.Protocol("exist"))
	hf2 := testMux.Match(mux.Protocol("do not exist"))

	assert.NotNil(t, hf)
	assert.Nil(t, hf2)
}

func TestMux_ServeRequest(t *testing.T) {

	//when
	testMux := mux.NewMux()

	testMux.Handle(mux.Protocol("exist"), func(message conn.OutterMessage) {
		assert.Equal(t, message.Data, []byte("hello"))
	})

	message := conn.OutterMessage{}
	message.Data = []byte("hello")
	message.Envelope = &pb.Envelope{Protocol: "exist"}

	testMux.ServeRequest(message)
}
