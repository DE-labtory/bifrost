package bifrost

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/key"
)

type ConnID = string

type innerMessage struct {
	Envelope  *pb.Envelope
	OnErr     func(error)
	OnSuccess func(interface{})
}

type Message struct {
	Envelope *pb.Envelope
	Data     []byte
	Conn     Connection
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
func (m *Message) Respond(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	m.Conn.Send(envelope, successCallBack, errCallBack)
}

type ReceivedMessageHandler interface {
	ServeRequest(msg Message)
	ServeError(conn Connection, err error)
}

type Connection interface {
	Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error))
	Close()
	GetAddress() Address
	GetPeerKey() key.PubKey
	GetID() ConnID
	Start() error
}

type GrpcConnection struct {
	ID            ConnID
	peerKey       key.PubKey
	address       Address
	streamWrapper StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandler
	outChannl     chan *innerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
}

func NewConnection(address Address, pubkey key.PubKey, streamWrapper StreamWrapper, handle ReceivedMessageHandler) (Connection, error) {

	if streamWrapper == nil || handle == nil || pubkey == nil {
		return nil, errors.New("fail to create connection streamWrapper or handle is nil")
	}

	return &GrpcConnection{
		ID:            FromPubKey(pubkey),
		peerKey:       pubkey,
		address:       address,
		streamWrapper: streamWrapper,
		handle:        handle,
		outChannl:     make(chan *innerMessage, 200),
		readChannel:   make(chan *pb.Envelope, 200),
		stopChannel:   make(chan struct{}, 1),
	}, nil
}

func (conn *GrpcConnection) GetAddress() Address {
	return conn.address
}
func (conn *GrpcConnection) GetPeerKey() key.PubKey {
	return conn.peerKey
}
func (conn *GrpcConnection) GetID() ConnID {
	return conn.ID
}

func (conn *GrpcConnection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *GrpcConnection) Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	conn.Lock()
	defer conn.Unlock()

	m := &innerMessage{
		Envelope:  envelope,
		OnErr:     errCallBack,
		OnSuccess: successCallBack,
	}

	conn.outChannl <- m
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
				conn.handle.ServeError(conn, err)
			}
			return err
		case message := <-conn.readChannel:
			if conn.handle != nil {
				m := Message{Envelope: message, Conn: conn, Data: message.Payload}
				conn.handle.ServeRequest(m)
			}
		}
	}

	return nil
}

//
//func NewConnInfo(id string, address Address, pubKey key.PubKey) ConnInfo {
//	return ConnInfo{
//		Id:      id,
//		Address: address,
//		PeerKey: pubKey,
//	}
//}

//
//type PublicConnInfo struct {
//	Id        string
//	Address   Address
//	Pubkey    []byte
//	KeyType   key.KeyType
//	KeyGenOpt key.KeyGenOpts
//}
