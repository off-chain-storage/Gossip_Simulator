package state

import "flag-example/async/event"

// 느낌은 State Update를 Consumer에게 제공하는 것으로 이해되고 있다.
// 어떤 State ??? <- 이걸 알아보자
type Notifier interface {
	StateFeed() *event.Feed
}
