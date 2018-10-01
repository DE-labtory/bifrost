package server_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mocks"
	"github.com/it-chain/bifrost/pb"
	"github.com/stretchr/testify/assert"
)

func TestServer_validateRequestPeerInfo_whenValidPeerInfo(t *testing.T) {
	// given
	keyOpt := mocks.NewMockKeyOpts()

	peerInfo := &bifrost.PeerInfo{
		IP:          "127.0.0.1",
		PubKeyBytes: keyOpt.PubKey.ToByte(),
		IsPrivate:   keyOpt.PubKey.IsPrivate(),
		KeyGenOpt:   keyOpt.PubKey.KeyGenOpt(),
	}

	mockServer := mocks.NewMockServer()

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	// when
	flag, ip, peerKey := mockServer.ValidateResponsePeerInfo(envelope)

	// then
	assert.True(t, flag)
	assert.Equal(t, peerInfo.IP, ip)
	assert.Equal(t, keyOpt.PubKey, peerKey)
}

func TestServer_BifrostStream(t *testing.T) {
	// given
	s := mocks.NewMockServer()

	keyOpt := mocks.NewMockKeyOpts()

	peerInfo := &bifrost.PeerInfo{
		IP:          "127.0.0.1",
		PubKeyBytes: keyOpt.PubKey.ToByte(),
		IsPrivate:   keyOpt.PubKey.IsPrivate(),
		KeyGenOpt:   keyOpt.PubKey.KeyGenOpt(),
	}

	mockStreamServer := mocks.NewMockStreamServer(*peerInfo)

	// when
	err := s.BifrostStream(mockStreamServer)

	// then
	assert.NoError(t, err)
}

func TestServer_Listen(t *testing.T) {
	// given
	s := mocks.NewMockServer()

	// when
	go s.Listen("127.0.0.1:7777")

	time.Sleep(3 * time.Second)
}

func TestServer_Stop(t *testing.T) {
	// given
	s := mocks.NewMockServer()
	go s.Listen("127.0.0.1:7778")

	time.Sleep(3 * time.Second)

	// when
	s.Stop()
}
