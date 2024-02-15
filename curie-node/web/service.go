package c_web

import (
	"context"
	"flag-example/curie-node/db"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

type Config struct {
	Host   string
	Port   string
	Router *fiber.App
	DB     db.AccessRedisDB
}

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *Config
	host   string
	port   string
	router *fiber.App
}

func NewService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)

	server := &Service{
		ctx:    ctx,
		cancel: cancel,
		host:   cfg.Host,
		port:   cfg.Port,
		router: cfg.Router,
	}

	// Register Proposer Web Server's Router
	if err := server.InitializeRoutes(); err != nil {
		log.WithError(err).Fatal("Could not initialize routes")
	}

	return server
}

func (s *Service) Start() {
	// Setup the Http Server Address
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	// Start the Web Server
	go func() {
		s.router.Listen(address)
	}()

	log.WithField("address", address).Info("http listening on address")
}

func (s *Service) InitializeRoutes() error {
	// Register Proposer Web Server's Router
	if s.router == nil {
		return errors.New("no fiber router on server")
	}

	// Register all routes api
	api := s.router.Group("/curie")
	api.Post("/pub_key", s.StoreProposerPubKey)
	log.Info("Initialize Proposer REST API Routes")

	return nil
}

func (s *Service) Stop() error {
	defer s.cancel()

	return nil
}