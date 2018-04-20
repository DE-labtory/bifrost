package key

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestKeyStorer_Store(t *testing.T) {
	var keyGenTester = RSAKeyGenerator{bits: 2048}
	var keyGenOption = KeyGenOpts(RSA2048)

	var pri, pub, _ = keyGenTester.Generate(keyGenOption)
	var keyStoreTester = keyStorer{path: "./.testKeys"}

	defer os.RemoveAll("./.testKeys")

	err := keyStoreTester.Store(pri, pub)
	assert.NoError(t, err)
}
