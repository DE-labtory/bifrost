package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var RSAKeyBits = RSAKeyGenerator{bits: 2048}
var RSAKeyGenOption = KeyGenOpts(RSA2048)

func TestRSAKeyPairGeneration(t *testing.T) {
	pri, pub, err := RSAKeyBits.Generate(RSAKeyGenOption)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)
}

func TestRSAKeyPairSKI(t *testing.T) {
	pri, pub, _ := RSAKeyBits.Generate(RSAKeyGenOption)

	priSki := pri.SKI()
	assert.NotNil(t, priSki)

	pubSki := pub.SKI()
	assert.NotNil(t, pubSki)
}

func TestRSAGetAlgorithm(t *testing.T) {
	pri, pub, _ := RSAKeyBits.Generate(RSAKeyGenOption)

	priKeyOption := pri.Algorithm()
	assert.NotNil(t, priKeyOption)

	pubKeyOption := pub.Algorithm()
	assert.NotNil(t, pubKeyOption)
}

func TestRSAGetPublicKey(t *testing.T) {
	pri, _, _ := RSAKeyBits.Generate(RSAKeyGenOption)

	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)
}

func TestRSAKeyToPEM(t *testing.T) {
	pri, pub, _ := RSAKeyBits.Generate(RSAKeyGenOption)

	priPEM, err := pri.ToPEM()
	assert.NoError(t, err)
	assert.NotNil(t, priPEM)

	pubPEM, err := pub.ToPEM()
	assert.NoError(t, err)
	assert.NotNil(t, pubPEM)
}

func TestGetRSAKeyType(t *testing.T) {
	pri, pub, _ := RSAKeyBits.Generate(RSAKeyGenOption)

	priType := pri.Type()
	assert.Equal(t, priType, PRIVATE_KEY, "They should be equal")

	pubType := pub.Type()
	assert.Equal(t, pubType, PUBLIC_KEY, "They should be equal")
}
