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
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultTimeout = time.Second * 3

type Address struct {
	IP string
}

func NewClientConn(address Address, tslEnabled bool, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption

	if tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(defaultTimeout))
	conn, err := grpc.Dial(address.IP, opts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}
