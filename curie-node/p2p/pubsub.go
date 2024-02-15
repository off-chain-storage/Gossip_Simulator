package p2p

import (
	"context"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
)

// const (
// 	// overlay parameters
// 	gossipSubD   = 8  // topic stable mesh target count
// 	gossipSubDlo = 6  // topic stable mesh low watermark
// 	gossipSubDhi = 12 // topic stable mesh high watermark

// 	// gossip parameters
// 	gossipSubMcacheLen    = 6   // number of windows to retain full messages in cache for `IWANT` responses
// 	gossipSubMcacheGossip = 3   // number of windows to gossip about
// 	gossipSubSeenTTL      = 550 // number of heartbeat intervals to retain message IDs

// 	// heartbeat interval
// 	gossipSubHeartbeatInterval = 700 * time.Millisecond // frequency of heartbeat, milliseconds
// )

// JoinTopic will join PubSub topic, if not already joined.
func (s *Service) JoinTopic(topic string, opts ...pubsub.TopicOpt) (*pubsub.Topic, error) {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	// 이미 토픽에 Join 되어 있는지 확인
	if _, ok := s.joinedTopics[topic]; !ok {
		topicHandle, err := s.pubsub.Join(topic, opts...)
		if err != nil {
			return nil, err
		}
		s.joinedTopics[topic] = topicHandle
	}

	return s.joinedTopics[topic], nil
}

func (s *Service) LeaveTopic(topic string) error {
	s.joinedTopicsLock.Lock()
	defer s.joinedTopicsLock.Unlock()

	if t, ok := s.joinedTopics[topic]; ok {
		if err := t.Close(); err != nil {
			return err
		}
		delete(s.joinedTopics, topic)
	}

	return nil
}

func (s *Service) PublishToTopic(ctx context.Context, topic string, data []byte, opts ...pubsub.PubOpt) error {
	// Topic Join by msg type
	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return err
	}

	for {
		// 토픽에 피어가 있어야만 Publish, 최소 동기화 피어 수에 대한 것은 일단 주석 처리
		if len(topicHandle.ListPeers()) > 0 /* || flags.Get().MinimumSyncPeers == 0 */ {
			log.WithField("topic", topic).Debug("publishing message to topic")
			return topicHandle.Publish(ctx, data, opts...)
		}
		select {
		// ctx.Done()이 호출되면 에러를 반환
		case <-ctx.Done():
			return errors.Wrapf(ctx.Err(), "unable to find requisite number of peers for topic %s, 0 peers found to publish to", topic)
		// 100ms 동안 대기
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// SubscribeToTopic joins (if necessary) and subscribes to PubSub topic.
func (s *Service) SubscribeToTopic(topic string, opts ...pubsub.SubOpt) (*pubsub.Subscription, error) {
	topicHandle, err := s.JoinTopic(topic)
	if err != nil {
		return nil, err
	}

	log.WithField("topic", topic).Debug("subscribing to topic")

	return topicHandle.Subscribe(opts...)
}

func (s *Service) pubsubOptions() []pubsub.Option {
	psOpts := []pubsub.Option{
		// 메시지 서명 첨부하지 않고, 서명이 있는 메세지 거부
		// pubsub.WithMessageSignaturePolicy(pubsub.StrictNoSign),

		// 메세지 Author 정보 제거 - 익명성 향상
		pubsub.WithNoAuthor(),

		// 메세지 ID 함수 설정
		// pubsub.WithMessageIdFn(func(pmsg *pubsubpb.Message) string {
		// 	return MsgID(s.genesisValidatorsRoot, pmsg)
		// }),

		// 토픽 구독과 관련하여 필터링 옵션 - 일단 제외
		// pubsub.WithSubscriptionFilter(s),

		// 피어 아웃바운드 큐 크기 설정 - 일단 제외
		// pubsub.WithPeerOutboundQueueSize(int(s.cfg.QueueSize)),

		// 최대 메세지 크기 설정 - 일단 제외
		pubsub.WithMaxMessageSize(10 * 1 << 20),

		// 메세지 검증 큐 크기 설정 - 일단 제외
		// pubsub.WithValidateQueueSize(int(s.cfg.QueueSize)),

		// 피어 점수 파라미터 설정 - 일단 제외
		// pubsub.WithPeerScore(peerScoringParams()),

		// 피어 점수 검사 - 일단 제외
		// pubsub.WithPeerScoreInspect(s.peerInspector, time.Minute),

		// GossipSub 파라미터 설정
		// pubsub.WithGossipSubParams(pubsubGossipParam()),

		// Raw Tracer 설정 - 일단 제외
		// pubsub.WithRawTracer(gossipTracer{host: s.host}),
	}
	return psOpts
}

// creates a custom gossipsub parameter set.
// 커스텀 gossipsub 파라미터 셋 생성
// func pubsubGossipParam() pubsub.GossipSubParams {
// 	gParams := pubsub.DefaultGossipSubParams()
// 	gParams.Dlo = gossipSubDlo
// 	gParams.D = gossipSubD
// 	gParams.HeartbeatInterval = gossipSubHeartbeatInterval
// 	gParams.HistoryLength = gossipSubMcacheLen
// 	gParams.HistoryGossip = gossipSubMcacheGossip
// 	return gParams
// }

// // 이게 뭔지 모르겠다
// func setPubSubParameters() {
// 	pubsub.TimeCacheDuration = 550 * gossipSubHeartbeatInterval
// }
