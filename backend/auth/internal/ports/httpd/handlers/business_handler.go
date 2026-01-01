package handlers

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/authn"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
)

type BusinessHandler struct {
	businessSrv service.BusinessService
}

func NewBusinessHandler(businessSrv service.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessSrv: businessSrv,
	}
}

func (h *BusinessHandler) Create(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var payload service.CreateBusinessPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	resp, err := h.businessSrv.Create(c.Context(), user.UserID, payload)
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "controller.create", fiber.Map{"Entity": "Business"}), resp, nil))
}
