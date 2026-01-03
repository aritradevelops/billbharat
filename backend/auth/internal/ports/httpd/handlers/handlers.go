package handlers

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/database"
)

type Handler struct {
	db       database.Database
	Auth     *AuthHandler
	User     *UserHandler
	Business *BusinessHandler
}

func New(db database.Database, service *service.Service, environment string) *Handler {
	return &Handler{
		db:       db,
		Auth:     NewAuthHandler(service.Auth, environment),
		User:     NewUserHandler(service.User),
		Business: NewBusinessHandler(service.Business, environment),
	}
}
