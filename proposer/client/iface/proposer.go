package iface

import (
	"context"
	"flag-example/proposer/keymanager"
)

type Proposer interface {
	// ProposeBlock for Original Gossip
	ProposeCurieBlockForOG(ctx context.Context, blockData []byte) error
	// ProposeBlock for New Gossip
	ProposeCurieBlockForNG(ctx context.Context, blockData []byte) error
	// Send PubKey to Curie Node
	SendPubKeyToCurieNode(ctx context.Context) error
	// Initialize Proposer's KeyManager
	WaitForKeyManagerInitialization(ctx context.Context) error
	// Get KeyManager
	KeyManager() (keymanager.IKeymanager, error)
}
