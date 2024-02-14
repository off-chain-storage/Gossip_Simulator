package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"net"

	"flag-example/config/features"
	ecdsacurie "flag-example/crypto/ecdsa"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

// IP와 PORT 받아서 Multiaddr 형태로 변환해주는 함수
func MultiAddressBuilder(ipAddr string, port uint) (ma.Multiaddr, error) {
	parsedIP := net.ParseIP(ipAddr)
	// IPv4 or 16바이트 주소가 아니면 에러 반환
	if parsedIP.To4() == nil && parsedIP.To16() == nil {
		return nil, errors.Errorf("invalid ip address provided")
	}
	if parsedIP.To4() != nil {
		return ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ipAddr, port))
	}
	return ma.NewMultiaddr(fmt.Sprintf("/ip6/%s/tcp/%d", ipAddr, port))
}

// libp2p host 빌드를 위한 옵션
func (s *Service) buildOptions(ip net.IP, priKey *ecdsa.PrivateKey) []libp2p.Option {
	cfg := s.cfg
	listen, err := MultiAddressBuilder("0.0.0.0", cfg.TCPPort)
	// listen, err := MultiAddressBuilder(ip.String(), cfg.TCPPort)
	if err != nil {
		log.WithError(err).Fatal("Failed to p2p listen")
	}
	// 호스트의 Local IP가 비어있지 않을 때
	// if cfg.LocalIP != "" {
	// 	// IP Addr이 유효하지 않을 경우 nil 반환
	// 	if net.ParseIP(cfg.LocalIP) == nil {
	// 		log.Fatalf("Invalid local ip provided: %s", cfg.LocalIP)
	// 	}
	// 	// Local IP로 Multi Addr 빌드하기
	// 	listen, err = MultiAddressBuilder(cfg.LocalIP, cfg.TCPPort)
	// 	if err != nil {
	// 		log.WithError(err).Fatal("Failed to p2p listen")
	// 	}
	// }

	// ecdsa private key를 libp2p crypto private key 형태로 변환
	ifaceKey, err := ecdsacurie.ConvertToInterfacePrivkey(priKey)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve private key")
	}
	// Public Key에 해당하는 Peer ID 반환
	id, err := peer.IDFromPublicKey(ifaceKey.GetPublic())
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve peer id")
	}
	log.Infof("Running node with peer id of %s ", id.String())

	// libp2p 옵션 제정
	options := []libp2p.Option{
		// libp2p private key
		privKeyOption(priKey),
		// 해당 MultiAddr로 listen 하도록 설정
		libp2p.ListenAddrs(listen),
		// libp2p 통신에 사용될 프로토콜 지정 - TCP
		libp2p.Transport(tcp.NewTCPTransport),
		// 기본 멀티 플렉서 사용 - 단일 연결을 통해 여러 데이터 스트림을 동시에 전송할 수 있게 하는 기술
		libp2p.DefaultMuxers,
		// // 특정 멀티 플렉서 사용 - 지원 X
		// libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
		// libp2p.UserAgent(version.BuildData()),
		// libp2p.ConnectionGater(s),
	}

	// 옵션 배열에 libp2p 보안 메커니즘 추가
	// noise.ID -> Noise Protocol Identifier
	// noise.New -> Noise Protocol 초기화하고 구현, libp2p 연결 시 Noise 프로토콜을 사용하여 연결을 암호화하고 인증
	// options = append(options, libp2p.Security(noise.ID, noise.New))

	options = append(options, libp2p.DisableRelay())

	// if cfg.HostAddress != "" {
	// 	options = append(options, libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
	// 		external, err := MultiAddressBuilder(cfg.HostAddress, cfg.TCPPort)
	// 		if err != nil {
	// 			log.WithError(err).Error("Unable to create external multiaddress")
	// 		} else {
	// 			addrs = append(addrs, external)
	// 		}
	// 		return addrs
	// 	}))
	// }

	// ping service 금지
	options = append(options, libp2p.Ping(false))
	if features.Get().DisableResourceManager {
		// libp2p Resource Manager을 NULL로 설정 -> Resource 무제한으로 사용하게끔 옵션 설정
		options = append(options, libp2p.ResourceManager(&network.NullResourceManager{}))
	}

	return options
}

func multiAddressBuilderWithID(ipAddr, protocol string, port uint, id peer.ID) (ma.Multiaddr, error) {
	parsedIP := net.ParseIP(ipAddr)
	if parsedIP.To4() == nil && parsedIP.To16() == nil {
		return nil, errors.Errorf("invalid ip address provided: %s", ipAddr)
	}
	if id.String() == "" {
		return nil, errors.New("empty peer id given")
	}

	return ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/%s/%d/p2p/%s", ipAddr, protocol, port, id.String()))
}

func privKeyOption(privkey *ecdsa.PrivateKey) libp2p.Option {
	return func(cfg *libp2p.Config) error {
		// ecdsa private key를 libp2p crypto private key 형태로 변환
		ifaceKey, err := ecdsacurie.ConvertToInterfacePrivkey(privkey)
		if err != nil {
			return err
		}
		log.Debug("ECDSA private key generated")
		// libp2p option key 반환
		return cfg.Apply(libp2p.Identity(ifaceKey))
	}
}
