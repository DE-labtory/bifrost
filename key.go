package bifrost

import (
	"github.com/it-chain/heimdall/key"
	b58 "github.com/jbenet/go-base58"
)

func FromPubKey(key key.PubKey) string {

	encoded := b58.Encode(key.SKI())
	return encoded
}

//Create ID from Pri Key
func FromPriKey(key key.PriKey) string {

	pub, _ := key.PublicKey()
	return FromPubKey(pub)
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
