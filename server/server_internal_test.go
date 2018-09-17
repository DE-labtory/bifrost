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

package server

import (
	"testing"

	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestGetKeyOpts(t *testing.T) {
	// when
	keyOpt := GetKeyOpts()

	// then
	assert.NotNil(t, keyOpt.PriKey)
	assert.NotNil(t, keyOpt.PubKey)
}

func TestGetServer(t *testing.T) {
	// when
	s := GetServer()

	// then
	assert.NotNil(t, s.pubKey)
	assert.NotNil(t, s.Crypto)
	assert.Nil(t, s.onConnectionHandler)
	assert.Nil(t, s.onErrorHandler)
	assert.Empty(t, s.ip)
	assert.Nil(t, s.lis)
}

func TestServer_OnConnection(t *testing.T) {
	// given
	s := GetServer()

	// when
	s.OnConnection(func(conn bifrost.Connection) {

	})

	// then
	assert.NotNil(t, s.onConnectionHandler)
}

func TestServer_OnError(t *testing.T) {
	// given
	s := GetServer()

	// when
	s.OnError(func(err error) {

	})

	// then
	assert.NotNil(t, s.onErrorHandler)
}

func TestServer_validateRequestPeerInfo_whenInValidPeerInfo(t *testing.T) {
	// given
	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   []byte("123"),
		CurveOpt: 2,
	}

	mockServer := GetServer()

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	// when
	flag, _, _ := mockServer.validateRequestPeerInfo(envelope)

	// then
	assert.False(t, flag)
}
