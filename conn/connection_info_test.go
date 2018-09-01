package conn

import (
	"testing"

	"github.com/it-chain/heimdall"
	"github.com/stretchr/testify/assert"
)


func TestFromPublicConnInfo(t *testing.T) {

	pri, err := heimdall.GenerateKey(heimdall.SECP384R1)
	pub := &pri.PublicKey

	b := heimdall.PubKeyToBytes(pub)

	pci := PublicConnInfo{}
	pci.Id = "test1"
	pci.Address = Address{IP: "127.0.0.1"}
	pci.Pubkey = b
	pci.CurveOpt = heimdall.CurveToCurveOpt(pub.Curve)

	//when
	connInfo, err := FromPublicConnInfo(pci)

	//then
	assert.NoError(t, err)
	assert.Equal(t, pub, connInfo.PubKey)
	assert.Equal(t, pci.Id, string(connInfo.Id))
	assert.Equal(t, pci.Address, connInfo.Address)
}
