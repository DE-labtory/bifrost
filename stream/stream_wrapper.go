package stream

import (
	"context"

	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
)

type StreamWrapper interface {
	Close()
	GetStream() Stream
}

type CStreamWrapper struct {
	conn         *grpc.ClientConn
	client       pb.StreamServiceClient
	clientStream pb.StreamService_StreamClient
	cancel       context.CancelFunc
}

//client stream wrapper
func NewClientStreamWrapper(conn *grpc.ClientConn) (StreamWrapper, error) {

	ctx, cf := context.WithCancel(context.Background())
	streamServiceClient := pb.NewStreamServiceClient(conn)
	clientStream, err := streamServiceClient.Stream(ctx)

	if err != nil {
		return nil, err
	}

	return &CStreamWrapper{
		cancel:       cf,
		conn:         conn,
		clientStream: clientStream,
		client:       streamServiceClient,
	}, nil
}

func (csw *CStreamWrapper) GetStream() Stream {
	return csw.clientStream
}

func (csw *CStreamWrapper) Close() {
	csw.conn.Close()
	csw.clientStream.CloseSend()
	csw.cancel()
}

//server stream wrapper
type SStreamWrapper struct {
	serverStream pb.StreamService_StreamServer
	cancel       context.CancelFunc
}

func NewServerStreamWrapper(serverStream pb.StreamService_StreamServer, cancel context.CancelFunc) StreamWrapper {
	return &SStreamWrapper{
		cancel:       cancel,
		serverStream: serverStream,
	}
}

func (ssw *SStreamWrapper) GetStream() Stream {
	return ssw.serverStream
}

func (ssw *SStreamWrapper) Close() {
	ssw.cancel()
}
