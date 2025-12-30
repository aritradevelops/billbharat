package handlers

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/database"
)

type Handler struct {
	db   database.Database
	Auth *AuthHandler
	User *UserHandler
}

func New(db database.Database, service *service.Service) *Handler {
	return &Handler{
		db:   db,
		Auth: NewAuthHandler(service.Auth),
		User: NewUserHandler(service.User),
	}
}
