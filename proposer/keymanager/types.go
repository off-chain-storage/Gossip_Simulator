package keymanager

import (
	"context"
	"crypto/ecdsa"

	curieecdsa "flag-example/crypto/ecdsa"
	curiepb "flag-example/proto"
)

type IKeymanager interface {
	PublicKeysFetcher
	Signer
	Generator
}

type PublicKeysFetcher interface {
	FetchValidatingPublicKeys() (*ecdsa.PublicKey, error)
}

type Generator interface {
	GenerateKey() error
}

type Signer interface {
	Sign(context.Context, *curiepb.SignRequest) (curieecdsa.Signature, error)
}
