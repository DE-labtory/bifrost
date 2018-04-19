package conn

import (
	"errors"
	"regexp"
	"strings"

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
	PubKey  key.PubKey
}

func NewConnInfo(id ID, address Address, pubKey key.PubKey) ConnInfo {
	return ConnInfo{
		Id:      id,
		Address: address,
		PubKey:  pubKey,
	}
}

type HostInfo struct {
	ConnInfo
	PriKey key.PriKey
}

func NewHostInfo(id ID, address Address, pubKey key.PubKey, priKey key.PriKey) HostInfo {

	return HostInfo{
		ConnInfo: NewConnInfo(id, address, pubKey),
		PriKey:   priKey,
	}
}

type PublicConnInfo struct {
	Id        ID
	Address   Address
	Pubkey    []byte
	KeyType   key.KeyType
	KeyGenOpt key.KeyGenOpts
}

func (hostInfo HostInfo) GetPublicInfo() *PublicConnInfo {

	publicConnInfo := &PublicConnInfo{}
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
