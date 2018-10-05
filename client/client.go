package client

import (
	"time"

	"context"
	"encoding/json"
	"errors"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"

	"github.com/it-chain/iLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// handshake 과정에서 올바르지 않은 메세지 타입이 올 경우 발생하는 에러
var ErrNotExpectedMessage = errors.New("wrong message type")

// defaultDialTimeout
const (
	defaultDialTimeout = 3 * time.Second
)

// Server 와 연결시 사용되는 Client option
type ClientOpts struct {
	Ip     string
	PubKey bifrost.Key
}

// Server 와 연결시 사용되는 grpc option.
type GrpcOpts struct {
	TlsEnabled bool
	Creds      credentials.TransportCredentials
}

// 서버와 연결 요청. 실패시 err. handshake 과정을 거침.
func Dial(serverIp string, clientOpts ClientOpts, grpcOpts GrpcOpts, crypto bifrost.Crypto) (bifrost.Connection, error) {

	var opts []grpc.DialOption //required options

	if grpcOpts.TlsEnabled {
		opts = append(opts, grpc.WithTransportCredentials(grpcOpts.Creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	dialContext, _ := context.WithTimeout(context.Background(), defaultDialTimeout)
	gconn, err := grpc.DialContext(dialContext, serverIp, opts...)

	if err != nil {
		return nil, err
	}

	// create stream wrapper
	// inside stream wrapper, call main rpc service method BifrostStream()
	streamWrapper, err := bifrost.NewClientStreamWrapper(gconn)

	if err != nil {
		return nil, err
	}

	serverPubKey, err := handShake(streamWrapper, clientOpts, crypto.KeyRecoverer)

	if err != nil {
		return nil, err
	}

	conn, err := bifrost.NewConnection(serverIp, serverPubKey, streamWrapper, crypto)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// handshake 함수, return : serverPubKey, err
func handShake(streamWrapper bifrost.StreamWrapper, clientOpts ClientOpts, keyRecoverer bifrost.KeyRecoverer) (bifrost.Key, error) {

	err := waitServer(streamWrapper)

	if err != nil {
		iLogger.Infof(nil, "[Bifrost] Waiting server failed [%s]", err.Error())
		streamWrapper.Close()
		return nil, err
	}

	err = sendInfo(streamWrapper, clientOpts)
	if err != nil {
		iLogger.Infof(nil, "[Bifrost] Send info failed [%s]", err.Error())
		streamWrapper.Close()
		return nil, err
	}

	serverPubKey, err := getServerInfo(streamWrapper, keyRecoverer)

	if err != nil {
		iLogger.Infof(nil, "[Bifrost] Get server info failed [%s]", err.Error())
		streamWrapper.Close()
		return nil, err
	}

	iLogger.Info(nil, "[Bifrost] Handshake success")

	return serverPubKey, nil
}

// handshake 첫번째 과정 함수. server 의 request peer info 메세지를 기다린다.
func waitServer(streamWrapper bifrost.StreamWrapper) error {
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
func sendInfo(streamWrapper bifrost.StreamWrapper, clientOpts ClientOpts) error {
	env, err := bifrost.BuildResponsePeerInfo(clientOpts.PubKey)

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
func getServerInfo(streamWrapper bifrost.StreamWrapper, keyRecoverer bifrost.KeyRecoverer) (bifrost.Key, error) {
	env, err := bifrost.RecvWithTimeout(3*time.Second, streamWrapper)

	if err != nil {
		return nil, err
	}

	peerInfo := &bifrost.PeerInfo{}

	err = json.Unmarshal(env.Payload, peerInfo)

	if err != nil {
		return nil, err
	}

	serverPubKey, err := keyRecoverer.RecoverKeyFromByte(peerInfo.PubKeyBytes, peerInfo.IsPrivate)

	if err != nil {
		return nil, err
	}

	return serverPubKey, nil
}
