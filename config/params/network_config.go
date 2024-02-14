package params

import "github.com/mohae/deepcopy"

type NetworkConfig struct {
	BootstrapNodes []string
}

// Singleton Instance
var networkConfig = curieNetworkConfig

// Using Singleton Pattern
func CurieNetworkConfig() *NetworkConfig {
	return networkConfig
}

// 새로운 부트스트랩 노드 정보 받아와서 덮어쓰기
func OverrideCurieNetworkConfig(cfg *NetworkConfig) {
	networkConfig = cfg.Copy()
}

// 새로운 인스턴스로 덮어 쓰기 위한 Copy Function
func (c *NetworkConfig) Copy() *NetworkConfig {
	config, ok := deepcopy.Copy(*c).(NetworkConfig)
	if !ok {
		config = *networkConfig
	}
	return &config
}
