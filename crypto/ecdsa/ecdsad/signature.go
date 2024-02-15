package ecdsad

import (
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
	curiepb "flag-example/proto"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	sigPublicKey, err := crypto.Ecrecover(msg, s.sig)
	if err != nil {
		logrus.WithError(err).Error("Failed to recover public key")
		return false
	}

	recoveredPubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		log.Fatalf("Failed to unmarshal public key: %v", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)
	originalAddr := crypto.PubkeyToAddress(*pubKey)
	if recoveredAddr.Hex() == originalAddr.Hex() {
		return true
	} else {
		return false
	}
}

func (s *Signature) Marshal() []byte {
	return s.sig
}

func (s *Signature) Proto() (proto.Message, error) {
	return &curiepb.Signature{
		Sig: s.sig,
	}, nil
}
