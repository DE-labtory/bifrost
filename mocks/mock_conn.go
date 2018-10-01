package mocks

import (
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/pb"
)

func NewMockConnection(targetIP string) (bifrost.Connection, error) {
	keyOpts := NewMockKeyOpts()
	mockCrypto := NewMockCrypto()

	mockStreamWrapper := MockStreamWrapper{}
	mockStreamWrapper.CloseCallBack = func() {

	}
	mockStreamWrapper.SendCallBack = func(envelope *pb.Envelope) {

	}

	conn, err := bifrost.NewConnection(targetIP, keyOpts.PubKey, mockStreamWrapper, mockCrypto)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type SendCallBack func(envelope *pb.Envelope)
type CloseCallBack func()

type MockStreamWrapper struct {
	SendCallBack  SendCallBack
	CloseCallBack CloseCallBack
}

func (msw MockStreamWrapper) Send(envelope *pb.Envelope) error {
	msw.SendCallBack(envelope)
	return nil
}

func (MockStreamWrapper) Recv() (*pb.Envelope, error) {
	panic("implement me")
}

func (msw MockStreamWrapper) Close() {
	msw.CloseCallBack()
}

func (MockStreamWrapper) GetStream() bifrost.Stream {
	panic("implement me")
}
