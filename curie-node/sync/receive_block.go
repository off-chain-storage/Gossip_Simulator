package sync

import (
	"context"
	"flag-example/blocks/interfaces"
)

type BlockReceiver interface {
	ReceiveOGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error
	ReceiveNGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error
}

func (s *Service) ReceiveOGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error {
	/* Check Received Data for Validation */

	// Hashing Received Data
	hash := block.Hash()
	log.Info(hash)

	// Decryption Signature && Compare Hashing and Decryption Signature
	// 1. 수신 데이터로부터 서명 데이터 추출하기
	sig := block.Signature()

	// 2. 공개키 이용하여 Verify() 함수 호출하기
	if sig.Verify(s.pubKey, hash) {
		// If it is valid, Send Normal ACK to Check Node
		log.Info("Received Data is Valid")
	} else {
		// If it is invalid, Send NACK to Check Node
		log.Info("Received Data is Non-Valid")
	}

	return nil
}

func ReceiveNGBlock(ctx context.Context, block interfaces.SignedCurieBlock) error {
	/* Check Received Data for Validation */

	// Hashing Received Data
	// hash := block.Hash()

	// Decryption Signature && Compare Hashing and Decryption Signature
	// 1. 수신 데이터로부터 서명 데이터 추출하기
	// sig := block.Signature()

	// 2. 공개키 이용하여 Verify() 함수 호출하기

	// If it is valid, Send Normal ACK to Check Node

	// If it is invalid, Send NACK to Check Node

	return nil
}
