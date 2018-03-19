package peer

import (
	"testing"
	"github.com/it-chain/heimdall"
	"fmt"
	"os"
)

func TestFromPubkey(t *testing.T) {
	cryp, _ := heimdall.NewCryptoImpl(".myKeys", &heimdall.ECDSAKeyGenOpts{})
	defer os.RemoveAll(".myKeys")

	pri, _, _ := cryp.GetKey()
	fmt.Println(FromPubkey(pri))
}