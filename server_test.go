package bifrost

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	"github.com/stretchr/testify/assert"
)

func TestServer_OnConnection(t *testing.T) {

	//given
	path := "./key"
	defer os.RemoveAll(path)

	s := getServer(path)

	assert.Nil(t, s.onConnectionHandler)

	//when
	s.OnConnection(func(conn Connection) {
		log.Printf("asd")
	})

	//then
	assert.NotNil(t, s.onConnectionHandler)
}

func TestServer_OnError(t *testing.T) {

	//given
	path := "./key"
	defer os.RemoveAll(path)

	s := getServer(path)

	assert.Nil(t, s.onErrorHandler)

	//when
	s.OnError(func(err error) {
		log.Printf("asd")
	})

	//then
	assert.NotNil(t, s.onErrorHandler)
}

func TestServer_validateRequestPeerInfo_whenValidPeerInfo(t *testing.T) {

	//given
	path := "./key"
	km, err := key.NewKeyManager(path)
	defer os.RemoveAll(path)

	if err != nil {
		log.Fatal(err.Error())
	}

	_, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	b, _ := pub.ToPEM()

	peerInfo := &PeerInfo{
		Ip:        "127.0.0.1",
		Pubkey:    b,
		KeyGenOpt: pub.Algorithm(),
	}

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_REQUEST_PEERINFO
	envelope.Payload = payload

	//when
	flag, ip, peerKey := ValidateRequestPeerInfo(envelope)

	//then
	assert.True(t, flag)
	assert.Equal(t, peerInfo.Ip, ip)
	assert.Equal(t, pub, peerKey)
}

func TestServer_validateRequestPeerInfo_whenInValidPeerInfo(t *testing.T) {

	//given
	peerInfo := &PeerInfo{
		Ip:        "127.0.0.1",
		Pubkey:    []byte("123"),
		KeyGenOpt: key.RSA2048,
	}

	payload, _ := json.Marshal(peerInfo)

	envelope := &pb.Envelope{}
	envelope.Type = pb.Envelope_REQUEST_PEERINFO
	envelope.Payload = payload

	//when
	flag, _, _ := ValidateRequestPeerInfo(envelope)

	//then
	assert.False(t, flag)
}

func TestServer_BifrostStream(t *testing.T) {

	path := "./key"
	path2 := "./key2"

	km, err := key.NewKeyManager(path)
	s := getServer(path2)

	defer os.RemoveAll(path)
	defer os.RemoveAll(path2)

	_, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	b, _ := pub.ToPEM()

	peerInfo := &PeerInfo{
		Ip:        "127.0.0.1",
		Pubkey:    b,
		KeyGenOpt: pub.Algorithm(),
	}

	mockStreamServer := NewMockStreamServer(*peerInfo)

	err = s.BifrostStream(mockStreamServer)

	assert.NoError(t, err)
}

func TestServer_Listen(t *testing.T) {

	path := "./key"

	defer os.RemoveAll(path)

	s := getServer(path)

	go s.Listen("127.0.0.1:7777")

	time.Sleep(3 * time.Second)
}
