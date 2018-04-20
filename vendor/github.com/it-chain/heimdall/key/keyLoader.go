// This file provides components that is needed for loading key process.

package key

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// keyLoader contains path as a string.
type keyLoader struct {
	path string
}

// keyInfos contains key related information such as id, key generation option and key type.
type keyInfos struct {
	id         string
	keyGenOpts KeyGenOpts
	keyType    KeyType
}

// Load loads private key and public key from stored file.
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
				key, err := loader.loadKey(keyInfos.id, keyInfos.keyGenOpts, keyInfos.keyType)
				if err != nil {
					return nil, nil, err
				}

				pri, err = MatchPrivateKeyOpt(key, keyInfos.keyGenOpts)
				if err != nil {
					return nil, nil, err
				}

			case PUBLIC_KEY:
				key, err := loader.loadKey(keyInfos.id, keyInfos.keyGenOpts, keyInfos.keyType)
				if err != nil {
					return nil, nil, err
				}

				pub, err = MatchPublicKeyOpt(key, keyInfos.keyGenOpts)
				if err != nil {
					return nil, nil, err
				}
			}
		}
	}

	if pri == nil || pub == nil {
		return nil, nil, errors.New("Failed to load Key")
	}

	return pri, pub, nil

}

// loadKey reads key from file and changes the format from PEM to key.
func (loader *keyLoader) loadKey(alias string, keyGenOpt KeyGenOpts, keyType KeyType) (key interface{}, err error) {

	if len(alias) == 0 {
		return nil, errors.New("Input value should not be blank")
	}

	path, err := loader.getFullPath(alias, keyGenOpt.String(), string(keyType))
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

// getKeyInfos gets key information that are id, key generation option and key type.
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
		id:         datas[0],
		keyGenOpts: keyGenOpts,
		keyType:    keyType,
	}

	return infos, true

}

// getFullPath gets full (absolute) path of a key file.
func (loader *keyLoader) getFullPath(alias, keyGenOpt string, suffix string) (string, error) {
	if _, err := os.Stat(loader.path); os.IsNotExist(err) {
		err = os.MkdirAll(loader.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(loader.path, alias+"_"+keyGenOpt+"_"+suffix), nil
}
