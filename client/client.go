package client

import (
	"time"

	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"context"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost"
	"errors"
	"encoding/json"
)

// handshake 과정에서 올바르지 않은 메세지 타입이 올 경우 발생하는 에러
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

	serverPubKey, err := handShake(streamWrapper, clientOpts)

	if err != nil {
		return nil, err
	}

	conn, err := bifrost.NewConnection(serverIp, clientOpts.priKey, serverPubKey, streamWrapper)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// handshake 함수, return : serverPubKey, err
func handShake(streamWrapper bifrost.StreamWrapper, clientOpts ClientOpts) (key.PubKey, error) {
	err := handShakeWaitServer(streamWrapper)
	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	err = handShakeSendInfo(streamWrapper, clientOpts)
	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	serverPubKey, err := handShakeGetServerInfo(streamWrapper)
	if err != nil {
		streamWrapper.Close()
		return nil, err
	}

	return serverPubKey, nil

}

// handshake 첫번째 과정 함수. server 의 request peer info 메세지를 기다린다.
func handShakeWaitServer(streamWrapper bifrost.StreamWrapper) error {
	env, err := bifrost.RecvWithTimeout(10*time.Second, streamWrapper)
	if err != nil {
		return err
	}

	if env.GetType() != pb.Envelope_REQUEST_PEERINFO {
		return ErrNotExpectedMessage
	}
	return nil
}

// handShake 두번째 과정 함수. client 의 peer info 메세지를 server 에게 전달한다.
func handShakeSendInfo(streamWrapper bifrost.StreamWrapper, clientOpts ClientOpts) error {
	env, err := bifrost.BuildRequestPeerInfo(clientOpts.ip, clientOpts.pubKey)

	if err != nil {
		return err
	}
	err = streamWrapper.Send(env)

	if err != nil {
		return err
	}
	return nil
}

// handShake 세번째 과정 함수. server 의 peer info 메세지를 기다린다(Get 한다).
func handShakeGetServerInfo(streamWrapper bifrost.StreamWrapper) (*key.PubKey, error) {
	env, err := bifrost.RecvWithTimeout(3*time.Second, streamWrapper)

	if err != nil {
		return nil, err
	}

	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)

	if err != nil {
		return nil, err
	}

	serverPubKey, err := bifrost.ByteToPubKey(peerInfo.Pubkey, peerInfo.KeyGenOpt)

	if err != nil {
		return nil, err
	}

	return &serverPubKey, nil
}

func connect(conn *grpc.ClientConn) (bifrost.StreamWrapper, error) {

	streamWrapper, err := bifrost.NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
