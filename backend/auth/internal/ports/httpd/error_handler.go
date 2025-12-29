package httpd

import (
	"fmt"

	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/handlers"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		logger.Error().Type("type", err).Err(err).Msg("request failed")

		if e, ok := err.(*fiber.Error); ok {
			c.Status(e.Code)
			return c.JSON(handlers.NewResponse(translation.Localize(c, fmt.Sprintf("errors.%d", e.Code)), nil, err))
		}

		if e, ok := err.(*service.ServiceError); ok {
			c.Status(e.HttpErrorCode)
			return c.JSON(handlers.NewResponse(translation.Localize(c, e.Short), nil, err))
		}

		if e, ok := err.(*validation.ValidationErrors); ok {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(handlers.NewResponse(translation.Localize(c, "errors.422"), nil, e))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(handlers.NewResponse(translation.Localize(c, "errors.500"), nil, err))
	}
}
