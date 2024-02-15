package p2p

import (
	"context"
	"crypto/ecdsa"
	"sync"
	"time"

	ecdsacurie "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"
	"flag-example/curie-node/p2p/cnode"
	"flag-example/curie-node/p2p/peers"
	curienetwork "flag-example/network"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"go.opencensus.io/trace"
)

// 피어 제한에 도달한 경우 새 피어 검색을 중단하고 대신 아래에 정의된 기간 동안 현재 피어 제한 상태를 폴링합니다.
// var pollingPeriod = 6 * time.Second

// maxBadResponses는 통신을 중단하기 전에 피어로부터의 잘못된 응답의 최대 수입니다.
// const maxBadResponses = 5

// maxDialTimeout은 단일 피어 다이얼에 대한 시간 초과입니다.
// var maxDialTimeout = params.BeaconConfig().RespTimeoutDuration()
const maxDialTimeout = time.Duration(30) * time.Second

// Service for managing peer to peer (p2p) networking.
type Service struct {
	started          bool
	startupErr       error
	ctx              context.Context
	cfg              *Config
	cancel           context.CancelFunc
	privKey          *ecdsa.PrivateKey
	pubsub           *pubsub.PubSub
	peers            *peers.Status
	joinedTopics     map[string]*pubsub.Topic
	joinedTopicsLock sync.RWMutex
	host             host.Host
	dht              *dht.IpfsDHT

	// 단순 P2P에서 필요없는 옵션
	// pingMethod func(ctx context.Context, id peer.ID) error
	// dv5Listener           Listener
	// metaData              metadata.Metadata
	// isPreGenesis          bool
	// addrFilter            *multiaddr.Filters
	// ipLimiter             *leakybucket.Collector
	// subnetsLock           map[uint64]*sync.RWMutex
	// subnetsLockLock       sync.Mutex // Lock access to subnetsLock
	// initializationLock    sync.Mutex
	// genesisTime           time.Time
	// genesisValidatorsRoot []byte
	// activeValidatorCount  uint64
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	_ = cancel // govet fix for lost cancel. Cancel is handled in service.Stop().

	s := &Service{
		ctx:          ctx,
		cancel:       cancel,
		cfg:          cfg,
		joinedTopics: make(map[string]*pubsub.Topic, len(gossipTopicMappings)),
		// isPreGenesis: true,
		// subnetsLock:  make(map[uint64]*sync.RWMutex),
	}

	// cfg 검증
	s.cfg = validateConfig(s.cfg)

	// // 부트 스트랩 노드 cfg에서 파싱 후 Service Struct에 붙이기
	// dv5Nodes := parseBootStrapAddrs(s.cfg.BootstrapNodeAddr)
	// cfg.Discv5BootStrapAddr = dv5Nodes

	// 뭔가 내 네트워크 정보 가져오는 것 같은데 자세히는 모르겠다,,
	ipAddr := curienetwork.IPAddr()

	// private key 생성하기
	s.privKey, err = privKey(s.cfg)
	if err != nil {
		log.WithError(err).Error("Failed to generate p2p private key")
		return nil, err
	}

	// // Publisher 이면 PublicKey Redis에 등록하기
	// if s.cfg.IsPublisher {
	// 	ecdsaPubKey, err := s.getPubKeyFromPrivKey()
	// 	if err != nil {
	// 		log.WithError(err).Error("Failed to get p2p Public key")
	// 		return nil, err
	// 	}
	// 	if err := s.cfg.DB.SetDataToRedis("Publisher", ecdsaPubKey); err != nil {
	// 		log.Error(err)
	// 		return nil, err
	// 	}
	// }

	// libp2p options 설정
	opts := s.buildOptions(ipAddr, s.privKey)
	h, err := libp2p.New(opts...)
	if err != nil {
		log.WithError(err).Error("Failed to create p2p host")
		return nil, err
	}

	// libp2p host
	s.host = h

	// Set pubsub option
	// psOpts := s.pubsubOptions()

	// 뭔지 몰라,,,
	// setPubSubParameters()

	// Create GossipSub Instance
	gs, err := pubsub.NewGossipSub(s.ctx, s.host)
	if err != nil {
		log.WithError(err).Error("Failed to start pubsub")
		return nil, err
	}
	s.pubsub = gs

	// Peer를 위한 새로운 Status Entity 생성
	s.peers = peers.NewStatus(ctx, &peers.StatusConfig{
		PeerLimit: int(s.cfg.MaxPeers),
	})

	return s, nil
}

func (s *Service) Start() {
	// 피어가 이미 시작됐는지 확인
	if s.started {
		log.Error("Attempted to start p2p service when it was already started")
		return
	}

	// ** 여기다가 Proposer 노드의 Public Key 저장 또는 가져오는 함수 추가 하기 **
	if !s.cfg.IsPublisher {
		// Subscriber
		pubKey, err := s.cfg.DB.GetDataFromRedis("Proposer")
		if err != nil {
			log.WithError(err).Error("Failed to get Proposer's Public Key from DB")
		}

		log.Info("Proposer's PubKey from RedisDB is ", pubKey)

		ecdsaPubKey, err := ecdsacurie.ConvertToEcdsaPubKeyString(pubKey)
		if err != nil {
			log.WithError(err).Error("Failed to convert *ecdsa.Publickey from string")
		}

		// Singleton Pattern for storing pubKey
		ecdsad.PublicKeyFromProposer(ecdsaPubKey)
	}

	// Peer Discovery를 위한 DHT Init 함수 실행
	if err := s.startDHT(); err != nil {
		log.WithError(err).Fatal("Failed to start discovery")
		s.startupErr = err
		return
	}

	s.started = true
}

func (s *Service) Stop() error {
	defer s.cancel()
	s.started = false

	return nil
}

func (s *Service) connectWithAllPeers(multiAddrs []multiaddr.Multiaddr) {
	addrInfos, err := peer.AddrInfosFromP2pAddrs(multiAddrs...)
	if err != nil {
		log.WithError(err).Error("Could not convert to peer address info's from multiaddresses")
		return
	}

	var wg sync.WaitGroup
	for _, info := range addrInfos {
		wg.Add(1)
		go func(info peer.AddrInfo) {
			if err := s.connectWithPeer(s.ctx, info); err != nil {
				log.WithError(err).Tracef("Could not connect with peer %s", info.String())
			}
		}(info)
	}
	wg.Wait()
}

func (s *Service) connectWithPeer(ctx context.Context, info peer.AddrInfo) error {
	ctx, span := trace.StartSpan(ctx, "p2p.connectWithPeer")
	defer span.End()

	if info.ID == s.host.ID() {
		return nil
	}

	// 이미 연결된 ID라면 연결 시도 X
	if s.host.Network().Connectedness(info.ID) == network.Connected {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, maxDialTimeout)
	defer cancel()

	if err := s.host.Connect(ctx, info); err != nil {
		log.Warn(err)
		return err
	} else {
		log.Infof("Connection established with node: [%q] [%q]", info.Addrs, info.ID.String())
	}

	return nil
}

func (s *Service) connectToBootnodes() error {
	nodes := make([]*cnode.Node, 0, len(s.cfg.BootstrapNodeAddr))

	for _, peerAddr := range s.cfg.BootstrapNodeAddr {
		bootNode, err := cnode.Parse(peerAddr)
		if err != nil {
			return err
		}

		nodes = append(nodes, bootNode)
	}
	multiAddresses := s.convertToMultiAddr(nodes)

	s.connectWithAllPeers(multiAddresses)
	return nil
}

// Peers returns the peer status interface.
func (s *Service) Peers() *peers.Status {
	return s.peers
}

func (s *Service) PubSub() *pubsub.PubSub {
	return s.pubsub
}

func (s *Service) PeerID() peer.ID {
	return s.host.ID()
}

func (s *Service) PublisherPeer() bool {
	return s.cfg.IsPublisher
}
