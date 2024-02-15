package rpc

import (
	"context"
	"flag-example/curie-node/db"
	"flag-example/curie-node/monitor"
	"flag-example/curie-node/p2p"
	"flag-example/curie-node/rpc/proposer"
	curiepb "flag-example/proto"
	"fmt"
	"net"

	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Host         string
	Port         string
	CertFlag     string
	KeyFlag      string
	MaxMsgSize   int
	Broadcaster  p2p.Broadcaster
	PeersFetcher p2p.PeersProvider
	Monitor      monitor.Monitor
	DB           db.AccessRedisDB
}

type Service struct {
	cfg                 *Config
	ctx                 context.Context
	cancel              context.CancelFunc
	listener            net.Listener
	grpcServer          *grpc.Server
	connectedRPCClients map[net.Addr]bool
}

func NewService(ctx context.Context, cfg *Config) *Service {
	ctx, cancel := context.WithCancel(ctx)
	s := &Service{
		cfg:                 cfg,
		ctx:                 ctx,
		cancel:              cancel,
		connectedRPCClients: make(map[net.Addr]bool),
	}

	// Register RPC Server
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.WithError(err).Fatalln("Could not listen to port in Start()", address)
	}
	s.listener = lis
	log.WithField(
		"address", address,
	).Info("gRPC server listening on port")

	// Register gRPC Server's Option
	opts := []grpc.ServerOption{
		// gRPC 서버에서 OpenCensus의 통계 처리 기능을 활성화하는 옵션
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		// gRPC 서버에서 최대 메시지 크기를 설정하는 옵션
		grpc.MaxRecvMsgSize(s.cfg.MaxMsgSize),
	}
	if s.cfg.CertFlag != "" && s.cfg.KeyFlag != "" {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFlag, s.cfg.KeyFlag)
		if err != nil {
			log.WithError(err).Fatal("Could not load TLS keys")
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC server. Please use TLS keys")
	}
	s.grpcServer = grpc.NewServer(opts...)

	return s
}

func (s *Service) Start() {
	proposerServer := &proposer.Server{
		Ctx:     s.ctx,
		DB:      s.cfg.DB,
		P2P:     s.cfg.Broadcaster,
		Monitor: s.cfg.Monitor,
	}

	curiepb.RegisterCurieNodeProposerServer(s.grpcServer, proposerServer)

	reflection.Register(s.grpcServer)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				log.WithError(err).Errorf("Could not serve gRPC")
			}
		}
	}()
}

func (s *Service) Stop() error {
	s.cancel()
	if s.listener != nil {
		s.grpcServer.GracefulStop()
		log.Debug("Initiated graceful stop of gRPC server")
	}
	return nil
}
