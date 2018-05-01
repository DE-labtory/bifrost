// This file provides interfaces of Key, Private key and Public key.

package key

// A KeyType represents key type such as 'public' or 'private'.
type KeyType string

const (
	PRIVATE_KEY KeyType = "pri"
	PUBLIC_KEY  KeyType = "pub"
)

// A Key represents a cryptographic key.
type Key interface {
	// SKI provides name of file that will be store a key
	SKI() (ski []byte)

	// Algorithm returns key generation option such as 'rsa2048'.
	Algorithm() KeyGenOpts

	// ToPEM makes a key to PEM format.
	ToPEM() ([]byte, error)

	// Type returns type of the key.
	Type() KeyType
}

// PriKey represents a private key by implementing Key interface.
type PriKey interface {
	Key

	// PublicKey returns public key of key pair
	PublicKey() (PubKey, error)
}

// PubKey represents a public key by implementing Key interface.
type PubKey interface {
	Key
}
