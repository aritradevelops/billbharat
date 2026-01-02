package handlers

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/authn"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSrv service.UserService
}

type UpdateDPPayload struct {
	Dp *string `json:"dp"`
}

type InvitePayload struct {
	Name        string `json:"name"`
	Email       string `json:"email" `
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone" `
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

func (h *UserHandler) Invite(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var payload InvitePayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	response, err := h.userSrv.Invite(c.Context(), user.UserID, user.BusinessID, service.InvitePayload{
		Name:        payload.Name,
		Email:       payload.Email,
		CountryCode: payload.CountryCode,
		Phone:       payload.Phone,
		Origin:      c.Get("Origin"),
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "user.invite", nil), response, nil))
}
