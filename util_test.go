package bifrost

import (
	"testing"

	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestBuildResponsePeerInfo(t *testing.T) {
	//given
	mockGenerator := MockGenerator{}
	pri, err := mockGenerator.GenerateKey()
	assert.NoError(t, err)

	mockFormatter := MockFormatter{}

	//when
	envelope, err := BuildResponsePeerInfo(&pri.PublicKey, &mockFormatter)
	assert.NoError(t, err)

	//then
	assert.NoError(t, err)
	assert.Equal(t, envelope.Type, pb.Envelope_RESPONSE_PEERINFO)
}
