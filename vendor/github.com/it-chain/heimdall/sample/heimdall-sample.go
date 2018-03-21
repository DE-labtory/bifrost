package main

import (
	"log"
	"crypto/rsa"
	"crypto"
	"crypto/sha256"
	"fmt"
	"os"
	"github.com/it-chain/heimdall"
)

/*
This sample shows data to be transmitted
is signed and verified by RSA Key.
*/

func main() {

	cryp, err := heimdall.NewCryptoImpl(".myKeys", &heimdall.RSAKeyGenOpts{})
	if err != nil {
		log.Fatalln(err)
	}

	defer os.RemoveAll("./.myKeys")

	opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthEqualsHash, Hash:crypto.SHA256}
	rawData := []byte("This is raw Data with []byte")
	signature, err := cryp.Sign(rawData, opts)

	hash := sha256.New()
	hash.Write(rawData)
	digest := hash.Sum(nil)

	_, pub, err := cryp.GetKey()
	if err != nil {
		log.Fatalln(err)
	}

	ok, err := cryp.Verify(pub, signature, digest, opts)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(ok)

}
