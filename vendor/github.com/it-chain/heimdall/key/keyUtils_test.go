package key

import (
	"testing"

	"os"

	"encoding/pem"

	"github.com/stretchr/testify/assert"
)

func TestPEMToPrivateKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	pri, _, _ := testKeyManager.GenerateKey(keyGenOption)

	priPEM, _ := pri.ToPEM()

	testPri, err := PEMToPrivateKey(priPEM, keyGenOption)
	assert.NotNil(t, testPri)
	assert.NoError(t, err)

	defer os.RemoveAll("./.testKeys")
}

func TestPEMToPublicKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	_, pub, _ := testKeyManager.GenerateKey(keyGenOption)

	pubPEM, _ := pub.ToPEM()

	testPub, err := PEMToPublicKey(pubPEM, keyGenOption)
	assert.NotNil(t, testPub)
	assert.NoError(t, err)

	defer os.RemoveAll("./.testKeys")
}

func TestDERToPrivateKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	pri, _, _ := testKeyManager.GenerateKey(keyGenOption)

	priPEM, _ := pri.ToPEM()
	block, _ := pem.Decode(priPEM)

	myPri, err := DERToPrivateKey(block.Bytes)
	assert.NotNil(t, myPri)
	assert.NoError(t, err)

	defer os.RemoveAll("./.testKeys")
}

func TestDERToPublicKey(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	_, pub, _ := testKeyManager.GenerateKey(keyGenOption)

	pubPEM, _ := pub.ToPEM()
	block, _ := pem.Decode(pubPEM)

	myPub, err := DERToPublicKey(block.Bytes)
	assert.NotNil(t, myPub)
	assert.NoError(t, err)

	defer os.RemoveAll("./.testKeys")
}

func TestMatchPrivateKeyOpt(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	pri, _, _ := testKeyManager.GenerateKey(keyGenOption)

	priPEM, _ := pri.ToPEM()

	block, _ := pem.Decode(priPEM)

	testPri, _ := DERToPrivateKey(block.Bytes)

	myPri, err := MatchPrivateKeyOpt(testPri, keyGenOption)
	assert.NoError(t, err)
	assert.NotNil(t, myPri)

	defer os.RemoveAll("./.testKeys")
}

func TestMatchPublicKeyOpt(t *testing.T) {
	var keyGenOption = KeyGenOpts(RSA2048)

	testKeyManager, _ := NewKeyManager("./.testKeys")
	_, pub, _ := testKeyManager.GenerateKey(keyGenOption)

	pubPEM, _ := pub.ToPEM()

	block, _ := pem.Decode(pubPEM)

	testPub, _ := DERToPublicKey(block.Bytes)

	myPub, err := MatchPublicKeyOpt(testPub, keyGenOption)
	assert.NoError(t, err)
	assert.NotNil(t, myPub)

	defer os.RemoveAll("./.testKeys")
}
