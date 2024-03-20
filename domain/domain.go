package domain

import "github.com/gofiber/fiber/v2"

type Service interface {
	AddUser(AddUserBody) (int, fiber.Map)
	GetAllUsers() (int, fiber.Map)
	GetUser(string) (int, fiber.Map)
	UpdateScore(UpdateScore) (int, fiber.Map)
}

type Repositories interface {
	AddUser(User) error
	GetAllUsers() (*[]User, error)
	GetUser(string) (*User, error)
	UpdateScore(user User) error
}
