package bifrost_test

import (
	"testing"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mocks"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestBuildResponsePeerInfo(t *testing.T) {
	//given
	ip := "127.0.0.1:2323"
	keyOpt := mocks.NewMockKeyOpts()

	//when
	envelope, err := bifrost.BuildResponsePeerInfo(ip, keyOpt.PubKey, nil)
	assert.NoError(t, err)

	//then
	assert.NoError(t, err)
	assert.Equal(t, envelope.Type, pb.Envelope_RESPONSE_PEERINFO)
}
