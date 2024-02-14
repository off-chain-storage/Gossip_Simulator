package peerdata

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type StoreConfig struct {
	MaxPeers int
}

type Store struct {
	sync.RWMutex
	ctx    context.Context
	config *StoreConfig
	peers  map[peer.ID]*PeerData
	// 여기서 부터 필요없는 데이터
	// trustedPeers map[peer.ID]bool
}

// PeerData aggregates protocol and application level info about a single peer.
type PeerData struct {
	// Network related data.
	Address   ma.Multiaddr
	Direction network.Direction
	// ConnState PeerConnectionState

	// 여기서 부터는 관련 없는 데이터
	// Network related data. - 일단 제외
	// Enr           *enr.Record
	// NextValidTime time.Time
	// Chain related data. - 일단 제외
	// MetaData                  metadata.Metadata
	// ChainState                *ethpb.Status
	// ChainStateLastUpdated     time.Time
	// ChainStateValidationError error
	// Scorers internal data. - 일단 제외
	// BadResponses         int
	// ProcessedBlocks      uint64
	// BlockProviderUpdated time.Time
	// Gossip Scoring data. - 일단 제외
	// TopicScores      map[string]*ethpb.TopicScoreSnapshot
	// GossipScore      float64
	// BehaviourPenalty float64
}

// NewStore creates new peer data store.
func NewStore(ctx context.Context, config *StoreConfig) *Store {
	return &Store{
		ctx:    ctx,
		config: config,
		peers:  make(map[peer.ID]*PeerData),
		// trustedPeers: make(map[peer.ID]bool),
	}
}
