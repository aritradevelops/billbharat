package handlers

import (
	"github.com/aritradevelops/billbharat/backend/product/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/database"
)

type Handler struct {
	db       database.Database
	Category *ProductCategoryHandler
}

func New(db database.Database, service *service.Service, environment string) *Handler {
	return &Handler{
		db:       db,
		Category: NewProductCategoryHandler(service.Category),
	}
}
