package sync

import (
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type subTopicHandler struct {
	sync.RWMutex
	subTopics map[string]*pubsub.Subscription
	digestMap map[[4]byte]int
}

func newSubTopicHandler() *subTopicHandler {
	return &subTopicHandler{
		subTopics: map[string]*pubsub.Subscription{},
		digestMap: map[[4]byte]int{},
	}
}

func (s *subTopicHandler) addTopic(topic string, sub *pubsub.Subscription) {
	s.Lock()
	defer s.Unlock()
	s.subTopics[topic] = sub
}

func (s *subTopicHandler) topicExists(topic string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.subTopics[topic]
	return ok
}
