package sync

import (
	"flag-example/curie-node/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var errNilPubsubMessage = errors.New("nil pubsub message")

// var errInvalidTopic = errors.New("invalid topic format")

func (s *Service) decodePubsubMessage(msg *pubsub.Message) (proto.Message, error) {
	if msg == nil || msg.Topic == nil || *msg.Topic == "" {
		return nil, errNilPubsubMessage
	}

	topic := *msg.Topic

	base := p2p.GossipTopicMappings(topic)
	if base == nil {
		return nil, p2p.ErrMessageNotMapped
	}

	// m, ok := proto.Clone(base).(proto.Message)
	// if !ok {
	// 	return nil, errors.Errorf("message of %T does not support marshaller interface", base)
	// }

	return base, nil
}
