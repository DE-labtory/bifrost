package key

type KeyType string

const (
	PRIVATE_KEY KeyType = "pri"
	PUBLIC_KEY	KeyType = "pub"
)

type Key interface {

	SKI() (ski []byte)

	Algorithm() KeyGenOpts

	ToPEM() ([]byte,error)

	Type() (KeyType)

}

type PriKey interface {

	Key

	PublicKey() (PubKey, error)

}

type PubKey interface { Key }
