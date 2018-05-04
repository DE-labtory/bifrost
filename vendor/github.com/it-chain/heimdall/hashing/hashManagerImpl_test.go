package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashManager_Hash(t *testing.T) {

	rawData := []byte("This data will be hashed by hashManager")

	// normal case
	digest, err := Hash(rawData, nil, SHA512)
	assert.NoError(t, err)
	assert.NotNil(t, digest)

	// compare between hashed data by the same hash function
	anotherDigest, err := Hash(rawData, nil, SHA512)
	assert.Equal(t, digest, anotherDigest)

	// compare between hashed data by the different hash function
	anotherDigest, err = Hash(rawData, nil, SHA256)
	assert.NotEqual(t, digest, anotherDigest)

}
