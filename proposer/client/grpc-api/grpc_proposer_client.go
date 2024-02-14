package grpcapi

import (
	"context"
	"flag-example/proposer/client/iface"
	curiepb "flag-example/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type grpcProposerClient struct {
	curieNodeProposerClient curiepb.CurieNodeProposerClient
}

func (c *grpcProposerClient) GetBlock(ctx context.Context, in *empty.Empty) (*curiepb.CurieBlock, error) {
	return c.curieNodeProposerClient.GetBlock(ctx, in)
}

func (c *grpcProposerClient) ProposeCurieBlockForOG(ctx context.Context, in *curiepb.SignedCurieBlockForOG) (*curiepb.ProposeResponse, error) {
	return c.curieNodeProposerClient.ProposeCurieBlockForOG(ctx, in)
}

func (c *grpcProposerClient) ProposeCurieBlockForNG(ctx context.Context, in *curiepb.SignedCurieBlockForNG) (*curiepb.ProposeResponse, error) {
	return c.curieNodeProposerClient.ProposeCurieBlockForNG(ctx, in)
}

func (c *grpcProposerClient) SendProposerPublicKey(ctx context.Context, in *curiepb.ProposerPublicKeyRequest) (*curiepb.ProposeResponse, error) {
	return c.curieNodeProposerClient.SendProposerPublicKey(ctx, in)
}

func NewGrpcProposerClient(cc grpc.ClientConnInterface) iface.ProposerClient {
	return &grpcProposerClient{curiepb.NewCurieNodeProposerClient(cc)}
}
