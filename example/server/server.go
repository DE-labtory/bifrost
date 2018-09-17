package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/example"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/server"
	"github.com/it-chain/engine/common/logger"
)

var ip = "127.0.0.1:7777"

var DefaultMux *mux.DefaultMux

func main() {

	generator := example.SimpleGenerator{Curve: elliptic.P384(), Rand: rand.Reader}
	pri, err := generator.GenerateKey()

	if err != nil {
		logger.Fatalf(nil, err.Error())
	}

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		logger.Info(nil, fmt.Sprintf("[Bifrost] %s", message.Data))
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		logger.Info(nil, fmt.Sprintf("[Bifrost] %s", message.Data))
	})

	formatter := example.SimpleFormatter{}
	idGetter := example.SimpleIdGetter{IDPrefix: "ITTEST", Formatter: &formatter}
	err = generator.StoreKey(pri, "", "./.key", idGetter.GetID(&pri.PublicKey))
	if err != nil {
		logger.Fatalf(nil, err.Error())
	}

	keyLoader := example.SimpleKeyLoader{KeyDirPath: "./.key", KeyID: idGetter.GetID(&pri.PublicKey)}
	signer := example.SimpleSigner{KeyLoader: &keyLoader}
	verifier := example.SimpleVerifier{}
	crypto := bifrost.Crypto{IDGetter: &idGetter, Verifier: &verifier, Signer: &signer, Formatter: &formatter}

	s := server.New(bifrost.KeyOpts{PriKey: pri, PubKey: &pri.PublicKey}, crypto)

	s.OnConnection(OnConnection)
	s.OnError(OnError)

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			os.RemoveAll("./.key")
			os.Exit(0)
		}
	}()

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
	logger.Fatalf(nil, err.Error())
}
