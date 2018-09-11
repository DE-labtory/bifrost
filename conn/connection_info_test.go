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

package conn_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/it-chain/bifrost/conn"
	"github.com/stretchr/testify/assert"
)

func TestID_ToString(t *testing.T) {
	id := conn.ID("testID")
	strID := id.ToString()
	assert.Equal(t, strID, "testID")
}

func TestNewConnInfo(t *testing.T) {
	id := "testID"
	address := conn.Address{IP: "127.0.0.1:1111"}
	pubKey := new(ecdsa.PublicKey)

	connInfo := conn.NewConnInfo(id, address, pubKey)
	assert.Equal(t, connInfo.Id, conn.ID(id))
	assert.Equal(t, connInfo.Address, address)
	assert.Equal(t, connInfo.PubKey, pubKey)
}

func TestToAddress(t *testing.T) {
	ipv4Addr := "192.168.0.0:1234"
	addrNoPort := "127.0.0.114324"
	addrTooManyNumInIP := "127.0.0.114324"
	addrTooManyNumInPort := "127.0.0.1:555555"

	address, err := conn.ToAddress(ipv4Addr)
	assert.NoError(t, err)
	assert.Equal(t, ipv4Addr, address.IP)

	address, err = conn.ToAddress(addrNoPort)
	assert.Error(t, err)

	address, err = conn.ToAddress(addrTooManyNumInIP)
	assert.Error(t, err)

	address, err = conn.ToAddress(addrTooManyNumInPort)
	assert.Error(t, err)
}
