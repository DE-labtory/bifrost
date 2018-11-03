package bifrost

import (
	"context"
	"encoding/json"
	"time"

	"github.com/it-chain/bifrost/pb"
)

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
	PriKey Key
	PubKey Key
}

func BuildResponsePeerInfo(ip string, pubKey Key, metaData map[string]string) (*pb.Envelope, error) {
	b, err := pubKey.ToByte()

	if err != nil {
		return nil, err
	}

	pi := &PeerInfo{
		IP:          ip,
		PubKeyBytes: b,
		IsPrivate:   pubKey.IsPrivate(),
		MetaData:    metaData,
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
