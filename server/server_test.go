package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/engine/common/logger"
	"github.com/stretchr/testify/assert"
)

func TestServer_OnConnection(t *testing.T) {

	//given
	s := GetServer()

	assert.Nil(t, s.onConnectionHandler)

	//when
	s.OnConnection(func(conn bifrost.Connection) {
		logger.Infof(nil, "test server on connection")
	})

	//then
	assert.NotNil(t, s.onConnectionHandler)
}

func TestServer_OnError(t *testing.T) {

	//given
	s := GetServer()

	assert.Nil(t, s.onErrorHandler)

	//when
	s.OnError(func(err error) {
		logger.Infof(nil, "test server on error")
	})

	//then
	assert.NotNil(t, s.onErrorHandler)
}

func TestServer_validateRequestPeerInfo_whenValidPeerInfo(t *testing.T) {

	//given
	mockGenerator := MockGenerator{}
	pri, err := mockGenerator.GenerateKey()

	if err != nil {
		logger.Fatal(nil, err.Error())
	}

	pub := &pri.PublicKey
	mockFormatter := MockFormatter{}
	b := mockFormatter.ToByte(pub)

	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   b,
		CurveOpt: mockFormatter.GetCurveOpt(pub),
	}

	mockServer := GetServer()

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	//when
	flag, ip, peerKey := mockServer.ValidateResponsePeerInfo(envelope)

	//then
	assert.True(t, flag)
	assert.Equal(t, peerInfo.IP, ip)
	assert.Equal(t, pub, peerKey)
}

func TestServer_validateRequestPeerInfo_whenInValidPeerInfo(t *testing.T) {

	//given
	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   []byte("123"),
		CurveOpt: 2,
	}

	mockServer := GetServer()

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_RESPONSE_PEERINFO
	envelope.Payload = payload

	//when
	flag, _, _ := mockServer.validateRequestPeerInfo(envelope)

	//then
	assert.False(t, flag)
}

func TestServer_BifrostStream(t *testing.T) {

	s := GetServer()

	mockGenerator := MockGenerator{}
	pri, err := mockGenerator.GenerateKey()

	if err != nil {
		logger.Fatal(nil, err.Error())
	}

	pub := &pri.PublicKey
	mockFormatter := MockFormatter{}
	b := mockFormatter.ToByte(pub)

	peerInfo := &bifrost.PeerInfo{
		IP:       "127.0.0.1",
		Pubkey:   b,
		CurveOpt: mockFormatter.GetCurveOpt(pub),
	}

	mockStreamServer := NewMockStreamServer(*peerInfo)

	err = s.BifrostStream(mockStreamServer)

	assert.NoError(t, err)
}

func TestServer_Listen(t *testing.T) {

	s := GetServer()

	go s.Listen("127.0.0.1:7777")

	time.Sleep(3 * time.Second)
}
