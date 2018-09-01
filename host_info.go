package bifrost

import (
	"crypto/ecdsa"

	"github.com/it-chain/bifrost/conn"
	"github.com/it-chain/heimdall"
)

//Identitiy of Connection
type ID string

//Create ID from Public Key
func FromPubKey(key *ecdsa.PublicKey) ID {
	return ID(heimdall.PubKeyToKeyID(key))
}

func (id ID) String() string {
	return string(id)
}

type HostInfo struct {
	conn.ConnInfo
	PriKey *ecdsa.PrivateKey
}

func NewHostInfo(address conn.Address, pubKey *ecdsa.PublicKey, priKey *ecdsa.PrivateKey) HostInfo {

	id := FromPubKey(pubKey)

	return HostInfo{
		ConnInfo: conn.NewConnInfo(id.String(), address, pubKey),
		PriKey:   priKey,
	}
}

func (hostInfo HostInfo) GetPublicInfo() *conn.PublicConnInfo {

	publicConnInfo := &conn.PublicConnInfo{}
	publicConnInfo.Id = hostInfo.Id.ToString()
	publicConnInfo.Address = hostInfo.Address

	b := heimdall.PubKeyToBytes(hostInfo.PubKey)

	publicConnInfo.Pubkey = b
	publicConnInfo.CurveOpt = heimdall.CurveToCurveOpt(hostInfo.PubKey.Curve)

	return publicConnInfo
}
