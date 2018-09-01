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
