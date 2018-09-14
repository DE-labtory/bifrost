package main

import (
	"crypto/elliptic"
	"crypto/rand"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/example"
	"github.com/it-chain/bifrost/logger"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/server"
)

var ip = "127.0.0.1:7777"

var DefaultMux *mux.DefaultMux

func main() {

	generator := example.SimpleGenerator{Curve: elliptic.P384(), Rand: rand.Reader}
	pri, err := generator.GenerateKey()

	if err != nil {
		logger.Fatal(nil, err.Error())
	}

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		logger.Infof(nil, "%s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		logger.Infof(nil, "%s", message.Data)
	})

	formatter := example.SimpleFormatter{}
	idGetter := example.SimpleIdGetter{IDPrefix: "ITTEST", Formatter: &formatter}
	signer := example.SimpleSigner{PriKey: pri}
	verifier := example.SimpleVerifier{}
	s := server.New(bifrost.KeyOpts{PriKey: pri, PubKey: &pri.PublicKey}, &idGetter, &formatter, &signer, &verifier)

	s.OnConnection(OnConnection)
	s.OnError(OnError)

	s.Listen(ip)
}

func OnConnection(connection bifrost.Connection) {

	connection.Handle(DefaultMux)
	defer connection.Close()

	if err := connection.Start(); err != nil {
		connection.Close()
	}
}

func OnError(err error) {
	logger.Fatal(nil, err.Error())
}
