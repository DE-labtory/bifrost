package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DE-labtory/bifrost"
	"github.com/DE-labtory/bifrost/mocks"
	"github.com/DE-labtory/bifrost/mux"
	"github.com/DE-labtory/bifrost/server"
	"github.com/DE-labtory/iLogger"
)

var ip = "127.0.0.1:7777"

var DefaultMux *mux.DefaultMux
var testServerKeyDir = "./.test_server_key"

func main() {

	keyPair := mocks.NewMockKeyOpts()

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		iLogger.Infof(nil, "[Bifrost] %s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		iLogger.Infof(nil, "[Bifrost] %s", message.Data)
	})

	err := mocks.MockStoreKey(keyPair.PriKey, testServerKeyDir)
	if err != nil {
		iLogger.Fatal(nil, err.Error())
	}

	signer := mocks.MockECDSASigner{KeyID: keyPair.PubKey.ID(), KeyDirPath: testServerKeyDir}
	verifier := mocks.MockECDSAVerifier{}
	recoverer := mocks.MockECDSAKeyRecoverer{}
	crypto := bifrost.Crypto{Verifier: &verifier, Signer: &signer, KeyRecoverer: &recoverer}

	s := server.New(keyPair, crypto, nil)

	s.OnConnection(OnConnection)
	s.OnError(OnError)

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			os.RemoveAll(testServerKeyDir)
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
	iLogger.Fatalf(nil, err.Error())
}
