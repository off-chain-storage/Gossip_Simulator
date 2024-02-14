package iface

import (
	"context"

	curiepb "flag-example/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

type ProposerClient interface {
	GetBlock(ctx context.Context, in *empty.Empty) (*curiepb.CurieBlock, error)
	ProposeCurieBlockForOG(ctx context.Context, in *curiepb.SignedCurieBlockForOG) (*curiepb.ProposeResponse, error)
	ProposeCurieBlockForNG(ctx context.Context, in *curiepb.SignedCurieBlockForNG) (*curiepb.ProposeResponse, error)
	SendProposerPublicKey(ctx context.Context, in *curiepb.ProposerPublicKeyRequest) (*curiepb.ProposeResponse, error)
}
