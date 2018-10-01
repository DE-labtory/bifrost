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

package mocks_test

import (
	"testing"

	"os"

	"github.com/it-chain/bifrost/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewMockKeyPair(t *testing.T) {
	// when
	pri, pub, err := mocks.NewMockKeyPair()

	// then
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)
}

func TestNewMockCrypto(t *testing.T) {
	// when
	crypto := mocks.NewMockCrypto()

	// then
	assert.NotNil(t, crypto)
}

func TestMockECDSASigner_Sign(t *testing.T) {
	// given
	pri, _, err := mocks.NewMockKeyPair()
	keyDirPath := "./.test_key"
	wrongPath := "./.key"
	err = mocks.MockStoreKey(pri, keyDirPath)
	defer os.RemoveAll(keyDirPath)

	mockSigner := mocks.MockECDSASigner{
		KeyDirPath: keyDirPath,
		KeyID:      pri.ID(),
	}

	invalidSigner := mocks.MockECDSASigner{
		KeyDirPath: wrongPath,
		KeyID:      pri.ID(),
	}

	message := []byte("test message")

	// when
	signature, noErr := mockSigner.Sign(message)
	nilSignature, err := invalidSigner.Sign(message)

	// then
	assert.NoError(t, noErr)
	assert.NotNil(t, signature)

	assert.Error(t, err)
	assert.Nil(t, nilSignature)
}

// todo: 나머지 테스트 코드 작성
