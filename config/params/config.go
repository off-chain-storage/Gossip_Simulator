package params

// CurieNode의 Config 관련 사항 모음
type CurieNodeConfig struct {
	// Compression 되지 않은 Gossip msg의 최대 크기
	GossipMaxSize uint64
	// Compression 되지 않은 Req/Resp의 최대 크기
	MaxChunkSize uint64
}
