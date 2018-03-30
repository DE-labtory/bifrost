package key

import (
	"encoding/hex"
	"os"
	"io/ioutil"
	"errors"
	"path/filepath"
)

type keyStorer struct {
	path string
}

func (storer *keyStorer) Store(keys... Key) (err error) {

	if len(keys) == 0 {
		return errors.New("Input values should not be NIL")
	}

	for _, key := range keys {
		err := storer.storeKey(key)

		if err != nil{
			return err
		}
	}

	return nil

}

func (storer *keyStorer) storeKey(key Key) (error) {

	var data []byte
	var err error

	if key == nil{
		return errors.New("No Key Errors")
	}

	data, err = key.ToPEM()

	if err != nil {
		return err
	}

	path, err := storer.getFullPath(hex.EncodeToString(key.SKI()), key.Algorithm().String() + "_" + string(key.Type()))

	if err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = ioutil.WriteFile(path, data, 0700)
		if err != nil {
			return err
		}
	}

	return nil

}

func (storer *keyStorer) getFullPath(alias, suffix string) (string, error) {

	if _, err := os.Stat(storer.path); os.IsNotExist(err) {
		err = os.MkdirAll(storer.path, 0755)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(storer.path, alias + "_" + suffix), nil

}