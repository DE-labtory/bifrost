package bifrost

import "sync"

type ConnectionID string

type ConnectionStore struct {
	sync.RWMutex
	connMap map[ID]Connection
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connMap: make(map[ID]Connection),
	}
}

func (connStore ConnectionStore) AddConnection(conn Connection) {
	connStore.Lock()
	defer connStore.Unlock()

	connID := ID(conn.GetConnInfo().Id)

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return
	}

	connStore.connMap[connID] = conn
}

func (connStore ConnectionStore) DeleteConnection(connID ID) {
	connStore.Lock()
	defer connStore.Unlock()

	conn := connStore.GetConnection(connID)

	if conn == nil {
		return
	}

	conn.Close()
	delete(connStore.connMap, connID)
}

func (connStore ConnectionStore) GetConnection(connID ID) Connection {

	conn, ok := connStore.connMap[connID]

	//exist
	if ok {
		return conn
	}

	return nil
}
