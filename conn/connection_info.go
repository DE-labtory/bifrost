package conn

import (
	"errors"
	"regexp"
	"strings"

	"github.com/it-chain/heimdall/key"
	b58 "github.com/jbenet/go-base58"
)

type ID string

func FromRsaPubKey(key key.PubKey) ID {
	encoded := b58.Encode(key.SKI())
	return ID(encoded)
}

func FromRsaPriKey(key key.PriKey) ID {
	pub, _ := key.PublicKey()

	return FromRsaPubKey(pub)
}

func (id ID) String() string {
	return string(id)
}

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

type ConnenctionInfo struct {
	Id      ID
	Address Address
	PubKey  key.PubKey
}

func NewConnenctionInfo(id ID, address Address, pubKey key.PubKey) ConnenctionInfo {
	return ConnenctionInfo{
		Id:      id,
		Address: address,
		PubKey:  pubKey,
	}
}

type MyConnectionInfo struct {
	ConnenctionInfo
	PriKey key.PriKey
}

func NewMyConnectionInfo(id ID, address Address, pubKey key.PubKey, priKey key.PriKey) MyConnectionInfo {

	return MyConnectionInfo{
		ConnenctionInfo: NewConnenctionInfo(id, address, pubKey),
		PriKey:          priKey,
	}
}

func (myConnectionInfo MyConnectionInfo) GetPublicInfo() ConnenctionInfo {
	return myConnectionInfo.ConnenctionInfo
}
