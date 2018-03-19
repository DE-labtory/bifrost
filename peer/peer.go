package peer

import (
	"github.com/it-chain/heimdall"
	b58 "github.com/jbenet/go-base58"
)

type ID string

func FromPubkey(key heimdall.Key) ID{
	encoded := b58.Encode(key.SKI())
	return ID(encoded)
}

func (id ID) String() string{
	return string(id)
}

type Address struct{
	Ip string
	Port string
}

func ToAddress(ipv4 string) Address{
	return Address{

	}
}

type Peer struct{
	Id ID
	Address Address
}