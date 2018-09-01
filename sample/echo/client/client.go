package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"encoding/json"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall"
	"crypto/ecdsa"
)

func CreateHost(ip string, mux *mux.Mux, pub *ecdsa.PublicKey, pri *ecdsa.PrivateKey) *bifrost.BifrostHost {

	myconnectionInfo := bifrost.NewHostInfo(conn.Address{IP: ip}, pub, pri)

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
	}

	return bifrost.New(myconnectionInfo, mux, OnConnectionHandler)
}

func ReadFromConsole() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter text: ")
	text, _ := reader.ReadString('\n')
	return text
}

func BuildEnvelope(protocol mux.Protocol, data interface{}) *pb.Envelope {

	payload, _ := json.Marshal(data)
	envelope := &pb.Envelope{}
	envelope.Protocol = string(protocol)
	envelope.Payload = payload

	return envelope
}

func Sign(envelope *pb.Envelope, priKey *ecdsa.PrivateKey) *pb.Envelope {
	envelope.Signature, _ = heimdall.Sign(priKey, envelope.Payload, nil, heimdall.SHA384)

	return envelope
}

func main() {


	defer os.RemoveAll("~/key")

	priv, err := heimdall.GenerateKey(heimdall.SECP384R1)

	var protocol mux.Protocol
	protocol = "/echo/1.0"
	mux := mux.NewMux()

	mux.Handle(protocol, func(message conn.OutterMessage) {
		//log.Printf("Echoed [%s]", string(message.Data))
		fmt.Println(fmt.Sprintf("%s", message.Data[:]))
	})

	address := "127.0.0.1:7777"
	host := CreateHost(address, mux, &priv.PublicKey, priv)

	conn, err := host.ConnectToPeer(bifrost.NewAddress("127.0.0.1:8888"))

	defer conn.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		input := ReadFromConsole()
		envelope := Sign(BuildEnvelope(protocol, input), priv)
		conn.Send(envelope, nil, nil)
	}
}
