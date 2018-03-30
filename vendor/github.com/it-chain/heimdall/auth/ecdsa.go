package auth

import (
	"math/big"
	"encoding/asn1"
	"crypto/ecdsa"
	"errors"
	"crypto/rand"
	"github.com/it-chain/heimdall/key"
)

type ecdsaSignature struct {
	R, S *big.Int
}

type ECDSASigner struct {}

func marshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ecdsaSignature{r, s})
}

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

type ECDSAVerifier struct {}

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