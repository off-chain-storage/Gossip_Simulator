package client

import (
	"context"
	"flag-example/proposer/client/iface"
)

func run(ctx context.Context, p iface.Proposer) {
	// Initialize proposer keyManager before send to curie node
	if err := p.WaitForKeyManagerInitialization(ctx); err != nil {
		log.WithError(err).Fatal("Could not initialize KeyManager")
	}

	// Send PubKey To CurieNode
	if err := p.SendPubKeyToCurieNode(ctx); err != nil {
		log.WithError(err).Fatal("Could not send public key to Curie Node")
	}
}
