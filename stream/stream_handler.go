package stream

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/it-chain/bifrost/pb"
)

type ReceivedMessageHandler interface {
	ServeRequest(msg OutterMessage)
	ServeError(err error)
}

type StreamHandler interface {
	Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error))
	Start() error
	Close()
}

type StreamHandlerImpl struct {
	streamWrapper StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandler
	outChannl     chan *InnerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
}

func SetStreamHandler(streamWrapper StreamWrapper, handle ReceivedMessageHandler) (StreamHandler, error) {

	if streamWrapper == nil || handle == nil {
		return nil, errors.New("fail to create streamHandler streamWrapper or handle is nil")
	}

	sh := &StreamHandlerImpl{
		streamWrapper: streamWrapper,
		handle:        handle,
		outChannl:     make(chan *InnerMessage, 200),
		readChannel:   make(chan *pb.Envelope, 200),
		stopChannel:   make(chan struct{}, 1),
	}

	return sh, nil
}

func (sh *StreamHandlerImpl) toDie() bool {
	return atomic.LoadInt32(&(sh.stopFlag)) == int32(1)
}

func (sh *StreamHandlerImpl) writeStream() {

	for !sh.toDie() {

		select {

		case m := <-sh.outChannl:
			err := sh.streamWrapper.Send(m.Envelope)
			if err != nil {
				if m.OnErr != nil {
					go m.OnErr(err)
				}
			} else {
				if m.OnSuccess != nil {
					go m.OnSuccess("")
				}
			}
		case stop := <-sh.stopChannel:
			sh.stopChannel <- stop
			return
		}
	}
}

func (sh *StreamHandlerImpl) readStream(errChan chan error) {

	defer func() {
		recover()
	}()

	for !sh.toDie() {

		envelope, err := sh.streamWrapper.Recv()

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

func (sh *StreamHandlerImpl) Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	sh.Lock()
	defer sh.Unlock()

	m := &InnerMessage{
		Envelope:  envelope,
		OnErr:     errCallBack,
		OnSuccess: successCallBack,
	}

	sh.outChannl <- m
}

func (sh *StreamHandlerImpl) Close() {

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

func (sh *StreamHandlerImpl) Start() error {

	errChan := make(chan error, 1)

	go sh.readStream(errChan)
	go sh.writeStream()

	for !sh.toDie() {
		select {
		case stop := <-sh.stopChannel:
			sh.stopChannel <- stop
			return nil
		case err := <-errChan:
			if sh.handle != nil {
				sh.handle.ServeError(err)
			}
			return err
		case message := <-sh.readChannel:
			if sh.handle != nil {
				sh.handle.ServeRequest(OutterMessage{Envelope: message, Stream: sh})
			}
		}
	}

	return nil
}
