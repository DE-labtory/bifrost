package heimdall

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"encoding/hex"
	"crypto/rsa"
)

func TestKeyManager_StoreKey(t *testing.T) {
	km := &keyManager{}
	km.Init("")

	defer os.RemoveAll(km.path)

	rsaRawKey, err := rsa.GenerateKey(rand.Reader, 1024)

	rsaPriKey := &RsaPrivateKey{rsaRawKey}
	err = km.Store(rsaPriKey)
	assert.NoError(t, err)

	rsaPubKey := &RsaPublicKey{&rsaPriKey.priv.PublicKey}
	err = km.Store(rsaPubKey)
	assert.NoError(t, err)

	ecdsaRawKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	ecdsaPriKey := &EcdsaPrivateKey{ecdsaRawKey}
	err = km.Store(ecdsaPriKey)
	assert.NoError(t, err)

	ecdsaPubKey := &EcdsaPublicKey{&ecdsaPriKey.priv.PublicKey}
	err = km.Store(ecdsaPubKey)
	assert.NoError(t, err)

	// check whether file is exist
	path := filepath.Join(km.path, hex.EncodeToString(rsaPriKey.SKI()) + "_pri")
	assert.FileExists(t, path)

	path = filepath.Join(km.path, hex.EncodeToString(rsaPubKey.SKI()) + "_pub")
	assert.FileExists(t, path)

	path = filepath.Join(km.path, hex.EncodeToString(ecdsaPriKey.SKI()) + "_pri")
	assert.FileExists(t, path)

	path = filepath.Join(km.path, hex.EncodeToString(ecdsaPubKey.SKI()) + "_pub")
	assert.FileExists(t, path)

	// Invalid input for Store
	err = km.Store(nil)
	assert.Error(t, err)
}

//test fail 이유?!!
func TestKeyManager_StoreInvalidInput(t *testing.T) {

	km := &keyManager{}
	km.Init("")

	defer os.RemoveAll(km.path)

	rsaRawKey, err := rsa.GenerateKey(rand.Reader, 1024)

	rsaPriKey := &RsaPrivateKey{rsaRawKey}
	err = km.storeKey(rsaPriKey)
	assert.Error(t, err)

	rsaPubKey := &RsaPublicKey{&rsaPriKey.priv.PublicKey}
	err = km.storeKey(rsaPubKey)
	assert.Error(t, err)

	ecdsaRawKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	ecdsaPriKey := &EcdsaPrivateKey{ecdsaRawKey}
	err = km.storeKey(ecdsaPriKey)
	assert.Error(t, err)

	ecdsaPubKey := &EcdsaPublicKey{&ecdsaPriKey.priv.PublicKey}
	err = km.storeKey(ecdsaPubKey)
	assert.Error(t, err)
}

func TestKeyManager_LoadKey(t *testing.T) {

	km := &keyManager{}
	km.Init("")

	defer os.RemoveAll(km.path)

	rsaRawKey, err := rsa.GenerateKey(rand.Reader, 1024)

	rsaPriKey := &RsaPrivateKey{rsaRawKey}
	err = km.Store(rsaPriKey)
	assert.NoError(t, err)

	rsaPubKey := &RsaPublicKey{&rsaPriKey.priv.PublicKey}
	err = km.Store(rsaPubKey)
	assert.NoError(t, err)

	pri, pub, err := km.Load()
	assert.NoError(t, err)
	assert.NotNil(t, pri)
	assert.NotNil(t, pub)

	assert.Equal(t, rsaPriKey, pri.(*RsaPrivateKey))
	assert.Equal(t, rsaPubKey, pub.(*RsaPublicKey))

}