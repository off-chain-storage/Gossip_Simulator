package main

import (
	"context"
	"fmt"
	"os"

	"flag-example/cmd"
	"flag-example/cmd/curie-node/flags"
	"flag-example/curie-node/node"

	golog "github.com/ipfs/go-log/v2"

	"github.com/golang/gddo/log"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var appFlags = []cli.Flag{
	// Common Flag
	cmd.RPCMaxPageSizeFlag,
	cmd.GrpcMaxCallRecvMsgSizeFlag,
	cmd.DataDirFlag,
	cmd.VerbosityFlag,
	cmd.LogFormat,
	cmd.LogFileName,

	// P2P Flag
	cmd.BootstrapNode,
	cmd.BootstrapNodeRegion,
	cmd.NoDiscovery,
	cmd.P2PUDPPort,
	cmd.P2PTCPPort,

	// cmd.P2PIP,
	cmd.P2PHost,
	cmd.P2PHostDNS,
	cmd.P2PMaxPeers,
	cmd.P2PPrivKey,
	cmd.PubsubQueueSize,
	cmd.IsPublisher,

	// RedisDB Flag
	cmd.RedisDBAddrFlag,
	cmd.PoolSizeFlag,
	cmd.MaxIdleConnsFlag,

	// RPC Flag
	flags.RPCHost,
	flags.RPCPort,

	// Monitor Flag
	cmd.MonitorUDPAddrFlag,

	// HTTP Flag
	cmd.HTTPPort,
	cmd.HTTPHost,

	// 필요 없는 것
	// cmd.BackupWebhookOutputDir,
	// cmd.MinimalConfigFlag,
	// cmd.E2EConfigFlag,
	// cmd.StaticPeers,
	// cmd.RelayNode,
	// cmd.P2PStaticID,
	// cmd.P2PMetadata,
	// cmd.P2PAllowList,
	// cmd.P2PDenyList,
	// cmd.EnableTracingFlag,
	// cmd.TracingProcessNameFlag,
	// cmd.TracingEndpointFlag,
	// cmd.TraceSampleFractionFlag,
	// cmd.MonitoringHostFlag,
	// cmd.DisableMonitoringFlag,
	// cmd.ClearDB,
	// cmd.ForceClearDB,
	// cmd.MaxGoroutines,
	// cmd.EnableUPnPFlag,
	// cmd.ConfigFileFlag,
	// cmd.ChainConfigFileFlag,
	// cmd.AcceptTosFlag,
	// cmd.RestoreSourceFileFlag,
	// cmd.RestoreTargetDirFlag,
	// cmd.ValidatorMonitorIndicesFlag,
	// cmd.ApiTimeoutFlag,
}

func main() {
	rctx, cancel := context.WithCancel(context.Background())

	app := cli.App{}
	app.Name = "Curie-node"
	app.Usage = "this is a Curie-node implementation for Propagation Experiment"
	app.Action = func(ctx *cli.Context) error {
		if err := startNode(ctx, cancel); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}

	app.Flags = appFlags

	if err := app.RunContext(rctx, os.Args); err != nil {
		log.Error(rctx, err.Error())
	}
}

func startNode(ctx *cli.Context, cancel context.CancelFunc) error {
	verbosity := ctx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	// Set libp2p logger to only panic logs for the info level.
	golog.SetAllLoggers(golog.LevelPanic)

	curie, err := node.New(ctx, cancel)
	if err != nil {
		return fmt.Errorf("unable to start curie node: %w", err)
	}
	curie.Start()
	return nil
}
