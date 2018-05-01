package key

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyManager(t *testing.T) {
	testKeyManager, err := NewKeyManager("./.testKeys")
	assert.NoError(t, err)
	assert.NotNil(t, testKeyManager)
	defer os.RemoveAll("./.testKeys")
}

func TestKeyManagerImpl_GenerateKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	pri, pub, err := testKeyManager.GenerateKey(keyGenOption)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	defer os.RemoveAll("./.testKeys")
}

func TestKeyManagerImpl_GetKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	testKeyManager.GenerateKey(keyGenOption)

	pri, pub, err := testKeyManager.GetKey()
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	defer os.RemoveAll("./.testKeys")
}

func TestKeyManagerImpl_RemoveKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	testKeyManager.GenerateKey(keyGenOption)

	err := testKeyManager.RemoveKey()
	assert.NoError(t, err)

	defer os.RemoveAll("./.testKeys")
}

func TestKeyManagerImpl_GetPath(t *testing.T) {
	testKeyManager, _ := NewKeyManager("./.testKeys")
	originPath := "./.testKeys/.keys"
	path := testKeyManager.GetPath()
	assert.Equal(t, originPath, path, "They should be equal")

	defer os.RemoveAll("./.testKeys")
}
