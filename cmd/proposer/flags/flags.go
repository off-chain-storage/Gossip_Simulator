package flags

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	CuriePRCProviderFlag = &cli.StringFlag{
		Name:  "curie-rpc-provider",
		Usage: "The RPC provider endpoint for the curie network",
	}
	// HTTPHost specifies a http host for the proposer client.
	HTTPHost = &cli.StringFlag{
		Name:  "http-server-host",
		Usage: "The host on which the HTTP server runs on",
	}
	// HTTPPort enables a http to be exposed for the proposer client.
	HTTPPort = &cli.IntFlag{
		Name:  "http-server-port",
		Usage: "Enable HTTP Server for JSON requests",
	}
	// GrpcRetriesFlag defines the number of times to retry a failed gRPC request.
	GrpcRetriesFlag = &cli.UintFlag{
		Name:  "grpc-retries",
		Usage: "Number of attempts to retry gRPC requests",
		Value: 5,
	}
	// GrpcRetryDelayFlag defines the interval to retry a failed gRPC request.
	GrpcRetryDelayFlag = &cli.DurationFlag{
		Name:  "grpc-retry-delay",
		Usage: "The amount of time between gRPC retry requests.",
		Value: 1 * time.Second,
	}
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = &cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
)
