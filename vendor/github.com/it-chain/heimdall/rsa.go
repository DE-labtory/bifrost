package heimdall

import (
	"crypto/rsa"
	"errors"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type RsaSigner struct{}

func (s *RsaSigner) Sign(key Key, digest []byte, opts SignerOpts) ([]byte, error) {

	if opts == nil {
		return nil, errors.New("invalid options")
	}

	return key.(*RsaPrivateKey).priv.Sign(rand.Reader, digest, opts)

}

type RsaVerifier struct{}

func (v *RsaVerifier) Verify(key Key, signature, digest []byte, opts SignerOpts) (bool, error) {

	if opts == nil {
		return false, errors.New("invalid options")
	}

	switch opts.(type) {
	case *rsa.PSSOptions:
		err := rsa.VerifyPSS(key.(*RsaPublicKey).pub,
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

type rsaKeyMarshalOpt struct {
	N *big.Int
	E int
}

type RsaPrivateKey struct {
	priv *rsa.PrivateKey
}

func (key *RsaPrivateKey) SKI() ([]byte) {

	if key.priv == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		key.priv.N, key.priv.E,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)

}

func (key *RsaPrivateKey) Algorithm() string {
	return RSA
}

func (key *RsaPrivateKey) PublicKey() (pub *RsaPublicKey, err error) {
	return &RsaPublicKey{&key.priv.PublicKey}, nil
}

func (key *RsaPrivateKey) ToPEM() ([]byte,error) {
	keyData := x509.MarshalPKCS1PrivateKey(key.priv)

	return pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: keyData,
		},
	), nil
}

func (key *RsaPrivateKey) Type() (keyType){
	return PRIVATE_KEY
}

type RsaPublicKey struct {
	pub *rsa.PublicKey
}

func (key *RsaPublicKey) SKI() ([]byte) {

	if key.pub == nil {
		return nil
	}

	data, _ := asn1.Marshal(rsaKeyMarshalOpt{
		big.NewInt(123), 57,
	})

	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func (key *RsaPublicKey) Algorithm() string {
	return RSA
}

func (key *RsaPublicKey) ToPEM() ([]byte,error) {
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

func (key *RsaPublicKey) Type() (keyType){
	return PUBLIC_KEY
}
