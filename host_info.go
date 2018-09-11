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
)

//Identitiy of Connection
type ID string

func (id ID) String() string {
	return string(id)
}

type HostInfo struct {
	conn.ConnInfo
	PriKey *ecdsa.PrivateKey
}

func NewHostInfo(address conn.Address, priKey *ecdsa.PrivateKey, idGetter IDGetter) HostInfo {
	id := idGetter.GetID(&priKey.PublicKey)

	return HostInfo{
		ConnInfo: conn.NewConnInfo(id.String(), address, &priKey.PublicKey),
		PriKey:   priKey,
	}
}
