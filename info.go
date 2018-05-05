package bifrost

import "github.com/it-chain/heimdall/key"

//type ID string
//
//func (id ID) ToString() string {
//	return string(id)
//}
//
//
type PeerInfo struct {
	Ip        string
	Pubkey    []byte
	KeyGenOpt key.KeyGenOpts
}

//
//func NewConnInfo(id string, address Address, pubKey key.PubKey) ConnInfo {
//	return ConnInfo{
//		Id:      ID(id),
//		Address: address,
//		PeerKey: pubKey,
//	}
//}
//
//type PublicConnInfo struct {
//	Id        string
//	Address   Address
//	Pubkey    []byte
//	KeyType   key.KeyType
//	KeyGenOpt key.KeyGenOpts
//}
//
//func FromPublicConnInfo(publicConnInfo PublicConnInfo) (*ConnInfo, error) {
//
//	pubKey, err := ByteToPubKey(publicConnInfo.Pubkey, publicConnInfo.KeyGenOpt, publicConnInfo.KeyType)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &ConnInfo{
//		Id:      ID(publicConnInfo.Id),
//		Address: publicConnInfo.Address,
//		PeerKey: pubKey,
//	}, nil
//}
//
//func ByteToPubKey(byteKey []byte, keyGenOpt key.KeyGenOpts, keyType key.KeyType) (key.PubKey, error) {
//
//	k, err := key.PEMToPublicKey(byteKey)
//
//	if err != nil {
//		return nil, err
//	}
//
//	pub, err := key.MatchPublicKeyOpt(k, keyGenOpt)
//	if err != nil {
//		return nil, err
//	}
//
//	return pub, nil
//}
