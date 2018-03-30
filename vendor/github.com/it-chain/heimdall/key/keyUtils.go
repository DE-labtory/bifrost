package key

import (
	"encoding/pem"
	"crypto/x509"
	"errors"
)

func PEMToPublicKey(data []byte) (interface{}, error) {

	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Failed to decode data")
	}

	key, err := DERToPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to convert PEM data to public key")
	}

	return key, nil

}

func PEMToPrivateKey(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("Failed to decode data")
	}

	key, err := DERToPrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("Failed to convert PEM data to private key")
	}

	return key, nil

}

func DERToPublicKey(data []byte) (interface{}, error) {

	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, errors.New("Failed to Parse data")
	}

	return key, nil

}

func DERToPrivateKey(data []byte) (interface{}, error) {

	var key interface{}
	var err error

	if len(data) == 0 {
		return nil, errors.New("Input data should not be NIL")
	}

	if key, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return key, err
	}

	if key, err = x509.ParseECPrivateKey(data); err == nil {
		return key, err
	}

	return nil, errors.New("Unspported Private Key Type")

}