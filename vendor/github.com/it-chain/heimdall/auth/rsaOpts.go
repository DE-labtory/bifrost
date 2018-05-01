// This file provides RSA signing options.

package auth

import (
	"crypto"
	"crypto/rsa"
)

// RSASignerOpts represents RSA signing options as integer.
type RSASignerOpts int

const (
	AUTO_SHA224 RSASignerOpts = iota
	AUTO_SHA256
	AUTO_SHA384
	AUTO_SHA512

	EQUAL_SHA224
	EQUAL_SHA256
	EQUAL_SHA384
	EQUAL_SHA512

	UNKNOWN
)

// SignerOptsToPSSOptions parse the RSASignerOpts(RSA signer option) to PSS option.
func (opts RSASignerOpts) SignerOptsToPSSOptions() *rsa.PSSOptions {

	switch opts {
	case AUTO_SHA224:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512_224}
	case AUTO_SHA256:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512_256}
	case AUTO_SHA384:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA384}
	case AUTO_SHA512:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512}

	case EQUAL_SHA224:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA512_224}
	case EQUAL_SHA256:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA512_256}
	case EQUAL_SHA384:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA384}
	case EQUAL_SHA512:
		return &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA512}

	default:
		return nil

	}

}
