package main

import (
	"log"

	"crypto/elliptic"
	"crypto/rand"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/example"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/server"
)

var ip = "127.0.0.1:7777"

var DefaultMux *mux.DefaultMux

func main() {

	generator := example.SimpleGenerator{Curve: elliptic.P384(), Rand: rand.Reader}
	pri, err := generator.GenerateKey()

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

	formatter := example.SimpleFormatter{}
	idGetter := example.SimpleIdGetter{IDPrefix: "ITTEST", Formatter: &formatter}
	err = generator.StoreKey(pri, "", "./.key", idGetter.GetID(&pri.PublicKey))
	if err != nil {
		log.Fatal(err.Error())
	}

	keyLoader := example.SimpleKeyLoader{KeyDirPath: "./.key", KeyID: idGetter.GetID(&pri.PublicKey)}
	signer := example.SimpleSigner{KeyLoader: &keyLoader}
	verifier := example.SimpleVerifier{}
	crypto := bifrost.Crypto{IDGetter: &idGetter, Verifier: &verifier, Signer: &signer, Formatter: &formatter}

	s := server.New(bifrost.KeyOpts{PriKey: pri, PubKey: &pri.PublicKey}, crypto)

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
