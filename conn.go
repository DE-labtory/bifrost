package bifrost

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"crypto/ecdsa"

	"github.com/it-chain/bifrost/logger"
	"github.com/it-chain/bifrost/pb"
)

type ConnID = string

type PeerInfo struct {
	IP       string
	Pubkey   []byte
	CurveOpt int
}

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
func (m *Message) Respond(data []byte, protocol string, successCallBack func(interface{}), errCallBack func(error)) {

	m.Conn.Send(data, protocol, successCallBack, errCallBack)
}

type Handler interface {
	ServeRequest(msg Message)
}

type Connection interface {
	Send(data []byte, protocol string, successCallBack func(interface{}), errCallBack func(error))
	Close()
	GetIP() Address
	GetPeerKey() *ecdsa.PublicKey
	GetID() ConnID
	Start() error
	Handle(handler Handler)
}

type GrpcConnection struct {
	ID            ConnID
	peerKey       *ecdsa.PublicKey
	ip            Address
	streamWrapper StreamWrapper
	stopFlag      int32
	handler       Handler
	outChannl     chan *innerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
	Crypto
}

func NewConnection(ip string, peerKey *ecdsa.PublicKey, streamWrapper StreamWrapper,
	crypto Crypto) (Connection, error) {

	if streamWrapper == nil || peerKey == nil {
		return nil, errors.New("fail to create connection streamWrapper or peerKey is nil")
	}

	ipAddr, err := ToAddress(ip)
	if err != nil {
		return nil, err
	}

	return &GrpcConnection{
		ID:            ConnID(crypto.GetID(peerKey)),
		peerKey:       peerKey,
		ip:            ipAddr,
		streamWrapper: streamWrapper,
		outChannl:     make(chan *innerMessage, 200),
		readChannel:   make(chan *pb.Envelope, 200),
		stopChannel:   make(chan struct{}, 1),
		Crypto:        crypto,
	}, nil
}

func (conn *GrpcConnection) GetIP() Address {
	return conn.ip
}
func (conn *GrpcConnection) GetPeerKey() *ecdsa.PublicKey {
	return conn.peerKey
}
func (conn *GrpcConnection) GetID() ConnID {
	return conn.ID
}

func (conn *GrpcConnection) toDie() bool {
	return atomic.LoadInt32(&(conn.stopFlag)) == int32(1)
}

func (conn *GrpcConnection) Handle(handler Handler) {
	conn.handler = handler
}

func (conn *GrpcConnection) Send(payload []byte, protocol string, successCallBack func(interface{}), errCallBack func(error)) {

	conn.Lock()
	defer conn.Unlock()

	signedEnvelope, err := conn.build(protocol, payload)

	if err != nil {
		go errCallBack(errors.New(fmt.Sprintf("fail to sign envelope [%s]", err.Error())))
		return
	}

	m := &innerMessage{
		Envelope:  signedEnvelope,
		OnErr:     errCallBack,
		OnSuccess: successCallBack,
	}

	conn.outChannl <- m
}

func (conn *GrpcConnection) build(protocol string, payload []byte) (*pb.Envelope, error) {

	sig, err := conn.Sign(payload)
	if err != nil {
		return nil, err
	}

	envelope := &pb.Envelope{}
	envelope.Signature = sig
	envelope.Payload = payload
	envelope.Type = pb.Envelope_NORMAL
	envelope.Protocol = protocol
	envelope.Pubkey = []byte("key")

	return envelope, nil
}

func (conn *GrpcConnection) Verify(envelope *pb.Envelope) bool {
	flag, err := conn.Crypto.Verify(conn.peerKey, envelope.Signature, envelope.Payload)

	if err != nil {
		logger.Info(nil, fmt.Sprintf("[Bifrost] %s", err.Error()))
		return false
	}

	return flag
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
			return err
		case message := <-conn.readChannel:
			if conn.Verify(message) {
				if conn.handler != nil {
					m := Message{Envelope: message, Conn: conn, Data: message.Payload}
					conn.handler.ServeRequest(m)
				}
			} else {
				// todo: verify 결과 false인 경우
			}
		}
	}

	return nil
}
