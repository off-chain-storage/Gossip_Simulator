package common

import (
	"google.golang.org/protobuf/proto"
)

type PrivateKey interface {
	PublicKey() PublicKey
	Sign(msg []byte) Signature
}

type PublicKey interface {
	Copy() PublicKey
}

type Signature interface {
	Verify(pubKey PublicKey, msg []byte) bool
	Marshal() ([]byte, []byte, error)
	Proto() (proto.Message, error)
}
