package auth

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/heimdall/key"
	"crypto/sha512"
)

func TestRSA_SignVerify(t *testing.T) {

	ecdsaCurve := elliptic.P521()
	reader := rand.Reader

	generatedKey, err := ecdsa.GenerateKey(ecdsaCurve, reader)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	pri := &key.ECDSAPrivateKey{generatedKey}
	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	rawData := []byte("ECDSA test data!")

	hash := sha512.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	signer := &ECDSASigner{}
	signature, err := signer.Sign(pri, digest, nil)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	verifier := &ECDSAVerifier{}
	ok, err := verifier.Verify(pub, signature, digest, nil)
	assert.NoError(t, err)
	assert.True(t, ok)

}