package db

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) buildOptions() *redis.Options {
	cfg := s.cfg

	cfg.PoolFIFO = true

	cfg.Dialer = func(ctx context.Context, network string, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, 5*time.Second)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	options := &redis.Options{
		Addr:         cfg.DbAddr,
		PoolSize:     int(cfg.PoolSize),
		PoolFIFO:     cfg.PoolFIFO,
		Dialer:       cfg.Dialer,
		MaxIdleConns: int(cfg.MaxIdleConns),
	}

	return options
}
