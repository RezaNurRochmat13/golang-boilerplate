package user

import (
	"golang-boilerplate-example/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(c *fiber.Ctx) error {

	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	err := h.service.Register(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user registered",
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	user, err := h.service.Login(c.Context(), req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	token, err := auth.GenerateToken(uint(user.ID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "token generation failed",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
