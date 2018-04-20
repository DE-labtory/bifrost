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

	log.Println("Listen... on: [%s]", host.info.Address.IP)
	func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			s.Stop()
			lis.Close()
		}
	}()
}
