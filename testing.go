package bifrost

import (
	"fmt"
	"io"
	"net"

	"crypto/ecdsa"

	"crypto/elliptic"
	"crypto/rand"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/iLogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockGenerator struct {
}

func (generator MockGenerator) GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
}

type MockSigner struct {
}

func (signer *MockSigner) Sign(message []byte) ([]byte, error) {
	return []byte("signature"), nil
}

type MockVerifier struct {
}

func (verifier *MockVerifier) Verify(key *ecdsa.PublicKey, signature []byte, message []byte) (bool, error) {
	return true, nil
}

type MockFormatter struct {
}

func (formatter *MockFormatter) ToByte(*ecdsa.PublicKey) []byte {
	return []byte("byte format of ecdsa public key")
}

func (formatter *MockFormatter) FromByte([]byte, int) *ecdsa.PublicKey {
	return new(ecdsa.PublicKey)
}

func (formatter *MockFormatter) GetCurveOpt(pubKey *ecdsa.PublicKey) int {
	return *new(int)
}

type MockIdGetter struct {
}

func (idGetter *MockIdGetter) GetID(key *ecdsa.PublicKey) KeyID {
	return *new(KeyID)
}

type MockConnectionHandler func(stream pb.StreamService_BifrostStreamServer)
type MockRecvHandler func(envelope *pb.Envelope)
type MockCloseHandler func()

type MockServer struct {
	Rh  MockRecvHandler
	Ch  MockConnectionHandler
	Clh MockCloseHandler
}

type MockHandler struct{}

func (h MockHandler) ServeRequest(message Message) {

}

func (h MockHandler) ServeError(conn Connection, err error) {

}

func (ms MockServer) BifrostStream(stream pb.StreamService_BifrostStreamServer) error {

	if ms.Ch != nil {
		ms.Ch(stream)
	}

	for {
		envelope, err := stream.Recv()

		//fmt.Printf(err.Error())

		if err == io.EOF {
			return nil
		}

		if err != nil {
			if ms.Clh != nil {
				ms.Clh()
			}
			return err
		}

		if ms.Rh != nil {
			ms.Rh(envelope)
		}
	}
}

func ListenMockServer(mockServer pb.StreamServiceServer, ipAddress string) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", ipAddress)

	if err != nil {
		iLogger.Fatalf(nil, "Failed to listen: %v", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, mockServer)
	reflection.Register(s)

	fmt.Printf("listen..")

	go func() {
		if err := s.Serve(lis); err != nil {
			iLogger.Fatalf(nil, "Failed to serve: %v", err.Error())
			s.Stop()
			lis.Close()
		}
	}()

	return s, lis
}

func GetKeyOpts() KeyOpts {

	geneartor := MockGenerator{}

	pri, err := geneartor.GenerateKey()

	if err != nil {
		iLogger.Fatalf(nil, err.Error())
	}

	return KeyOpts{
		PubKey: &pri.PublicKey,
		PriKey: pri,
	}
}

func GetMockCrypto() Crypto {
	mockIDGetter := MockIdGetter{}
	mockFormatter := MockFormatter{}
	mockSigner := MockSigner{}
	mockVerifier := MockVerifier{}

	return Crypto{IDGetter: &mockIDGetter, Formatter: &mockFormatter, Signer: &mockSigner, Verifier: &mockVerifier}
}

func GetMockConnection(targetIP string) (Connection, error) {
	keyOpts := GetKeyOpts()
	mockCrypto := GetMockCrypto()

	mockStreamWrapper := MockStreamWrapper{}
	mockStreamWrapper.CloseCallBack = func() {

	}
	mockStreamWrapper.SendCallBack = func(envelope *pb.Envelope) {

	}

	conn, err := NewConnection(targetIP, keyOpts.PubKey, mockStreamWrapper, mockCrypto)
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

func (MockStreamWrapper) GetStream() Stream {
	panic("implement me")
}
