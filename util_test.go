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
	keyOpt := mocks.NewMockKeyOpts()

	//when
	envelope, err := bifrost.BuildResponsePeerInfo(keyOpt.PubKey)
	assert.NoError(t, err)

	//then
	assert.NoError(t, err)
	assert.Equal(t, envelope.Type, pb.Envelope_RESPONSE_PEERINFO)
}
