package host

import (
	"testing"

	"google.golang.org/grpc"
)

type MockHello struct {
}

func (m MockHello) Register(*grpc.Server) {

}

type MockHello2 struct {
}

func (m MockHello2) Register(*grpc.Server) {

}

func Test_HostInterface(t *testing.T) {

}
