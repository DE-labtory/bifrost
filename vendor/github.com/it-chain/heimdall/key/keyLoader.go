package key

import (
	"os"
	"io/ioutil"
	"strings"
	"crypto/rsa"
	"crypto/ecdsa"
	"errors"
	"path/filepath"
)

type keyLoader struct {
	path string
}

type keyInfos struct {
	id			string
	keyGenOpts	KeyGenOpts
	keyType		KeyType
}

func (loader *keyLoader) Load() (pri PriKey, pub PubKey, err error) {

	if _, err := os.Stat(loader.path); os.IsNotExist(err) {

		return nil, nil, errors.New("Keys are not exist")
	}

	files, err := ioutil.ReadDir(loader.path)
	if err != nil {
		return nil, nil, errors.New("Failed to read key repository directory")
	}

	for _, file := range files {

		keyInfos, ok := loader.getKeyInfos(file.Name())

		if ok {
			switch keyInfos.keyType {
			case PRIVATE_KEY:
				key, err := loader.loadKey(keyInfos.id, keyInfos.keyType)
				if err != nil {
					return nil, nil, err
				}

				switch key.(type) {
				case *rsa.PrivateKey:
					pri = &RSAPrivateKey{PrivKey: key.(*rsa.PrivateKey), bits: KeyGenOptsToRSABits(keyInfos.keyGenOpts)}
				case *ecdsa.PrivateKey:
					pri = &ECDSAPrivateKey{PrivKey: key.(*ecdsa.PrivateKey)}
				default:
					return nil, nil, errors.New("Failed to load Key")
				}

			case PUBLIC_KEY:
				key, err := loader.loadKey(keyInfos.id, keyInfos.keyType)
				if err != nil {
					return nil, nil, err
				}

				switch key.(type) {
				case *rsa.PublicKey:
					pub = &RSAPublicKey{PubKey: key.(*rsa.PublicKey), bits: KeyGenOptsToRSABits(keyInfos.keyGenOpts)}
				case *ecdsa.PublicKey:
					pub = &ECDSAPublicKey{key.(*ecdsa.PublicKey)}
				default:
					return nil, nil, errors.New("Failed to load Key")
				}
			}
		}
	}

	if pri == nil || pub == nil {
		return nil, nil, errors.New("Failed to load Key")
	}

	return pri, pub, nil

}

func (loader *keyLoader) loadKey(alias string, keyType KeyType) (key interface{}, err error) {

	if len(alias) == 0 {
		return nil, errors.New("Input value should not be blank")
	}

	path, err := loader.getFullPath(alias, string(keyType))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch keyType {
	case PRIVATE_KEY:
		key, err = PEMToPrivateKey(data)
	case PUBLIC_KEY:
		key, err = PEMToPublicKey(data)
	}

	if err != nil {
		return nil, err
	}

	return key, nil

}

func (loader *keyLoader) getKeyInfos(name string) (*keyInfos, bool) {

	datas := strings.Split(name, "_")

	if len(datas) != 3 {
		return nil, false
	}

	keyGenOpts, ok := StringToKeyGenOpts(datas[1])
	if !ok {
		return nil, false
	}

	keyType := KeyType(datas[2])
	if !(keyType == PRIVATE_KEY || keyType == PUBLIC_KEY) {
		return nil, false
	}

	infos := &keyInfos{
		id:datas[0],
		keyGenOpts:keyGenOpts,
		keyType:keyType,
	}

	return infos, true

}

func (loader *keyLoader) getFullPath(alias, suffix string) (string, error) {
	if _, err := os.Stat(loader.path); os.IsNotExist(err) {
		err = os.MkdirAll(loader.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(loader.path, alias + "_" + suffix), nil
}
