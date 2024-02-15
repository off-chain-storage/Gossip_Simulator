package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"

	"github.com/mohae/deepcopy"
)

var publicKey common.PublicKey

type PublicKey struct {
	p *ecdsa.PublicKey
}

func GetPublicKey() common.PublicKey {
	return publicKey
}

func PublicKeyFromProposer(pubKey *ecdsa.PublicKey) common.PublicKey {
	publicKey = &PublicKey{p: pubKey}

	return publicKey
}

func (p *PublicKey) Copy() common.PublicKey {
	config, ok := deepcopy.Copy(*p).(common.PublicKey)
	if !ok {
		config = publicKey
	}
	return config
}
