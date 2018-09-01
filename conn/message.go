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
	"sync"

	"github.com/it-chain/bifrost/pb"
)

type InnerMessage struct {
	Envelope  *pb.Envelope
	OnErr     func(error)
	OnSuccess func(interface{})
}

type OutterMessage struct {
	Envelope *pb.Envelope
	Data     []byte
	Conn     Connection
	sync.Mutex
}

// Respond sends a msg to the source that sent the ReceivedMessageImpl
func (m *OutterMessage) Respond(envelope *pb.Envelope, successCallBack func(interface{}), errCallBack func(error)) {

	m.Conn.Send(envelope, successCallBack, errCallBack)
}
