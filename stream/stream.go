package stream

import (
	"github.com/it-chain/bifrost/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Stream interface{
	Send(*pb.Envelope) error
	Recv() (*pb.Envelope, error)
}

type StreamWrapper interface{
	Close()
	GetStream() Stream
}

type CStreamWrapper struct {
	conn         *grpc.ClientConn
	client       pb.StreamServiceClient
	clientStream pb.StreamService_StreamClient
	cancel       context.CancelFunc
}


func NewClientStreamWrapper(conn *grpc.ClientConn,client pb.StreamServiceClient, clientStream pb.StreamService_StreamClient,cancel context.CancelFunc) StreamWrapper{
	return &CStreamWrapper{
		cancel:cancel,
		conn: conn,
		clientStream: clientStream,
		client: client,
	}
}

func (csw *CStreamWrapper) GetStream() Stream{
	return csw.clientStream
}

func (csw *CStreamWrapper) Close(){
	csw.conn.Close()
	csw.clientStream.CloseSend()
	csw.cancel()
}

type SStreamWrapper struct{
	serverStream   pb.StreamService_StreamServer
	cancel         context.CancelFunc
}


func NewServerStreamWrapper(serverStream pb.StreamService_StreamServer, cancel context.CancelFunc) StreamWrapper{
	return &SStreamWrapper{
		cancel:cancel,
		serverStream:serverStream,
	}
}

func (ssw *SStreamWrapper) GetStream() Stream{
	return ssw.serverStream
}

func (ssw *SStreamWrapper) Close(){
	ssw.cancel()
}