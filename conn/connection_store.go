package conn

import "sync"

type ConnectionID string

type ConnectionStore struct {
	sync.RWMutex
	connMap map[ConnectionID]Connection
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		connMap: make(map[ConnectionID]Connection),
	}
}

func (connStore ConnectionStore) AddConnection(conn Connection) {
	connStore.Lock()
	defer connStore.Unlock()

	connID := ConnectionID(conn.GetConnInfo().Id)

	_, ok := connStore.connMap[connID]

	//exist
	if ok {
		return
	}

	connStore.connMap[connID] = conn
}

func (connStore ConnectionStore) DeleteConnection(connID ConnectionID) {
	connStore.Lock()
	defer connStore.Unlock()

	conn := connStore.GetConnection(connID)

	if conn == nil {
		return
	}

	conn.Close()
	delete(connStore.connMap, connID)
}

func (connStore ConnectionStore) GetConnection(connID ConnectionID) Connection {
	connStore.Lock()
	defer connStore.Unlock()

	conn, ok := connStore.connMap[connID]

	//exist
	if ok {
		return conn
	}

	return nil
}
