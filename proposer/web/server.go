package web

import (
	"context"
	"flag-example/proposer/client"
	"flag-example/proposer/client/iface"
	"fmt"
	"net"

	"github.com/gofiber/fiber/v3"
	"github.com/off-chain-storage/GoSphere/sdk"
	"github.com/pkg/errors"
)

type Config struct {
	Host            string
	Port            string
	CertFlag        string
	KeyFlag         string
	Router          *fiber.App
	ProposerService *client.ProposerService
}

type Server struct {
	ctx               context.Context
	cancel            context.CancelFunc
	host              string
	port              string
	withCert          string
	withKey           string
	router            *fiber.App
	proposerService   *client.ProposerService
	curieNodeProposer iface.Proposer
	pmanager          *sdk.PManager

	// UDP Server - temp service
	udpServer *net.UDPAddr
	conn      *net.UDPConn
}

func NewServer(ctx context.Context, cfg *Config) *Server {
	ctx, cancel := context.WithCancel(ctx)

	server := &Server{
		ctx:             ctx,
		cancel:          cancel,
		host:            cfg.Host,
		port:            cfg.Port,
		withCert:        cfg.CertFlag,
		withKey:         cfg.KeyFlag,
		router:          cfg.Router,
		proposerService: cfg.ProposerService,
	}

	// Register Proposer Web Server's Router
	if err := server.InitializeRoutes(); err != nil {
		log.WithError(err).Fatal("Could not initialize routes")
	}

	// ** Register PManager Client in Web Service ** //
	pm, err := sdk.NewPManager(ctx)
	if err != nil {
		log.WithError(err).Fatal("Could not create PManager")
	}

	server.pmanager = pm

	// UDP Server - temp service
	if err := server.buildUDPAddr(); err != nil {
		log.WithError(err).Fatal("Could not build UDP address")
	}

	if err := server.Conn(); err != nil {
		log.WithError(err).Fatal("Could not start UDP listener")
	}

	return server
}

func (s *Server) buildUDPAddr() error {
	udpServer, err := net.ResolveUDPAddr("udp4", "43.200.145.206:30004")
	if err != nil {
		return err
	}

	s.udpServer = udpServer
	return nil
}

func (s *Server) Conn() error {
	conn, err := net.DialUDP("udp4", nil, s.udpServer)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *Server) SendUDPMessage(msg string) error {
	_, err := s.conn.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Start() {
	// Setup the Http Server Address
	address := fmt.Sprintf("%s:%s", s.host, s.port)

	// Register Proposer Client in Web Service
	s.curieNodeProposer = s.proposerService.Proposer()

	// Start the Web Server
	go func() {
		s.router.Listen(address)
	}()

	log.WithField("address", address).Info("http listening on address")
}

func (s *Server) InitializeRoutes() error {
	// Register Proposer Web Server's Router
	if s.router == nil {
		return errors.New("no fiber router on server")
	}

	// Register all routes api
	api := s.router.Group("/curie/proposer")
	api.Post("/original", s.ProposeCurieBlockForOG)
	api.Post("/new", s.ProposeCurieBlockForNG)

	log.Info("Initialize Proposer REST API Routes")

	return nil
}

func (s *Server) Stop() error {
	defer s.cancel()

	return nil
}
