package service

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/repository"
)

type Service struct {
	Auth AuthService
	User UserService
}

func New(repository repository.Repository, jwtManager *jwtutil.JwtManager) *Service {
	return &Service{
		Auth: NewAuthService(repository, jwtManager),
		User: NewUserService(repository),
	}
}
