package bifrost

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultTimeout = time.Second * 3

//type Address struct {
//	IP string
//}

func NewClientConn(ip string, tslEnabled bool, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption

	if tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	conn, err := grpc.Dial(ip, opts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func Connect(conn *grpc.ClientConn) (StreamWrapper, error) {

	streamWrapper, err := NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
