package mux

import (
	"sync"

	"github.com/it-chain/bifrost/stream"
)

type Protocol string

type BiMux struct {
	sync.RWMutex
	registerHandled map[Protocol]*Handle
}

type Handle struct {
	Handler stream.ReceivedMessageHandler
}
