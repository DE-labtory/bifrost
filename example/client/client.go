package main

import (
	"bufio"
	"fmt"
	"os"

	"os/signal"
	"syscall"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/bifrost/mocks"
	"github.com/it-chain/bifrost/mux"
	"github.com/it-chain/iLogger"
)

var clientIp = "127.0.0.1:7778"
var serverIp = "127.0.0.1:7777"
var DefaultMux *mux.DefaultMux
var testClientKeyDirPath = "./.test_client_key"

//todo 아직깔끔하지않음 여러 수정필요
func main() {

	keyPair := mocks.NewMockKeyOpts()

	DefaultMux := mux.New()

	DefaultMux.Handle("chat", func(message bifrost.Message) {
		iLogger.Infof(nil, "[Bifrost] %s", message.Data)
	})

	DefaultMux.Handle("join", func(message bifrost.Message) {
		iLogger.Infof(nil, "[Bifrost] %s", message.Data)
	})

	clientOpt := client.ClientOpts{
		Ip:     clientIp,
		PubKey: keyPair.PubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	err := mocks.MockStoreKey(keyPair.PriKey, testClientKeyDirPath)
	if err != nil {
		iLogger.Fatalf(nil, err.Error())
	}

	signer := mocks.MockECDSASigner{KeyID: keyPair.PubKey.ID(), KeyDirPath: testClientKeyDirPath}
	verifier := mocks.MockECDSAVerifier{}
	recoverer := mocks.MockECDSAKeyRecoverer{}
	crypto := bifrost.Crypto{Signer: &signer, Verifier: &verifier, KeyRecoverer: &recoverer}

	conn, err := client.Dial(serverIp, nil, clientOpt, grpcOpt, crypto)
	if err != nil {
		iLogger.Fatalf(nil, err.Error())
	}

	conn.Handle(DefaultMux)

	go func() {
		if err := conn.Start(); err != nil {
			iLogger.Info(nil, "[Bifrost] Conn close")
			conn.Close()
		}
	}()

	conn.Send([]byte("client join!!"), "join", nil, nil)

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			os.RemoveAll(testClientKeyDirPath)
			os.Exit(0)
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		conn.Send([]byte(text), "chat", nil, nil)
	}

}
