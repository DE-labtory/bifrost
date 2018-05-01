package bifrost

import (
	"os"
	"testing"

	"github.com/it-chain/heimdall/key"
	"github.com/stretchr/testify/assert"
)

func TestByteToPubKey(t *testing.T) {

	//given
	km, err := key.NewKeyManager("~/key")
	defer os.RemoveAll("~/key")

	_, pub, err := km.GenerateKey(key.RSA4096)
	assert.NoError(t, err)

	b, err := pub.ToPEM()
	assert.NoError(t, err)

	keyGenOpt := pub.Algorithm()
	keyType := pub.Type()

	//when
	pubK, err := ByteToPubKey(b, keyGenOpt, keyType)

	//then
	assert.NoError(t, err)
	assert.Equal(t, pubK, pub)
}
