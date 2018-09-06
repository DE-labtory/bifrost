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

package bifrost

import (
	"log"
	"net"

	"github.com/it-chain/bifrost/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Listen(host *BifrostHost) {

	lis, err := net.Listen("tcp", host.info.Address.IP)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, host)
	reflection.Register(s)

	log.Printf("Listen... on: [%v]", host.info.Address.IP)
	func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()
}
