package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"

	"github.com/mohae/deepcopy"
	"github.com/sirupsen/logrus"
)

var publicKey *PublicKey

type PublicKey struct {
	p *ecdsa.PublicKey
}

func PublicKeyFromProposer(pubKey *ecdsa.PublicKey) {
	publicKey = &PublicKey{p: pubKey}

	logrus.Info("Public key set", publicKey.p)
}

func GetPublicKey() *ecdsa.PublicKey {
	return publicKey.p
}

func (p *PublicKey) Copy() common.PublicKey {
	config, ok := deepcopy.Copy(*p).(common.PublicKey)
	if !ok {
		config = publicKey
	}
	return config
}
