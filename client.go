package bifrost

import (
	"time"

	"github.com/it-chain/heimdall/key"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultTimeout = time.Second * 3

type Client struct {
}

type GrpcOpts struct {
	ip         string
	tslEnabled bool
	creds      credentials.TransportCredentials
}

type KeyOpts struct {
	priKey key.PriKey
	pubKey key.PubKey
}

func Dial(grpcOpts GrpcOpts, key KeyOpts) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption

	if grpcOpts.tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(grpcOpts.creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	gconn, err := grpc.Dial(grpcOpts.ip, opts...)

	if err != nil {
		return nil, err
	}

	connect(gconn)

	streamWrapper, err := NewClientStreamWrapper(gconn)

	if err != nil {
		return nil, err
	}

	conn := NewConnection()

	return conn, err
}

func connect(conn *grpc.ClientConn) (StreamWrapper, error) {

	streamWrapper, err := NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
