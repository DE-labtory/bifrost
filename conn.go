package Bifrost

import (
	"sync"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type Connection interface{

}

type ConnectionImpl struct {
	conn           *grpc.ClientConn
	cancel         context.CancelFunc
	stopFlag       int32
	connectionID   string
	handle         ReceiveMessageHandle
	outChannl      chan *msg.InnerMessage
	readChannel    chan *message.Envelope
	stopChannel    chan struct{}
	sync.RWMutex
}