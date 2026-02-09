package note

import (
	"errors"
	"golang-boilerplate-example/queue/email"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

type Handler struct {
	Service *Service
	Queue   *asynq.Client
}

func NewHandler(service *Service, queue *asynq.Client) *Handler {
	return &Handler{
		Service: service,
		Queue:   queue,
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

func (h *Handler) SendEmail(c *fiber.Ctx) error {
	type Payload struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var body Payload
	if err := c.BodyParser(&body); err != nil {
		return handleError(c, err)
	}

	task, err := email.NewSendEmailTask(email.Payload{
		To:      body.To,
		Subject: body.Subject,
		Body:    body.Body,
	})
	if err != nil {
		return handleError(c, err)
	}

	_, err = h.Queue.Enqueue(
		task,
		asynq.MaxRetry(5),
		asynq.Timeout(30),
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"message": "email queued"})
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
