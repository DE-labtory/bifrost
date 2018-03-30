package key

type KeyManager interface {

	GenerateKey(opts KeyGenOpts) (pri PriKey, pub PubKey, err error)

	GetKey() (pri PriKey, pub PubKey, err error)

	RemoveKey() (error)

}

