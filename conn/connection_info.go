package conn

import (
	"errors"
	"regexp"
	"strings"

	"crypto/ecdsa"

	"github.com/it-chain/heimdall"
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
	PubKey  *ecdsa.PublicKey
}

func NewConnInfo(id string, address Address, pubKey *ecdsa.PublicKey) ConnInfo {
	return ConnInfo{
		Id:      ID(id),
		Address: address,
		PubKey:  pubKey,
	}
}

type PublicConnInfo struct {
	Id       string
	Address  Address
	Pubkey   []byte
	CurveOpt heimdall.CurveOpts
}

func FromPublicConnInfo(publicConnInfo PublicConnInfo) (*ConnInfo, error) {

	pubKey, err := heimdall.BytesToPubKey(publicConnInfo.Pubkey, publicConnInfo.CurveOpt)

	if err != nil {
		return nil, err
	}

	return &ConnInfo{
		Id:      ID(publicConnInfo.Id),
		Address: publicConnInfo.Address,
		PubKey:  pubKey,
	}, nil
}
