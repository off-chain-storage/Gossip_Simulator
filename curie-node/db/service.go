package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	started     bool
	startupErr  error
	ctx         context.Context
	cfg         *Config
	cancel      context.CancelFunc
	redisClient *redis.Client
	conn        *redis.Conn
}

func NewRedisClient(ctx context.Context, cfg *Config) (*Service, error) {
	// var err error
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel

	s := &Service{
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
	}

	opts := s.buildOptions()
	r := redis.NewClient(opts)

	s.redisClient = r

	return s, nil
}

func (s *Service) Start() {
	if s.started {
		log.Error("Attempted to start RedisDB Service when it was already started")
		return
	}

	if s.conn == nil {
		log.Warnf("Attempted to make a connection with Redis DB")
		s.SetRedisConn()
	}

	s.started = true
}

func (s *Service) Stop() error {
	defer s.cancel()
	defer s.conn.Close()
	s.started = false

	return nil
}

func (s *Service) RedisClient() redis.Client { return *s.redisClient }
