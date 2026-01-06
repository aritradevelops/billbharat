package service

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/product/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/google/uuid"
)

type ProductCategoryService interface {
	CreateProductCategory(ctx context.Context, payload CreateProductCategoryPayload) (CreateProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, payload UpdateProductCategoryPayload) (UpdateProductCategoryResponse, error)
	ListProductCategories(ctx context.Context, payload ListProductCategoryPayload) ([]ListProductCategoryResponse, error)
}

type ListProductCategoryPayload struct {
	BusinessID uuid.UUID `json:"business_id" validate:"required,uuid"`
	Page       int       `json:"page" validate:"required"`
	Limit      int       `json:"limit" validate:"required"`
}

type ListProductCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateProductCategoryPayload struct {
	Name       string    `json:"name" validate:"required,min=3,max=100"`
	BusinessID uuid.UUID `json:"business_id" validate:"required,uuid"`
	Initiator  uuid.UUID `json:"created_by" validate:"required,uuid"`
}

type CreateProductCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateProductCategoryPayload struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid"`
	Name       string    `json:"name" validate:"required,min=3,max=100"`
	BusinessID uuid.UUID `json:"business_id" validate:"required,uuid"`
	Initiator  uuid.UUID `json:"updated_by" validate:"required,uuid"`
}

type UpdateProductCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
type productCategoryService struct {
	repository repository.Repository
}

func NewProductCategoryService(repository repository.Repository) ProductCategoryService {
	return &productCategoryService{
		repository: repository,
	}
}

func (s *productCategoryService) CreateProductCategory(ctx context.Context, payload CreateProductCategoryPayload) (CreateProductCategoryResponse, error) {
	var response CreateProductCategoryResponse

	if errs := validation.Validate(payload); errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	category, err := s.repository.CreateProductCategory(ctx, dao.CreateProductCategoryParams{
		Name:       payload.Name,
		BusinessID: payload.BusinessID,
		CreatedBy:  payload.Initiator,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create product category")
		return response, InternalError
	}

	response.ID = category.ID
	response.Name = category.Name

	return response, nil
}

func (s *productCategoryService) UpdateProductCategory(ctx context.Context, payload UpdateProductCategoryPayload) (UpdateProductCategoryResponse, error) {
	var response UpdateProductCategoryResponse

	if errs := validation.Validate(payload); errs != nil {
		return response, errs
	}

	category, err := s.repository.SetProductCategoryNameByID(ctx, dao.SetProductCategoryNameByIDParams{
		ID:        payload.ID,
		Name:      payload.Name,
		UpdatedBy: &payload.Initiator,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update product category")
		return response, err
	}

	response.ID = category.ID
	response.Name = category.Name

	return response, nil
}

func (s *productCategoryService) ListProductCategories(ctx context.Context, payload ListProductCategoryPayload) ([]ListProductCategoryResponse, error) {
	response := []ListProductCategoryResponse{}

	if errs := validation.Validate(payload); errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}
	// TODO: bring it from constants
	if payload.Limit == 0 {
		payload.Limit = 10
	}
	if payload.Page == 0 {
		payload.Page = 1
	}

	categories, err := s.repository.ListProductCategoriesByBusinessID(ctx, dao.ListProductCategoriesByBusinessIDParams{
		BusinessID: payload.BusinessID,
		Limit:      payload.Limit,
		Offset:     (payload.Page - 1) * payload.Limit,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to list product categories")
		return response, err
	}

	for _, category := range categories {
		response = append(response, ListProductCategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return response, nil
}
