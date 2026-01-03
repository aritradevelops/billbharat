package handlers

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/authn"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
	environment string
}

func NewAuthHandler(authService service.AuthService, environment string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		environment: environment,
	}
}

type RegisterPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
}

type VerifyEmailPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type VerifyPhonePayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendEmailVerificationRequestPayload struct {
	Email string `json:"email"`
}

type SendPhoneVerificationRequestPayload struct {
	Email string `json:"email"`
}

type ForgotPasswordPayload struct {
	Email string `json:"email"`
}

type ResetPasswordPayload struct {
	Email           string `json:"email"`
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangePasswordPayload struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
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

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var payload VerifyEmailPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.VerifyEmail(c.Context(), service.VerifyEmailPayload{
		Email: payload.Email,
		Code:  payload.Code,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.verify_email", nil), nil, nil))
}
func (h *AuthHandler) VerifyPhone(c *fiber.Ctx) error {
	var payload VerifyPhonePayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.VerifyPhone(c.Context(), service.VerifyPhonePayload{
		Email: payload.Email,
		Code:  payload.Code,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.verify_phone", nil), nil, nil))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload LoginPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	response, err := h.authService.Login(c.Context(), service.LoginPayload{
		Email:     payload.Email,
		Password:  payload.Password,
		UserIP:    c.IP(),
		UserAgent: c.Get("User-Agent"),
	})
	if err != nil {
		return err
	}

	cookieSameSite := fiber.CookieSameSiteLaxMode
	secure := true
	if h.environment != "production" {
		cookieSameSite = fiber.CookieSameSiteStrictMode
		secure = false
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    response.AccessToken,
		HTTPOnly: true,
		Expires:  response.AccessTokenLifetime,
		SameSite: cookieSameSite,
		Secure:   secure,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    response.RefreshToken,
		HTTPOnly: true,
		Expires:  response.RefreshTokenLifetime,
		SameSite: cookieSameSite,
		Secure:   secure,
	})

	return c.JSON(NewResponse(translation.Localize(c, "auth.login", nil), response, nil))
}

func (h *AuthHandler) SendEmailVerificationRequest(c *fiber.Ctx) error {
	var payload VerifyEmailPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.SendEmailVerificationRequest(c.Context(), service.SendEmailVerificationRequestPayload{
		Email: payload.Email,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.send_email_verification_request", nil), nil, nil))
}

func (h *AuthHandler) SendPhoneVerificationRequest(c *fiber.Ctx) error {
	var payload VerifyEmailPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.SendPhoneVerificationRequest(c.Context(), service.SendPhoneVerificationRequestPayload{
		Email: payload.Email,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.send_phone_verification_request", nil), nil, nil))
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var payload ForgotPasswordPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.ForgotPassword(c.Context(), service.ForgotPasswordPayload{
		Email: payload.Email,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.forgot_password", nil), nil, nil))
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var payload ResetPasswordPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	_, err := h.authService.ResetPassword(c.Context(), service.ResetPasswordPayload{
		Email:           payload.Email,
		Code:            payload.Code,
		Password:        payload.Password,
		ConfirmPassword: payload.ConfirmPassword,
	})
	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.reset_password", nil), nil, nil))
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var payload ChangePasswordPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	_, err = h.authService.ChangePassword(c.Context(), user.UserID, service.ChangePasswordPayload{
		Email:           user.Email,
		CurrentPassword: payload.CurrentPassword,
		NewPassword:     payload.NewPassword,
		ConfirmPassword: payload.ConfirmPassword,
	})

	if err != nil {
		return err
	}
	return c.JSON(NewResponse(translation.Localize(c, "auth.change_password", nil), nil, nil))
}
