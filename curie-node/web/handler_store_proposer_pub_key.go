package c_web

import (
	ecdsacurie "flag-example/crypto/ecdsa"
	"flag-example/crypto/ecdsa/ecdsad"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) StoreProposerPubKey(c fiber.Ctx) error {
	// Subscriber
	pubKey, err := s.cfg.DB.GetDataFromRedis("Proposer")
	if err != nil {
		log.WithError(err).Error("Failed to get Proposer's Public Key from DB")
	}

	log.Info("Proposer's PubKey from RedisDB is ", pubKey)

	// Convert string to *ecdsa.PublicKey (Geth-secp256k1)
	ecdsaPubKey, err := ecdsacurie.ConvertToEcdsaPubKeyString(pubKey)
	if err != nil {
		log.WithError(err).Error("Failed to convert *ecdsa.Publickey from string")
	}

	// Singleton Pattern for storing pubKey
	ecdsad.PublicKeyFromProposer(ecdsaPubKey)

	return nil
}
