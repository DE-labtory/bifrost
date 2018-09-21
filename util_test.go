package bifrost

import (
	"os"
	"testing"

	"github.com/it-chain/bifrost/pb"
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

	//when
	pubK, err := ByteToPubKey(b, keyGenOpt)

	//then
	assert.NoError(t, err)
	assert.Equal(t, pubK, pub)
}

func TestBuildResponsePeerInfo(t *testing.T) {

	//given
	km, err := key.NewKeyManager("~/key")
	assert.NoError(t, err)
	defer os.RemoveAll("~/key")

	_, pub, err := km.GenerateKey(key.RSA4096)
	assert.NoError(t, err)

	//when
	envelope, err := BuildResponsePeerInfo(pub, nil)
	assert.NoError(t, err)

	//then
	assert.NoError(t, err)
	assert.Equal(t, envelope.Type, pb.Envelope_RESPONSE_PEERINFO)
}
