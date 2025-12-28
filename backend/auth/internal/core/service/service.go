package service

import (
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradeveops/billbharat/backend/shared/eventbroker"
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
