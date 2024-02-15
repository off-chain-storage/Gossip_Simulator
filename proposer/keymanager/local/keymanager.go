package local

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"

	curieecdsa "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"
	curiepb "flag-example/proto"
	"sync"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	privKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return errors.Wrap(err, "failed to generate key")
	}
	pubKey := &privKey.PublicKey

	logrus.Info("Key Generation, ECDSA PubKey is ", pubKey)
	logrus.Info("Key Generation, String PubKey is ", curieecdsa.ConvertToStringEcdsaPubKey(pubKey))

	logrus.Info("Key Generation, ECDSA PrivKey is ", privKey)

	lock.Lock()
	publicKey = pubKey
	privateKeyCache[curieecdsa.ConvertToStringEcdsaPubKey(pubKey)] =
		ecdsad.PrivateKeyFromBytes(curieecdsa.ConvertToByteEcdsaPrivKey(privKey))
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

	logrus.Info("Before Signing, ECDSA PrivKey is ", privateKey)

	return privateKey.Sign(req.SigningMsg), nil
}

func (*KeyManager) FetchValidatingPublicKeys() (*ecdsa.PublicKey, error) {
	lock.RLock()
	result := publicKey
	lock.RUnlock()

	return result, nil
}
