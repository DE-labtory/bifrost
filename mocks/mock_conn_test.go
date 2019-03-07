/*
 * Copyright 2018 DE-labtory
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

package mocks_test

import (
	"testing"

	"github.com/DE-labtory/bifrost/mocks"
	"github.com/DE-labtory/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestNewMockConnection(t *testing.T) {
	// when
	conn, err := mocks.NewMockConnection("127.0.0.1:1234")

	// then
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestMockStreamWrapper_Send(t *testing.T) {
	// given
	var isSendCallBackCalled = false
	mockStreamWrapper := mocks.MockStreamWrapper{
		SendCallBack: func(envelope *pb.Envelope) {
			isSendCallBackCalled = true
		},
	}

	// when
	envelope := new(pb.Envelope)
	mockStreamWrapper.Send(envelope)

	// then
	assert.True(t, isSendCallBackCalled)
}

func TestMockStreamWrapper_Close(t *testing.T) {
	// given
	var isCloseCallBackCalled = false
	mockStreamWrapper := mocks.MockStreamWrapper{
		CloseCallBack: func() {
			isCloseCallBackCalled = true
		},
	}

	// when
	mockStreamWrapper.Close()

	// then
	assert.True(t, isCloseCallBackCalled)
}
