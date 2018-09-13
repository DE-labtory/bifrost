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
)

type Generator interface {
	GenerateKey() (*ecdsa.PrivateKey, error)
}

type Keystore interface {
	StoreKey(priKey *ecdsa.PrivateKey, pwd string) error
	LoadKey(keyID string, pwd string) (*ecdsa.PrivateKey, error)
}

type IDGetter interface {
	GetID(*ecdsa.PublicKey) ConnID
}

type Signer interface {
	Sign(message []byte) ([]byte, error)
}

type Verifier interface {
	Verify(pubKey *ecdsa.PublicKey, signature, message []byte) (bool, error)
}

type Formatter interface {
	FromByte(pubKey []byte, curveOpt int) *ecdsa.PublicKey
	ToByte(pubKey *ecdsa.PublicKey) []byte
	GetCurveOpt(pubKey *ecdsa.PublicKey) int
}
