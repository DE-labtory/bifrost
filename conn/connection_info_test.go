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
	"testing"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/heimdall"
	"github.com/stretchr/testify/assert"
)

func TestFromPublicConnInfo(t *testing.T) {

	pri, err := heimdall.GenerateKey(heimdall.SECP384R1)
	pub := &pri.PublicKey

	b := heimdall.PubKeyToBytes(pub)

	pci := conn.PublicConnInfo{}
	pci.Id = "test1"
	pci.Address = conn.Address{IP: "127.0.0.1"}
	pci.Pubkey = b
	pci.CurveOpt = heimdall.CurveToCurveOpt(pub.Curve)

	//when
	connInfo, err := conn.FromPublicConnInfo(pci)

	//then
	assert.NoError(t, err)
	assert.Equal(t, pub, connInfo.PubKey)
	assert.Equal(t, pci.Id, string(connInfo.Id))
	assert.Equal(t, pci.Address, connInfo.Address)
}
