package node

import (
	"context"
	"flag-example/cmd"
	"flag-example/cmd/curie-node/flags"
	"flag-example/curie-node/db"
	"flag-example/curie-node/monitor"
	"flag-example/curie-node/node/registration"
	"flag-example/curie-node/p2p"
	"flag-example/curie-node/rpc"
	regularsync "flag-example/curie-node/sync"
	c_web "flag-example/curie-node/web"
	"fmt"

	"flag-example/runtime"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/urfave/cli/v2"
)

type CurieNode struct {
	cliCtx              *cli.Context
	ctx                 context.Context
	cancel              context.CancelFunc
	services            *runtime.ServiceRegistry
	lock                sync.RWMutex
	db                  db.RedisDB
	stop                chan struct{}
	initialSyncComplete chan struct{}
}

// 새로운 Curie Node instance 만들기 (w/ config options, register other things)
func New(cliCtx *cli.Context, cancel context.CancelFunc) (*CurieNode, error) {
	// Curie Node의 Network Config 설정
	configureNetwork(cliCtx)

	// golang에서 reflect, type 등을 활용하는데 사용 되는 문법 <- 추가로 공부할 것
	registry := runtime.NewServiceRegistry()

	ctx := cliCtx.Context
	curie := &CurieNode{
		cliCtx:   cliCtx,
		ctx:      ctx,
		cancel:   cancel,
		services: registry,
	}

	curie.initialSyncComplete = make(chan struct{})

	// Register Redis DB for storing sender's public key
	log.Debugln("Starting Redis DB")
	if err := curie.startRedisDB(cliCtx); err != nil {
		return nil, err
	}

	// Register P2P Service for Gossip
	log.Debugln("Registering P2P Service")
	if err := curie.registerP2P(cliCtx); err != nil {
		return nil, err
	}

	// Register Web Service for Storing Public Key
	log.Debugln("Registering Web Service")
	router := newRouter(cliCtx)
	if err := curie.registerWebService(router); err != nil {
		return nil, err
	}

	// Register Monitor Service for ACK
	log.Debugln("Registering Monitor Service")
	if err := curie.registerMonitoringService(curie.initialSyncComplete); err != nil {
		return nil, err
	}

	// Register RPC Service for Connection with Validator Node
	log.Debugln("Registering RPC Service")
	if err := curie.registerRPCService(); err != nil {
		return nil, err
	}

	// Register Sync Service for Syncing
	log.Debugln("Registering Sync Service")
	if err := curie.registerSyncService(curie.initialSyncComplete); err != nil {
		return nil, err
	}

	return curie, nil
}

func newRouter(cliCtx *cli.Context) *fiber.App {
	r := fiber.New()
	return r
}

func (c *CurieNode) Start() {
	// Mutex Lock
	c.lock.Lock()

	log.Info("Starting cuire node")

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

func (c *CurieNode) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Info("Stopping curie node")
	c.services.StopAll()
	// 여기에 디비 닫는거 있어야 함
	c.cancel()

	// 여기서 Stop Signal 날려주면 stop chan 시그널 날려줌
	close(c.stop)
}

func (c *CurieNode) startRedisDB(cliCtx *cli.Context) error {
	dbAddr := cliCtx.String(cmd.RedisDBAddrFlag.Name)
	poolSize := cliCtx.Uint(cmd.P2PMaxPeers.Name)
	maxIdleConns := cliCtx.Uint(cmd.P2PMaxPeers.Name)

	svc, err := db.NewRedisClient(c.ctx, &db.Config{
		DbAddr:       dbAddr,
		PoolSize:     poolSize,
		MaxIdleConns: maxIdleConns,
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect database")
		return err
	}

	svc.SetRedisConn()

	// 여기 좀 수정해야함;;
	c.db = svc

	log.WithField("database-addr", dbAddr).Info("Connecting DB")

	return c.services.RegisterService(svc)
}

func (c *CurieNode) registerP2P(cliCtx *cli.Context) error {
	// 부트 스트랩 노드 등록 하기
	bootstrapNodeAddrs, dataDir, err := registration.P2PPreregistration(cliCtx)
	if err != nil {
		return err
	}

	// Register New P2P Service
	svc, err := p2p.NewService(c.ctx, &p2p.Config{
		BootstrapNodeAddr: bootstrapNodeAddrs,
		DataDir:           dataDir,
		NoDiscovery:       cliCtx.Bool(cmd.NoDiscovery.Name),
		IsPublisher:       cliCtx.Bool(cmd.IsPublisher.Name),
		HostAddress:       cliCtx.String(cmd.P2PHost.Name),
		HostDNS:           cliCtx.String(cmd.P2PHostDNS.Name),
		PrivateKey:        cliCtx.String(cmd.P2PPrivKey.Name),
		TCPPort:           cliCtx.Uint(cmd.P2PTCPPort.Name),
		UDPPort:           cliCtx.Uint(cmd.P2PUDPPort.Name),
		MaxPeers:          cliCtx.Uint(cmd.P2PMaxPeers.Name),
		DB:                c.db,
		// LocalIP:           cliCtx.String(cmd.P2PIP.Name),
	})
	if err != nil {
		return err
	}
	return c.services.RegisterService(svc)
}

func (c *CurieNode) registerWebService(router *fiber.App) error {
	// Register Proposer Web Server's Router
	httpHost := c.cliCtx.String(cmd.HTTPHost.Name)
	httpPort := c.cliCtx.Int(cmd.HTTPPort.Name)

	webServer := c_web.NewService(c.ctx, &c_web.Config{
		Host:   httpHost,
		Port:   fmt.Sprintf("%d", httpPort),
		Router: router,
		DB:     c.db,
	})

	return c.services.RegisterService(webServer)
}

func (c *CurieNode) registerSyncService(initialSyncComplete chan struct{}) error {
	rs := regularsync.NewService(
		c.ctx,
		regularsync.WithP2P(c.fetchP2P()),
		regularsync.WithDatabase(c.db),
		regularsync.WithMonitor(c.fetchMonitor()),
		regularsync.WithInitialSyncComplete(initialSyncComplete),
	)

	log.Info("Register Sync Service")

	return c.services.RegisterService(rs)
}

func (c *CurieNode) registerRPCService() error {
	rpcService := rpc.NewService(c.ctx, &rpc.Config{
		Host:        c.cliCtx.String(flags.RPCHost.Name),
		Port:        c.cliCtx.String(flags.RPCPort.Name),
		MaxMsgSize:  c.cliCtx.Int(cmd.GrpcMaxCallRecvMsgSizeFlag.Name),
		Broadcaster: c.fetchP2P(),
		Monitor:     c.fetchMonitor(),
		DB:          c.db,
	})

	return c.services.RegisterService(rpcService)
}

func (c *CurieNode) registerMonitoringService(complete chan struct{}) error {
	udpAddr := c.cliCtx.String(cmd.MonitorUDPAddrFlag.Name)

	ms, err := monitor.NewService(c.ctx, &monitor.Config{
		UDPAddr:             udpAddr,
		InitialSyncComplete: complete,
	})
	if err != nil {
		log.WithError(err).Error("Failed to register monitor service")
		return err
	}

	return c.services.RegisterService(ms)
}

func (c *CurieNode) fetchP2P() p2p.P2P {
	var p *p2p.Service
	if err := c.services.FetchService(&p); err != nil {
		panic(err)
	}
	return p
}

func (c *CurieNode) fetchMonitor() monitor.Monitor {
	var m *monitor.Service
	if err := c.services.FetchService(&m); err != nil {
		panic(err)
	}
	return m
}
