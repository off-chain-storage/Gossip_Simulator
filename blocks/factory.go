package blocks

import (
	"flag-example/blocks/interfaces"
	curiepb "flag-example/proto"

	"github.com/pkg/errors"
)

var (
	ErrUnsupportedSignedCurieBlock = errors.New("unsupported signed curie block")
)

func NewSignedBlock(i interface{}) (interfaces.SignedCurieBlock, error) {
	switch b := i.(type) {
	case *curiepb.SignedCurieBlockForOG:
		return initSignedBlockForOGFromProto(b)
	case *curiepb.SignedCurieBlockForNG:
		return initSignedBlockForNGFromProto(b)
	default:
		return nil, errors.Wrapf(ErrUnsupportedSignedCurieBlock, "unable to create block from type %T", i)
	}
}

func NewCurieBlock(cb *curiepb.CurieBlock, blockData []byte) (interfaces.ReadOnlyCurieBlock, error) {
	cb.DummyData = blockData
	return initBlockFromProto(cb)
}

func BuildSignedCurieBlockForOG(blk interfaces.ReadOnlyCurieBlock, sig_r []byte, sig_s []byte) (*curiepb.SignedCurieBlockForOG, error) {
	pb, err := blk.Proto()
	if err != nil {
		return nil, err
	}

	b, ok := pb.(*curiepb.CurieBlock)
	if !ok {
		return nil, err
	}

	// Build SignedBlockData
	signedBlock := &curiepb.SignedCurieBlockForOG{
		Body: b,
		Signature: &curiepb.Signature{
			SigR: sig_r,
			SigS: sig_s,
		},
	}

	return signedBlock, nil
}

func BuildSignedCurieBlockForNG(sig_r []byte, sig_s []byte) (*curiepb.SignedCurieBlockForNG, error) {
	// Build SignedBlockData
	signedBlock := &curiepb.SignedCurieBlockForNG{
		Signature: &curiepb.Signature{
			SigR: sig_r,
			SigS: sig_s,
		},
	}

	return signedBlock, nil
}
