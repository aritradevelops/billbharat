package handlers

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradeveops/billbharat/backend/auth/internal/ports/httpd/authn"
	"github.com/aritradeveops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSrv service.UserService
}

type UpdateDPPayload struct {
	Dp string `form:"dp"`
}

func NewUserHandler(userSrv service.UserService) *UserHandler {
	return &UserHandler{userSrv: userSrv}
}
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	response, err := h.userSrv.Profile(c.Context(), user.UserID, service.ProfilePayload{
		ID: c.Params("id"),
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "user.profile", nil), response, nil))
}

func (h *UserHandler) UpdateDP(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var payload UpdateDPPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	response, err := h.userSrv.UpdateDP(c.Context(), user.UserID, service.UpdateDPPayload{
		Dp: payload.Dp,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "user.update_dp", nil), response, nil))
}
