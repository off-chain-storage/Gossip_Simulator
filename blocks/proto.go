package blocks

import (
	"flag-example/crypto/ecdsa/common"
	curieecdsad "flag-example/crypto/ecdsa/ecdsad"
	"flag-example/crypto/hash"
	curiepb "flag-example/proto"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Implement interfaces for CurieBlock
func (cb *CurieBlock) Proto() (proto.Message, error) {
	return &curiepb.CurieBlock{
		DummyData: cb.dummyData,
	}, nil
}

func (cb *CurieBlock) Hash() []byte {
	return hash.Hash(cb.dummyData)
}

func (cb *CurieBlock) Body() []byte {
	return cb.dummyData
}

// Implement interfaces for SignedCurieBlockForOG
func (sbo *SignedCurieBlockForOG) Proto() (proto.Message, error) {
	blockMsg, err := sbo.body.Proto()
	if err != nil {
		return nil, err
	}

	block, ok := blockMsg.(*curiepb.CurieBlock)
	if !ok {
		return nil, errors.New("Failed to convert to curiepb.CurieBlock")
	}

	sigRawData, err := sbo.signature.Proto()
	if err != nil {
		return nil, err
	}

	signature, ok := sigRawData.(*curiepb.Signature)
	if !ok {
		return nil, errors.New("Failed to convert to curiepb.CurieBlock")
	}

	return &curiepb.SignedCurieBlockForOG{
		Body:      block,
		Signature: signature,
	}, nil
}

func (sbo *SignedCurieBlockForOG) Hash() []byte {
	return hash.Hash(sbo.body.dummyData)
}

// ReceiveBlock() 함수에서 블록 검증을 위해 서명을 빼오는 함수가 추가되어야함,
// 이때 common.Signature.Verify() 함수 호출을 위해 common.Signature 타입으로 반환해야함
func (sbo *SignedCurieBlockForOG) Signature() common.Signature {
	return sbo.signature
}

// Implement interfaces for SignedCurieBlockForNG
func (sbn *SignedCurieBlockForNG) Proto() (proto.Message, error) {
	sigRawData, err := sbn.signature.Proto()
	if err != nil {
		return nil, err
	}

	signature, ok := sigRawData.(*curiepb.Signature)
	if !ok {
		return nil, errors.New("Failed to convert to curiepb.Signature")
	}

	return &curiepb.SignedCurieBlockForNG{
		Signature: signature,
	}, nil
}

func (sbn *SignedCurieBlockForNG) Hash() []byte {
	sigHash := make([]byte, 0, 64)
	return hash.Hash(sigHash)
}

func (sbn *SignedCurieBlockForNG) Signature() common.Signature {
	return sbn.signature
}

func initBlockFromProto(pb *curiepb.CurieBlock) (*CurieBlock, error) {
	blk := &CurieBlock{
		dummyData: pb.DummyData,
	}

	return blk, nil
}

func initSignedBlockForOGFromProto(pb *curiepb.SignedCurieBlockForOG) (*SignedCurieBlockForOG, error) {
	block, err := initBlockFromProto(pb.Body)
	if err != nil {
		return nil, err
	}

	sig, err := curieecdsad.InitSignFromProto(pb.Signature)
	if err != nil {
		return nil, err
	}

	b := &SignedCurieBlockForOG{
		body:      block,
		signature: sig,
	}

	return b, nil
}

func initSignedBlockForNGFromProto(pb *curiepb.SignedCurieBlockForNG) (*SignedCurieBlockForNG, error) {
	sig, err := curieecdsad.InitSignFromProto(pb.Signature)
	if err != nil {
		return nil, err
	}

	b := &SignedCurieBlockForNG{
		signature: sig,
	}

	return b, nil
}
