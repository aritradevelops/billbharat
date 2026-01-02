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

func (h *BusinessHandler) List(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	resp, err := h.businessSrv.List(c.Context(), user.UserID)
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "controller.list", fiber.Map{"Entity": "Business"}), resp, nil))
}

func (h *BusinessHandler) Select(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	response, err := h.businessSrv.Select(c.Context(), user.UserID, c.Params("business_id"), service.SwitchBusinessPayload{
		UserIP:    c.IP(),
		UserAgent: c.Get("User-Agent"),
	})
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    response.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  response.AccessTokenLifetime,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		Path:     "/",
		HTTPOnly: true,
		Expires:  response.RefreshTokenLifetime,
	})
	return c.JSON(NewResponse(translation.Localize(c, "auth.login"), response, nil))
}
