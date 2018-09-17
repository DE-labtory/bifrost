package server_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/logger"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/server"
	"github.com/stretchr/testify/assert"
)

func TestServer_validateRequestPeerInfo_whenValidPeerInfo(t *testing.T) {
	// given
	mockGenerator := server.MockGenerator{}
	pri, err := mockGenerator.GenerateKey()
	if err != nil {
		logger.Fatalf(nil, err.Error())
	}

	pub := &pri.PublicKey
	mockFormatter := server.MockFormatter{}
	b := mockFormatter.ToByte(pub)

	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   b,
		CurveOpt: mockFormatter.GetCurveOpt(pub),
	}

	mockServer := server.GetServer()

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	// when
	flag, ip, peerKey := mockServer.ValidateResponsePeerInfo(envelope)

	// then
	assert.True(t, flag)
	assert.Equal(t, peerInfo.IP, ip)
	assert.Equal(t, pub, peerKey)
}

func TestServer_BifrostStream(t *testing.T) {
	// given
	s := server.GetServer()

	mockGenerator := bifrost.MockGenerator{}
	pri, err := mockGenerator.GenerateKey()

	if err != nil {
		logger.Fatalf(nil, err.Error())
	}

	pub := &pri.PublicKey
	mockFormatter := bifrost.MockFormatter{}
	b := mockFormatter.ToByte(pub)

	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   b,
		CurveOpt: mockFormatter.GetCurveOpt(pub),
	}

	mockStreamServer := server.NewMockStreamServer(*peerInfo)

	// when
	err = s.BifrostStream(mockStreamServer)

	// then
	assert.NoError(t, err)
}

func TestServer_Listen(t *testing.T) {
	// given
	s := server.GetServer()

	// when
	go s.Listen("127.0.0.1:7777")

	time.Sleep(3 * time.Second)
}

func TestServer_Stop(t *testing.T) {
	// given
	s := server.GetServer()
	go s.Listen("127.0.0.1:7778")

	time.Sleep(3 * time.Second)

	// when
	s.Stop()
}
