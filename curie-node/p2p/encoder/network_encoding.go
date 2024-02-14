package encoder

import (
	"io"

	ssz "github.com/prysmaticlabs/fastssz"
)

// Curie-Node의 P2P 네트워킹시 인코딩/디코딩 등을 지원하는 인터페이스
type NetworkEncoding interface {
	// 수신 받은 Gossip Msg의 디코딩(역직렬화)을 지원
	DecodeGossip([]byte, ssz.Unmarshaler) error
	// DecodeWithMaxLength는 reader에서 varint prefix를 포함한 바이트를 읽어 디코딩(역직렬화)
	DecodeWithMaxLength(io.Reader, ssz.Unmarshaler) error
	// EncodeGossip는 주어진 writer에 주어진 ssz.Marshaler를 Gossip 메시지로 인코딩(직렬화)
	EncodeGossip(io.Writer, ssz.Marshaler) (int, error)
	// EncodeWithMaxLength는 주어진 writer에게 varint prefix를 포함한 바이트를 인코딩(직렬화)
	EncodeWithMaxLength(io.Writer, ssz.Marshaler) (int, error)
	// ?
	ProtocolSuffix() string
}
