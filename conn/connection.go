package conn

import (
	"sync"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/stream"
)

type Connection interface {
	Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error))
}

type GrpcConnection struct {
	connInfo      ConnenctionInfo
	streamHandler stream.StreamHandler
	sync.RWMutex
}

func NewConnection(connInfo ConnenctionInfo, streamHandler stream.StreamHandler) Connection {

	streamHandler.Start()

	return &GrpcConnection{
		connInfo:      connInfo,
		streamHandler: streamHandler,
	}
}

func (conn *GrpcConnection) Close() {
	conn.streamHandler.Close()
}

func (conn *GrpcConnection) Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {
	conn.streamHandler.Send(envelope, successCallBack, errCallBack)
}
