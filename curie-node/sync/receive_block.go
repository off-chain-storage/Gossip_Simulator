package sync

import (
	"context"
	"flag-example/blocks/interfaces"
	"flag-example/crypto/hash"

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
		// log.Info("Received Data is Valid")
		return nil
	} else {
		log.Error("Received Data is Non-Valid")
		return errors.New("Received Data is Non-Valid")
	}
}

func (s *Service) ReceiveNGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error {
	/* Check Received Data for Validation */
	// Decryption Signature && Compare Hashing and Decryption Signature
	sig := block.Signature()

	for {
		msg, err := s.pmSub.ReadMessage(context.TODO())
		if err != nil {
			s.pmSub.Cancel()
			return errors.New("Failed to read message")
		}
		hash := hash.Hash(msg.Data)

		if sig.Verify(s.pubKey, hash) {
			// log.Info("Received Data is Valid")
			return nil
		} else {
			log.Error("Received Data is Non-Valid")
			return errors.New("Received Data is Non-Valid")
		}
	}
}
