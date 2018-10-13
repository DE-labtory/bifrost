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
	"encoding/json"
	"testing"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestServer_OnConnection(t *testing.T) {
	// given
	keyOpt := bifrost.KeyOpts{
		PubKey: *new(bifrost.Key),
		PriKey: *new(bifrost.Key),
	}

	crypto := bifrost.Crypto{}

	mockServer := New(keyOpt, crypto, nil)

	// when
	mockServer.OnConnection(func(conn bifrost.Connection) {

	})

	// then
	assert.NotNil(t, mockServer.onConnectionHandler)
}

func TestServer_OnError(t *testing.T) {
	// given
	keyOpt := bifrost.KeyOpts{
		PubKey: *new(bifrost.Key),
		PriKey: *new(bifrost.Key),
	}

	crypto := bifrost.Crypto{}

	mockServer := New(keyOpt, crypto, nil)

	// when
	mockServer.OnError(func(err error) {

	})

	// then
	assert.NotNil(t, mockServer.onErrorHandler)
}

func TestServer_validateRequestPeerInfo_whenInValidPeerInfo(t *testing.T) {
	// given
	peerInfo := &bifrost.PeerInfo{
		IP:          "127.0.0.1",
		PubKeyBytes: []byte("123"),
		IsPrivate:   false,
		KeyGenOpt:   "P-384",
	}

	keyOpt := bifrost.KeyOpts{
		PubKey: *new(bifrost.Key),
		PriKey: *new(bifrost.Key),
	}

	crypto := bifrost.Crypto{}

	mockServer := New(keyOpt, crypto, nil)

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	// when
	flag, _, _ := mockServer.validateRequestPeerInfo(envelope)

	// then
	assert.False(t, flag)
}
