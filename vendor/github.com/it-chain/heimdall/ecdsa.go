package heimdall

import (
	"crypto/ecdsa"
	"math/big"
	"encoding/asn1"
	"errors"
	"crypto/rand"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type EcdsaSignature struct {
	R, S *big.Int
}

type EcdsaSigner struct{}

func marshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(EcdsaSignature{r, s})
}

func unmarshalECDSASignature(signature []byte) (*big.Int, *big.Int, error) {
	ecdsaSig := new(EcdsaSignature)
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

func (signer *EcdsaSigner) Sign(key Key, digest []byte, opts SignerOpts) ([]byte, error) {

	r, s, err := ecdsa.Sign(rand.Reader, key.(*EcdsaPrivateKey).priv, digest)
	if err != nil {
		return nil, err
	}

	signature, err := marshalECDSASignature(r, s)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

type EcdsaVerifier struct{}

func (v *EcdsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	r, s, err := unmarshalECDSASignature(signature)
	if err != nil {
		return false, err
	}

	valid := ecdsa.Verify(key.(*EcdsaPublicKey).pub, digest, r, s)
	if !valid {
		return valid, errors.New("failed to verify")
	}

	return valid, nil
}

type EcdsaPrivateKey struct {
	priv *ecdsa.PrivateKey
}

func (key *EcdsaPrivateKey) SKI() ([]byte) {

	if key.priv == nil {
		return nil
	}

	data := elliptic.Marshal(key.priv.Curve, key.priv.PublicKey.X, key.priv.PublicKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *EcdsaPrivateKey) Algorithm() string {
	return ECDSA
}

func (key *EcdsaPrivateKey) PublicKey() (Key, error) {
	return &EcdsaPublicKey{&key.priv.PublicKey}, nil
}

func (key *EcdsaPrivateKey) ToPEM() ([]byte,error){
	keyData, err := x509.MarshalECPrivateKey(key.priv)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type: "ECDSA PRIVATE KEY",
			Bytes: keyData,
		},
	), nil
}

func (key *EcdsaPrivateKey) Type() (keyType){
	return PRIVATE_KEY
}

type EcdsaPublicKey struct {
	pub *ecdsa.PublicKey
}

func (key *EcdsaPublicKey) SKI() ([]byte) {

	if key.pub == nil {
		return nil
	}

	data := elliptic.Marshal(key.pub.Curve, key.pub.X, key.pub.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *EcdsaPublicKey) Algorithm() string {
	return ECDSA
}

func (key *EcdsaPublicKey) ToPEM() ([]byte,error){
	keyData, err := x509.MarshalPKIXPublicKey(key.pub)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PUBLIC KEY",
			Bytes: keyData,
		},
	), nil
}

func (key *EcdsaPublicKey) Type() (keyType){
	return PUBLIC_KEY
}