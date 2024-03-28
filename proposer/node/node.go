package node

import (
	"context"
	"flag-example/cmd"
	"flag-example/cmd/proposer/flags"
	"flag-example/proposer/client"
	"flag-example/proposer/web"
	"flag-example/runtime"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/off-chain-storage/GoSphere/sdk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ProposerClient struct {
	cliCtx   *cli.Context
	ctx      context.Context
	cancel   context.CancelFunc
	services *runtime.ServiceRegistry
	lock     sync.RWMutex
	stop     chan struct{}
}

func NewProposerClient(cliCtx *cli.Context) (*ProposerClient, error) {
	verbosity := cliCtx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(level)

	registry := runtime.NewServiceRegistry()
	ctx, cancel := context.WithCancel(cliCtx.Context)
	proposerClient := &ProposerClient{
		cliCtx:   cliCtx,
		ctx:      ctx,
		cancel:   cancel,
		services: registry,
		stop:     make(chan struct{}),
	}

	router := newRouter(cliCtx)
	if err := proposerClient.initialize(cliCtx, router); err != nil {
		return nil, err
	}

	return proposerClient, nil
}

func newRouter(cliCtx *cli.Context) *fiber.App {
	r := fiber.New()
	return r
}

func (c *ProposerClient) initialize(cliCtx *cli.Context, router *fiber.App) error {
	/* RPC Client - Propose Data to Curie Node by call curie-node's propagation rpc method */
	// Signature Module
	// Create Block Data Module
	if err := c.registerProposerService(cliCtx); err != nil {
		return err
	}

	/* Web Server - User send request by HTTP POST API */
	if err := c.registerWebService(router); err != nil {
		return err
	}

	httpHost := cliCtx.String(flags.HTTPHost.Name)
	httpPort := cliCtx.Int(flags.HTTPPort.Name)
	httpAddress := fmt.Sprintf("http://%s:%d", httpHost, httpPort)
	log.WithField("address", httpAddress).Info(
		"Starting Proposer HTTP Server on address",
	)
	return nil
}

func (c *ProposerClient) registerProposerService(cliCtx *cli.Context) error {
	// curie-node rpc endpoint
	endpoint := cliCtx.String(flags.CuriePRCProviderFlag.Name)
	// grpc Recv Msg Size
	maxCallRecvMsgSize := c.cliCtx.Int(cmd.GrpcMaxCallRecvMsgSizeFlag.Name)
	// grpc Options
	grpcRetries := c.cliCtx.Uint(flags.GrpcRetriesFlag.Name)
	grpcRetryDelay := c.cliCtx.Duration(flags.GrpcRetryDelayFlag.Name)
	// grpc Cert
	cert := c.cliCtx.String(flags.CertFlag.Name)

	// Create New Proposer Client Instance
	p, err := client.NewProposerClient(c.cliCtx.Context, &client.Config{
		Endpoint:                   endpoint,
		GrpcMaxCallRecvMsgSizeFlag: maxCallRecvMsgSize,
		GrpcRetriesFlag:            grpcRetries,
		GrpcRetryDelay:             grpcRetryDelay,
		CertFlag:                   cert,
	})
	if err != nil {
		return errors.Wrap(err, "could not initialize validator service")
	}

	return c.services.RegisterService(p)
}

func (c *ProposerClient) registerWebService(router *fiber.App) error {
	var ps *client.ProposerService
	if err := c.services.FetchService(&ps); err != nil {
		return err
	}

	// HTTP Server address
	httpHost := c.cliCtx.String(flags.HTTPHost.Name)
	httpPort := c.cliCtx.Int(flags.HTTPPort.Name)

	// Create New Web Server Instance
	webServer := web.NewServer(c.ctx, &web.Config{
		Host:            httpHost,
		Port:            fmt.Sprintf("%d", httpPort),
		Router:          router,
		ProposerService: ps,
	})

	return c.services.RegisterService(webServer)
}

func (c *ProposerClient) Start() {
	c.lock.Lock()
	log.Info("Starting Proposer Node")

	// Setup Propagation Module
	go sdk.SetupPropagationModule()

	c.services.StartAll()

	stop := c.stop
	c.lock.Unlock()

	go func() {
		// OS Signal을 수신할 수 있는 채널 생성 (버퍼 크기 : 1)
		sigc := make(chan os.Signal, 1)
		// sigc 채널이 SIGINT, SIGTERM 신호를 받을 수 있도록 설정
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		// Start() 함수가 종료될 때 signal.Stop(sigc) 함수 실행
		defer signal.Stop(sigc)
		// sigc 채널에서 신호 수신 대기하기
		<-sigc
		log.Info("Got interrupt, shutting down...")
		// 일단 생략
		// debug.Exit(c.cliCtx)

		go c.Close()
		// 추가적인 SIGINT, SIGTERM 신호 처리
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.WithField("times", i-1).Info("Already shutting down, interrupt more to panic")
			}
		}
		panic("Panic closing the beacon node")
	}()

	<-stop
}

func (c *ProposerClient) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.services.StopAll()
	log.Info("Stopping Curie Proposer")
	c.cancel()

	// 여기서 Stop Signal 날려주면 stop chan 시그널 날려줌
	close(c.stop)

}
