package bifrost

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/auth"
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
	GetIP() string
	GetPeerKey() key.PubKey
	GetID() ConnID
	Start() error
}

type GrpcConnection struct {
	ID            ConnID
	key           key.PriKey
	peerKey       key.PubKey
	ip            string
	streamWrapper StreamWrapper
	stopFlag      int32
	handle        ReceivedMessageHandler
	outChannl     chan *innerMessage
	readChannel   chan *pb.Envelope
	stopChannel   chan struct{}
	sync.RWMutex
}

func NewConnection(ip string, priKey key.PriKey, peerKey key.PubKey, streamWrapper StreamWrapper) (Connection, error) {

	if streamWrapper == nil || peerKey == nil || priKey == nil {
		return nil, errors.New("fail to create connection streamWrapper or handle is nil")
	}

	return &GrpcConnection{
		ID:            FromPubKey(peerKey),
		key:           priKey,
		peerKey:       peerKey,
		ip:            ip,
		streamWrapper: streamWrapper,
		outChannl:     make(chan *innerMessage, 200),
		readChannel:   make(chan *pb.Envelope, 200),
		stopChannel:   make(chan struct{}, 1),
	}, nil
}

func (conn *GrpcConnection) GetIP() string {
	return conn.ip
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

	signedEnvelope, err := sign(envelope, conn.key)

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

//todo signer opts from config
func sign(envelope *pb.Envelope, priKey key.PriKey) (*pb.Envelope, error) {

	hash := sha512.New()
	hash.Write(envelope.Payload)
	digest := hash.Sum(nil)

	sig, err := auth.Sign(priKey, digest, auth.EQUAL_SHA512.SignerOptsToPSSOptions())

	if err != nil {
		return nil, err
	}

	envelope.Signature = sig

	return envelope, nil
}

//todo signer opts from config
func verify(envelope *pb.Envelope, pubkey key.PubKey) bool {

	b, _ := pubkey.ToPEM()

	if !bytes.Equal(envelope.Pubkey, b) {
		return false
	}

	hash := sha512.New()
	hash.Write(envelope.Payload)
	digest := hash.Sum(nil)

	flag, err := auth.Verify(pubkey, envelope.Signature, digest, auth.EQUAL_SHA512.SignerOptsToPSSOptions())

	if err != nil {
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
			if conn.handle != nil {
				conn.handle.ServeError(conn, err)
			}
			return err
		case message := <-conn.readChannel:
			if verify(message, conn.peerKey) {
				if conn.handle != nil {
					m := Message{Envelope: message, Conn: conn, Data: message.Payload}
					conn.handle.ServeRequest(m)
				}
			} else {
				//
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
