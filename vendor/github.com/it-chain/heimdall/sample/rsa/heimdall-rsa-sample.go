package main

import (
	"fmt"
	"log"
	"os"

	"reflect"

	"github.com/it-chain/heimdall/auth"
	"github.com/it-chain/heimdall/hashing"
	"github.com/it-chain/heimdall/key"
)

/*
This sample shows data to be transmitted
is signed and verified by RSA Key.
*/

func main() {

	keyManager, err := key.NewKeyManager("")
	errorCheck(err)

	defer os.RemoveAll("./.heimdall")

	// Generate key pair with RSA algorithm.
	pri, pub, err := keyManager.GenerateKey(key.RSA4096)
	errorCheck(err)

	// Get key from memory of keyManager or from key files in key path of keyManager.
	pri, pub, err = keyManager.GetKey()
	errorCheck(err)

	// Convert key to PEM(byte) format.
	bytePriKey, err := pri.ToPEM()
	bytePubKey, err := pub.ToPEM()

	// Reconstruct key pair from bytes to key.
	recPri, err := key.PEMToPrivateKey(bytePriKey, key.RSA4096)
	recPub, err := key.PEMToPublicKey(bytePubKey, key.RSA4096)
	errorCheck(err)

	// Compare reconstructed key pair with original key pair.
	if reflect.DeepEqual(pri, recPri) && reflect.DeepEqual(pub, recPub) {
		print("reconstruct complete!\n")
	}

	sampleData := []byte("This is sample data from heimdall.")

	// Convert raw data to digest(hash value) by using SHA512 function.
	digest, err := hashing.Hash(sampleData, nil, hashing.SHA512)
	errorCheck(err)

	// The option will be used in signing process in case of using RSA key.
	signerOpts := auth.EQUAL_SHA512.SignerOptsToPSSOptions()

	// AuthManager makes digest(hash value) to signature with private key.
	signature, err := auth.Sign(pri, digest, signerOpts)
	errorCheck(err)

	/* --------- After data transmitted --------- */

	// AuthManager verify that received data has any forgery during transmitting process by digest.
	// and verify that the received data is surely from the expected sender by public key.
	ok, err := auth.Verify(pub, signature, digest, signerOpts)
	errorCheck(err)

	fmt.Println(ok)

}

func errorCheck(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
