package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/heimdall/key"
)

func CreateHost(ip string, mux *mux.Mux) *bifrost.BifrostHost {

	km, err := key.NewKeyManager("~/key")

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer os.RemoveAll("~/key")

	priv, pub, err := km.GenerateKey(key.RSA4096)

	myconnectionInfo := bifrost.NewHostInfo(conn.Address{IP: ip}, pub, priv)

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
	}

	return bifrost.New(myconnectionInfo, mux, OnConnectionHandler)
}

func ReadFromConsole() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {

	mux := mux.NewMux()

	mux.Handle("/echo/1.0", func(message conn.OutterMessage) {
		log.Printf("Echoed [%s]", string(message.Data))
	})

	address := "127.0.0.1:8888"
	host := CreateHost(address, mux)

	conn, err := host.ConnectToPeer(bifrost.NewAddress("127.0.0.1:7777"))

	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		input := ReadFromConsole()
		conn.Send(input)
	}
}
