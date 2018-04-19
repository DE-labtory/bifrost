package bifrost

import (
	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/heimdall/key"
	b58 "github.com/jbenet/go-base58"
)

//Identitiy of Connection
type ID string

//Create ID from Public Key
func FromPubKey(key key.PubKey) ID {

	encoded := b58.Encode(key.SKI())
	return ID(encoded)
}

//Create ID from Pri Key
func FromPriKey(key key.PriKey) ID {

	pub, _ := key.PublicKey()
	return FromPubKey(pub)
}

func (id ID) String() string {
	return string(id)
}

type HostInfo struct {
	conn.ConnInfo
	PriKey key.PriKey
}

func NewHostInfo(id ID, address conn.Address, pubKey key.PubKey, priKey key.PriKey) HostInfo {

	return HostInfo{
		ConnInfo: conn.NewConnInfo(id, address, pubKey),
		PriKey:   priKey,
	}
}

func (hostInfo HostInfo) GetPublicInfo() *conn.PublicConnInfo {

	publicConnInfo := &conn.PublicConnInfo{}
	publicConnInfo.Id = hostInfo.Id
	publicConnInfo.Address = hostInfo.Address

	b, err := hostInfo.PubKey.ToPEM()

	if err != nil {
		return nil
	}

	publicConnInfo.Pubkey = b
	publicConnInfo.KeyType = hostInfo.PubKey.Type()
	publicConnInfo.KeyGenOpt = hostInfo.PubKey.Algorithm()

	return publicConnInfo
}
