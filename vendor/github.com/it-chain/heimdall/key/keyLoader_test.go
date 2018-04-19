package key

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyLoader_Load(t *testing.T) {
	var keyGenTester = RSAKeyGenerator{bits: 2048}
	var keyGenOption = KeyGenOpts(RSA2048)

	var pri, pub, _ = keyGenTester.Generate(keyGenOption)
	var keyStoreTester = keyStorer{path: "./.testKeys"}

	_ = keyStoreTester.Store(pri, pub)

	var keyLoadTester = keyLoader{path: "./.testKeys"}

	loadedPri, loadedPub, err := keyLoadTester.Load()
	assert.NoError(t, err)
	assert.NotNil(t, loadedPri)
	assert.NotNil(t, loadedPub)

	defer os.RemoveAll("./.testKeys")
}
