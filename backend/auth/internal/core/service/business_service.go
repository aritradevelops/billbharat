package service

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/google/uuid"
)

type BusinessService interface {
	CreateBusiness(ctx context.Context, initiator string, payload CreateBusinessPayload) (CreateBusinessResponse, error)
}

type CreateBusinessPayload struct {
	Name            string   `json:"name" validate:"required,min=3,max=255"`
	Description     *string  `json:"description" validate:"min=50,max=255"`
	Logo            *string  `json:"logo" validate:"url"`
	Industry        string   `json:"industry" validate:"required,oneof=IT,Healthcare,Education,Finance,Manufacturing,Retail,Travel,Entertainment,Other"`
	PrimaryCurrency string   `json:"primary_currency" validate:"required"`
	OwnerID         string   `json:"owner_id" validate:"required"`
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

type businessService struct {
	repository   repository.Repository
	eventManager events.EventManager
}

func NewBusinessService(repository repository.Repository, eventManager events.EventManager) BusinessService {
	return &businessService{repository: repository, eventManager: eventManager}
}

func (s *businessService) CreateBusiness(ctx context.Context, initiator string, payload CreateBusinessPayload) (CreateBusinessResponse, error) {
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
		OwnerID:         uuid.MustParse(payload.OwnerID),
		Currencies:      payload.Currencies,
		CreatedBy:       uuid.MustParse(initiator),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create business")
		return response, InternalError
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
		return response, InternalError
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
