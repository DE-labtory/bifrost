package client

import (
	"time"

	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"context"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost"
	"errors"
	"encoding/json"
)

var ErrNotExpectedMessage = errors.New("wrong message type")

// defaultDialTimeout
const (
	defaultDialTimeout = 3 * time.Second
)

// Server 와 연결시 사용되는 Client option
type ClientOpts struct {
	ip     string
	priKey key.PriKey
	pubKey key.PubKey
	mux    mux.Mux
}

// Server 와 연결시 사용되는 grpc option.
type GrpcOpts struct {
	tlsEnabled bool
	creds      credentials.TransportCredentials
}

// 서버와 연결 요청. 실패시 err. handshake 과정을 거침.
func Dial(serverIp string, clientOpts ClientOpts, grpcOpts GrpcOpts) (bifrost.Connection, error) {

	var opts []grpc.DialOption

	if grpcOpts.tlsEnabled {
		opts = append(opts, grpc.WithTransportCredentials(grpcOpts.creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	dialContext, _ := context.WithTimeout(context.Background(), defaultDialTimeout)
	gconn, err := grpc.DialContext(dialContext, serverIp, opts...)

	if err != nil {
		return nil, err
	}

	streamWrapper, err := connect(gconn)

	if err != nil {
		return nil, err
	}

	serverPubKey, err := handShake(streamWrapper,clientOpts)

	if err != nil {
		return nil, err
	}

	conn, err := bifrost.NewConnection(serverIp, clientOpts.priKey, serverPubKey, streamWrapper)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// handshake func, return : serverPubKey, err
func handShake(streamWrapper bifrost.StreamWrapper, clientOpts ClientOpts) (key.PubKey, error) {
	env, err := bifrost.RecvWithTimeout(10*time.Second, streamWrapper)
	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	if env.GetType() != pb.Envelope_REQUEST_PEERINFO {
		streamWrapper.Close()
		return nil, ErrNotExpectedMessage
	}

	env, err = bifrost.BuildRequestPeerInfo(clientOpts.ip, clientOpts.pubKey)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	err = streamWrapper.Send(env)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	env, err = bifrost.RecvWithTimeout(3 * time.Second, streamWrapper)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	/////////////
	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)

	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	serverPubKey, err := bifrost.ByteToPubKey(peerInfo.Pubkey, peerInfo.KeyGenOpt)


	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	return serverPubKey, nil


}

func connect(conn *grpc.ClientConn) (bifrost.StreamWrapper, error) {

	streamWrapper, err := bifrost.NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
