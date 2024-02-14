package p2p

import (
	"crypto/ecdsa"

	ecdsacurie "flag-example/crypto/ecdsa"
	"flag-example/curie-node/p2p/cnode"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

func (s *Service) startDHT() error {
	var options []dht.Option

	// NoDiscovery가 True일 경우 부트스트랩 모드
	if s.cfg.NoDiscovery {
		options = append(options, dht.Mode(dht.ModeAutoServer))
	}

	kdht, err := dht.New(s.ctx, s.host, options...)
	if err != nil {
		return err
	}

	// NoDiscovery cmd가 false이면 진입 -> 즉, 부트스트랩과 연결될거라는 뜻
	if !s.cfg.NoDiscovery {
		// BootStrap Node와의 연결
		err = s.connectToBootnodes()
		if err != nil {
			log.WithError(err).Error("Could not add bootnode to the exclusion list")
			s.startupErr = err
			return err
		}

		// Connection with New Peer - Go Routine
		go s.listenForNewNodes(OriginalTopicFormat)
		go s.listenForNewNodes(NewApproachTopicFormat)
	}

	if err = kdht.Bootstrap(s.ctx); err != nil {
		return err
	}

	s.dht = kdht

	return nil
}

func (s *Service) listenForNewNodes(topic string) {
	var routingDiscovery = drouting.NewRoutingDiscovery(s.dht)

	dutil.Advertise(s.ctx, routingDiscovery, topic)

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			peers, err := routingDiscovery.FindPeers(s.ctx, topic)
			if err != nil {
				log.WithError(err).Tracef("Could not find peers for topic %s", topic)
				continue
			}

			for p := range peers {
				go func(peer *peer.AddrInfo) {
					if err := s.connectWithPeer(s.ctx, *peer); err != nil {
						log.WithError(err).Tracef("Could not connect with peer %s for topic %s", peer.String(), topic)
					}
				}(&p)
			}
		}
	}
}

func (s *Service) convertToMultiAddr(nodes []*cnode.Node) []ma.Multiaddr {
	var multiAddrs []ma.Multiaddr
	for _, node := range nodes {
		if node.IP() == nil {
			continue
		}
		multiAddr, err := s.convertToSingleMultiAddr(node)
		if err != nil {
			log.WithError(err).Error("could not convert to multiAddr")
			continue
		}
		multiAddrs = append(multiAddrs, multiAddr)
	}
	return multiAddrs
}

func (s *Service) getPubKeyFromPrivKey() (string, error) {
	// pubkey := s.privKey.PublicKey
	pubkey := s.privKey.Public()
	ecdsaPubKey, ok := pubkey.(*ecdsa.PublicKey)
	if !ok { // 실패할 경우(Secp256k1 타입이 아닐 경우) 오류 메세지 반환
		return "", errors.New("could not cast to ecdsaPublicKey")
	}

	stringPubKey := ecdsacurie.ConvertToStringEcdsaPubKey(ecdsaPubKey)

	return stringPubKey, nil
}

func (s *Service) convertToSingleMultiAddr(node *cnode.Node) (ma.Multiaddr, error) {
	return multiAddressBuilderWithID(node.IP().String(), "tcp", uint(node.TCP()), node.PeerID)
}
