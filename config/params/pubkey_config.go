package params

import (
	"flag-example/crypto/ecdsa/common"

	"github.com/mohae/deepcopy"
)

type ProposerConfig struct {
	ProposerPubKey *common.PublicKey
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
