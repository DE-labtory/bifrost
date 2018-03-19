package Bifrost

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
)

type Host interface{
	Register(*grpc.Server)
}

type Address struct{
	Ip string
}

func NewAddress(ipAddress string) Address{
	// validate ip pattern
	return Address{
		Ip:ipAddress,
	}
}

type Bifrost struct{
	server *grpc.Server
}

func NewHost(address Address) Bifrost{
	lis, err := net.Listen("tcp", address.Ip)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	s.RegisterService()
	reflection.Register(s)

	return Bifrost{
		server:s,
	}
}
