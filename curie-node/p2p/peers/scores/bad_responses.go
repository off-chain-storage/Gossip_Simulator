package scores

import (
	"flag-example/curie-node/p2p/peers/peerdata"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var _ Scorer = (*BadResponsesScorer)(nil)

const (
	DefaultBadResponsesThreshold = 6

	DefaultBadResponsesDecayInterval = time.Hour

	DefaultBadResponsesPenaltyFactor = 10
)

type BadResponsesScorer struct {
	config *BadResponsesScorerConfig
	store  *peerdata.Store
}

type BadResponsesScorerConfig struct {
	Threshold     int
	DecayInterval time.Duration
}

func newBadResponsesScorer(store *peerdata.Store, config *BadResponsesScorerConfig) *BadResponsesScorer {
	if config == nil {
		config = &BadResponsesScorerConfig{}
	}
	scorer := &BadResponsesScorer{
		config: config,
		store:  store,
	}
	if scorer.config.Threshold == 0 {
		scorer.config.Threshold = DefaultBadResponsesThreshold
	}
	if scorer.config.DecayInterval == 0 {
		scorer.config.DecayInterval = DefaultBadResponsesDecayInterval
	}
	return scorer
}

func (s *BadResponsesScorer) Score(pid peer.ID) float64 {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.scoreNoLock(pid)
}

func (s *BadResponsesScorer) scoreNoLock(pid peer.ID) float64 {
	if s.isBadPeerNoLock(pid) {
		return BadPeerScore
	}
	score := float64(0)
	peerData, ok := s.store.PeerData(pid)
	if !ok {
		return score
	}
	if peerData.BadResponses > 0 {
		score = float64(peerData.BadResponses) / float64(s.config.Threshold)
		// Since score represents a penalty, negate it and multiply
		// it by a factor.
		score *= -DefaultBadResponsesPenaltyFactor
	}
	return score
}

func (s *BadResponsesScorer) Increment(pid peer.ID) {
	s.store.Lock()
	defer s.store.Unlock()

	peerData, ok := s.store.PeerData(pid)
	if !ok {
		s.store.SetPeerData(pid, &peerdata.PeerData{
			BadResponses: 1,
		})
		return
	}
	peerData.BadResponses++
}

func (s *BadResponsesScorer) IsBadPeer(pid peer.ID) bool {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.isBadPeerNoLock(pid)
}

func (s *BadResponsesScorer) isBadPeerNoLock(pid peer.ID) bool {
	if peerData, ok := s.store.PeerData(pid); ok {
		return peerData.BadResponses >= s.config.Threshold
	}
	return false
}

// BadPeers returns the peers that are considered bad.
func (s *BadResponsesScorer) BadPeers() []peer.ID {
	s.store.RLock()
	defer s.store.RUnlock()

	badPeers := make([]peer.ID, 0)
	for pid := range s.store.Peers() {
		if s.isBadPeerNoLock(pid) {
			badPeers = append(badPeers, pid)
		}
	}
	return badPeers
}
