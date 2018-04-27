package main

import (
	"log"
	"os"
	"runtime/pprof"

	"bufio"
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"github.com/google/gops/agent"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/key"
)

func CreateHost(ip string, mux *mux.Mux, pub key.PubKey, pri key.PriKey) *bifrost.BifrostHost {

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

func Sign(envelope *pb.Envelope, priKey key.PriKey) *pb.Envelope {

	au, _ := auth.NewAuth()

	hash := sha512.New()
	hash.Write(envelope.Payload)
	digest := hash.Sum(nil)

	sig, _ := au.Sign(priKey, digest, auth.EQUAL_SHA512.SignerOptsToPSSOptions())
	envelope.Signature = sig

	return envelope
}

func main() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)

	km, err := key.NewKeyManager("~/key")

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer os.RemoveAll("~/key")

	priv, pub, err := km.GenerateKey(key.RSA4096)

	var protocol mux.Protocol
	protocol = "/echo/1.0"
	mux := mux.NewMux()

	mux.Handle(protocol, func(message conn.OutterMessage) {
		log.Printf("Echoed [%s]", string(message.Envelope.Payload))
		envelope := Sign(BuildEnvelope(protocol, string(message.Data)), priv)
		message.Respond(envelope, nil, nil)
	})

	address := "127.0.0.1:8888"
	host := CreateHost(address, mux, pub, priv)

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	//time.Sleep(time.Hour)
	bifrost.Listen(host)
}
