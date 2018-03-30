package auth

import (
	"crypto/rsa"
	"crypto/rand"
	"errors"
	"github.com/it-chain/heimdall/key"
)

type RSASigner struct{}

func (s *RSASigner) Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error) {

	if opts == nil {
		return nil, errors.New("invalid options")
	}

	return priKey.(*key.RSAPrivateKey).PrivKey.Sign(rand.Reader, digest, opts)

}

type RSAVerifier struct{}

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