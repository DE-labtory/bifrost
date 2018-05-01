package auth

import (
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/heimdall/key"
	"crypto/sha512"
	"testing"
	"crypto/rand"
	"crypto/rsa"
)

func TestECDSA_SignVerify(t *testing.T) {

	rsaBits := 4096
	reader := rand.Reader

	generatedKey, err := rsa.GenerateKey(reader, rsaBits)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	pri := &key.RSAPrivateKey{generatedKey, rsaBits}
	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	rawData := []byte("RSA test data!")

	hash := sha512.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	signer := &RSASigner{}
	signature, err := signer.Sign(pri, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	verifier := &RSAVerifier{}
	ok, err := verifier.Verify(pub, signature, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.NoError(t, err)
	assert.True(t, ok)

}