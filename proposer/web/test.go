package web

import (
	"github.com/gofiber/fiber/v3"
)

func (s *Server) Test(c fiber.Ctx) error {
	log.Info("HI")
	id := c.Params("cid")
	log.Info(id)

	return c.Status(fiber.StatusOK).JSON("id")
}
