/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conn

import (
	"errors"
	"regexp"
	"strings"

	"crypto/ecdsa"
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

	// just ip address xxx.xxx.xxx.xxx
	//re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	// ip with port number xxx.xxx.xxx.xxx:xxxx(x)  -> port number can be 0~99999 (real port numbers are in 0~65535 -> unsigned short 2bytes)
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]){1}([:][0-9][0-9][0-9][0-9][0-9]?)$`)
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
	CurveOpt int
}
