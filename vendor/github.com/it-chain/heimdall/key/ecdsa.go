// This file implement ECDSA key and its generation.

package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// An ECDSAKeyGenerator contains elliptic curve for ECDSA.
type ECDSAKeyGenerator struct {
	curve elliptic.Curve
}

// Generate returns private key and public key for ECDSA using key generation option.
func (keygen *ECDSAKeyGenerator) Generate(opts KeyGenOpts) (pri PriKey, pub PubKey, err error) {

	if keygen.curve == nil {
		return nil, nil, errors.New("Curve value have not to be nil")
	}

	generatedKey, err := ecdsa.GenerateKey(keygen.curve, rand.Reader)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate ECDSA key : %s", err)
	}

	pri = &ECDSAPrivateKey{generatedKey}
	pub, err = pri.(*ECDSAPrivateKey).PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return pri, pub, nil

}

// ECDSAPrivateKey contains private key of ECDSA.
type ECDSAPrivateKey struct {
	PrivKey *ecdsa.PrivateKey
}

// SKI provides name of file that will be store a ECDSA private key.
func (key *ECDSAPrivateKey) SKI() (ski []byte) {

	if key.PrivKey == nil {
		return nil
	}

	data := elliptic.Marshal(key.PrivKey.Curve, key.PrivKey.PublicKey.X, key.PrivKey.PublicKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

// Algorithm returns key generation option of ECDSA.
func (key *ECDSAPrivateKey) Algorithm() KeyGenOpts {
	return ECDSACurveToKeyGenOpts(key.PrivKey.Curve)
}

// PublicKey returns ECDSA public key of key pair.
func (key *ECDSAPrivateKey) PublicKey() (PubKey, error) {
	return &ECDSAPublicKey{&key.PrivKey.PublicKey}, nil
}

// ToPEM makes a ECDSA private key to PEM format.
func (key *ECDSAPrivateKey) ToPEM() ([]byte, error) {
	keyData, err := x509.MarshalECPrivateKey(key.PrivKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "ECDSA PRIVATE KEY",
			Bytes: keyData,
		},
	), nil
}

// Type returns type of the ECDSA private key.
func (key *ECDSAPrivateKey) Type() KeyType {
	return PRIVATE_KEY
}

// ECDSAPublicKey contains components of a public key.
type ECDSAPublicKey struct {
	PubKey *ecdsa.PublicKey
}

// SKI provides name of file that will be store a ECDSA public key.
func (key *ECDSAPublicKey) SKI() (ski []byte) {

	if key.PubKey == nil {
		return nil
	}

	data := elliptic.Marshal(key.PubKey.Curve, key.PubKey.X, key.PubKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

// Algorithm returns ECDSA public key generation option.
func (key *ECDSAPublicKey) Algorithm() KeyGenOpts {
	return ECDSACurveToKeyGenOpts(key.PubKey.Curve)
}

// ToPEM makes a ECDSA public key to PEM format.
func (key *ECDSAPublicKey) ToPEM() ([]byte, error) {
	keyData, err := x509.MarshalPKIXPublicKey(key.PubKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "ECDSA PUBLIC KEY",
			Bytes: keyData,
		},
	), nil
}

// Type returns type of the ECDSA public key.
func (key *ECDSAPublicKey) Type() KeyType {
	return PUBLIC_KEY
}
