package key

import (
	"testing"

	"crypto/elliptic"

	"github.com/stretchr/testify/assert"
)

var ecdsaCurve = ECDSAKeyGenerator{curve: elliptic.P521()}
var ECDSAkeyGenOption = KeyGenOpts(ECDSA521)

func TestECDSAKeyPairGeneration(t *testing.T) {
	pri, pub, err := ecdsaCurve.Generate(ECDSAkeyGenOption)
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)
}

func TestECDSAKeyPairSKI(t *testing.T) {
	pri, pub, _ := ecdsaCurve.Generate(ECDSAkeyGenOption)

	priSki := pri.SKI()
	assert.NotNil(t, priSki)

	pubSki := pub.SKI()
	assert.NotNil(t, pubSki)
}

func TestECDSAGetAlgorithm(t *testing.T) {
	pri, pub, _ := ecdsaCurve.Generate(ECDSAkeyGenOption)

	priKeyOption := pri.Algorithm()
	assert.NotNil(t, priKeyOption)

	pubKeyOption := pub.Algorithm()
	assert.NotNil(t, pubKeyOption)
}

func TestECDSAGetPublicKey(t *testing.T) {
	pri, _, _ := ecdsaCurve.Generate(ECDSAkeyGenOption)

	pub, err := pri.PublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, pub)
}

func TestECDSAKeyToPEM(t *testing.T) {
	pri, pub, _ := ecdsaCurve.Generate(ECDSAkeyGenOption)

	priPEM, err := pri.ToPEM()
	assert.NoError(t, err)
	assert.NotNil(t, priPEM)

	pubPEM, err := pub.ToPEM()
	assert.NoError(t, err)
	assert.NotNil(t, pubPEM)
}

func TestGetECDSAKeyType(t *testing.T) {
	pri, pub, _ := ecdsaCurve.Generate(ECDSAkeyGenOption)

	priType := pri.Type()
	assert.Equal(t, priType, PRIVATE_KEY, "They should be equal")

	pubType := pub.Type()
	assert.Equal(t, pubType, PUBLIC_KEY, "They should be equal")
}
