package api

import (
	
	"popwa/domain"
	"popwa/logs"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PopWaHandler interface {
	AddUser(c *fiber.Ctx) error
	GetAllUsers(c *websocket.Conn)
	GetUser(c *fiber.Ctx) error
	UpdateScore(c *fiber.Ctx) error
}

type popWaHandler struct {
	popWaService domain.Service
	clients      map[*websocket.Conn]struct{}
	mu           sync.Mutex
}

func NewPopWaHandler(popWaService domain.Service) PopWaHandler {
	return &popWaHandler{
		popWaService: popWaService,
		clients:      make(map[*websocket.Conn]struct{}),
	}
}

func (h *popWaHandler) AddUser(c *fiber.Ctx) error {
	var newUser domain.AddUserBody
	if err := c.BodyParser(&newUser); err != nil {
		logs.Info("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	statusCode, response := h.popWaService.AddUser(newUser)
	return c.Status(statusCode).JSON(response)
}

func (h *popWaHandler) GetAllUsers(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	for {
		_, data := h.popWaService.GetAllUsers()
	
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteJSON(data)
			if err != nil {
				logs.Error("Error sending data to client:", zap.Error(err))
				delete(h.clients, client)
				client.Close()
			}
		}
		h.mu.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func (h *popWaHandler) GetUser(c *fiber.Ctx) error {
	username := c.Params("Username")

	statusCode, response := h.popWaService.GetUser(username)
	return c.Status(statusCode).JSON(response)
}

func (h *popWaHandler) UpdateScore(c *fiber.Ctx) error {
	var updateScore domain.UpdateScore
	if err := c.BodyParser(&updateScore); err != nil {
		logs.Info("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	statusCode, response := h.popWaService.UpdateScore(updateScore)
	return c.Status(statusCode).JSON(response)
}
