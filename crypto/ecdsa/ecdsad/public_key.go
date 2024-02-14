package ecdsad

import "crypto/ecdsa"

type PublicKey struct {
	p *ecdsa.PublicKey
}
