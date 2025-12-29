package service

import (
	"context"

	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService interface {
	Profile(ctx context.Context, initiator string, payload ProfilePayload) (ProfileResponse, error)
	UpdateDP(ctx context.Context, initiator string, payload UpdateDPPayload) (ProfileResponse, error)
}

type ProfilePayload struct {
	ID string `json:"id"`
}

type UpdateDPPayload struct {
	ID string `json:"id"`
	Dp string `json:"dp"`
}

type ProfileResponse struct {
	HumanID string `json:"human_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Dp      string `json:"dp"`
	Phone   string `json:"phone"`
}

type userService struct {
	repository repository.Repository
}

func NewUserService(repository repository.Repository) UserService {
	return &userService{repository: repository}
}

func (s *userService) Profile(ctx context.Context, initiator string, payload ProfilePayload) (ProfileResponse, error) {
	var response ProfileResponse
	if !checkAccess(initiator, payload.ID) {
		return response, UserNotFoundErr
	}
	user, err := s.repository.FindUserById(ctx, uuid.MustParse(payload.ID))
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user by id")
		return response, UserNotFoundErr
	}
	response = ProfileResponse{
		HumanID: user.HumanID,
		Email:   user.Email,
		Name:    user.Name,
		Dp:      user.Dp.String,
		Phone:   user.Phone,
	}
	return response, nil
}

func (s *userService) UpdateDP(ctx context.Context, initiator string, payload UpdateDPPayload) (ProfileResponse, error) {
	var response ProfileResponse
	if !checkAccess(initiator, payload.ID) {
		return response, UserNotFoundErr
	}

	user, err := s.repository.UpdateUserDP(ctx, dao.UpdateUserDPParams{
		ID: uuid.MustParse(payload.ID),
		Dp: pgtype.Text{
			String: payload.Dp,
			Valid:  true,
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update user dp")
		return response, UserNotFoundErr
	}
	response = ProfileResponse{
		HumanID: user.HumanID,
		Email:   user.Email,
		Name:    user.Name,
		Dp:      payload.Dp,
		Phone:   user.Phone,
	}
	return response, nil
}

func checkAccess(initiator string, target string) bool {
	return true
}
