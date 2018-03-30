package key

type keyGenerator interface {

	Generate(opts KeyGenOpts) (pri PriKey, pub PubKey, err error)

}
