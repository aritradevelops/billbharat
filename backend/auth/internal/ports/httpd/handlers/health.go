package handlers

import (
	"github.com/aritradeveops/billbharat/backend/shared/logger"
	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Overall  string `json:"overall"`
	Database string `json:"database"`
}

func (h *Handler) Health(c *fiber.Ctx) error {
	var response HealthResponse
	hasError := false
	err := h.db.Health()
	if err != nil {
		logger.Error().Err(err).Msg("database is not healthy")
		response.Database = "NOT OK"
		hasError = true
	} else {
		response.Database = "OK"
	}
	if hasError {
		response.Overall = "NOT OK"
	} else {
		response.Overall = "OK"
	}

	return c.JSON(response)
}
