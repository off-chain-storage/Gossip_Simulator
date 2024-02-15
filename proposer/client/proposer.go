package client

import (
	"context"
	"errors"
	"flag-example/blocks"
	"flag-example/blocks/interfaces"
	"flag-example/proposer/client/iface"
	"flag-example/proposer/keymanager"
	"flag-example/proposer/keymanager/local"

	curieecdsa "flag-example/crypto/ecdsa"
	"flag-example/crypto/hash"
	curiepb "flag-example/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type proposer struct {
	proposerClient iface.ProposerClient
	keyManager     keymanager.IKeymanager
}

// 여기다가 gRPC method 호출하는 함수 만들면 되고, 아래 함수를 HTTP API Handler에서 호출하면 됨
// 여기서 호출하는 gRPC method는 curie-node의 propagation rpc method를 호출하는 것이며
// 이 때, curie-node의 propagation rpc method는 flag-example/curie-node/rpc/curie/server.go와 propose.go에 구현되어 있음 (구현 예정)

// ProposeBlock for Original Gossip
func (p *proposer) ProposeCurieBlockForOG(ctx context.Context, blockData []byte) error {
	// Get Block from Curie Node
	b, err := p.proposerClient.GetBlock(ctx, &emptypb.Empty{})
	if err != nil {
		log.WithError(err).Error("Failed to get block from curie node")
		return err
	}

	// Get CurieBlock with Block Data - Original Data
	wb, err := blocks.NewCurieBlock(b, blockData)
	if err != nil {
		log.WithError(err).Error("Failed to wrap block")
		return err
	}

	// Sign Block with Proposer's Private Key
	sig, err := p.signData(ctx, wb)
	if err != nil {
		log.WithError(err).Error("Failed to sign block")
		return err
	}

	blk, err := blocks.BuildSignedCurieBlockForOG(wb, sig)
	if err != nil {
		log.WithError(err).Error("Failed to build signed curie block")
		return err
	}

	_, err = p.proposerClient.ProposeCurieBlockForOG(ctx, blk)
	if err != nil {
		log.WithError(err).Error("Failed to propose original gossip data")
		return err
	}

	log.Info("Successfully proposed original gossip data")

	return nil
}

// ProposeBlock for New Gossip
func (p *proposer) ProposeCurieBlockForNG(ctx context.Context, blockData []byte) error {
	// Get Block from Curie Node
	b, err := p.proposerClient.GetBlock(ctx, &emptypb.Empty{})
	if err != nil {
		log.WithError(err).Error("Failed to get block from curie node")
		return err
	}

	// Make Hash from Block Data
	h := hash.Hash(blockData)

	// Get CurieBlock with Block Data - Data Hash
	wb, err := blocks.NewCurieBlock(b, h)
	if err != nil {
		log.WithError(err).Error("Failed to wrap block")
		return err
	}

	// Sign Block with Proposer's Private Key
	sig, err := p.signData(ctx, wb)
	if err != nil {
		log.WithError(err).Error("Failed to sign block")
		return err
	}

	blk, err := blocks.BuildSignedCurieBlockForNG(sig)
	if err != nil {
		log.WithError(err).Error("Failed to build signed curie block")
		return err
	}

	_, err = p.proposerClient.ProposeCurieBlockForNG(ctx, blk)
	if err != nil {
		log.WithError(err).Error("Failed to propose original gossip data")
		return err
	}

	log.Info("Successfully proposed original gossip data")

	return nil
}

// Send PubKey to Curie Node
func (p *proposer) SendPubKeyToCurieNode(ctx context.Context) error {
	// 1. Generate Key
	if p.keyManager == nil {
		return errors.New("could not initialize KeyManager")
	}

	// 2. Convert *ecdsa.PubKey to string
	var pubKey string
	pk, err := p.keyManager.FetchValidatingPublicKeys()
	if err == nil {
		pubKey = curieecdsa.ConvertToStringEcdsaPubKey(pk)
	}

	// 3. Send PubKey to Curie Node
	_, err = p.proposerClient.SendProposerPublicKey(ctx, &curiepb.ProposerPublicKeyRequest{
		PublicKey: pubKey,
	})
	if err != nil {
		log.WithField("Proposer's PubKey", pubKey).WithError(err).Error("Failed to send pubkey to curie node")
		return err
	}

	log.WithField("Proposer's PubKey", pubKey).Info("Successfully sent pubkey to curie node")
	return nil
}

func (p *proposer) signData(ctx context.Context, b interfaces.ReadOnlyCurieBlock) ([]byte, error) {
	var pubKey string
	pk, err := p.keyManager.FetchValidatingPublicKeys()
	if err == nil {
		pubKey = curieecdsa.ConvertToStringEcdsaPubKey(pk)
	}

	sig, err := p.keyManager.Sign(ctx, &curiepb.SignRequest{
		PublicKey:  pubKey,
		SigningMsg: b.Hash(),
	})
	if err != nil {
		log.WithError(err).Error("Failed to sign block")
		return nil, err
	}

	return sig.Marshal(), nil
}

func (p *proposer) WaitForKeyManagerInitialization(ctx context.Context) error {
	k, err := local.NewKeyManager(ctx)
	if err != nil {
		return err
	}

	p.keyManager = k

	return nil
}

func (p *proposer) KeyManager() (keymanager.IKeymanager, error) {
	if p.keyManager == nil {
		return nil, errors.New("keymanager is not initialized")
	}
	return p.keyManager, nil
}
