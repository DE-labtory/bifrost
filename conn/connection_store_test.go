package conn

import (
	"testing"

	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/stream"
	"github.com/stretchr/testify/assert"
)

type MockStreamWrapper struct {
}

func (m MockStreamWrapper) Send(*pb.Envelope) error {
	return nil
}
func (m MockStreamWrapper) Recv() (*pb.Envelope, error) {
	return nil, nil
}

func (m MockStreamWrapper) Close() {

}
func (m MockStreamWrapper) GetStream() stream.Stream {
	return nil
}

type MockReceivedHandler struct {
}

func (m MockReceivedHandler) ServeRequest(msg OutterMessage) {

}

func (m MockReceivedHandler) ServeError(conn Connection, err error) {

}

func TestNewConnectionStore(t *testing.T) {
	connStore := NewConnectionStore()
	assert.NotNil(t, connStore.connMap)
}

func TestConnectionStore_AddConnection(t *testing.T) {
	//given
	connStore := NewConnectionStore()
	msw := MockStreamWrapper{}
	mrh := MockReceivedHandler{}

	connInfo := ConnInfo{}
	connInfo.Id = ID("ASD")

	conn, err := NewConnection(connInfo, msw, mrh)

	if err != nil {

	}

	//when
	connStore.AddConnection(conn)

	//then
	assert.Equal(t, 1, len(connStore.connMap))
}

func TestConnectionStore_DeleteConnection(t *testing.T) {
	//given
	connStore := NewConnectionStore()
	msw := MockStreamWrapper{}
	mrh := MockReceivedHandler{}

	connInfo := ConnInfo{}
	connInfo.Id = ID("ASD")

	conn, err := NewConnection(connInfo, msw, mrh)

	if err != nil {

	}
	connStore.AddConnection(conn)

	//when
	connStore.DeleteConnection(conn.GetConnInfo().Id)

	//then
	assert.Equal(t, 0, len(connStore.connMap))
}

func TestConnectionStore_GetConnection(t *testing.T) {
	//given
	connStore := NewConnectionStore()
	msw := MockStreamWrapper{}
	mrh := MockReceivedHandler{}

	connInfo := ConnInfo{}
	connInfo.Id = ID("ASD")

	conn, err := NewConnection(connInfo, msw, mrh)

	if err != nil {

	}

	connStore.AddConnection(conn)

	//when
	fconn := connStore.GetConnection(conn.GetConnInfo().Id)

	//then
	assert.Equal(t, conn.GetConnInfo(), fconn.GetConnInfo())
}
