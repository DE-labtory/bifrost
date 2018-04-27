package conn

import (
	"errors"
	"regexp"
	"strings"

	"github.com/it-chain/heimdall/key"
)

type ID string

func (id ID) ToString() string {
	return string(id)
}

//Address to connect other peer
type Address struct {
	IP string
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}

	return false
}

//format should be xxx.xxx.xxx.xxx:xxxx
func ToAddress(ipv4 string) (Address, error) {

	valid := validIP4(ipv4)

	if !valid {
		return Address{}, errors.New("invalid IP4 format")
	}

	return Address{
		IP: ipv4,
	}, nil
}

type ConnInfo struct {
	Id      ID
	Address Address
	PeerKey key.PubKey
}

func NewConnInfo(id string, address Address, pubKey key.PubKey) ConnInfo {
	return ConnInfo{
		Id:      ID(id),
		Address: address,
		PeerKey: pubKey,
	}
}

type PublicConnInfo struct {
	Id        string
	Address   Address
	Pubkey    []byte
	KeyType   key.KeyType
	KeyGenOpt key.KeyGenOpts
}

func FromPublicConnInfo(publicConnInfo PublicConnInfo) (*ConnInfo, error) {

	pubKey, err := ByteToPubKey(publicConnInfo.Pubkey, publicConnInfo.KeyGenOpt, publicConnInfo.KeyType)

	if err != nil {
		return nil, err
	}

	return &ConnInfo{
		Id:      ID(publicConnInfo.Id),
		Address: publicConnInfo.Address,
		PeerKey: pubKey,
	}, nil
}

func ByteToPubKey(byteKey []byte, keyGenOpt key.KeyGenOpts, keyType key.KeyType) (key.PubKey, error) {

	k, err := key.PEMToPublicKey(byteKey)

	if err != nil {
		return nil, err
	}

	pub, err := key.MatchPublicKeyOpt(k, keyGenOpt)
	if err != nil {
		return nil, err
	}

	return pub, nil
}
