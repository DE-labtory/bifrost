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

package mocks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"crypto/sha512"
	"encoding/asn1"

	"crypto/x509"

	"github.com/btcsuite/btcutil/base58"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/iLogger"
)

func NewMockKeyPair() (pri bifrost.Key, pub bifrost.Key, err error) {
	ecPri, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pri = &MockPriKey{internalPriKey: ecPri}
	pub = &MockPubKey{internalPubKey: &ecPri.PublicKey}

	return pri, pub, nil
}

func NewMockCrypto() bifrost.Crypto {
	mockSigner := MockECDSASigner{}
	mockVerifier := MockECDSAVerifier{}
	mockKeyRecoverer := MockECDSAKeyRecoverer{}

	return bifrost.Crypto{
		Signer:       &mockSigner,
		Verifier:     &mockVerifier,
		KeyRecoverer: &mockKeyRecoverer,
	}
}

type MockECDSASigner struct {
	KeyID      string
	KeyDirPath string
}

func (signer *MockECDSASigner) Sign(message []byte) ([]byte, error) {
	priKey, err := mockLoadKey(signer.KeyID, signer.KeyDirPath)
	if err != nil {
		return nil, err
	}

	// get hash from message
	hash := sha512.New384()
	hash.Write(message)
	digest := hash.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, priKey.(*MockPriKey).internalPriKey, digest)
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

type MockECDSAVerifier struct {
}

func (verifier *MockECDSAVerifier) Verify(peerKey bifrost.Key, signature []byte, message []byte) (bool, error) {
	hash := sha512.New384()
	hash.Write(message)
	digest := hash.Sum(nil)

	ecdsaSig := new(struct{ R, S *big.Int })
	rest, err := asn1.Unmarshal(signature, ecdsaSig)
	if err != nil {
		iLogger.Error(nil, "[Bifrost] Verify - unmarshal error")
		return false, err
	}

	if len(rest) != 0 {
		iLogger.Info(nil, "[Bifrost] Verify - rest error")
		return false, errors.New("invalid values follow signature")
	}

	return ecdsa.Verify(peerKey.(*MockPubKey).internalPubKey, digest, ecdsaSig.R, ecdsaSig.S), nil
}

type MockECDSAKeyRecoverer struct {
}

func (recoverer *MockECDSAKeyRecoverer) RecoverKeyFromByte(keyBytes []byte, isPrivateKey bool) (bifrost.Key, error) {
	switch isPrivateKey {
	case true:
		internalPriKey, err := x509.ParseECPrivateKey(keyBytes)
		if err != nil {
			return nil, err
		}

		pri := newMockPriKey(internalPriKey)

		return pri, nil

	case false:
		internalPubKey, err := x509.ParsePKIXPublicKey(keyBytes)
		if err != nil {
			return nil, err
		}

		pub := newMockPubKey(internalPubKey.(*ecdsa.PublicKey))

		return pub, nil
	}

	return nil, errors.New("isPrivateKey value should be true or false")
}

type mockKeyFile struct {
	KeyBytes     []byte
	IsPrivateKey bool
}

func MockStoreKey(key bifrost.Key, keyDirPath string) error {
	// check and make dir if not exist
	if _, err := os.Stat(keyDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(keyDirPath, 0755)
		if err != nil {
			return err
		}
	}

	// full path
	keyFilePath := filepath.Join(keyDirPath, key.ID())

	// byte formatted key
	keyBytes, err := key.ToByte()
	if err != nil {
		return err
	}

	isPrivateKey := key.IsPrivate()

	// make json marshal
	jsonKeyFile, err := json.Marshal(mockKeyFile{KeyBytes: keyBytes, IsPrivateKey: isPrivateKey})
	if err != nil {
		return err
	}

	// store json key file
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		err = ioutil.WriteFile(keyFilePath, jsonKeyFile, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func mockLoadKey(keyID bifrost.KeyID, keyDirPath string) (bifrost.Key, error) {
	var keyFile mockKeyFile

	if _, err := os.Stat(keyDirPath); os.IsNotExist(err) {
		return nil, err
	}

	keyFilePath := filepath.Join(keyDirPath, keyID)

	jsonKeyFile, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(jsonKeyFile, &keyFile); err != nil {
		return nil, err
	}

	mockKeyRecoverer := MockECDSAKeyRecoverer{}
	key, err := mockKeyRecoverer.RecoverKeyFromByte(keyFile.KeyBytes, keyFile.IsPrivateKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}

type MockPriKey struct {
	internalPriKey *ecdsa.PrivateKey
}

func newMockPriKey(internalKey *ecdsa.PrivateKey) *MockPriKey {
	return &MockPriKey{
		internalPriKey: internalKey,
	}
}

func (mockPriKey *MockPriKey) ID() bifrost.KeyID {
	// get keyBytes from key
	mockPub := &mockPriKey.internalPriKey.PublicKey
	keyBitString := elliptic.Marshal(mockPub.Curve, mockPub.X, mockPub.Y)

	// get ski from keyBytes
	hash := sha256.New()
	hash.Write(keyBitString)
	hashValue := hash.Sum(nil)

	ski := make([]byte, 20)
	ski = hashValue[:20]

	return "IT" + base58.Encode(ski)
}

func (mockPriKey *MockPriKey) ToByte() ([]byte, error) {
	return x509.MarshalECPrivateKey(mockPriKey.internalPriKey)
}

func (mockPriKey *MockPriKey) KeyGenOpt() string {
	mockPub := &mockPriKey.internalPriKey.PublicKey
	return mockPub.Curve.Params().Name
}

func (mockPriKey *MockPriKey) IsPrivate() bool {
	return true
}

type MockPubKey struct {
	internalPubKey *ecdsa.PublicKey
}

func newMockPubKey(internalKey *ecdsa.PublicKey) *MockPubKey {
	return &MockPubKey{
		internalPubKey: internalKey,
	}
}

func (mockPubKey *MockPubKey) ID() bifrost.KeyID {
	// get keyBytes from key
	keyBitString := elliptic.Marshal(mockPubKey.internalPubKey.Curve, mockPubKey.internalPubKey.X, mockPubKey.internalPubKey.Y)

	// get ski from keyBytes
	hash := sha256.New()
	hash.Write(keyBitString)
	hashValue := hash.Sum(nil)

	ski := make([]byte, 20)
	ski = hashValue[:20]

	return "IT" + base58.Encode(ski)
}

func (mockPubKey *MockPubKey) ToByte() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(mockPubKey.internalPubKey)
}

func (mockPubKey *MockPubKey) KeyGenOpt() string {
	return mockPubKey.internalPubKey.Curve.Params().Name
}

func (mockPubKey *MockPubKey) IsPrivate() bool {
	return false
}
