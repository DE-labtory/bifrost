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

	"io/ioutil"
	"path/filepath"

	"errors"
	"strings"

	"os"

	"crypto/x509"

	"github.com/btcsuite/btcutil/base58"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/iLogger"
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

func (generator *SimpleGenerator) GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(generator.Curve, generator.Rand)
}

func (generator *SimpleGenerator) StoreKey(priKey *ecdsa.PrivateKey, pwd string, keyDirPath string, keyID bifrost.KeyID) error {
	priKeyBytes, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return err
	}

	if _, err := os.Stat(keyDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(keyDirPath, 0755)
		if err != nil {
			return err
		}
	}

	keyFilePath := filepath.Join(keyDirPath, string(keyID))
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		err = ioutil.WriteFile(keyFilePath, priKeyBytes, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

type SimpleSigner struct {
	bifrost.KeyLoader
}

func (signer *SimpleSigner) Sign(message []byte) ([]byte, error) {
	hash := sha512.New384()
	hash.Write(message)
	digest := hash.Sum(nil)

	priKey, err := signer.LoadKey("")
	if err != nil {
		iLogger.Infof(nil, "[Bifrost] Load error %s", err.Error())
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, priKey, digest)
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
	bifrost.KeyLoader
}

func (verifier *SimpleVerifier) Verify(peerKey *ecdsa.PublicKey, signature, message []byte) (bool, error) {

	hash := sha512.New384()
	hash.Write(message)
	digest := hash.Sum(nil)

	ecdsaSig := new(struct{ R, S *big.Int })
	rest, err := asn1.Unmarshal(signature, ecdsaSig)
	if err != nil {
		iLogger.Info(nil, "[Bifrost] Verify - unmarshal error")
		return false, err
	}

	if len(rest) != 0 {
		iLogger.Info(nil, "[Bifrost] Verify - rest error")
		return false, errors.New("invalid values follow signature")
	}

	return ecdsa.Verify(peerKey, digest, ecdsaSig.R, ecdsaSig.S), nil
}

type SimpleIdGetter struct {
	IDPrefix string
	bifrost.Formatter
}

func (idGetter *SimpleIdGetter) GetID(pubKey *ecdsa.PublicKey) bifrost.KeyID {
	hash := sha512.New512_256()
	hash.Write(idGetter.ToByte(pubKey))

	return bifrost.KeyID(idGetter.IDPrefix + base58.Encode(hash.Sum(nil)))
}

type SimpleKeyLoader struct {
	KeyDirPath string
	KeyID      bifrost.KeyID
}

func (keyLoader *SimpleKeyLoader) LoadKey(pwd string) (*ecdsa.PrivateKey, error) {
	var keyPath string
	var loadedKey *ecdsa.PrivateKey

	files, err := ioutil.ReadDir(keyLoader.KeyDirPath)
	if err != nil {
		return nil, errors.New("invalid keystore path - failed to read directory path")
	}

	for _, file := range files {
		if strings.Compare(file.Name(), string(keyLoader.KeyID)) == 0 {
			keyPath = filepath.Join(keyLoader.KeyDirPath, file.Name())
			break
		}
	}

	priKeyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	loadedKey, err = x509.ParseECPrivateKey(priKeyBytes)
	if err != nil {
		return nil, err
	}

	return loadedKey, nil
}
