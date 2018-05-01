// This file provides certificate formatting and Sign and Verify functions for ECDSA.

package auth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"errors"
	"math/big"

	"github.com/it-chain/heimdall/key"
)

// ecdsaSignature contains ECDSA signature components that are two integers, R and S.
type ecdsaSignature struct {
	R, S *big.Int
}

//ECDSASigner represents subject of ECDSA signing process.
type ECDSASigner struct{}

// marshalECDSASignature returns encoding format (ASN.1) of signature.
func marshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ecdsaSignature{r, s})
}

// unmarshalECDSASignature parses the ASN.1 structure to ECDSA signature.
func unmarshalECDSASignature(signature []byte) (*big.Int, *big.Int, error) {
	ecdsaSig := new(ecdsaSignature)
	_, err := asn1.Unmarshal(signature, ecdsaSig)
	if err != nil {
		return nil, nil, errors.New("failed to unmarshal")
	}

	if ecdsaSig.R == nil {
		return nil, nil, errors.New("invalid signature")
	}
	if ecdsaSig.S == nil {
		return nil, nil, errors.New("invalid signature")
	}

	if ecdsaSig.R.Sign() != 1 {
		return nil, nil, errors.New("invalid signature")
	}
	if ecdsaSig.S.Sign() != 1 {
		return nil, nil, errors.New("invalid signature")
	}

	return ecdsaSig.R, ecdsaSig.S, nil

}

// Sign signs a digest(hash) using priKey(private key), and returns signature.
func (signer *ECDSASigner) Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error) {

	r, s, err := ecdsa.Sign(rand.Reader, priKey.(*key.ECDSAPrivateKey).PrivKey, digest)
	if err != nil {
		return nil, err
	}

	signature, err := marshalECDSASignature(r, s)
	if err != nil {
		return nil, err
	}

	return signature, nil

}

//ECDSAVerifier represents subject of ECDSA verifying process.
type ECDSAVerifier struct{}

// Verify verifies the signature using pubKey(public key) and digest of original message, then returns boolean value.
func (v *ECDSAVerifier) Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	r, s, err := unmarshalECDSASignature(signature)
	if err != nil {
		return false, err
	}

	valid := ecdsa.Verify(pubKey.(*key.ECDSAPublicKey).PubKey, digest, r, s)
	if !valid {
		return valid, errors.New("failed to verify")
	}

	return valid, nil

}
