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

	"github.com/it-chain/bifrost"
	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	//when
	testMux := New()

	//then
	testMux.Handle(Protocol("test1"), func(message bifrost.Message) {

	})

	err := testMux.Handle(Protocol("test1"), func(message bifrost.Message) {

	})

	testMux.Handle(Protocol("test3"), func(message bifrost.Message) {

	})

	//result
	assert.Error(t, err, "Asd")
	assert.Equal(t, len(testMux.registerHandled), 2)
}

func TestMux_Handle(t *testing.T) {
	// given
	testMux := New()

	// when
	err := testMux.Handle(Protocol("exist"), func(message bifrost.Message) {

	})
	assert.NoError(t, err)

	hf := testMux.match(Protocol("exist"))
	hf2 := testMux.match(Protocol("do not exist"))

	// then
	assert.NotNil(t, hf)
	assert.Nil(t, hf2)
}

func TestMux_HandleError(t *testing.T) {
	// given
	testMux := New()

	// when
	testMux.HandleError(func(conn bifrost.Connection, err error) {

	})

	// then
	assert.NotNil(t, testMux.errorFunc)
}
