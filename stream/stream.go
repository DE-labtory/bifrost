package stream

import (
	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
)

type Stream interface {
	Send(*pb.Envelope) error
	Recv() (*pb.Envelope, error)
}

func Connect(conn *grpc.ClientConn, handle ReceivedMessageHandler) (StreamWrapper, error) {

	streamWrapper, err := NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	return streamWrapper, nil
}
