/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conn

import "sync"

type ConnectionID string

type ConnectionStore struct {
	sync.RWMutex
	ConnMap map[ID]Connection
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		ConnMap: make(map[ID]Connection),
	}
}

func (connStore ConnectionStore) AddConnection(conn Connection) {
	connStore.Lock()
	defer connStore.Unlock()

	connID := ID(conn.GetConnInfo().Id)

	_, ok := connStore.ConnMap[connID]

	//exist
	if ok {
		return
	}

	connStore.ConnMap[connID] = conn
}

func (connStore ConnectionStore) DeleteConnection(connID ID) {
	connStore.Lock()
	defer connStore.Unlock()

	conn := connStore.GetConnection(connID)

	if conn == nil {
		return
	}

	conn.Close()
	delete(connStore.ConnMap, connID)
}

func (connStore ConnectionStore) GetConnection(connID ID) Connection {

	conn, ok := connStore.ConnMap[connID]

	//exist
	if ok {
		return conn
	}

	return nil
}
