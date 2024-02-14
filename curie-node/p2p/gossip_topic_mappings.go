package p2p

import (
	curiepb "flag-example/proto"
	"reflect"

	"google.golang.org/protobuf/proto"
)

var gossipTopicMappings = map[string]proto.Message{
	OriginalTopicFormat:    &curiepb.SignedCurieBlockForOG{},
	NewApproachTopicFormat: &curiepb.SignedCurieBlockForNG{},
}

func GossipTopicMappings(topic string) proto.Message {
	return gossipTopicMappings[topic]
}

var GossipTypeMapping = make(map[reflect.Type]string, len(gossipTopicMappings))

func init() {
	for k, v := range gossipTopicMappings {
		GossipTypeMapping[reflect.TypeOf(v)] = k
	}

	GossipTypeMapping[reflect.TypeOf(&curiepb.SignedCurieBlockForOG{})] = OriginalTopicFormat
	GossipTypeMapping[reflect.TypeOf(&curiepb.SignedCurieBlockForNG{})] = NewApproachTopicFormat
}
