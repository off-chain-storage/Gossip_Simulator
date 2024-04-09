package peers

import (
	"context"
	"flag-example/curie-node/p2p/peers/peerdata"
	"flag-example/curie-node/p2p/peers/scores"

	"github.com/libp2p/go-libp2p/core/peer"
)

// const (
// 	// 피어에 연결되지 않았음을 의미
// 	// PeerDisconnected peerdata.PeerConnectionState = iota
// 	// 피어와 연결을 끊으려고 시도 중
// 	PeerDisconnecting
// 	// 피어에 연결되어 있음을 의미
// 	PeerConnected
// 	// 피어와 연결을 계속 시도 중
// 	PeerConnecting
// )

const (
	// ColocationLimit restricts how many peer identities we can see from a single ip or ipv6 subnet.
	ColocationLimit = 5

	// Additional buffer beyond current peer limit, from which we can store the relevant peer statuses.
	maxLimitBuffer = 150

	// InboundRatio is the proportion of our connected peer limit at which we will allow inbound peers.
	InboundRatio = float64(0.8)

	// MinBackOffDuration minimum amount (in milliseconds) to wait before peer is re-dialed.
	// When node and peer are dialing each other simultaneously connection may fail. In order, to break
	// of constant dialing, peer is assigned some backoff period, and only dialed again once that backoff is up.
	MinBackOffDuration = 100
	// MaxBackOffDuration maximum amount (in milliseconds) to wait before peer is re-dialed.
	MaxBackOffDuration = 5000
)

// 피어 상태 정보
type Status struct {
	ctx     context.Context
	store   *peerdata.Store
	scorers *scores.Service
	// 필요없는 데이터
	// ipTracker map[string]uint64
	// rand      *rand.Rand
}

// StatusConfig는 피어 상태 서비스 매개변수
type StatusConfig struct {
	PeerLimit    int
	ScoresParams *scores.Config
}

func NewStatus(ctx context.Context, config *StatusConfig) *Status {
	store := peerdata.NewStore(ctx, &peerdata.StoreConfig{
		MaxPeers: maxLimitBuffer + config.PeerLimit,
	})
	return &Status{
		ctx:     ctx,
		store:   store,
		scorers: scores.NewService(ctx, store, config.ScoresParams),
	}
}

func (p *Status) Scorers() *scores.Service {
	return p.scorers
}

func (p *Status) IsBad(pid peer.ID) bool {
	p.store.Lock()
	defer p.store.Unlock()
	return p.isBad(pid)
}

// isBad is the lock-free version of IsBad.
func (p *Status) isBad(pid peer.ID) bool {
	// Do not disconnect from trusted peers.
	if p.store.IsTrustedPeer(pid) {
		return false
	}
	// return p.isfromBadIP(pid) || p.scorers.IsBadPeerNoLock(pid)
	return p.scorers.IsBadPeerNoLock(pid)
}
