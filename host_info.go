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

package bifrost

import (
	"crypto/ecdsa"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/heimdall"
)

//Identitiy of Connection
type ID string

//Create ID from Public Key
func FromPubKey(key *ecdsa.PublicKey) ID {
	return ID(heimdall.PubKeyToKeyID(key))
}

func (id ID) String() string {
	return string(id)
}

type HostInfo struct {
	conn.ConnInfo
	PriKey *ecdsa.PrivateKey
}

func NewHostInfo(address conn.Address, pubKey *ecdsa.PublicKey, priKey *ecdsa.PrivateKey) HostInfo {

	id := FromPubKey(pubKey)

	return HostInfo{
		ConnInfo: conn.NewConnInfo(id.String(), address, pubKey),
		PriKey:   priKey,
	}
}

func (hostInfo HostInfo) GetPublicInfo() *conn.PublicConnInfo {

	publicConnInfo := &conn.PublicConnInfo{}
	publicConnInfo.Id = hostInfo.Id.ToString()
	publicConnInfo.Address = hostInfo.Address

	b := heimdall.PubKeyToBytes(hostInfo.PubKey)

	publicConnInfo.Pubkey = b
	publicConnInfo.CurveOpt = heimdall.CurveToCurveOpt(hostInfo.PubKey.Curve)

	return publicConnInfo
}
