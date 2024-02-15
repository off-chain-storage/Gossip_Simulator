package sync

import (
	"context"
	"flag-example/curie-node/p2p"
	"fmt"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const pubsubMessageTimeout = 30 * time.Second

type wrappedVal func(context.Context, peer.ID, *pubsub.Message) (pubsub.ValidationResult, error)

type subHandler func(context.Context, proto.Message) error

func (s *Service) msgValidator(_ context.Context, _ peer.ID, msg *pubsub.Message) (pubsub.ValidationResult, error) {
	m, err := s.decodePubsubMessage(msg)
	if err != nil {
		log.WithError(err).Error("Could not decode message")
		return pubsub.ValidationReject, nil
	}
	msg.ValidatorData = m

	return pubsub.ValidationAccept, nil
}

func (s *Service) registerSubscribers() {
	s.subscribe(
		p2p.OriginalTopicFormat,
		s.msgValidator,
		s.originalCurieBlockSubscriber,
	)
	s.subscribe(
		p2p.NewApproachTopicFormat,
		s.msgValidator,
		s.newCurieBlockSubscriber,
	)
}

func (s *Service) subscribe(topic string, validator wrappedVal, handle subHandler) *pubsub.Subscription {
	base := p2p.GossipTopicMappings(topic)
	if base == nil {
		panic(fmt.Sprintf("%s is not mapped to any message in GossipTopicMappings", topic))
	}

	return s.subscribeWithBase(topic, validator, handle)
}

func (s *Service) subscribeWithBase(topic string, validator wrappedVal, handle subHandler) *pubsub.Subscription {
	log := log.WithField("topic", topic)

	ok := s.subHandler.topicExists(topic)
	if ok {
		log.Debugf("Provided topic already has an active subscription running: %s", topic)
		return nil
	}

	if err := s.cfg.p2p.PubSub().RegisterTopicValidator(s.wrapAndReportValidation(topic, validator)); err != nil {
		log.WithError(err).Error("Could not register topic validator")
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

		if msg.ValidatorData == nil {
			log.Error("Received nil message on pubsub")
			return
		}

		if err := handle(ctx, msg.ValidatorData.(proto.Message)); err != nil {
			log.WithError(err).Error("Could not handle message")
			return
		}

		log.Info("Message handled")
	}

	messageLoop := func() {
		for {
			// Subscriber 쪽에서 메세지를 수신하더라도 여기 이후로 넘어가지 않음 - 해결
			log.Info("Subscriber is waiting for message...")
			msg, err := sub.Next(s.ctx)
			if err != nil {
				if err != pubsub.ErrSubscriptionCancelled {
					log.WithError(err).Warn("Subscription next failed")
				}
				sub.Cancel()
				return
			}

			log.Info("@@ This msg is from ", msg.ReceivedFrom.String(), "@@")

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

func (s *Service) wrapAndReportValidation(topic string, v wrappedVal) (string, pubsub.ValidatorEx) {
	return topic, func(ctx context.Context, pid peer.ID, msg *pubsub.Message) (res pubsub.ValidationResult) {
		b, err := v(ctx, pid, msg)
		if b == pubsub.ValidationReject {
			fields := logrus.Fields{
				"topic":   topic,
				"peer id": pid.String(),
			}
			log.WithError(err).WithFields(fields).Debugf("Gossip message was rejected")
		}
		if b == pubsub.ValidationIgnore {
			if err != nil {
				log.WithError(err).WithFields(logrus.Fields{
					"topic":   topic,
					"peer id": pid.String(),
				}).Debugf("Gossip message was ignored")
			}
		}

		return b
	}
}
