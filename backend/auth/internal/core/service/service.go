package service

import (
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/eventbroker"
)

type Service struct {
	Auth AuthService
	User UserService
}

func New(repository repository.Repository, jwtManager *jwtutil.JwtManager, eventBroker eventbroker.Producer) *Service {
	return &Service{
		Auth: NewAuthService(repository, jwtManager, eventBroker),
		User: NewUserService(repository),
	}
}
