package ecdsad

import (
	"bytes"
	"crypto/ecdsa"
	"flag-example/crypto/ecdsa/common"
	curiepb "flag-example/proto"

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

	logrus.Info("Original Public Key: ", len(sigPublicKey))
	k, err := crypto.DecompressPubkey(sigPublicKey)
	if err != nil {
		logrus.WithError(err).Error("Failed to decompress public key")
		return false
	}

	if bytes.Equal(crypto.CompressPubkey(k), crypto.CompressPubkey(pubKey)) {
		return true
	} else {
		logrus.Info("Original Public Key: ", crypto.CompressPubkey(pubKey))
		logrus.Info("Recovered Public Key: ", crypto.CompressPubkey(k))
		logrus.Info(bytes.Equal(crypto.CompressPubkey(k), crypto.CompressPubkey(pubKey)))
		logrus.Info("Signature: ", s.sig)
		return false
	}

	// if (sigPublicKey != nil) && bytes.Equal(crypto.CompressPubkey(pubKey), sigPublicKey) {
	// 	return true
	// } else {
	// 	logrus.Info("Original Public Key: ", crypto.CompressPubkey(pubKey))
	// 	logrus.Info("Recovered Public Key: ", sigPublicKey)
	// 	logrus.Info(bytes.Equal(crypto.CompressPubkey(pubKey), sigPublicKey))
	// 	logrus.Info("Signature: ", s.sig)
	// 	return false
	// }
}

func (s *Signature) Marshal() []byte {
	return s.sig
}

func (s *Signature) Proto() (proto.Message, error) {
	return &curiepb.Signature{
		Sig: s.sig,
	}, nil
}
