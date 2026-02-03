package note

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetAllNotes(c *fiber.Ctx) error {
	notes, err := h.Service.GetAllNotes()
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notes fetched successfully",
		"data":    notes,
	})
}

func (h *Handler) GetNote(c *fiber.Ctx) error {
	id := c.Params("id")

	note, err := h.Service.GetNote(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note fetched successfully",
		"data":    note,
	})
}

func (h *Handler) CreateNote(c *fiber.Ctx) error {
	var note Note
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := h.Service.CreateNote(&note); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note created successfully",
	})
}

func (h *Handler) UpdateNote(c *fiber.Ctx) error {
	var note Note
	if err := c.BodyParser(&note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if err := h.Service.UpdateNote(c.Params("id"), &note); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note updated successfully",
	})
}

func (h *Handler) DeleteNote(c *fiber.Ctx) error {
	if err := h.Service.DeleteNote(c.Params("id")); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note deleted successfully",
	})
}

func handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidPayload):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})

	case errors.Is(err, ErrNoteNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})

	default:
		// log error di sini kalau mau
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
}
