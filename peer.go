package bifrost

import "github.com/it-chain/heimdall/key"

type PeerInfo struct {
	Ip        string
	Pubkey    []byte
	KeyGenOpt key.KeyGenOpts
	MetaData  map[string]string
}
