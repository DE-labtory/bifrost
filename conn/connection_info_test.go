package conn

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

func TestFromPublicConnInfo(t *testing.T) {

	//given
	km, err := key.NewKeyManager("~/key")
	assert.NoError(t, err)

	defer os.RemoveAll("~/key")

	_, pub, err := km.GenerateKey(key.RSA4096)

	b, err := pub.ToPEM()
	assert.NoError(t, err)

	pci := PublicConnInfo{}
	pci.Id = "test1"
	pci.Address = Address{IP: "127.0.0.1"}
	pci.Pubkey = b
	pci.KeyGenOpt = pub.Algorithm()
	pci.KeyType = pub.Type()

	//when
	connInfo, err := FromPublicConnInfo(pci)

	//then
	assert.NoError(t, err)
	assert.Equal(t, pub, connInfo.PubKey)
	assert.Equal(t, pci.Id, string(connInfo.Id))
	assert.Equal(t, pci.Address, connInfo.Address)
}
