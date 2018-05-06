package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/heimdall/key"
)

var clientIp = "127.0.0.1:7778"
var serverIp = "127.0.0.1:7777"
var DefaultMux *mux.DefaultMux

//todo 아직깔끔하지않음 여러 수정필요
func main() {

	km, err := key.NewKeyManager("")

	if err != nil {
		log.Fatal(err.Error())
	}

	pri, pub, err := km.GenerateKey(key.RSA4096)

	if err != nil {
		log.Fatal(err.Error())
	}

	DefaultMux := mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		log.Printf("%s", message.Data)
	})

	clientOpt := client.ClientOpts{
		Ip:     clientIp,
		PriKey: pri,
		PubKey: pub,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	conn, err := client.Dial(serverIp, clientOpt, grpcOpt)

	if err != nil {
		log.Fatal(err.Error())
	}

	conn.Handle(DefaultMux)

	go func() {
		if err := conn.Start(); err != nil {
			conn.Close()
		}
	}()

	conn.Send([]byte("client join!!"), "join", nil, nil)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		conn.Send([]byte(text), "chat", nil, nil)
	}
}
