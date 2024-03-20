package services

import (
	"popwa/logs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *popWaService) GetUser(Username string) (int, fiber.Map) {
	user, err := s.popWaRepo.GetUser(Username)
	if err != nil {
		logs.Info("Error retrieving users", zap.Error(err))
		return fiber.StatusInternalServerError, fiber.Map{
			"message": "internal server error",
		}
	}

	return fiber.StatusOK, fiber.Map{
		"message": "user get success",
		"users":   user,
	}
}
