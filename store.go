package bifrost

import (
	"errors"
	"sync"
)

var ErrConnAlreadyExist = errors.New("connection already exist in store")
var ErrConnNotExist = errors.New("connection not exist in store")

type ConnectionID string

type ConnectionStore struct {
	sync.RWMutex
	connMap map[ConnID]Connection
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connMap: make(map[ConnID]Connection),
	}
}

func (connStore ConnectionStore) AddConnection(conn Connection) error {
	connStore.Lock()
	defer connStore.Unlock()

	connID := conn.GetID()

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return ErrConnAlreadyExist
	}

	connStore.connMap[connID] = conn

	return nil
}

func (connStore ConnectionStore) DeleteConnection(connID ConnID) error {
	connStore.Lock()
	defer connStore.Unlock()

	conn, err := connStore.GetConnection(connID)

	if conn == nil {
		return err
	}

	conn.Close()
	delete(connStore.connMap, connID)

	return nil
}

func (connStore ConnectionStore) GetConnection(connID ConnID) (Connection, error) {

	conn, ok := connStore.connMap[connID]

	//exist
	if ok {
		return conn, nil
	}

	return nil, ErrConnNotExist
}
