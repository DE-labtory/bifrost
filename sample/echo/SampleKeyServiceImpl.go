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

package echo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"encoding/asn1"
	"io"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/it-chain/bifrost"
)

type SimpleFormatter struct {
}

func (formatter *SimpleFormatter) ToByte(pubKey *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
}

func (formatter *SimpleFormatter) FromByte(pubKey []byte, curveOpt int) *ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(elliptic.P384(), pubKey)

	return &ecdsa.PublicKey{X: x, Y: y, Curve: elliptic.P384()}
}

func (formatter *SimpleFormatter) GetCurveOpt(pubKey *ecdsa.PublicKey) int {
	switch pubKey.Curve {
	case elliptic.P224():
		return 0
	case elliptic.P256():
		return 1
	case elliptic.P384():
		return 2
	case elliptic.P521():
		return 3
	default:
		return 4
	}
}

type SimpleGenerator struct {
	Curve elliptic.Curve
	Rand  io.Reader
}

func (generator SimpleGenerator) GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(generator.Curve, generator.Rand)
}

type SimpleSigner struct {
	Message []byte
	PriKey  *ecdsa.PrivateKey
}

func (signer *SimpleSigner) Sign() ([]byte, error) {
	hash := sha512.New384()
	hash.Write(signer.Message)
	digest := hash.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, signer.PriKey, digest)
	if err != nil {
		return nil, err
	}

	signature, err := asn1.Marshal(struct {
		R, S *big.Int
	}{r, s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

type SimpleIdGetter struct {
	IDPrefix   string
	PubKeyByte []byte
}

func (idGetter *SimpleIdGetter) GetID(key *ecdsa.PublicKey) bifrost.ID {
	hash := sha512.New512_256()
	hash.Write(idGetter.PubKeyByte)

	return bifrost.ID(idGetter.IDPrefix + base58.Encode(hash.Sum(nil)))
}
