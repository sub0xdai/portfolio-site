package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sub0x/resume-ai/internal/knowledge"
	"github.com/sub0x/resume-ai/internal/types"
)

type Handler struct {
	service *knowledge.Service
}

func NewHandler(service *knowledge.Service) *Handler {
	return &Handler{service: service}
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Chunks []types.Chunk `json:"chunks"`
}

func (h *Handler) HandleQuery(c *fiber.Ctx) error {
	var req QueryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	chunks, err := h.service.Query(req.Query, 5)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to query: %v", err),
		})
	}

	return c.JSON(QueryResponse{Chunks: chunks})
}

func (h *Handler) HandleGetTags(c *fiber.Ctx) error {
	tags := h.service.GetAllTags()
	return c.JSON(fiber.Map{
		"tags": tags,
	})
}

func (h *Handler) HandleGetNotesByTag(c *fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tag parameter is required",
		})
	}

	notes := h.service.GetNotesByTag(tag)
	return c.JSON(fiber.Map{
		"notes": notes,
	})
}
