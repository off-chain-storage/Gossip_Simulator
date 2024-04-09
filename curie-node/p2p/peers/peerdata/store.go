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
	ctx          context.Context
	config       *StoreConfig
	peers        map[peer.ID]*PeerData
	trustedPeers map[peer.ID]bool
}

// PeerData aggregates protocol and application level info about a single peer.
type PeerData struct {
	// Network related data.
	Address   ma.Multiaddr
	Direction network.Direction
	// Peer Score
	BadResponses int
}

// NewStore creates new peer data store.
func NewStore(ctx context.Context, config *StoreConfig) *Store {
	return &Store{
		ctx:          ctx,
		config:       config,
		peers:        make(map[peer.ID]*PeerData),
		trustedPeers: make(map[peer.ID]bool),
	}
}

func (s *Store) IsTrustedPeer(p peer.ID) bool {
	return s.trustedPeers[p]
}

func (s *Store) PeerData(pid peer.ID) (*PeerData, bool) {
	peerData, ok := s.peers[pid]
	return peerData, ok
}

func (s *Store) SetPeerData(pid peer.ID, data *PeerData) {
	s.peers[pid] = data
}

func (s *Store) Peers() map[peer.ID]*PeerData {
	return s.peers
}
