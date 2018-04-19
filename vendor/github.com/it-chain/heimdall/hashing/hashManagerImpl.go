// This file implements hash manager for hashing process.

package hashing

import (
	"crypto/sha1"
	"crypto/sha512"
	"errors"
	"hash"
)

// hashManagerImpl represents a instance of hashManager.
type hashManagerImpl struct{}

// NewHashManager initialize a hash manager.
func NewHashManager() (HashManager, error) {

	hm := &hashManagerImpl{}

	return hm, nil

}

// Hash hashes the input data.
func (hm *hashManagerImpl) Hash(data []byte, tail []byte, opts HashOpts) ([]byte, error) {

	if data == nil {
		return nil, errors.New("Data should not be NIL")
	}

	var hash hash.Hash

	switch opts {
	case SHA1:
		hash = sha1.New()
	case SHA224:
		hash = sha512.New512_224()
	case SHA256:
		hash = sha512.New512_256()
	case SHA384:
		hash = sha512.New384()
	case SHA512:
		hash = sha512.New()
	default:
		return nil, errors.New("Invalid hash opts")
	}

	hash.Write(data)
	return hash.Sum(tail), nil

}
