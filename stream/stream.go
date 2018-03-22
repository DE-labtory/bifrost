package stream

import (
	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type Stream interface{
	Send(*pb.Envelope) error
	Recv() (*pb.Envelope, error)
}

type StreamWrapper interface{
	Close()
	GetStream() Stream
}

type CStreamWrapper struct{
	conn *grpc.ClientConn
	client         pb.StreamServiceClient
	clientStream   pb.StreamService_StreamClient
	cancel         context.CancelFunc
}

type SStreamWrapper struct{
	serverStream   pb.StreamService_StreamServer
	cancel         context.CancelFunc
}
