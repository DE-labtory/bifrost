// This file provides authentication interface.

package auth

import (
	"github.com/it-chain/heimdall/key"
)

// Auth represents authentication including sign and verify process.
type Auth interface {
	// Sign signs a digest(hash) using priKey(private key), and returns signature.
	Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error)

	// Verify verifies the signature using pubKey(public key) and digest of original message, then returns boolean value.
	Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error)
}
