package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"

	"github.com/ethereum/go-ethereum/crypto"
)

type ecdsaPrivateKey struct {
	p *ecdsa.PrivateKey
}

func PrivateKeyFromBytes(privKey *ecdsa.PrivateKey) common.PrivateKey {
	wrappedKey := &ecdsaPrivateKey{p: privKey}

	return wrappedKey
}

func (p *ecdsaPrivateKey) Sign(msg []byte) common.Signature {
	sig, err := crypto.Sign(msg, p.p)
	if err != nil {
		panic(err)
	}

	return &Signature{sig: sig}
}

func (p *ecdsaPrivateKey) PublicKey() common.PublicKey {
	return &PublicKey{p: &p.p.PublicKey}
}
