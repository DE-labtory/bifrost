package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
	"crypto/rand"
	"github.com/it-chain/heimdall/key"
	"crypto/sha512"
	"crypto/elliptic"
	"crypto/ecdsa"
)

func TestNewAuth(t *testing.T) {

	authManager, err := NewAuth()
	assert.NoError(t, err)
	assert.NotNil(t, authManager)

}

func TestAuth_RSASignVerify(t *testing.T) {

	rsaKeyBits := 4096

	generatedKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	pri := &key.RSAPrivateKey{generatedKey, rsaKeyBits}
	assert.NotNil(t, pri)

	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	diffGeneratedKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)

	diffPri := &key.RSAPrivateKey{diffGeneratedKey, rsaKeyBits}
	diffPub, err := diffPri.PublicKey()

	rawData := []byte("RSA Sign test data!!!")

	hash := sha512.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	// hash for generating wrong type of digest
	hash = sha512.New512_256()
	hash.Write(rawData)
	wrongDigest := hash.Sum(nil)

	authManager, err := NewAuth()
	assert.NoError(t, err)
	assert.NotNil(t, authManager)

	signature, err := authManager.Sign(pri, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	diffSignature, err := authManager.Sign(diffPri, digest, EQUAL_SHA512.SignerOptsToPSSOptions())

	// normal case
	ok, err := authManager.Verify(pub, signature, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.NoError(t, err)
	assert.True(t, ok)

	// passing different signature case
	_, err = authManager.Verify(pub, diffSignature, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// public key missing case
	_, err = authManager.Verify(nil, signature, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// passing different public key case
	_, err = authManager.Verify(diffPub, signature, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// signature missing case
	_, err = authManager.Verify(pub, nil, digest, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// digest missing case
	_, err = authManager.Verify(pub, signature, nil, EQUAL_SHA512.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// passing wrong digest case
	_, err = authManager.Verify(pub, signature, wrongDigest, EQUAL_SHA256.SignerOptsToPSSOptions())
	assert.Error(t, err)

	// passing wrong signer option case
	ok, err = authManager.Verify(pub, signature, digest, EQUAL_SHA256.SignerOptsToPSSOptions())
	assert.Error(t, err)
	assert.False(t, ok)

}

func TestAuth_ECDSASignVerify(t *testing.T) {

	ecdsaCurve := elliptic.P521()

	generatedKey, err := ecdsa.GenerateKey(ecdsaCurve, rand.Reader)
	assert.NoError(t, err)
	assert.NotNil(t, generatedKey)

	pri := &key.ECDSAPrivateKey{PrivKey:generatedKey}
	assert.NotNil(t, pri)

	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)

	diffGeneratedKey, err := ecdsa.GenerateKey(ecdsaCurve, rand.Reader)

	diffPri := &key.ECDSAPrivateKey{diffGeneratedKey}
	diffPub, err := diffPri.PublicKey()

	rawData := []byte("ECDSA Sign test data!!!")

	hash := sha512.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	hash = sha512.New512_256()
	hash.Write(rawData)
	wrongDigest := hash.Sum(nil)

	authManager, err := NewAuth()
	assert.NoError(t, err)
	assert.NotNil(t, authManager)

	signature, err := authManager.Sign(pri, digest, nil)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	diffSignature, err := authManager.Sign(diffPri, digest, nil)

	// normal case
	ok, err := authManager.Verify(pub, signature, digest, nil)
	assert.NoError(t, err)
	assert.True(t, ok)

	// passing different signature case
	_, err = authManager.Verify(pub, diffSignature, digest, nil)
	assert.Error(t, err)

	// public key missing case
	_, err = authManager.Verify(nil, signature, digest, nil)
	assert.Error(t, err)

	// passing different public key case
	_, err = authManager.Verify(diffPub, signature, digest, nil)
	assert.Error(t, err)

	// signature missing case
	_, err = authManager.Verify(pub, nil, digest, nil)
	assert.Error(t, err)

	// digest missing case
	_, err = authManager.Verify(pub, signature, nil, nil)
	assert.Error(t, err)

	// passing wrong digest case
	_, err = authManager.Verify(pub, signature, wrongDigest, nil)
	assert.Error(t, err)

}