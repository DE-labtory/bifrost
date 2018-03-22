package stream

import (
	"sync"
	"sync/atomic"
	"github.com/it-chain/bifrost/msg"
	"errors"
)

type ReceivedMessageHandle func(message msg.OutterMessage)

type StreamHandler struct {
	streamWrapper StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandle
	outChannl     chan msg.InnerMessage
	readChannel   chan interface{}
	stopChannel   chan struct{}
	sync.RWMutex
}



func (sh *StreamHandler) toDie() bool {
	return atomic.LoadInt32(&(sh.stopFlag)) == int32(1)
}

func (sh *StreamHandler) WriteStream(){

	for !sh.toDie() {

		stream := sh.streamWrapper.GetStream()

		if stream == nil {
			return
		}

		select{

		case m := <-sh.outChannl:
			err := stream.Send(m.Envelope)
			if err != nil {
				if m.OnErr != nil{
					go m.OnErr(err)
				}
			}else{
				if m.OnSuccess != nil{
					go m.OnSuccess("")
				}
			}
		case stop := <-sh.stopChannel:
			sh.stopChannel <- stop
			return
		}
	}
}

func (sh *StreamHandler) ReadStream(errChan chan error){

	defer func() {
		recover()
	}()

	for !sh.toDie() {

		stream := sh.streamWrapper.GetStream()

		if stream == nil {
			errChan <- errors.New("Stream is nil")
			return
		}

		envelope, err := stream.Recv()

		if sh.toDie() {
			return
		}

		if err != nil {
			errChan <- err
			return
		}

		sh.readChannel <- envelope
	}
}