package features

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "flags")

const enabledFeatureFlag = "Enabled feature flag"
const disabledFeatureFlag = "Disabled feature flag"

type Flags struct {
	DisableResourceManager bool // Disables running the node with libp2p's resource manager.
}

var featureConfig *Flags
var featureConfigLock sync.RWMutex

func Get() *Flags {
	featureConfigLock.RLock()
	defer featureConfigLock.RUnlock()

	if featureConfig == nil {
		return &Flags{}
	}
	return featureConfig
}
