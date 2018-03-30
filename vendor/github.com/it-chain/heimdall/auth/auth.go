package auth

import (
	"github.com/it-chain/heimdall/key"
)

type Auth interface {

	Sign(priKey key.Key, digest []byte, opts SignerOpts) ([]byte, error)

	Verify(pubKey key.Key, signature, digest []byte, opts SignerOpts) (bool, error)

}