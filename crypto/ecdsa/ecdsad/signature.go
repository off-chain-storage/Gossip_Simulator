package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
	curiepb "flag-example/proto"
	"math/big"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

func (s *Signature) Verify(pubKey *ecdsa.PublicKey, msg []byte) bool {
	logrus.Info(pubKey)
	logrus.Info(len(msg))
	logrus.Info(s.sig.r)
	logrus.Info(s.sig.s)

	return ecdsa.Verify(pubKey, msg, s.sig.r, s.sig.s)
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
