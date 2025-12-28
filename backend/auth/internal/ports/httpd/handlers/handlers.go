package handlers

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/database"
	"github.com/aritradeveops/billbharat/backend/shared/eventbroker"
)

type Handler struct {
	db            database.Database
	eventProducer eventbroker.Producer
	Auth          *AuthHandler
	User          *UserHandler
}

func New(db database.Database, service *service.Service, eventProducer eventbroker.Producer) *Handler {
	return &Handler{
		db:            db,
		eventProducer: eventProducer,
		Auth:          NewAuthHandler(service.Auth),
		User:          NewUserHandler(service.User),
	}
}
