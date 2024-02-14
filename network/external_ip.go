package network

import (
	"net"
	"sort"

	"github.com/sirupsen/logrus"
)

// External IPV4 주소를 libp2p 형식으로 변환.
func IPAddr() net.IP {
	ip, err := ExternalIP()
	if err != nil {
		panic(err)
	}
	return net.ParseIP(ip)
}

// 사용 가능한 첫번째 IPv4/v6를 반환
func ExternalIP() (string, error) {
	ips, err := ipAddrs()
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "127.0.0.1", nil
	}
	return ips[0].String(), nil
}

// 사용 가능한 모든 유효한 IP를 반환
func ipAddrs() ([]net.IP, error) {
	// 시스템에 설치된 모든 네트워크 인터페이스 정보 가져옴
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// // 인터페이스가 활성화되어 있지 않으면 다음 인터페이스로 넘어감
		// if iface.Flags&net.FlagUp == 0 {
		// 	continue
		// }
		// // 루프백 인터페이스인 경우 다음 인터페이스로 넘어감
		// if iface.Flags&net.FlagLoopback != 0 {
		// 	continue
		// }
		// 인터페이스에 할당된 IP 주소 목록을 가져옴
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		// IP 주소 목록을 순회하면서 유효한 IP 주소만 추출
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				logrus.Info(v.IP)
			}
		}
	}

	var ipAddrs []net.IP
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue
			}
			ipAddrs = append(ipAddrs, ip)
		}
	}
	return SortAddresses(ipAddrs), nil
}

func SortAddresses(ipAddrs []net.IP) []net.IP {
	sort.Slice(ipAddrs, func(i, j int) bool {
		return ipAddrs[i].To4() != nil && ipAddrs[j].To4() == nil
	})
	return ipAddrs
}
