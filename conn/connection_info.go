package conn

import (
	"errors"
	"regexp"
	"strings"

	"github.com/it-chain/bifrost"
	"github.com/it-chain/heimdall/key"
)

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
	Id      bifrost.ID
	Address Address
	PubKey  key.PubKey
}

func NewConnInfo(id bifrost.ID, address Address, pubKey key.PubKey) ConnInfo {
	return ConnInfo{
		Id:      id,
		Address: address,
		PubKey:  pubKey,
	}
}

type PublicConnInfo struct {
	Id        bifrost.ID
	Address   Address
	Pubkey    []byte
	KeyType   key.KeyType
	KeyGenOpt key.KeyGenOpts
}

func FromPublicConnInfo(publicConnInfo PublicConnInfo) ConnInfo {

}
