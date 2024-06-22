package api

import (
	"encoding/json"
	"popwa/domain"
	"popwa/logs"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PopWaHandler interface {
	AddUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	HandleWebSocket(c *websocket.Conn)
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

func (h *popWaHandler) GetUser(c *fiber.Ctx) error {
	username := c.Params("Username")

	statusCode, response := h.popWaService.GetUser(username)
	return c.Status(statusCode).JSON(response)
}

func (h *popWaHandler) HandleWebSocket(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
		c.Close()
	}()

	for {
		var message struct {
			Action string          `json:"action"`
			Data   json.RawMessage `json:"data"`
		}

		err := c.ReadJSON(&message)
		if err != nil {
			logs.Error("Error reading JSON from client:", zap.Error(err))
			break
		}

		switch message.Action {
		case "getLeaderboard":
			_, users := h.popWaService.GetAllUsers()
			err := h.SendWebSocketMessage(c, "updateLeaderboard", users)
			if err != nil {
				logs.Error("Error sending data to client:", zap.Error(err))
			}

		case "updateScore":
			var updateScore domain.UpdateScore
			if err := json.Unmarshal(message.Data, &updateScore); err != nil {
				logs.Info("Invalid updateScore request", zap.Error(err))
				c.WriteJSON(fiber.Map{
					"error": "Invalid request",
				})
				continue
			}

			_, response := h.popWaService.UpdateScore(updateScore)
			_, users := h.popWaService.GetAllUsers()
			err := h.SendWebSocketMessage(c, "updateScore", response)
			if err != nil {
				logs.Error("Error sending response to client:", zap.Error(err))
			}
			err = h.SendWebSocketMessage(c, "updateLeaderboard", users)
			if err != nil {
				logs.Error("Error sending users to client:", zap.Error(err))
			}

			h.mu.Lock()
			for client := range h.clients {
				err := h.SendWebSocketMessage(client, "updateScore", response)
				if err != nil {
					logs.Error("Error broadcasting update to client:", zap.Error(err))
					delete(h.clients, client)
					client.Close()
				}
				err = h.SendWebSocketMessage(client, "updateLeaderboard", users)
				if err != nil {
					logs.Error("Error broadcasting users to client:", zap.Error(err))
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *popWaHandler) SendWebSocketMessage(c *websocket.Conn, action string, data interface{}) error {
	message := struct {
		Action string      `json:"action"`
		Data   interface{} `json:"data"`
	}{
		Action: action,
		Data:   data,
	}

	return c.WriteJSON(message)
}
