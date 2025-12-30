package service

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
)

type Service struct {
	Auth AuthService
	User UserService
}

func New(repository repository.Repository, jwtManager *jwtutil.JwtManager, eventManger events.EventManager) *Service {
	return &Service{
		Auth: NewAuthService(repository, jwtManager, eventManger),
		User: NewUserService(repository, eventManger),
	}
}
