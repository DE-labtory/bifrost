/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"os"
	"runtime/pprof"

	"bufio"
	"encoding/json"
	"fmt"

	"crypto/ecdsa"

	"crypto/elliptic"
	"crypto/rand"

	"github.com/google/gops/agent"
	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/sample/echo"
)

func CreateHost(ip string, mux *mux.Mux, pri *ecdsa.PrivateKey) *bifrost.BifrostHost {
	formatter := echo.SimpleFormatter{}
	idGetter := echo.SimpleIdGetter{IDPrefix: "ITTEST", PubKeyByte: formatter.ToByte(&pri.PublicKey)}
	signer := echo.SimpleSigner{PriKey: pri, Message: nil}

	myconnectionInfo := bifrost.NewHostInfo(conn.Address{IP: ip}, pri, &idGetter)

	var ErrorHandler = func(conn conn.Connection, err error) {
		log.Println(err.Error())
		log.Println(fmt.Sprintf("Connection is closing... [%s]", conn.GetConnInfo()))
		conn.Close()
	}

	mux.HandleError(ErrorHandler)

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
	}

	return bifrost.New(myconnectionInfo, mux, OnConnectionHandler, &signer, &formatter)
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

func main() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)

	defer os.RemoveAll("~/key")

	generator := echo.SimpleGenerator{Curve: elliptic.P384(), Rand: rand.Reader}
	priv, err := generator.GenerateKey()

	if err != nil {
		log.Fatal(err)
	}

	var protocol mux.Protocol
	protocol = "/echo/1.0"
	serverMux := mux.NewMux()

	address := "127.0.0.1:8888"
	host := CreateHost(address, serverMux, priv)

	host.Mux.Handle(protocol, func(message conn.OutterMessage) {
		log.Printf("Echoed [%s]", string(message.Envelope.Payload))
		envelope := BuildEnvelope(protocol, string(message.Data))
		host.Signer.(*echo.SimpleSigner).Message = envelope.Payload
		envelope.Signature, err = host.Signer.Sign()
		if err != nil {
			log.Fatalln(err.Error())
		}

		message.Respond(envelope, nil, nil)
	})

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}
	//time.Sleep(time.Hour)

	bifrost.Listen(host)
}
