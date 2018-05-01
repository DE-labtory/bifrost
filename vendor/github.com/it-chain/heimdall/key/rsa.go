// This file implement RSA key and its generation.

package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

// An RSAKeyGenerator contains RSA key length.
type RSAKeyGenerator struct {
	bits int
}

// Generate returns private key and public key for RSA using key generation option.
func (keygen *RSAKeyGenerator) Generate(opts KeyGenOpts) (pri PriKey, pub PubKey, err error) {

	if keygen.bits <= 0 {
		return nil, nil, errors.New("Bits length should be bigger than 0")
	}

	generatedKey, err := rsa.GenerateKey(rand.Reader, keygen.bits)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate RSA key : %s", err)
	}

	pri = &RSAPrivateKey{PrivKey: generatedKey, Bits: keygen.bits}
	pub, err = pri.(*RSAPrivateKey).PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return pri, pub, nil

}

// rsaKeyMarshalOpt contains N and E that are RSA key's components
type rsaKeyMarshalOpt struct {
	N *big.Int
	E int
}

// RSAPrivateKey contains private key of RSA.
type RSAPrivateKey struct {
	PrivKey *rsa.PrivateKey
	Bits    int
}

// SKI provides name of file that will be store a RSA private key.
func (key *RSAPrivateKey) SKI() (ski []byte) {

	if key.PrivKey == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		key.PrivKey.N, key.PrivKey.E,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

// Algorithm returns key generation option of RSA.
func (key *RSAPrivateKey) Algorithm() KeyGenOpts {
	return RSABitsToKeyGenOpts(key.Bits)
}

// PublicKey returns RSA public key of key pair.
func (key *RSAPrivateKey) PublicKey() (pub PubKey, err error) {
	return &RSAPublicKey{PubKey: &key.PrivKey.PublicKey, Bits: key.Bits}, nil
}

// ToPEM makes a RSA private key to PEM format.
func (key *RSAPrivateKey) ToPEM() ([]byte, error) {
	keyData := x509.MarshalPKCS1PrivateKey(key.PrivKey)

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyData,
		},
	), nil
}

// Type returns type of the RSA private key.
func (key *RSAPrivateKey) Type() KeyType {
	return PRIVATE_KEY
}

// RSAPublicKey contains components of a public key.
type RSAPublicKey struct {
	PubKey *rsa.PublicKey
	Bits   int
}

// SKI provides name of file that will be store a RSA public key.
func (key *RSAPublicKey) SKI() (ski []byte) {

	if key.PubKey == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		big.NewInt(123), 57,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// Algorithm returns RSA public key generation option.
func (key *RSAPublicKey) Algorithm() KeyGenOpts {
	return RSABitsToKeyGenOpts(key.Bits)
}

// ToPEM makes a RSA public key to PEM format.
func (key *RSAPublicKey) ToPEM() ([]byte, error) {

	keyData, err := x509.MarshalPKIXPublicKey(key.PubKey)

	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyData,
		},
	), nil

}

// Type returns type of the RSA public key.
func (key *RSAPublicKey) Type() KeyType {
	return PUBLIC_KEY
}
