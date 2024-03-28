package sync

import (
	"context"
	"flag-example/blocks/interfaces"
	"flag-example/crypto/hash"

	"github.com/off-chain-storage/GoSphere/sdk"
	"github.com/pkg/errors"
)

type BlockReceiver interface {
	ReceiveOGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error
	ReceiveNGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error
}

func (s *Service) ReceiveOGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error {
	/* Check Received Data for Validation */
	// Hashing Received Data
	hash := block.Hash()

	// Decryption Signature && Compare Hashing and Decryption Signature
	sig := block.Signature()

	if sig.Verify(s.pubKey, hash) {
		log.Info("Received Data is Valid")
	} else {
		log.Error("Received Data is Non-Valid")
		return errors.New("Received Data is Non-Valid")
	}

	return nil
}

func (s *Service) ReceiveNGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error {
	/* Check Received Data for Validation */
	// Decryption Signature && Compare Hashing and Decryption Signature
	sig := block.Signature()

	msgChan := sdk.ReadMessage(ctx)

	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return errors.New("Failed to read message")
			}

			hash := hash.Hash(msg)

			if sig.Verify(s.pubKey, hash) {
				log.Info("Received Data is Valid")
				return nil
			} else {
				log.Error("Received Data is Non-Valid")
				return errors.New("Received Data is Non-Valid")
			}

		case <-ctx.Done():
			return nil
		}
	}
}
