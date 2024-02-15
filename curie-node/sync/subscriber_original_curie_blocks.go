package sync

import (
	"context"
	"flag-example/blocks"

	"google.golang.org/protobuf/proto"
)

func (s *Service) originalCurieBlockSubscriber(ctx context.Context, msg proto.Message) error {
	log.Info("@@ STEP_1 @@")
	signed, err := blocks.NewSignedBlock(msg)
	if err != nil {
		return err
	}

	log.Info("@@ STEP_2 @@")
	if err := s.cfg.receiveModule.ReceiveOGBlock(ctx, signed, s.pubKey); err != nil {
		return err
	}

	log.Info("@@ STEP_3 @@")
	if err := s.cfg.monitor.SendUDPMessage(s.cfg.p2p.PeerID().String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) newCurieBlockSubscriber(ctx context.Context, msg proto.Message) error {
	// signed, err := blocks.NewSignedBlock(msg)
	// if err != nil {
	// 	return err
	// }

	// 여기서 ReceiveBlock 함수 추가 및 정리
	// if err := s.cfg.

	return nil
}
