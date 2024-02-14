package params

import (
	"crypto/ecdsa"

	"github.com/mohae/deepcopy"
)

type ProposerConfig struct {
	ProposerPubKey *ecdsa.PublicKey
}

var proposerConfig ProposerConfig

func CurieProposerConfig() *ProposerConfig {
	return &proposerConfig
}

func OverrideCuriePublisherConfig(cfg *ProposerConfig) {
	proposerConfig = *cfg.Copy()
}

func (p *ProposerConfig) Copy() *ProposerConfig {
	config, ok := deepcopy.Copy(*p).(ProposerConfig)
	if !ok {
		config = proposerConfig
	}
	return &config
}
