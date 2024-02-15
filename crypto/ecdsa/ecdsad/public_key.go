package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
)

type PublicKey struct {
	p *ecdsa.PublicKey
}

func GetPublicKey() common.PublicKey {
	return &PublicKey{}
}

func PublicKeyFromProposer(pubKey *ecdsa.PublicKey) *PublicKey {
	return &PublicKey{p: pubKey}
}

// Copy the public key to a new pointer reference.
func (p *PublicKey) Copy() common.PublicKey {
	np := *p.p
	return &PublicKey{p: &np}
}
