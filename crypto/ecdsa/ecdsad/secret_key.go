package ecdsad

import (
	"crypto/ecdsa"
	"crypto/rand"
	curieecdsa "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
)

type ecdsaPrivateKey struct {
	p *ecdsa.PrivateKey
}

func PrivateKeyFromBytes(privKey []byte) common.PrivateKey {
	pk := curieecdsa.ConvertToEcdsaPrivKeyByte(privKey)

	wrappedKey := &ecdsaPrivateKey{p: pk}

	return wrappedKey
}

func (p *ecdsaPrivateKey) Sign(msg []byte) common.Signature {
	r, s, err := ecdsa.Sign(rand.Reader, p.p, msg)
	if err != nil {
		panic(err)
	}
	return &Signature{sig: &signature{r: r, s: s}}
}

func (p *ecdsaPrivateKey) PublicKey() common.PublicKey {
	return &PublicKey{p: &p.p.PublicKey}
}
