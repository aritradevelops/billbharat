package service

import (
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/repository"
)

type Service struct {
	Category ProductCategoryService
}

func New(repository repository.Repository) *Service {
	return &Service{
		Category: &productCategoryService{
			repository: repository,
		},
	}
}
