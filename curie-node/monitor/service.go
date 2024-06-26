package monitor

import (
	"context"
	"net"
)

type Config struct {
	UDPAddr             string
	InitialSyncComplete chan struct{}
}

type Service struct {
	cfg       *Config
	ctx       context.Context
	cancel    context.CancelFunc
	udpServer *net.UDPAddr
	conn      *net.UDPConn
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	s := &Service{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	if err := s.buildUDPAddr(); err != nil {
		log.WithError(err).Error("Failed to build UDP address")
		return nil, err
	}

	s.SetUDPConn()

	return s, nil
}

func (s *Service) Start() {
	if s.conn == nil {
		s.SetUDPConn()
	}

	close(s.cfg.InitialSyncComplete)

	log.Info("Start UDP listener")
}

func (s *Service) Stop() error {
	defer s.cancel()
	defer s.conn.Close()

	return nil
}
