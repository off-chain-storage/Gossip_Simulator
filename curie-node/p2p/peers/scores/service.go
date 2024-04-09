package scores

import (
	"context"
	"flag-example/curie-node/p2p/peers/peerdata"
	"math"

	"github.com/libp2p/go-libp2p/core/peer"
)

const BadPeerScore = -100.0

const ScoreRoundingFactor = 10000

type Scorer interface {
	Score(pid peer.ID) float64
	IsBadPeer(pid peer.ID) bool
	BadPeers() []peer.ID
}

type Service struct {
	store   *peerdata.Store
	scorers struct {
		badResponsesScorer *BadResponsesScorer
		// peerStatusScorer   *PeerStatusScorer
	}
	weights     map[Scorer]float64
	totalWeight float64
}

type Config struct {
	BadResponsesScorerConfig *BadResponsesScorerConfig
}

func NewService(ctx context.Context, store *peerdata.Store, config *Config) *Service {
	s := &Service{
		store:   store,
		weights: make(map[Scorer]float64),
	}
	s.scorers.badResponsesScorer = newBadResponsesScorer(store, config.BadResponsesScorerConfig)
	s.setScorerWeight(s.scorers.badResponsesScorer, 0.3)

	return s
}

func (s *Service) BadResponsesScorer() *BadResponsesScorer {
	return s.scorers.badResponsesScorer
}

func (s *Service) Score(pid peer.ID) float64 {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.ScoreNoLock(pid)
}

func (s *Service) ScoreNoLock(pid peer.ID) float64 {
	score := float64(0)
	if _, ok := s.store.PeerData(pid); !ok {
		return 0
	}
	score += s.scorers.badResponsesScorer.scoreNoLock(pid) * s.scorerWeight(s.scorers.badResponsesScorer)
	return math.Round(score*ScoreRoundingFactor) / ScoreRoundingFactor
}

func (s *Service) IsBadPeer(pid peer.ID) bool {
	s.store.RLock()
	defer s.store.RUnlock()
	return s.IsBadPeerNoLock(pid)
}

func (s *Service) IsBadPeerNoLock(pid peer.ID) bool {
	if s.scorers.badResponsesScorer.isBadPeerNoLock(pid) {
		return true
	}
	// if s.scorers.peerStatusScorer.isBadPeerNoLock(pid) {
	// 	return true
	// }
	// if features.Get().EnablePeerScorer {
	// 	if s.scorers.gossipScorer.isBadPeerNoLock(pid) {
	// 		return true
	// 	}
	// }
	return false
}

// setScorerWeight adds scorer to map of known scorers.
func (s *Service) setScorerWeight(scorer Scorer, weight float64) {
	s.weights[scorer] = weight
	s.totalWeight += s.weights[scorer]
}

// scorerWeight calculates contribution percentage of a given scorer in total score.
func (s *Service) scorerWeight(scorer Scorer) float64 {
	return s.weights[scorer] / s.totalWeight
}
