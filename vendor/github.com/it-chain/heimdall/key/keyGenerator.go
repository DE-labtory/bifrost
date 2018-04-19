// This file provides key generator interface.

package key

// keyGenerator represents a key generator.
type keyGenerator interface {

	// Generate generates public and private key that match the input key generation option.
	Generate(opts KeyGenOpts) (pri PriKey, pub PubKey, err error)
}
