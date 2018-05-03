package bifrost

import (
	"time"

	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"context"
	"github.com/it-chain/bifrost/mux"
)

// defaultDialTimeout
const (
	defaultDialTimeout   = 3 * time.Second
	REQUEST_CONNINFO     = "/requestConnInfo"
	CONNECTION_ESTABLISH = "/connectionEstablish"
)

// Client struct. 하나의 연결당 하나의 client 가 필요.
type Client struct {
	ServerIp  string
	ServerKey key.PubKey
	clientKey key.PriKey
	conn      Connection
	mux       mux.Mux
}

// Client struct 를 생성하는 함수. 연결할 ip, 클라이언트의 private key, message 처리 mux 가 필요함
func NewClient(serverIp string, clientKey key.PriKey, mux mux.Mux) *Client {
	return &Client{
		ServerIp:  serverIp,
		clientKey: clientKey,
		mux: mux,
	}
}

// Server 와 연결시 사용되는 option.
type GrpcOpts struct {
	tlsEnabled bool
	creds      credentials.TransportCredentials
}

// 서버와 연결 요청. 실패시 err. handshake 과정을 거침.
func (c Client) ConnectToServer(grpcOpts GrpcOpts) error {

	var opts []grpc.DialOption

	if grpcOpts.tlsEnabled {
		opts = append(opts, grpc.WithTransportCredentials(grpcOpts.creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	dialContext, _ := context.WithTimeout(context.Background(), defaultDialTimeout)
	gconn, err := grpc.DialContext(dialContext, c.ServerIp, opts...)

	if err != nil {
		return err
	}

	streamWrapper, err := connect(gconn)

	if err != nil {
		return err
	}

	
	return err
}

func connect(conn *grpc.ClientConn) (StreamWrapper, error) {

	streamWrapper, err := NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
