// This file provides key generation options.

package key

import (
	"crypto/elliptic"
)

// KeyGenOpts represents key generation options by integer number.
type KeyGenOpts int

const (
	RSA1024 KeyGenOpts = iota
	RSA2048
	RSA4096

	ECDSA224
	ECDSA256
	ECDSA384
	ECDSA521

	UNKNOWN
)

var optsArr = [...]string{
	"rsa1024",
	"rsa2048",
	"rsa4096",

	"ecdsa224",
	"ecdsa256",
	"ecdsa384",
	"ecdsa521",

	"unknown",
}

// ValidCheck checks the input key generation option is valid or not.
func (opts KeyGenOpts) ValidCheck() bool {

	if opts < 0 || opts >= KeyGenOpts(len(optsArr)) {
		return false
	}

	return true

}

// String coverts format of key generation option from KeyGenOpts to string.
func (opts KeyGenOpts) String() string {

	if !opts.ValidCheck() {
		return "unknown"
	}

	return optsArr[opts]

}

// StringToKeyGenOpts converts format of key generation option from string to KeyGenOpts
func StringToKeyGenOpts(rawOpts string) (KeyGenOpts, bool) {

	for idx, opts := range optsArr {
		if rawOpts == opts {
			return KeyGenOpts(idx), true
		}
	}

	return -1, false

}

// ECDSACurveToKeyGenOpts converts format of ECDSA elliptic curve from elliptic.Curve to KeyGenOpts.
func ECDSACurveToKeyGenOpts(curve elliptic.Curve) KeyGenOpts {

	switch curve {
	case elliptic.P224():
		return ECDSA224
	case elliptic.P256():
		return ECDSA256
	case elliptic.P384():
		return ECDSA384
	case elliptic.P521():
		return ECDSA521
	default:
		return UNKNOWN
	}

}

// KetGenOptsToECDSACurve converts format of ECDSA elliptic curve from KeyGenOpts to elliptic.Curve.
func KeyGenOptsToECDSACurve(opts KeyGenOpts) elliptic.Curve {

	switch opts {
	case ECDSA224:
		return elliptic.P224()
	case ECDSA256:
		return elliptic.P256()
	case ECDSA384:
		return elliptic.P384()
	case ECDSA521:
		return elliptic.P521()
	default:
		return nil
	}

}

// RSABitsToKeyGenOpts converts format of RSA bits from bit length to KeyGenOpts.
func RSABitsToKeyGenOpts(bits int) KeyGenOpts {

	switch bits {
	case 1024:
		return RSA1024
	case 2048:
		return RSA2048
	case 4096:
		return RSA4096
	default:
		return UNKNOWN
	}

}

// KeyGenOptsToRSABits converts format of RSA bits from KeyGenOpts to bit length.
func KeyGenOptsToRSABits(opts KeyGenOpts) int {

	switch opts {
	case RSA1024:
		return 1024
	case RSA2048:
		return 2048
	case RSA4096:
		return 4096
	default:
		return -1
	}

}
