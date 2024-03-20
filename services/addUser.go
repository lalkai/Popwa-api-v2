package services

import (
	"popwa/domain"
	"popwa/logs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *popWaService) AddUser(newUser domain.AddUserBody) (int, fiber.Map) {

	user := domain.User{
		UserName: newUser.UserName,
		Score: 0,
	}
	
	err := s.popWaRepo.AddUser(user)
	if err != nil {
		logs.Error("Error saving profile:", zap.Error(err))
		return fiber.StatusInternalServerError, fiber.Map{
			"error": "failed to create",
		}
	}

	return fiber.StatusOK, fiber.Map{
		"message": "saved successfully",
	}
}
