package main

import (
	"log"

	"os"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/heimdall/key"
)

var keyPath = "./SERVER_KEY_PATH"
var ip = "127.0.0.1:7777"

var DefaultMux *mux.DefaultMux

func main() {

	km, err := key.NewKeyManager(keyPath)
	defer os.RemoveAll(keyPath)

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	DefaultMux = mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	s := bifrost.NewServer(bifrost.KeyOpts{PriKey: pri, PubKey: pub})

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
	log.Fatalln(err.Error())
}
