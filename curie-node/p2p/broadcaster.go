package p2p

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// ErrMessageNotMapped는 Msg가 GossipTypeMapping에 정의되지 않은 경우 에러 처리
var ErrMessageNotMapped = errors.New("message type is not mapped to a PubSub topic")

// Broadcast() - msg를 매핑하여 지정된 Topic에 브로드캐스트 될 수 있도록 함
func (s *Service) Broadcast(ctx context.Context, msg proto.Message) error {
	// GossipTypeMapping은 Gossip 메시지 유형을 PubSub 토픽에 매핑한다.
	topic, ok := GossipTypeMapping[reflect.TypeOf(msg)]
	if !ok {
		return ErrMessageNotMapped
	}

	// 주어진 topic에 대해 주어진 object를 broadcast한다.
	return s.broadcastObject(ctx, msg, topic)
}

// broadcastObject는 주어진 topic에 대해 주어진 object를 broadcast한다.
func (s *Service) broadcastObject(ctx context.Context, obj proto.Message, topic string) error {
	objPb, err := proto.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "could not marshal object")
	}

	// 주어진 topic에 대해 주어진 obj([]byte)를 Publish
	if err := s.PublishToTopic(ctx, topic, objPb); err != nil {
		err := errors.Wrap(err, "could not publish message")
		return err
	}

	return nil
}
