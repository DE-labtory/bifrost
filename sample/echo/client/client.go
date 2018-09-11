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
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/conn"

	"crypto/ecdsa"

	"crypto/elliptic"

	"crypto/rand"

	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/bifrost/pb"
	"github.com/it-chain/bifrost/sample/echo"
)

func CreateHost(ip string, mux *mux.Mux, pri *ecdsa.PrivateKey) *bifrost.BifrostHost {
	formatter := echo.SimpleFormatter{}
	idGetter := echo.SimpleIdGetter{IDPrefix: "ITTEST", PubKeyByte: formatter.ToByte(&pri.PublicKey)}
	signer := echo.SimpleSigner{PriKey: pri, Message: nil}

	myconnectionInfo := bifrost.NewHostInfo(conn.Address{IP: ip}, pri, &idGetter)

	var OnConnectionHandler = func(connection conn.Connection) {
		log.Printf("New connections are connected [%s]", connection)
	}

	return bifrost.New(myconnectionInfo, mux, OnConnectionHandler, &signer, &formatter)
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

func main() {

	defer os.RemoveAll("~/key")

	generator := echo.SimpleGenerator{Curve: elliptic.P384(), Rand: rand.Reader}
	priv, err := generator.GenerateKey()

	var protocol mux.Protocol
	protocol = "/echo/1.0"
	mux := mux.NewMux()

	mux.Handle(protocol, func(message conn.OutterMessage) {
		//log.Printf("Echoed [%s]", string(message.Data))
		fmt.Println(fmt.Sprintf("%s", message.Data[:]))
	})

	address := "127.0.0.1:7777"
	host := CreateHost(address, mux, priv)

	conn, err := host.ConnectToPeer(bifrost.NewAddress("127.0.0.1:8888"))

	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		input := ReadFromConsole()

		envelope := BuildEnvelope(protocol, input)

		host.Signer.(*echo.SimpleSigner).Message = envelope.Payload
		envelope.Signature, err = host.Signer.Sign()
		if err != nil {
			log.Fatalln(err.Error())
		}

		conn.Send(envelope, nil, nil)
	}

	defer conn.Close()
}
