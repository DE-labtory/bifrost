package main

import (
	"log"
	"os"
	"runtime/pprof"

	"bufio"
	"encoding/json"
	"fmt"

	"github.com/google/gops/agent"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall"
	"crypto/ecdsa"
)

func CreateHost(ip string, mux *mux.Mux, pub *ecdsa.PublicKey, pri *ecdsa.PrivateKey) *bifrost.BifrostHost {

	myconnectionInfo := bifrost.NewHostInfo(conn.Address{IP: ip}, pub, pri)

	var ErrorHandler = func(conn conn.Connection, err error) {
		log.Println(err.Error())
		log.Println(fmt.Sprintf("Connection is closing... [%s]", conn.GetConnInfo()))
		conn.Close()
	}

	mux.HandleError(ErrorHandler)

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
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)

	defer os.RemoveAll("~/key")

	priv, err := heimdall.GenerateKey(heimdall.SECP384R1)
	if err != nil {
		log.Fatal(err)
	}

	var protocol mux.Protocol
	protocol = "/echo/1.0"
	mux := mux.NewMux()

	mux.Handle(protocol, func(message conn.OutterMessage) {
		log.Printf("Echoed [%s]", string(message.Envelope.Payload))
		envelope := Sign(BuildEnvelope(protocol, string(message.Data)), priv)
		message.Respond(envelope, nil, nil)
	})

	address := "127.0.0.1:8888"
	host := CreateHost(address, mux, &priv.PublicKey, priv)

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	//time.Sleep(time.Hour)

	bifrost.Listen(host)
}
