package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
	curiepb "flag-example/proto"
	"math/big"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type signature struct {
	r, s *big.Int
}

type Signature struct {
	sig *signature
}

func InitSignFromProto(i interface{}) (common.Signature, error) {
	switch pb := i.(type) {
	case *curiepb.Signature:
		sig := &signature{
			r: new(big.Int).SetBytes(pb.SigR),
			s: new(big.Int).SetBytes(pb.SigS),
		}
		return &Signature{sig: sig}, nil
	default:
		return nil, errors.Wrapf(errors.New("unsupported signed curie block"), "unable to create block from type %T", i)
	}
}

func (s *Signature) Verify(pubKey common.PublicKey, msg []byte) bool {
	return ecdsa.Verify(pubKey.(*PublicKey).p, msg, s.sig.r, s.sig.s)
}

func (s *Signature) Marshal() ([]byte, []byte, error) {
	return s.sig.r.Bytes(), s.sig.s.Bytes(), nil
}

func (s *Signature) Proto() (proto.Message, error) {
	return &curiepb.Signature{
		SigR: s.sig.r.Bytes(),
		SigS: s.sig.s.Bytes(),
	}, nil
}
