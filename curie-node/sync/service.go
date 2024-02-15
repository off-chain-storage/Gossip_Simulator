package sync

import (
	"context"
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"
	"flag-example/curie-node/db"
	"flag-example/curie-node/monitor"
	"flag-example/curie-node/p2p"
)

type config struct {
	p2p     p2p.P2P
	curieDB db.ReadOnlyRedisDB
	monitor monitor.Monitor
}

type Service struct {
	cfg                 *config
	ctx                 context.Context
	cancel              context.CancelFunc
	subHandler          *subTopicHandler
	initialSyncComplete chan struct{}
	pubKey              *ecdsa.PublicKey
}

func NewService(ctx context.Context, opts ...Option) *Service {
	ctx, cancel := context.WithCancel(ctx)
	r := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    &config{},
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil
		}
	}

	// Get Proposer's Public Key
	// r.getPubKey()

	r.subHandler = newSubTopicHandler()

	return r
}

// Start the regular sync service
func (s *Service) Start() {
	log.Info("Start Sync Service")

	go s.registerHandlers()
}

func (s *Service) Stop() error {
	s.cancel()

	return nil
}

func (s *Service) registerHandlers() {
	select {
	case <-s.initialSyncComplete:
		// Register respective pubsub handlers at state synced event.
		s.registerSubscribers()
		return
	case <-s.ctx.Done():
		log.Debug("Context closed, exiting goroutine")
		return
	}
}

func (s *Service) getPubKey() {
	s.pubKey = ecdsad.GetPublicKey()
}
