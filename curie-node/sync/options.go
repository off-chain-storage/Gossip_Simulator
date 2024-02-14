package sync

import (
	"flag-example/curie-node/db"
	"flag-example/curie-node/monitor"
	"flag-example/curie-node/p2p"
)

type Option func(s *Service) error

func WithP2P(p2p p2p.P2P) Option {
	return func(s *Service) error {
		s.cfg.p2p = p2p
		return nil
	}
}

func WithDatabase(db db.ReadOnlyRedisDB) Option {
	return func(s *Service) error {
		s.cfg.curieDB = db
		return nil
	}
}

func WithMonitor(monitor monitor.Monitor) Option {
	return func(s *Service) error {
		s.cfg.monitor = monitor
		return nil
	}
}

func WithInitialSyncComplete(c chan struct{}) Option {
	return func(s *Service) error {
		s.initialSyncComplete = c
		return nil
	}
}
