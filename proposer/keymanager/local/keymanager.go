package local

import (
	"context"
	"crypto/ecdsa"

	curieecdsa "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"
	curiepb "flag-example/proto"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
)

var (
	lock            sync.RWMutex
	privateKeyCache                  = make(map[string]curieecdsa.PrivateKey)
	publicKey       *ecdsa.PublicKey = nil
)

type KeyManager struct{}

// KeyManager instance is Nil <- problem
func NewKeyManager(_ context.Context) (*KeyManager, error) {
	k := &KeyManager{}

	if err := k.GenerateKey(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize account store")
	}

	return k, nil
}

func (*KeyManager) GenerateKey() error {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "failed to generate key")
	}

	privKey.PublicKey.Curve = secp256k1.S256()

	pubKey := &privKey.PublicKey

	lock.Lock()
	publicKey = pubKey
	privateKeyCache[curieecdsa.ConvertToStringEcdsaPubKey(pubKey)] =
		ecdsad.PrivateKeyFromBytes(privKey)
	lock.Unlock()

	return nil
}

func (*KeyManager) Sign(ctx context.Context, req *curiepb.SignRequest) (curieecdsa.Signature, error) {
	lock.RLock()
	privateKey, ok := privateKeyCache[req.PublicKey]
	lock.RUnlock()
	if !ok {
		return nil, errors.New("secret key not found for public key")
	}

	return privateKey.Sign(req.SigningMsg), nil
}

func (*KeyManager) FetchValidatingPublicKeys() (*ecdsa.PublicKey, error) {
	lock.RLock()
	result := publicKey
	lock.RUnlock()

	return result, nil
}
