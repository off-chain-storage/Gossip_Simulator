package p2p

// 여기는 아직까지 뭘 위한 코드인지 잘 분별이 안 된다.

const (
	// Message Format
	GossipOriginalMessage       = "original"
	GossipNewPropagationMessage = "new_propagation"

	// Topic Format
	OriginalTopicFormat    = GossipOriginalMessage
	NewApproachTopicFormat = GossipNewPropagationMessage
)
