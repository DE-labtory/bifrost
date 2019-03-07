package bifrost_test

import (
	"testing"

	"github.com/DE-labtory/bifrost"
	"github.com/DE-labtory/bifrost/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConnectionStore_AddConnection(t *testing.T) {
	testConnStroe := bifrost.NewConnectionStore()
	testConn, err := mocks.NewMockConnection("127.0.0.1:1234")
	assert.NoError(t, err)

	err = testConnStroe.AddConnection(testConn)
	assert.NoError(t, err)
}

func TestConnectionStore_DeleteConnection(t *testing.T) {
	testConnStore := bifrost.NewConnectionStore()
	testConn, err := mocks.NewMockConnection("127.0.0.1:1234")
	assert.NoError(t, err)

	err = testConnStore.AddConnection(testConn)
	assert.NoError(t, err)

	err = testConnStore.DeleteConnection("wrong ID")
	assert.Error(t, err)

	err = testConnStore.DeleteConnection(testConn.GetID())
	assert.NoError(t, err)
}

func TestConnectionStore_GetConnection(t *testing.T) {
	testConnStore := bifrost.NewConnectionStore()
	testConn, err := mocks.NewMockConnection("127.0.0.1:1234")
	assert.NoError(t, err)

	err = testConnStore.AddConnection(testConn)
	assert.NoError(t, err)

	loadedConn, err := testConnStore.GetConnection(testConn.GetID())
	assert.NoError(t, err)
	assert.Equal(t, testConn, loadedConn)
}
