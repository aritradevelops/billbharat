package service

import "github.com/gofiber/fiber/v2"

type ServiceError struct {
	HttpErrorCode int    `json:"http_error_code"`
	DevErrorCode  string `json:"dev_error_code"`
	Short         string `json:"short"`
	Long          string `json:"long"`
}

func (e *ServiceError) Error() string {
	return e.Long
}

const (
	GeneralErrorCode = 1000
)

var (
	InternalError = &ServiceError{HttpErrorCode: fiber.StatusInternalServerError,
		DevErrorCode: "general_internal_error", Short: "Internal server error", Long: "Internal server error"}
)
