package p2p

import (
	"flag-example/curie-node/db"
)

type Config struct {
	NoDiscovery       bool
	BootstrapNodeAddr []string
	HostAddress       string
	HostDNS           string
	PrivateKey        string
	DataDir           string
	TCPPort           uint
	UDPPort           uint
	MaxPeers          uint
	QueueSize         uint
	IsPublisher       bool
	DB                db.AccessRedisDB

	// LocalIP           string
	// Discv5BootStrapAddr []string
	// EnableUPnP          bool
	// StaticPeerID        bool
	// StaticPeers         []string
	// RelayNodeAddr       string
	// MetaDataDir         string
	// AllowListCIDR       string
	// DenyListCIDR        []string
	// StateNotifier       statefeed.Notifier
	// DB                  db.ReadOnlyDatabase
	// ClockWaiter         startup.ClockWaiter
}

const defaultPubsubQueueSize = 600

func validateConfig(cfg *Config) *Config {
	if cfg.QueueSize == 0 {
		log.Warnf("Invalid pubsub queue size of %d initialized, setting the quese size as %d instead", cfg.QueueSize, defaultPubsubQueueSize)
		cfg.QueueSize = defaultPubsubQueueSize
	}
	return cfg
}
