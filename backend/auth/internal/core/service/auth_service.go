package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aritradeveops/billbharat/backend/auth/internal/core/cryptoutil"
	"github.com/aritradeveops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradeveops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradeveops/billbharat/backend/auth/internal/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost               = 10
	VerificationRequestExpiry = 15 * time.Minute
)

var (
	RootUserID    uuid.UUID = uuid.Nil
	UserExistsErr           = &ServiceError{HttpErrorCode: http.StatusConflict,
		DevErrorCode: "auth_001", Short: "user.exists", Long: "user already exists"}
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
	repository repository.Repository
}

func NewAuthService(repository repository.Repository) AuthService {
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

	tx, err := s.repository.StartTransaction(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to start transaction")
		return response, err
	}
	defer tx.Rollback(ctx)
	repo := s.repository.WithTx(tx)

	user, err := repo.CreateUser(ctx, dao.CreateUserParams{
		HumanID:       cryptoutil.HumanID("user"),
		Name:          payload.Name,
		Email:         payload.Email,
		Phone:         fmt.Sprintf("%s%s", payload.CountryCode, payload.Phone),
		EmailVerified: false,
		CreatedBy:     RootUserID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create user")
		return response, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate password hash")
		return response, InternalError
	}

	err = repo.CreatePassword(ctx, dao.CreatePasswordParams{
		UserID:    user.ID,
		Password:  string(hashedPassword),
		CreatedBy: user.ID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to create password")
		return response, err
	}

	otp, err := cryptoutil.GeneratOTP(6)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate otp")
		return response, err
	}

	err = repo.CreateVerificationRequest(ctx, dao.CreateVerificationRequestParams{
		UserID:    user.ID,
		Code:      otp,
		Type:      dao.VerificationTypeEmail,
		ExpiresAt: time.Now().Add(VerificationRequestExpiry),
		CreatedBy: user.ID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to create verification request")
		return response, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return response, err
	}

	return response, nil
}
