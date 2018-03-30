package auth

import (
	"crypto"
	"github.com/it-chain/heimdall/key"
)

type SignerOpts interface {
	crypto.SignerOpts
}

type signer interface {

	Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error)

}

type verifier interface {

	Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error)

}