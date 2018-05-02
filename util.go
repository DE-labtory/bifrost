package bifrost

import (
	"context"
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

func ByteToPubKey(byteKey []byte, keyGenOpt key.KeyGenOpts, keyType key.KeyType) (key.PubKey, error) {

	pubKey, err := key.PEMToPublicKey(byteKey, keyGenOpt)

	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

func recvWithTimeout(seconds int, stream Stream) (*pb.Envelope, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
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
