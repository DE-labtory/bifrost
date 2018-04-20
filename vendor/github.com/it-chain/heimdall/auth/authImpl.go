// This file provides implementation of the authentication interface.

package auth

import (
	"errors"
	"reflect"

	"github.com/it-chain/heimdall/key"
)

// authImpl contains signers and verifiers that is used for signing and verifying process.
type authImpl struct {
	signers   map[reflect.Type]signer
	verifiers map[reflect.Type]verifier
}

// NewAuth initialize the struct authImpl.
func NewAuth() (Auth, error) {

	signers := make(map[reflect.Type]signer)
	signers[reflect.TypeOf(&key.RSAPrivateKey{})] = &RSASigner{}
	signers[reflect.TypeOf(&key.ECDSAPrivateKey{})] = &ECDSASigner{}

	verifiers := make(map[reflect.Type]verifier)
	verifiers[reflect.TypeOf(&key.RSAPublicKey{})] = &RSAVerifier{}
	verifiers[reflect.TypeOf(&key.ECDSAPublicKey{})] = &ECDSAVerifier{}

	ai := &authImpl{
		signers:   signers,
		verifiers: verifiers,
	}

	return ai, nil

}

// Sign signs a digest(hash) using priKey(private key), and returns signature.
func (ai *authImpl) Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error) {

	var err error

	if len(digest) == 0 {
		return nil, errors.New("invalid data.")
	}

	if priKey == nil {
		return nil, errors.New("Private key is not exist.")
	}

	signer, found := ai.signers[reflect.TypeOf(priKey)]
	if !found {
		return nil, errors.New("unsupported key type.")
	}

	signature, err := signer.Sign(priKey, digest, opts)
	if err != nil {
		return nil, errors.New("signing error is occurred")
	}

	return signature, err

}

// Verify verifies the signature using pubKey(public key) and digest of original message, then returns boolean value.
func (ai *authImpl) Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if pubKey == nil {
		return false, errors.New("invalid key")
	}

	if len(signature) == 0 {
		return false, errors.New("invalid signature")
	}

	if len(digest) == 0 {
		return false, errors.New("invalid digest")
	}

	verifier, found := ai.verifiers[reflect.TypeOf(pubKey)]
	if !found {
		return false, errors.New("unsupported key type")
	}

	valid, err := verifier.Verify(pubKey, signature, digest, opts)
	if err != nil {
		return false, errors.New("verifying error is occurred")
	}

	return valid, nil

}
