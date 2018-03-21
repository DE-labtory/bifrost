package Bifrost

import (
	"sync"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/it-chain/bifrost/pb"
	"sync/atomic"
)

//type Stream interface{
//	Send(*pb.Envelope) error
//	Recv() (*pb.Envelope, error)
//}

type StreamHandler struct{
	stopFlag       int32
	//handle         ReceiveMessageHandle
	outChannl      chan interface{}
	readChannel    chan interface{}
	stopChannel    chan struct{}
	sync.RWMutex
}

type ClientStream struct{
	conn *grpc.ClientConn
	client         pb.StreamServiceClient
	clientStream   pb.StreamService_StreamClient
	cancel         context.CancelFunc
	StreamHandler
}

type ServerStream struct{
	serverStream   pb.StreamService_StreamServer
	cancel         context.CancelFunc
	StreamHandler
}

func (sh *StreamHandler) toDie() bool {
	return atomic.LoadInt32(&(sh.stopFlag)) == int32(1)
}

func (sh *StreamHandler) WriteStream(stream grpc.Stream){

	for !sh.toDie() {

		if stream == nil {
			return
		}

		select{

		case m := <-sh.outChannl:
			err := stream.SendMsg(m)
			//if err != nil {
			//	if m.OnErr != nil{
			//		go m.OnErr(err)
			//	}
			//	return
			//}else{
			//	if m.OnSuccess != nil{
			//		go m.OnSuccess("")
			//	}
			//}
		case stop := <-sh.stopChannel:
			sh.stopChannel <- stop
			return
		}
	}
}

//type ConnectionImpl struct {
//
//	cancel         context.CancelFunc
//	stopFlag       int32
//	connectionID   string
//	handle         ReceiveMessageHandle
//	outChannl      chan *msg.InnerMessage
//	readChannel    chan *message.Envelope
//	stopChannel    chan struct{}
//	sync.RWMutex
//}