// This file provides Sign and Verify functions for RSA.

package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"

	"github.com/it-chain/heimdall/key"
)

// RSASigner represents subject of RSA signing process.
type RSASigner struct{}

// Sign signs a digest(hash) using priKey(private key), and returns signature.
func (s *RSASigner) Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error) {

	if opts == nil {
		return nil, errors.New("invalid options")
	}

	return priKey.(*key.RSAPrivateKey).PrivKey.Sign(rand.Reader, digest, opts)

}

//RSAVerifier represents subject of RSA verifying process.
type RSAVerifier struct{}

// Verify verifies the signature using pubKey(public key) and digest of original message, then returns boolean value.
func (v *RSAVerifier) Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if opts == nil {
		return false, errors.New("invalid options")
	}

	switch opts.(type) {
	case *rsa.PSSOptions:
		err := rsa.VerifyPSS(pubKey.(*key.RSAPublicKey).PubKey,
			(opts.(*rsa.PSSOptions)).Hash,
			digest, signature, opts.(*rsa.PSSOptions))

		if err != nil {
			return false, errors.New("verifying error occurred")
		}

		return true, nil
	default:
		return false, errors.New("invalid options")
	}

}
