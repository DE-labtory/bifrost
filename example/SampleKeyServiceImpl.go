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

package example

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"encoding/asn1"
	"io"
	"math/big"

	"encoding/json"

	"io/ioutil"
	"os"

	"path/filepath"

	"errors"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/logger"
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
	PriKey *ecdsa.PrivateKey
}

func (signer *SimpleSigner) Sign(message []byte) ([]byte, error) {
	hash := sha512.New384()
	hash.Write(message)
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

type SimpleVerifier struct {
}

func (verifier *SimpleVerifier) Verify(pubKey *ecdsa.PublicKey, signature, message []byte) (bool, error) {
	hash := sha512.New384()
	hash.Write(message)
	digest := hash.Sum(nil)

	ecdsaSig := new(struct{ R, S *big.Int })
	rest, err := asn1.Unmarshal(signature, ecdsaSig)
	if err != nil {
		logger.Infof(nil, "Verify - unmarshal error")
		return false, err
	}

	if len(rest) != 0 {
		logger.Infof(nil, "Verify - rest error")
		return false, errors.New("invalid values follow signature")
	}

	return ecdsa.Verify(pubKey, digest, ecdsaSig.R, ecdsaSig.S), nil
}

type SimpleIdGetter struct {
	IDPrefix  string
	Formatter bifrost.Formatter
}

func (idGetter *SimpleIdGetter) GetID(pubKey *ecdsa.PublicKey) bifrost.ConnID {
	hash := sha512.New512_256()
	hash.Write(idGetter.Formatter.ToByte(pubKey))

	return bifrost.ConnID(idGetter.IDPrefix + base58.Encode(hash.Sum(nil)))
}

type SimpleKeystore struct {
	keyDirPath string
}

func (keystore *SimpleKeystore) StoreKey(priKey *ecdsa.PrivateKey, pwd string) error {
	priKeyBytes, err := json.Marshal(priKey)
	if err != nil {
		return err
	}

	keyFilePath := filepath.Join(keystore.keyDirPath, "simple_key")
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		err = ioutil.WriteFile(keyFilePath, priKeyBytes, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keystore *SimpleKeystore) LoadKey(keyID string, pwd string) (*ecdsa.PrivateKey, error) {
	var keyPath string

	files, err := ioutil.ReadDir(keystore.keyDirPath)
	if err != nil {
		return nil, errors.New("invalid keystore path - failed to read directory path")
	}

	for _, file := range files {
		if strings.Compare(file.Name(), keyID) == 0 {
			keyPath = filepath.Join(keystore.keyDirPath, file.Name())
			break
		}
	}

	jsonKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	loadedKey := new(ecdsa.PrivateKey)
	if err = json.Unmarshal(jsonKeyBytes, loadedKey); err != nil {
		return nil, err
	}

	return loadedKey, nil
}
