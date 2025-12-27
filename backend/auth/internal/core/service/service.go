package service

import "github.com/aritradeveops/billbharat/backend/auth/internal/persistence/repository"

type Service struct {
	Auth AuthService
}

func New(repository repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repository),
	}
}
