package handlers

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/database"
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
