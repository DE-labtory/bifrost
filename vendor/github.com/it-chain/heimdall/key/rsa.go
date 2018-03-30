package key

import (
	"crypto/rsa"
	"fmt"
	"errors"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type RSAKeyGenerator struct {
	bits int
}

func (keygen *RSAKeyGenerator) Generate(opts KeyGenOpts) (pri PriKey, pub PubKey, err error) {

	if keygen.bits <= 0 {
		return nil, nil, errors.New("Bits length should be bigger than 0")
	}

	generatedKey, err := rsa.GenerateKey(rand.Reader, keygen.bits)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to generate RSA key : %s", err)
	}

	pri = &RSAPrivateKey{PrivKey:generatedKey, bits:keygen.bits}
	pub, err = pri.(*RSAPrivateKey).PublicKey()
	if err != nil {
		return nil, nil, err
	}

	return pri, pub, nil

}

type rsaKeyMarshalOpt struct {
	N *big.Int
	E int
}

type RSAPrivateKey struct {
	PrivKey *rsa.PrivateKey
	bits int
}

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

func (key *RSAPrivateKey) Algorithm() KeyGenOpts {
	return RSABitsToKeyGenOpts(key.bits)
}

func (key *RSAPrivateKey) PublicKey() (pub PubKey, err error) {
	return &RSAPublicKey{PubKey: &key.PrivKey.PublicKey, bits: key.bits}, nil
}

func (key *RSAPrivateKey) ToPEM() ([]byte,error) {
	keyData := x509.MarshalPKCS1PrivateKey(key.PrivKey)

	return pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: keyData,
		},
	), nil
}

func (key *RSAPrivateKey) Type() (KeyType) {
	return PRIVATE_KEY
}

type RSAPublicKey struct {
	PubKey *rsa.PublicKey
	bits int
}

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

func (key *RSAPublicKey) Algorithm() KeyGenOpts {
	return RSABitsToKeyGenOpts(key.bits)
}

func (key *RSAPublicKey) ToPEM() ([]byte,error) {

	keyData, err := x509.MarshalPKIXPublicKey(key.PubKey)

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

func (key *RSAPublicKey) Type() (KeyType){
	return PUBLIC_KEY
}
