package main

import (
	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/key"
	"github.com/it-chain/heimdall/hashing"
	"log"
	"fmt"
	"os"
)

/*
This sample shows data to be transmitted
is signed and verified by ECDSA Key.
*/

func main() {

	keyManager, err := key.NewKeyManager("")
	errorCheck(err)

	pri, pub, err := keyManager.GenerateKey(key.ECDSA521)
	errorCheck(err)

	defer os.RemoveAll("./.keyRepository")

	sampleData := []byte("This is sample data from heimdall.")

	hashManager, err := hashing.NewHashManager()
	errorCheck(err)

	digest, err := hashManager.Hash(sampleData, nil, hashing.SHA512)
	errorCheck(err)

	authManager, err := auth.NewAuth()
	errorCheck(err)

	signature, err := authManager.Sign(pri, digest, nil)
	errorCheck(err)

	/* --------- After data transmitted --------- */
	ok, err := authManager.Verify(pub, signature, digest, nil)
	errorCheck(err)

	fmt.Println(ok)

}

func errorCheck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
