package params

var curieNetworkConfig = &NetworkConfig{
	BootstrapNodes: []string{
		// 여기다가 부트 스트랩 노드 값 추가하기
	},
}

var curieNodeConfig = &CurieNodeConfig{
	GossipMaxSize: 10 * 1 << 20, // 10 Mib
	MaxChunkSize:  10 * 1 << 20, // 10 Mib
}
