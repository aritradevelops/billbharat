package handlers

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradeveops/billbharat/backend/auth/internal/pkg/translation"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var payload RegisterPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.Register(c.Context(), service.RegisterPayload{
		Name:        payload.Name,
		Email:       payload.Email,
		CountryCode: payload.CountryCode,
		Phone:       payload.Phone,
		Password:    payload.Password,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.register", nil), nil, nil))
}
