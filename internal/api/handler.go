package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sub0x/resume-ai/internal/types"
)

type Handler struct {}

func NewHandler(_ interface{}) *Handler {
	return &Handler{}
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Answer string `json:"answer"`
}

func (h *Handler) HandleQuery(c *fiber.Ctx) error {
	var req QueryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Return a simple response without knowledge base integration
	return c.JSON(QueryResponse{
		Answer: fmt.Sprintf("You asked: %s. This is a simplified response without knowledge base integration.", req.Query),
	})
}

func (h *Handler) HandleGetTags(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"tags": []string{},
	})
}

func (h *Handler) HandleGetNotesByTag(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"notes": []types.Note{},
	})
}
