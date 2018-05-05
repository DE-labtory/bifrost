package bifrost

import (
	"context"
	"encoding/json"
	"time"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
	b58 "github.com/jbenet/go-base58"
)

func FromPubKey(key key.PubKey) string {

	encoded := b58.Encode(key.SKI())
	return encoded
}

//Create ID from Pri Key
func FromPriKey(key key.PriKey) string {

	pub, _ := key.PublicKey()
	return FromPubKey(pub)
}

func ByteToPubKey(byteKey []byte, keyGenOpt key.KeyGenOpts) (key.PubKey, error) {

	pubKey, err := key.PEMToPublicKey(byteKey, keyGenOpt)

	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func RecvWithTimeout(timeout time.Duration, stream Stream) (*pb.Envelope, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := make(chan *pb.Envelope, 1)
	errch := make(chan error, 1)

	go func() {
		envelope, err := stream.Recv()
		if err != nil {
			errch <- err
		}
		c <- envelope
	}()

	select {
	case <-ctx.Done():
		//timeoutted body
		return nil, ctx.Err()
	case err := <-errch:
		return nil, err
	case ok := <-c:
		//okay body
		return ok, nil
	}
}

type KeyOpts struct {
	PriKey key.PriKey
	PubKey key.PubKey
}

func BuildResponsePeerInfo(pubKey key.PubKey) (*pb.Envelope, error) {
	b, _ := pubKey.ToPEM()

	pi := &PeerInfo{
		Pubkey:    b,
		KeyGenOpt: pubKey.Algorithm(),
	}

	payload, err := json.Marshal(pi)

	if err != nil {
		return nil, err
	}

	return &pb.Envelope{
		Payload: payload,
		Type:    pb.Envelope_RESPONSE_PEERINFO,
	}, nil
}
