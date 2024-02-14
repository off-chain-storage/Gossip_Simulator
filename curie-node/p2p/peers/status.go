package peers

import (
	"context"
	"flag-example/curie-node/p2p/peers/peerdata"
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
	ctx   context.Context
	store *peerdata.Store

	// 필요없는 데이터
	// scorers   *scorers.Service
	// ipTracker map[string]uint64
	// rand      *rand.Rand
}

// StatusConfig는 피어 상태 서비스 매개변수
type StatusConfig struct {
	// 노드에 연결 가능한 최대 동시 피어 수 지정
	PeerLimit int
	// 피어 평판 관련 매개 변수 - 일단 제외
	// ScorerParams *scorers.Config
}

func NewStatus(ctx context.Context, config *StatusConfig) *Status {
	store := peerdata.NewStore(ctx, &peerdata.StoreConfig{
		MaxPeers: maxLimitBuffer + config.PeerLimit,
	})
	return &Status{
		ctx:   ctx,
		store: store,

		// 필요 없는 데이터
		// scorers:   scorers.NewService(ctx, store, config.ScorerParams),
		// ipTracker: map[string]uint64{},
		// rand:  rand.NewDeterministicGenerator(),
	}
}
