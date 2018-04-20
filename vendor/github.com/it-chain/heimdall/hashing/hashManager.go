// This file provides interface of hash function manager.

package hashing

// HashManager represents manager of hash function.
type HashManager interface {

	// Hash hashes the input data.
	Hash(data []byte, tail []byte, opts HashOpts) ([]byte, error)
}
