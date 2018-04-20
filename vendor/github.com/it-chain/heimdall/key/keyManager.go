// This file provides interface of key manager that used to key generation, getting key, and removing key.

package key

// KeyManager represents key manager.
type KeyManager interface {
	// GenerateKey generates public and private key pair that matches the input key generation option.
	GenerateKey(opts KeyGenOpts) (pri PriKey, pub PubKey, err error)

	// GetKey gets the key pair from keyManagerImpl struct.
	// if the keyManagerImpl doesn't have any key, then get keys from stored key files.
	GetKey() (pri PriKey, pub PubKey, err error)

	// RemoveKey removes key files.
	RemoveKey() error

	// Reconstruct key from bytes.
	ByteToKey(byteKey []byte, keyGenOpt KeyGenOpts, keyType KeyType) (err error)

	// GetPath returns path of key files
	GetPath() string
}
