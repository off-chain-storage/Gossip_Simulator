package interfaces

import (
	"google.golang.org/protobuf/proto"

	"flag-example/crypto/ecdsa/common"
)

// Signature interface
type Signature interface {
	Proto() (proto.Message, error)
}

// O.G 에서 Gossip 후 Verify 전 Hash 값 생성을 위한 Hash() 함수 추가
type SignedCurieBlock interface {
	Proto() (proto.Message, error)
	Hash() []byte
	Signature() common.Signature
}

// N.G 에서 Gossip 전 Hash 값 생성을 위한 Hash() 함수 추가
type ReadOnlyCurieBlock interface {
	Proto() (proto.Message, error)
	Hash() []byte
	Body() []byte
}
