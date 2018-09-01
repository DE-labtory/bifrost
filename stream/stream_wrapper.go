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

package stream

import (
	"context"

	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
)

type StreamWrapper interface {
	Stream
	Close()
	GetStream() Stream
}

type CStreamWrapper struct {
	conn         *grpc.ClientConn
	client       pb.StreamServiceClient
	clientStream pb.StreamService_StreamClient
	cancel       context.CancelFunc
}

//client stream wrapper
func NewClientStreamWrapper(conn *grpc.ClientConn) (StreamWrapper, error) {

	ctx, cf := context.WithCancel(context.Background())
	streamServiceClient := pb.NewStreamServiceClient(conn)
	clientStream, err := streamServiceClient.Stream(ctx)

	if err != nil {
		return nil, err
	}

	return &CStreamWrapper{
		cancel:       cf,
		conn:         conn,
		clientStream: clientStream,
		client:       streamServiceClient,
	}, nil
}

func (csw *CStreamWrapper) GetStream() Stream {
	return csw.clientStream
}

func (csw *CStreamWrapper) Send(envelope *pb.Envelope) error {
	return csw.clientStream.Send(envelope)
}

func (csw *CStreamWrapper) Recv() (*pb.Envelope, error) {
	return csw.clientStream.Recv()
}

func (csw *CStreamWrapper) Close() {
	csw.conn.Close()
	csw.clientStream.CloseSend()
	csw.cancel()
}

//server stream wrapper
type SStreamWrapper struct {
	serverStream pb.StreamService_StreamServer
	cancel       context.CancelFunc
}

func NewServerStreamWrapper(serverStream pb.StreamService_StreamServer, cancel context.CancelFunc) StreamWrapper {
	return &SStreamWrapper{
		cancel:       cancel,
		serverStream: serverStream,
	}
}

func (ssw *SStreamWrapper) GetStream() Stream {
	return ssw.serverStream
}

func (ssw *SStreamWrapper) Close() {
	ssw.cancel()
}

func (ssw *SStreamWrapper) Send(envelope *pb.Envelope) error {
	return ssw.serverStream.Send(envelope)
}

func (ssw *SStreamWrapper) Recv() (*pb.Envelope, error) {
	return ssw.serverStream.Recv()
}
