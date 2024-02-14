package web

import (
	file_util "flag-example/io/file"

	"github.com/gofiber/fiber/v3"
)

func (s *Server) ProposeCurieBlockForOG(c fiber.Ctx) error {

	if s.proposerService == nil {
		// http HandleError
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	f, err := file.Open()
	if err != nil {
		return err
	}

	rawBlkData, err := file_util.FileToBytes(f)
	if err != nil {
		return err
	}

	if err := s.curieNodeProposer.ProposeCurieBlockForOG(c.Context(), rawBlkData); err != nil {
		log.WithError(err).Error("Failed to propose block to curie node by Original Gossip")
		return err
	}

	return nil
}
