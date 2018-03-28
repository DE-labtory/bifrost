package stream

import (
	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
)

type Stream interface {
	Send(*pb.Envelope) error
	Recv() (*pb.Envelope, error)
}

func Connect(conn *grpc.ClientConn, handle ReceivedMessageHandler) (StreamHandler, error) {

	streamWrapper, err := NewClientStreamWrapper(conn)

	if err != nil {
		return nil, err
	}

	streamHandler, err := SetStreamHandler(streamWrapper, handle)

	go streamHandler.Start()

	return streamHandler, nil
}
