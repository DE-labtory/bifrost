package msg

import (
	"sync"
	"github.com/it-chain/bifrost/pb"
)

type InnerMessage struct{
	Envelope *pb.Envelope
	OnErr    func(error)
	OnSuccess func(interface{})
}

type OutterMessage struct {
	Envelope     *pb.Envelope
	Data         []byte
	ConnectionID string
	//Conn         connection
	sync.Mutex
}

type connection interface {
	Send(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error))
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
//func (m *OutterMessage) Respond(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {
//
//	m.Conn.Send(envelope, successCallBack, errCallBack)
//}