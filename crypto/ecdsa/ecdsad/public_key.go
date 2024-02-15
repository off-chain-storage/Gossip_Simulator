package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"

	"github.com/sirupsen/logrus"
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

	logrus.Info(p.p)

	np := *p.p
	return &PublicKey{p: &np}
}
