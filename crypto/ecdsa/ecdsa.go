package ecdsa

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/mohae/deepcopy"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type ProposerConfig struct {
	privKey *ecdsa.PrivateKey
	pubKey  *ecdsa.PublicKey
}

var proposerConfig ProposerConfig

func OverrideCurieProposerConfig(cfg *ProposerConfig) {
	proposerConfig = *cfg.Copy()
}

func (p *ProposerConfig) Copy() *ProposerConfig {
	config, ok := deepcopy.Copy(*p).(ProposerConfig)
	if !ok {
		config = proposerConfig
	}
	return &config
}

func CurieProposerConfig() *ProposerConfig {
	return &proposerConfig
}

func (pc *ProposerConfig) GenerateKey() error {
	privKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err == nil {
		pc.privKey = privKey
		pc.pubKey = &privKey.PublicKey
	}

	return err
}

func (pc *ProposerConfig) PrivKey() *ecdsa.PrivateKey { return pc.privKey }

func (pc *ProposerConfig) PubKey() *ecdsa.PublicKey { return pc.pubKey }
