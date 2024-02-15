package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
	curiepb "flag-example/proto"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Signature struct {
	sig []byte
}

func InitSignFromProto(i interface{}) (common.Signature, error) {
	switch pb := i.(type) {
	case *curiepb.Signature:
		return &Signature{sig: pb.Sig}, nil
	default:
		return nil, errors.Wrapf(errors.New("unsupported signed curie block"), "unable to create block from type %T", i)
	}
}

func (s *Signature) Verify(pubKey *ecdsa.PublicKey, msg []byte) bool {
	comPubKey := crypto.CompressPubkey(pubKey)
	return crypto.VerifySignature(comPubKey, msg, s.sig)
}

func (s *Signature) Marshal() []byte {
	return s.sig
}

func (s *Signature) Proto() (proto.Message, error) {
	return &curiepb.Signature{
		Sig: s.sig,
	}, nil
}
