package hashing

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewHashManager(t *testing.T) {

	hashManager, err := NewHashManager()
	assert.NoError(t, err)
	assert.NotNil(t, hashManager)

}

func TestHashManager_Hash(t *testing.T) {

	hashManager, err := NewHashManager()
	assert.NoError(t, err)
	assert.NotNil(t, hashManager)

	rawData := []byte("This data will be hashed by hashManager")

	// normal case
	digest, err := hashManager.Hash(rawData, nil, SHA512)
	assert.NoError(t, err)
	assert.NotNil(t, digest)

	// compare between hashed data by the same hash function
	anotherDigest, err := hashManager.Hash(rawData, nil, SHA512)
	assert.Equal(t, digest, anotherDigest)

	// compare between hashed data by the different hash function
	anotherDigest, err = hashManager.Hash(rawData, nil, SHA256)
	assert.NotEqual(t, digest, anotherDigest)

}