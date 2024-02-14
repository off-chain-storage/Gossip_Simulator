package p2p

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path"

	ecdsacurie "flag-example/crypto/ecdsa"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/pkg/errors"
)

const keyPath = "network-keys"

// Create Private Key for P2P Networking,
// If key isn't found, it'll be generates a new one
func privKey(cfg *Config) (*ecdsa.PrivateKey, error) {
	defaultKeyPath := path.Join(cfg.DataDir, keyPath)
	privateKeyPath := cfg.PrivateKey

	// 만약 Cli Flag에 Private Key Path가 있다면 이게 우선 순위가 가장 높다
	if privateKeyPath != "" {
		return privKeyFromFile(cfg.PrivateKey)
	}

	_, err := os.Stat(defaultKeyPath)
	defaultKeysExist := !os.IsNotExist(err)
	if err != nil && defaultKeysExist {
		return nil, err
	}

	// Default Key가 무엇인지는 모르겠다만
	// Cli Flag Priv Key Path 다음 우선순위를 가진다.
	if defaultKeysExist {
		return privKeyFromFile(defaultKeyPath)
	}

	// 여기까지 왔는데도 키가 없으면, 새롭게 만들자
	priv, _, err := crypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	return ecdsacurie.ConvertFromInterfacePrivKey(priv)
}

// Retrieves a private key from a file path in cli flag for p2p networking
func privKeyFromFile(path string) (*ecdsa.PrivateKey, error) {
	// 원본을 Path로부터 []byte 형태로 읽어오기
	src, err := os.ReadFile(path)
	if err != nil {
		log.WithError(err).Error("Error reading private key from file")
		return nil, err
	}

	// 원본 데이터를 디코딩 했을 때의 길이를 반환받고 그 길이만큼의 []byte를 생성
	// -> src hex data를 Decode 하기 위한 충분한 크기를 가지는 dst를 생성
	dst := make([]byte, hex.DecodedLen(len(src)))

	// 디코딩
	_, err = hex.Decode(dst, src)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hex string")
	}

	// byte 데이터에서 개인 키를 추출하여 반환
	unmarshalledKey, err := crypto.UnmarshalSecp256k1PrivateKey(dst)
	if err != nil {
		return nil, err
	}

	return ecdsacurie.ConvertFromInterfacePrivKey(unmarshalledKey)
}

// func extractIPFromAddr(ma multiaddr.Multiaddr) net.IP {
// 	ipComponent, err := ma.ValueForProtocol(multiaddr.P_IP4)
// 	if err != nil {
// 		return nil
// 	}

// 	ip := net.ParseIP(ipComponent)

// 	return ip
// }
