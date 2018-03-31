package conn

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/it-chain/heimdall/key"
	"github.com/stretchr/testify/assert"
)

func TestMarshalConnectionInfo(t *testing.T) {

	km, err := key.NewKeyManager("~/key")

	defer os.RemoveAll("~/key")

	_, pub, err := km.GenerateKey(key.RSA4096)

	connInfo := NewConnenctionInfo(FromRsaPubKey(pub), Address{IP: "127.0.0.1:8888"}, pub)

	b, err := json.Marshal(connInfo)

	fmt.Printf("[%s]", b)

	if err != nil {

	}
	//
	connectedConnInfo := &ConnenctionInfo{}
	err = json.Unmarshal(b, connectedConnInfo)

	assert.NoError(t, err)
}
