package services

import (
	"popwa/domain"
	"popwa/logs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *popWaService) UpdateScore(updateScore domain.UpdateScore) (int, fiber.Map) {

	user := domain.User{
		UserName: updateScore.UserName,
		Score:    updateScore.Score,
	}

	err := s.popWaRepo.UpdateScore(user)
	if err != nil {
		logs.Error("Error update user:", zap.Error(err))
		return fiber.StatusInternalServerError, fiber.Map{
			"error": "failed to update",
		}
	}

	return fiber.StatusOK, fiber.Map{
		"message": "update successfully",
	}
}