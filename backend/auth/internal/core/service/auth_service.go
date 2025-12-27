package service

import (
	"context"
	"fmt"

	"github.com/aritradeveops/billbharat/backend/auth/internal/core/cryptoutil"
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradeveops/billbharat/backend/auth/internal/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost          = 10
	AuthServiceErrorCode = 2000
)

var (
	UserExistsErr = &ServiceError{AuthServiceErrorCode + 1, "user.exists", "user already exists"}
)

type RegisterPayload struct {
	Name        string `json:"name" validate:"alphaspace,min=3,max=255"`
	Email       string `json:"email" validate:"email"`
	CountryCode string `json:"country_code" validate:"required"`
	Phone       string `json:"phone" validate:"numeric,min=10,max=16"`
	Password    string `json:"password" validate:"min=8,max=255"`
}

type RegisterResponse struct {
}

type AuthService interface {
	Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, error)
}

type authService struct {
	repository repository.Querier
}

func NewAuthService(repository repository.Querier) AuthService {
	return &authService{
		repository: repository,
	}
}

func (s *authService) Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, error) {
	var response RegisterResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}
	errs = validation.ValidatePassword(payload.Password)
	if errs != nil {
		logger.Error().Err(errs).Msg("password validation failed")
		return response, errs
	}
	if _, err := s.repository.FindUserByEmail(ctx, payload.Email); err == nil {
		logger.Error().Err(err).Msg("user already exists")
		return response, UserExistsErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate password hash")
		return response, InternalError
	}

	if err := s.repository.CreateUser(ctx, repository.CreateUserParams{
		HumanID:       cryptoutil.HumanID("user"),
		Name:          payload.Name,
		Email:         payload.Email,
		Phone:         fmt.Sprintf("%s%s", payload.CountryCode, payload.Phone),
		Password:      string(hashedPassword),
		EmailVerified: false,
		CreatedBy:     uuid.Nil,
	}); err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return response, err
	}
	// TODO: send verification email
	return response, nil
}
