package conn

import (
	"errors"
	"sync"

	"sync/atomic"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/stream"
)

type ReceivedMessageHandler interface {
	ServeRequest(msg OutterMessage)
	ServeError(err error)
}

type Connection interface {
	Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error))
	Close()
	GetConnInfo() ConnenctionInfo
}

type GrpcConnection struct {
	connInfo      ConnenctionInfo
	streamWrapper stream.StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandler
	outChannl     chan *InnerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
}

func NewConnection(connInfo ConnenctionInfo, streamWrapper stream.StreamWrapper, handle ReceivedMessageHandler) (Connection, error) {

	if streamWrapper == nil || handle == nil {
		return nil, errors.New("fail to create connection streamWrapper or handle is nil")
	}

	return &GrpcConnection{
		connInfo:      connInfo,
		streamWrapper: streamWrapper,
		handle:        handle,
		outChannl:     make(chan *InnerMessage, 200),
		readChannel:   make(chan *pb.Envelope, 200),
		stopChannel:   make(chan struct{}, 1),
	}, nil
}

func (conn *GrpcConnection) GetConnInfo() ConnenctionInfo {
	return conn.connInfo
}

func (conn *GrpcConnection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *GrpcConnection) writeStream() {

	for !conn.toDie() {

		select {

		case m := <-conn.outChannl:
			err := conn.streamWrapper.Send(m.Envelope)
			if err != nil {
				if m.OnErr != nil {
					go m.OnErr(err)
				}
			} else {
				if m.OnSuccess != nil {
					go m.OnSuccess("")
				}
			}
		case stop := <-conn.stopChannel:
			conn.stopChannel <- stop
			return
		}
	}
}

func (conn *GrpcConnection) readStream(errChan chan error) {

	defer func() {
		recover()
	}()

	for !conn.toDie() {

		envelope, err := conn.streamWrapper.Recv()

		if conn.toDie() {
			return
		}

		if err != nil {
			errChan <- err
			return
		}

		conn.readChannel <- envelope
	}
}

func (conn *GrpcConnection) Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	conn.Lock()
	defer conn.Unlock()

	m := &InnerMessage{
		Envelope:  envelope,
		OnErr:     errCallBack,
		OnSuccess: successCallBack,
	}

	conn.outChannl <- m
}

func (conn *GrpcConnection) Close() {

	if conn.toDie() {
		return
	}

	amIFirst := atomic.CompareAndSwapInt32(&conn.stopFlag, int32(0), int32(1))

	if !amIFirst {
		return
	}

	conn.stopChannel <- struct{}{}
	conn.Lock()

	conn.streamWrapper.Close()

	conn.Unlock()
}

func (conn *GrpcConnection) Start() error {

	errChan := make(chan error, 1)

	go conn.readStream(errChan)
	go conn.writeStream()

	for !conn.toDie() {
		select {
		case stop := <-conn.stopChannel:
			conn.stopChannel <- stop
			return nil
		case err := <-errChan:
			if conn.handle != nil {
				conn.handle.ServeError(err)
			}
			return err
		case message := <-conn.readChannel:
			if conn.handle != nil {
				conn.handle.ServeRequest(OutterMessage{Envelope: message, Conn: conn})
			}
		}
	}

	return nil
}
