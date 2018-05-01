package bifrost

import "sync"

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

func (connStore ConnectionStore) AddConnection(conn Connection) {
	connStore.Lock()
	defer connStore.Unlock()

	connID := conn.GetID()

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return
	}

	connStore.connMap[connID] = conn
}

func (connStore ConnectionStore) DeleteConnection(connID ConnID) {
	connStore.Lock()
	defer connStore.Unlock()

	conn := connStore.GetConnection(connID)

	if conn == nil {
		return
	}

	conn.Close()
	delete(connStore.connMap, connID)
}

func (connStore ConnectionStore) GetConnection(connID ConnID) Connection {

	conn, ok := connStore.connMap[connID]

	//exist
	if ok {
		return conn
	}

	return nil
}
