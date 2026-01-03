package service

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/aritradevelops/billbharat/backend/auth/internal/core/cryptoutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/google/uuid"
)

var (
	BusinessNotFoundErr = &ServiceError{
		HttpErrorCode: http.StatusNotFound, Short: "business.not_found", Long: "business not found",
		DevErrorCode: "business_001",
	}
	InvalidBusinessIdErr = &ServiceError{
		HttpErrorCode: http.StatusBadRequest, Short: "business.invalid_id", Long: "invalid business id",
		DevErrorCode: "business_002",
	}
)

type BusinessService interface {
	Create(ctx context.Context, initiator string, payload CreateBusinessPayload) (CreateBusinessResponse, error)
	List(ctx context.Context, initiator string) (ListBusinessesResponse, error)
	Select(ctx context.Context, initiator string, businessID string, payload SwitchBusinessPayload) (LoginResponse, error)
}

type CreateBusinessPayload struct {
	Name            string   `json:"name" validate:"required,min=3,max=255"`
	Description     *string  `json:"description" validate:"omitempty,min=50,max=255"`
	Logo            *string  `json:"logo" validate:"omitempty,url"`
	Industry        string   `json:"industry" validate:"required,oneof=IT Healthcare Education Finance Manufacturing Retail Travel Entertainment Other"`
	PrimaryCurrency string   `json:"primary_currency" validate:"required"`
	Currencies      []string `json:"currencies" validate:"required"`
}

type CreateBusinessResponse struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Description     *string  `json:"description"`
	Logo            *string  `json:"logo"`
	Industry        string   `json:"industry"`
	PrimaryCurrency string   `json:"primary_currency"`
	OwnerID         string   `json:"owner_id"`
	Currencies      []string `json:"currencies"`
}

type ListBusinessesResponse struct {
	Businesses []ListBusinessesResponseBusiness `json:"businesses"`
}

type ListBusinessesResponseBusiness struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type SwitchBusinessPayload struct {
	UserIP    string `json:"user_ip"`
	UserAgent string `json:"user_agent"`
}

type businessService struct {
	repository   repository.Repository
	eventManager events.EventManager
	jwtManager   *jwtutil.JwtManager
}

func NewBusinessService(repository repository.Repository, eventManager events.EventManager, jwtManager *jwtutil.JwtManager) BusinessService {
	return &businessService{repository: repository, eventManager: eventManager, jwtManager: jwtManager}
}

func (s *businessService) Create(ctx context.Context, initiator string, payload CreateBusinessPayload) (CreateBusinessResponse, error) {
	var response CreateBusinessResponse
	errs := validation.Validate(payload)
	if errs != nil {
		return response, errs
	}
	tx, err := s.repository.StartTransaction(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to start transaction")
		return response, err
	}
	defer tx.Rollback(ctx)
	repo := s.repository.WithTx(tx)

	business, err := repo.CreateBusiness(ctx, dao.CreateBusinessParams{
		Name:            payload.Name,
		Description:     payload.Description,
		Logo:            payload.Logo,
		Industry:        payload.Industry,
		PrimaryCurrency: payload.PrimaryCurrency,
		OwnerID:         uuid.MustParse(initiator),
		Currencies:      payload.Currencies,
		CreatedBy:       uuid.MustParse(initiator),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create business")
		return response, err
	}

	businessUser, err := repo.CreateBusinessUser(ctx, dao.CreateBusinessUserParams{
		UserID:     business.OwnerID,
		BusinessID: business.ID,
		// this has nothing to do with authorization
		Role:      "Owner",
		CreatedBy: uuid.MustParse(initiator),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create business user")
		return response, err
	}

	err = s.eventManager.EmitManageBusinessEvent(ctx, events.NewBusinessManageEvent("create", events.MangageBusinessEventPayload(business)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage business event")
		return response, InternalError
	}

	err = s.eventManager.EmitManageBusinessUserEvent(ctx, events.NewBusinessUserManageEvent("create", events.MangageBusinessUserEventPayload(businessUser)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage business user event")
		return response, InternalError
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return response, InternalError
	}

	response = CreateBusinessResponse{
		ID:              business.ID.String(),
		Name:            business.Name,
		Description:     business.Description,
		Logo:            business.Logo,
		Industry:        business.Industry,
		PrimaryCurrency: business.PrimaryCurrency,
		OwnerID:         business.OwnerID.String(),
		Currencies:      business.Currencies,
	}

	return response, nil
}

func (s *businessService) List(ctx context.Context, initiator string) (ListBusinessesResponse, error) {
	response := ListBusinessesResponse{
		Businesses: []ListBusinessesResponseBusiness{},
	}
	businesses, err := s.repository.FindBusinessesByUserID(ctx, uuid.MustParse(initiator))
	if err != nil {
		logger.Error().Err(err).Msg("failed to find businesses by user id")
		return response, err
	}

	for _, business := range businesses {
		response.Businesses = append(response.Businesses, ListBusinessesResponseBusiness{
			Name: business.Business.Name,
			ID:   business.Business.ID.String(),
		})
	}

	return response, nil
}

func (s *businessService) Select(ctx context.Context, initiator string, businessID string, payload SwitchBusinessPayload) (LoginResponse, error) {
	var response LoginResponse
	businesses, err := s.repository.FindBusinessesByUserID(ctx, uuid.MustParse(initiator))
	if err != nil {
		logger.Error().Err(err).Msg("failed to find businesses by user id")
		return response, InternalError
	}

	if !slices.ContainsFunc(businesses, func(business dao.FindBusinessesByUserIDRow) bool {
		return business.Business.ID.String() == businessID
	}) {
		return response, BusinessNotFoundErr
	}

	user, err := s.repository.FindUserById(ctx, uuid.MustParse(initiator))
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user by id")
		return response, InternalError
	}

	accessToken, err := s.jwtManager.Sign(jwtutil.JwtPayload{
		UserID:     user.ID.String(),
		Email:      user.Email,
		Name:       user.Name,
		Dp:         user.Dp,
		BusinessID: businessID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to sign access token")
		return response, InternalError
	}

	refreshToken, err := cryptoutil.GenerateRefreshToken()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate refresh token")
		return response, InternalError
	}

	err = s.repository.CreateSession(ctx, dao.CreateSessionParams{
		HumanID:      cryptoutil.HumanID("session"),
		UserID:       user.ID,
		UserIp:       payload.UserIP,
		UserAgent:    payload.UserAgent,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(SessionExpiry),
		CreatedBy:    user.ID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to create session")
		return response, InternalError
	}

	response.AccessToken = accessToken.Token
	response.AccessTokenLifetime = accessToken.Lifetime
	response.RefreshToken = refreshToken
	response.RefreshTokenLifetime = time.Now().Add(SessionExpiry)

	return response, nil
}
