package client

import (
	"context"
	"time"

	grpcapi "flag-example/proposer/client/grpc-api"
	"flag-example/proposer/client/iface"
	proposerHelper "flag-example/proposer/helpers"
	"flag-example/proposer/keymanager"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"go.opencensus.io/plugin/ocgrpc"
)

type Config struct {
	Proposer                   iface.Proposer
	CertFlag                   string
	Endpoint                   string
	GrpcMaxCallRecvMsgSizeFlag int
	GrpcRetriesFlag            uint
	GrpcRetryDelay             time.Duration
}

type ProposerService struct {
	ctx                context.Context
	cancel             context.CancelFunc
	withCert           string
	endpoint           string
	grpcRetryDelay     time.Duration
	grpcRetries        uint
	maxCallRecvMsgSize int
	conn               proposerHelper.NodeConnection
	proposer           iface.Proposer
}

func NewProposerClient(ctx context.Context, cfg *Config) (
	*ProposerService,
	error,
) {
	ctx, cancel := context.WithCancel(ctx)
	s := &ProposerService{
		ctx:                ctx,
		cancel:             cancel,
		endpoint:           cfg.Endpoint,
		withCert:           cfg.CertFlag,
		maxCallRecvMsgSize: cfg.GrpcMaxCallRecvMsgSizeFlag,
		grpcRetries:        cfg.GrpcRetriesFlag,
		grpcRetryDelay:     cfg.GrpcRetryDelay,
		proposer:           cfg.Proposer,
	}

	// gRPC Dial Options
	dialOpts := ConstructDialOptions(
		s.maxCallRecvMsgSize,
		s.withCert,
		s.grpcRetries,
		s.grpcRetryDelay,
	)
	if dialOpts == nil {
		return s, nil
	}

	// gRPC Dial Context to Curie gRPC Server
	grpcConn, err := grpc.DialContext(ctx, s.endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}

	if s.withCert != "" {
		log.Info("Established secure gRPC connection")
	} else {
		log.Info("Established insecure gRPC connection")
	}

	s.conn = proposerHelper.NewNodeConnection(grpcConn)

	return s, nil
}

func (p *ProposerService) Start() {
	// 1. Get gRPC Client Conn
	proposerClient := grpcapi.NewGrpcProposerClient(p.conn.GetGrpcClientConn())

	// 2. Create proposer struct instance with gRPC Client Conn
	proStruct := &proposer{
		proposerClient: proposerClient,
	}
	p.proposer = proStruct

	// 3. Start Proposer Service - Send PubKey to Curie Node
	go run(p.ctx, p.proposer)
}

func (p *ProposerService) Stop() error {
	p.cancel()
	log.Info("Stopping service")
	if p.conn != nil {
		return p.conn.GetGrpcClientConn().Close()
	}

	return nil
}

func (p *ProposerService) KeyManager() (keymanager.IKeymanager, error) {
	return p.proposer.KeyManager()
}

func (p *ProposerService) Proposer() iface.Proposer {
	return p.proposer
}

func ConstructDialOptions(
	maxCallRecvMsgSize int,
	withCert string,
	grpcRetries uint,
	grpcRetryDelay time.Duration,
	extraOpts ...grpc.DialOption,
) []grpc.DialOption {
	// Set SSL/TLS Cert for gRPC
	var transportSecurity grpc.DialOption
	if withCert != "" {
		// Create Credentials from file
		creds, err := credentials.NewClientTLSFromFile(withCert, "")
		if err != nil {
			log.WithError(err).Error("Could not get valid credentials")
			return nil
		}

		// Set gRPC SSL/TLS Option
		transportSecurity = grpc.WithTransportCredentials(creds)
	} else {
		// Don't use SSL/TLS for gRPC
		transportSecurity = grpc.WithInsecure()
	}

	// gRPC Max Message Size Option - 50MB
	if maxCallRecvMsgSize == 0 {
		maxCallRecvMsgSize = 10 * 5 << 20
	}

	// Set gRPC Dial Options
	dialOpts := []grpc.DialOption{
		// SSL/TLS Option
		transportSecurity,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize),
			grpcretry.WithMax(grpcRetries),
			grpcretry.WithBackoff(grpcretry.BackoffLinear(grpcRetryDelay)),
		),
		// 통계 핸들러??
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),

		// gRPC Unary, Stream Interceptor Option - X
		// gRPC Round Robin Option - X
	}

	dialOpts = append(dialOpts, extraOpts...)
	return dialOpts
}
