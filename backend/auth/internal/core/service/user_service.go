package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aritradevelops/billbharat/backend/auth/internal/core/cryptoutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
	"github.com/google/uuid"
)

const (
	InvitationLifetime = 15 * 24 * time.Hour // 15 days
)

type UserService interface {
	Profile(ctx context.Context, initiator string, payload ProfilePayload) (ProfileResponse, error)
	UpdateDP(ctx context.Context, initiator string, payload UpdateDPPayload) (ProfileResponse, error)
	Invite(ctx context.Context, initiator string, businessId string, payload InvitePayload) (InviteResponse, error)
}

type ProfilePayload struct {
	ID string `json:"id"`
}

type UpdateDPPayload struct {
	ID string  `json:"id"`
	Dp *string `json:"dp"`
}

type ProfileResponse struct {
	HumanID string  `json:"human_id"`
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Dp      *string `json:"dp"`
	Phone   string  `json:"phone"`
}

type InvitePayload struct {
	Name        string `json:"name" validate:"min=3,alphaspace,max=255"`
	Email       string `json:"email" validate:"email"`
	CountryCode string `json:"country_code" validate:"required"`
	Phone       string `json:"phone" validate:"numeric,min=10,max=16"`
	Origin      string `json:"origin" validate:"required"`
}

type InviteResponse struct {
}

type userService struct {
	repository   repository.Repository
	eventManager events.EventManager
}

func NewUserService(repository repository.Repository, eventManager events.EventManager) UserService {
	return &userService{repository: repository, eventManager: eventManager}
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
		Dp:      user.Dp,
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
		Dp: payload.Dp,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update user dp")
		return response, UserNotFoundErr
	}
	err = s.eventManager.EmitManageUserEvent(ctx, events.NewUserManageEvent("update", events.ManageUserEventPayload(user)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage user event")
		return response, InternalError
	}
	response = ProfileResponse{
		HumanID: user.HumanID,
		Email:   user.Email,
		Name:    user.Name,
		Dp:      user.Dp,
		Phone:   user.Phone,
	}
	return response, nil
}

func (s *userService) Invite(ctx context.Context, initiator string, businessId string, payload InvitePayload) (InviteResponse, error) {
	var response InviteResponse
	if err := validation.Validate(payload); err != nil {
		return response, err
	}

	invitationHash, err := cryptoutil.GenerateInvitationHash()

	if err != nil {
		return response, InternalError
	}

	invitation, err := s.repository.CreateInvitation(ctx, dao.CreateInvitationParams{
		Name:       payload.Name,
		Email:      payload.Email,
		Phone:      payload.CountryCode + payload.Phone,
		BusinessID: uuid.MustParse(businessId),
		Hash:       invitationHash,
		ExpiresAt:  time.Now().Add(InvitationLifetime),
	})
	go func(invitation dao.Invitation) {
		business, err := s.repository.FindBusinessById(ctx, invitation.BusinessID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to find business by id")
			return
		}

		err = s.eventManager.EmitManageNotificationEvent(ctx, events.NewNotificationManageEvent(events.ManageNotificationEventPayload{
			Event: notification.USER_INVITED,
			Kind:  notification.P2P,
			Payload: []events.NotificationChannelPayload{
				{Channel: notification.EMAIL, Data: notification.NewEmail(invitation.Email)},
				{Channel: notification.SMS, Data: notification.NewSMS(invitation.Phone)},
			},
			Tokens: map[string]string{
				"InvitationURL": fmt.Sprintf("%s/invites/%s", payload.Origin, invitationHash),
				"Email":         invitation.Email,
				"Name":          invitation.Name,
				"Phone":         invitation.Phone,
				"BusinessID":    invitation.BusinessID.String(),
				"ExpiresAt":     invitation.ExpiresAt.Format("2006-01-02 15:04:05"),
				"BusinessName":  business.Name,
			},
		}))
		if err != nil {
			logger.Info().Msgf("Failed to send notification to %s", invitation.Name)
		}
	}(invitation)

	return response, nil
}

func checkAccess(initiator string, target string) bool {
	return true
}
