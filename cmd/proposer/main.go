package main

import (
	"flag-example/cmd"
	"flag-example/cmd/proposer/flags"
	"flag-example/proposer/node"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var appFlags = []cli.Flag{
	// Common Flag
	flags.CuriePRCProviderFlag,
	flags.HTTPHost,
	flags.HTTPPort,
	flags.GrpcRetriesFlag,
	flags.GrpcRetryDelayFlag,
	flags.CertFlag,
	cmd.RPCMaxPageSizeFlag,
	cmd.GrpcMaxCallRecvMsgSizeFlag,
	cmd.DataDirFlag,
	cmd.VerbosityFlag,
	cmd.LogFormat,
	cmd.LogFileName,
}

func startNode(ctx *cli.Context) error {
	proposerClient, err := node.NewProposerClient(ctx)
	if err != nil {
		return err
	}

	proposerClient.Start()
	return nil
}

func main() {
	app := cli.App{}
	app.Name = "Proposer-node"
	app.Usage = "this is a Proposer-node for curie network"
	app.Action = func(ctx *cli.Context) error {
		if err := startNode(ctx); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}

	app.Flags = appFlags

	var log = logrus.WithField("prefix", "main")

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}
