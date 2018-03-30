package key

import (
	"crypto/elliptic"
	"crypto/ecdsa"
	"fmt"
	"errors"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"crypto/sha256"
)

type ECDSAKeyGenerator struct {
	curve elliptic.Curve
}

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

type ECDSAPrivateKey struct {
	PrivKey *ecdsa.PrivateKey
}

func (key *ECDSAPrivateKey) SKI() (ski []byte) {

	if key.PrivKey == nil {
		return nil
	}

	data := elliptic.Marshal(key.PrivKey.Curve, key.PrivKey.PublicKey.X, key.PrivKey.PublicKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *ECDSAPrivateKey) Algorithm() KeyGenOpts {
	return ECDSACurveToKeyGenOpts(key.PrivKey.Curve)
}

func (key *ECDSAPrivateKey) PublicKey() (PubKey, error) {
	return &ECDSAPublicKey{&key.PrivKey.PublicKey}, nil
}

func (key *ECDSAPrivateKey) ToPEM() ([]byte,error){
	keyData, err := x509.MarshalECPrivateKey(key.PrivKey)
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

func (key *ECDSAPrivateKey) Type() (KeyType){
	return PRIVATE_KEY
}

type ECDSAPublicKey struct {
	PubKey *ecdsa.PublicKey
}

func (key *ECDSAPublicKey) SKI() (ski []byte) {

	if key.PubKey == nil {
		return nil
	}

	data := elliptic.Marshal(key.PubKey.Curve, key.PubKey.X, key.PubKey.Y)

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *ECDSAPublicKey) Algorithm() KeyGenOpts {
	return ECDSACurveToKeyGenOpts(key.PubKey.Curve)
}

func (key *ECDSAPublicKey) ToPEM() ([]byte,error){
	keyData, err := x509.MarshalPKIXPublicKey(key.PubKey)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type: "ECDSA PUBLIC KEY",
			Bytes: keyData,
		},
	), nil
}

func (key *ECDSAPublicKey) Type() (KeyType){
	return PUBLIC_KEY
}
