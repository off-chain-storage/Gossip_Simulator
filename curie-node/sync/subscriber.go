package sync

import (
	"context"
	"flag-example/curie-node/p2p"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"
)

const pubsubMessageTimeout = 30 * time.Second

// type wrappedVal func(context.Context, peer.ID, *pubsub.Message) (pubsub.ValidationResult, error)

type subHandler func(context.Context, proto.Message) error

func (s *Service) registerSubscribers() {
	s.subscribe(
		p2p.OriginalTopicFormat,
		s.originalCurieBlockSubscriber,
	)
	s.subscribe(
		p2p.NewApproachTopicFormat,
		s.newCurieBlockSubscriber,
	)
}

func (s *Service) subscribe(topic string, handle subHandler) *pubsub.Subscription {
	base := p2p.GossipTopicMappings(topic)
	if base == nil {
		// Impossible condition as it would mean topic does not exist.
		panic(fmt.Sprintf("%s is not mapped to any message in GossipTopicMappings", topic))
	}

	return s.subscribeWithBase(topic, handle)
}

func (s *Service) subscribeWithBase(topic string, handle subHandler) *pubsub.Subscription {
	log := log.WithField("topic", topic)

	// Do not resubscribe already seen subscriptions.
	ok := s.subHandler.topicExists(topic)
	if ok {
		log.Debugf("Provided topic already has an active subscription running: %s", topic)
		return nil
	}

	sub, err := s.cfg.p2p.SubscribeToTopic(topic)
	if err != nil {
		log.WithError(err).Error("Could not subscribe topic")
		return nil
	}
	s.subHandler.addTopic(sub.Topic(), sub)

	pipeline := func(msg *pubsub.Message) {
		ctx, cancel := context.WithTimeout(s.ctx, pubsubMessageTimeout)
		defer cancel()

		var message proto.Message
		if err := proto.Unmarshal(msg.Data, message); err != nil {
			log.WithError(err).Error("Failed to unmarshal pubsub message")
			return
		}

		if err := handle(ctx, message); err != nil {
			log.WithError(err).Error("Could not handle message")
			return
		}
	}

	messageLoop := func() {
		for {
			// Subscriber 쪽에서 메세지를 수신하더라도 여기 이후로 넘어가지 않음
			msg, err := sub.Next(s.ctx)
			log.Info("msg is received")
			if err != nil {
				// context or subscription is cancelled.
				if err != pubsub.ErrSubscriptionCancelled {
					log.WithError(err).Warn("Subscription next failed")
				}
				sub.Cancel()
				return
			}

			if msg.ReceivedFrom == s.cfg.p2p.PeerID() {
				continue
			}

			go pipeline(msg)
		}
	}

	go messageLoop()
	log.WithField("topic", topic).Info("Subscribed to topic")
	return sub
}
