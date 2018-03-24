package stream

import (
	"sync"
	"sync/atomic"
	"github.com/it-chain/bifrost/msg"
	"errors"
	"github.com/it-chain/bifrost/pb"
)

type ReceivedMessageHandle func(message msg.OutterMessage)

type StreamHandler interface{
	Send(envelope *pb.Envelope, successCallBack func(interface{}),errCallBack func(error))
	Close()
}

type StreamHandlerImpl struct {
	streamWrapper StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandle
	outChannl     chan *msg.InnerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
}

func NewStreamHandler(streamWrapper StreamWrapper, handle ReceivedMessageHandle) (StreamHandler, error){

	if streamWrapper == nil || handle == nil{
		return nil, errors.New("fail to create streamHandler streamWrapper or handle is nil")
	}

	return &StreamHandlerImpl{
		streamWrapper: streamWrapper,
		handle: handle,
		outChannl: make(chan *msg.InnerMessage,200),
		readChannel: make(chan *pb.Envelope,200),
		stopChannel: make(chan struct{},1),
	}, nil
}

func (sh *StreamHandlerImpl) toDie() bool {
	return atomic.LoadInt32(&(sh.stopFlag)) == int32(1)
}

func (sh *StreamHandlerImpl) writeStream(){

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

func (sh *StreamHandlerImpl) readStream(errChan chan error){

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

func (sh *StreamHandlerImpl) Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)){

	sh.Lock()
	defer sh.Unlock()

	m := &msg.InnerMessage{
		Envelope: envelope,
		OnErr:    errCallBack,
		OnSuccess: successCallBack,
	}

	sh.outChannl <- m
}

func (sh *StreamHandlerImpl) Close(){

	if sh.toDie() {
		return
	}

	amIFirst := atomic.CompareAndSwapInt32(&sh.stopFlag, int32(0), int32(1))

	if !amIFirst {
		return
	}

	sh.stopChannel <- struct{}{}
	sh.Lock()

	sh.streamWrapper.Close()

	sh.Unlock()
}