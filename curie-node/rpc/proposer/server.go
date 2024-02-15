package proposer

import (
	"context"
	"flag-example/blocks"
	ecdsacurie "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"
	"flag-example/curie-node/db"
	"flag-example/curie-node/monitor"
	"flag-example/curie-node/p2p"
	curiepb "flag-example/proto"

	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Ctx     context.Context
	DB      db.AccessRedisDB
	P2P     p2p.Broadcaster
	Monitor monitor.Monitor
}

func (ps *Server) SendProposerPublicKey(ctx context.Context, req *curiepb.ProposerPublicKeyRequest) (*curiepb.ProposeResponse, error) {
	log.WithField("Proposer's PubKey", req.PublicKey).Info("Received Proposer's PubKey from Proposer Node")

	// Save Proposer's PubKey to Redis DB
	if err := ps.DB.SetDataToRedis("Proposer", req.PublicKey); err != nil {
		log.WithError(err).Error("Failed to save proposer's pubKey to redis db")

		return &curiepb.ProposeResponse{
			Message: "Failed received proposer's pubKey",
		}, err
	}

	// Store Proposer's PubKey
	ecdsaPubKey, err := ecdsacurie.ConvertToEcdsaPubKeyString(req.PublicKey)
	if err != nil {
		log.WithError(err).Error("Failed to convert *ecdsa.PublicKey from string")
	}

	// Singleton Pattern for storing pubKey
	ecdsad.PublicKeyFromProposer(ecdsaPubKey)

	// Return Propose Response to Proposer Node
	return &curiepb.ProposeResponse{
		Message: "Successfully received proposer's pubkey",
	}, nil
}

func (ps *Server) GetBlock(ctx context.Context, empty *empty.Empty) (*curiepb.CurieBlock, error) {
	log.Info("Received GetBlock Request from Proposer Node")

	// No Meaning - Just Dummy Data for Experiment
	var blockData []byte = []byte("Dummy Data")

	// Return Block to Proposer Node
	return &curiepb.CurieBlock{
		DummyData: blockData,
	}, nil
}

func (ps *Server) ProposeCurieBlockForOG(ctx context.Context, req *curiepb.SignedCurieBlockForOG) (*curiepb.ProposeResponse, error) {
	log.Info("Received Original Gossip Request from Proposer Node")

	if err := ps.Monitor.SendUDPMessage("Start Propagation"); err != nil {
		return nil, err
	}

	blk, err := blocks.NewSignedBlock(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not decode block: %v", err)
	}

	blkPb, err := blk.Proto()
	if err != nil {
		return nil, errors.Wrap(err, "could not get protobuf block")
	}

	if err := ps.P2P.Broadcast(ctx, blkPb); err != nil {
		return nil, fmt.Errorf("could not broadcast block: %v", err)
	}

	log.Info("Broadcasting Original Gossip Data")

	return &curiepb.ProposeResponse{
		Message: "Successfully broadcasted original gossip data",
	}, nil
}

func (ps *Server) ProposeCurieBlockForNG(ctx context.Context, req *curiepb.SignedCurieBlockForNG) (*curiepb.ProposeResponse, error) {
	log.Info("Received New Gossip Request from Proposer Node")

	blk, err := blocks.NewSignedBlock(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Could not decode block: %v", err)
	}

	blkPb, err := blk.Proto()
	if err != nil {
		return nil, errors.Wrap(err, "could not get protobuf block")
	}

	if err := ps.P2P.Broadcast(ctx, blkPb); err != nil {
		return nil, fmt.Errorf("could not broadcast block: %v", err)
	}

	log.Info("Broadcasting New Gossip Data")

	return &curiepb.ProposeResponse{
		Message: "Successfully broadcasted new gossip data",
	}, nil

}
