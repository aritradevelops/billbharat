package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aritradevelops/billbharat/backend/auth/internal/core/cryptoutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/validation"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost                 = 10
	VerificationRequestExpiry   = 15 * time.Minute
	VerificationRequestWaitTime = 1 * time.Minute
	SessionExpiry               = 30 * 24 * time.Hour
)

var (
	RootUserID    uuid.UUID = uuid.Nil
	UserExistsErr           = &ServiceError{HttpErrorCode: http.StatusConflict,
		DevErrorCode: "auth_001", Short: "user.exists", Long: "user already exists"}
	UserNotFoundErr = &ServiceError{HttpErrorCode: http.StatusNotFound,
		DevErrorCode: "auth_002", Short: "user.not_found", Long: "user not found"}
	VerificationRequestExpiredErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_003", Short: "verification_request.expired", Long: "verification request expired"}
	InvalidVerificationCodeErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_004", Short: "invalid_verification_code", Long: "invalid verification code"}
	UserEmailNotVerifiedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_005", Short: "user.email_not_verified", Long: "user email not verified"}
	UserPhoneNotVerifiedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_006", Short: "user.phone_not_verified", Long: "user phone not verified"}
	UserDeactivatedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_007", Short: "user.deactivated", Long: "user deactivated"}
	InvalidCredentialsErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_008", Short: "user.invalid_credentials", Long: "invalid credentials"}
	InvalidLoginMethodErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_009", Short: "user.invalid_login_method", Long: "invalid login method"}
	UserEmailVerifiedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_010", Short: "user.email_verified", Long: "user email verified"}
	TooManyVerificationRequestsErr = &ServiceError{HttpErrorCode: http.StatusTooManyRequests,
		DevErrorCode: "auth_011", Short: "auth.too_many_verification_requests", Long: "too many verification requests"}
	UserPhoneVerifiedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_012", Short: "user.phone_verified", Long: "user phone verified"}
	PasswordAlreadyUsedErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_013", Short: "password.already_used", Long: "password already used"}
	PasswordMismatchErr = &ServiceError{HttpErrorCode: http.StatusBadRequest,
		DevErrorCode: "auth_014", Short: "password.mismatch", Long: "password mismatch"}
)

type AuthService interface {
	Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, error)
	Login(ctx context.Context, payload LoginPayload) (LoginResponse, error)
	ForgotPassword(ctx context.Context, payload ForgotPasswordPayload) (ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, payload ResetPasswordPayload) (ResetPasswordResponse, error)
	ChangePassword(ctx context.Context, initiator string, payload ChangePasswordPayload) (ChangePasswordResponse, error)
	VerifyEmail(ctx context.Context, paylaod VerifyEmailPayload) (VerifyEmailResponse, error)
	VerifyPhone(ctx context.Context, payload VerifyPhonePayload) (VerifyPhoneResponse, error)
	SendEmailVerificationRequest(ctx context.Context, payload SendEmailVerificationRequestPayload) (SendEmailVerificationRequestResponse, error)
	SendPhoneVerificationRequest(ctx context.Context, payload SendPhoneVerificationRequestPayload) (SendPhoneVerificationRequestResponse, error)
}

type RegisterPayload struct {
	Name        string `json:"name" validate:"alphaspace,min=3,max=255"`
	Email       string `json:"email" validate:"email"`
	CountryCode string `json:"country_code" validate:"required"`
	Phone       string `json:"phone" validate:"numeric,min=10,max=16"`
	Password    string `json:"password" validate:"min=8,max=255"`
}

type RegisterResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"alphaspace,min=3,max=255"`
	Email       string    `json:"email" validate:"email"`
	CountryCode string    `json:"country_code" validate:"required"`
	Phone       string    `json:"phone" validate:"numeric,min=10,max=16"`
}

type VerifyEmailPayload struct {
	Email string `json:"email" validate:"email"`
	Code  string `json:"code" validate:"required,min=6,max=6"`
}

type VerifyEmailResponse struct {
}

type VerifyPhonePayload struct {
	Email string `json:"email" validate:"email"`
	Code  string `json:"code" validate:"required,min=6,max=6"`
}

type VerifyPhoneResponse struct {
}

type LoginPayload struct {
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=8,max=255"`
	UserIP    string `json:"user_ip" validate:"required"`
	UserAgent string `json:"user_agent" validate:"required"`
}

type LoginResponse struct {
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenLifetime  time.Time `json:"access_token_lifetime"`
	RefreshTokenLifetime time.Time `json:"refresh_token_lifetime"`
	BusinessFound        bool      `json:"business_found"`
}

type SendEmailVerificationRequestPayload struct {
	Email string `json:"email" validate:"email"`
}

type SendEmailVerificationRequestResponse struct {
}

type SendPhoneVerificationRequestPayload struct {
	Email string `json:"email" validate:"email"`
}

type SendPhoneVerificationRequestResponse struct {
}

type ForgotPasswordPayload struct {
	Email string `json:"email" validate:"email"`
}

type ForgotPasswordResponse struct {
}

type ResetPasswordPayload struct {
	Email           string `json:"email" validate:"email"`
	Code            string `json:"code" validate:"required,min=6,max=6"`
	Password        string `json:"password" validate:"min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password"`
}

type ResetPasswordResponse struct {
}

type ChangePasswordPayload struct {
	Email           string `json:"email" validate:"email"`
	CurrentPassword string `json:"current_password" validate:"min=8,max=255"`
	NewPassword     string `json:"new_password" validate:"min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=NewPassword"`
}

type ChangePasswordResponse struct {
}

type authService struct {
	repository   repository.Repository
	eventManager events.EventManager
	jwtManager   *jwtutil.JwtManager
}

func NewAuthService(repository repository.Repository, jwtManager *jwtutil.JwtManager, eventManager events.EventManager) AuthService {
	return &authService{
		repository:   repository,
		jwtManager:   jwtManager,
		eventManager: eventManager,
	}
}

// TODO: send otp to email and phone
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

	if err := tx.Commit(ctx); err != nil {
		logger.Error().Err(err).Msg("failed to commit transaction")
		return response, err
	}

	err = s.eventManager.EmitManageUserEvent(ctx, events.NewUserManageEvent("create", events.ManageUserEventPayload(user)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage user event")
		return response, err
	}

	response.ID = user.ID
	response.Name = user.Name
	response.Email = user.Email
	response.Phone = user.Phone
	return response, nil
}

// TODO: send welcome to email
func (s *authService) VerifyEmail(ctx context.Context, payload VerifyEmailPayload) (VerifyEmailResponse, error) {
	var response VerifyEmailResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	verificationRequest, err := s.repository.FindVerificationRequestByUserIdAndType(ctx, dao.FindVerificationRequestByUserIdAndTypeParams{
		UserID: user.ID,
		Type:   dao.VerificationTypeEmail,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find verification request")
		return response, VerificationRequestExpiredErr
	}

	if verificationRequest.Code != payload.Code {
		logger.Error().Err(err).Msg("invalid verification code")
		return response, InvalidVerificationCodeErr
	}

	user, err = s.repository.SetUserEmailVerified(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to set user email verified")
		return response, err
	}
	err = s.repository.SetVerificationRequestConsumedAt(ctx, verificationRequest.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to set verification request consumed at")
		return response, err
	}
	err = s.eventManager.EmitManageUserEvent(ctx, events.NewUserManageEvent("update", events.ManageUserEventPayload(user)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage user event")
		return response, InternalError
	}
	err = s.eventManager.EmitManageNotificationEvent(ctx, events.NewNotificationManageEvent(events.ManageNotificationEventPayload{
		Event: notification.EMAIL_VERIFIED,
		Kind:  notification.P2P,
		Payload: []events.NotificationChannelPayload{
			{
				Channel: notification.EMAIL,
				Data:    notification.NewEmail(user.Email),
			},
		},
	}))

	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage notification event")
		return response, InternalError
	}

	return response, nil
}

// TODO: send welcome to phone
func (s *authService) VerifyPhone(ctx context.Context, payload VerifyPhonePayload) (VerifyPhoneResponse, error) {
	var response VerifyPhoneResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	verificationRequest, err := s.repository.FindVerificationRequestByUserIdAndType(ctx, dao.FindVerificationRequestByUserIdAndTypeParams{
		UserID: user.ID,
		Type:   dao.VerificationTypePhone,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find verification request")
		return response, VerificationRequestExpiredErr
	}

	if verificationRequest.Code != payload.Code {
		logger.Error().Err(err).Msg("invalid verification code")
		return response, InvalidVerificationCodeErr
	}

	user, err = s.repository.SetUserPhoneVerified(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to set user phone verified")
		return response, err
	}
	err = s.repository.SetVerificationRequestConsumedAt(ctx, verificationRequest.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to set verification request consumed at")
		return response, err
	}
	err = s.eventManager.EmitManageUserEvent(ctx, events.NewUserManageEvent("update", events.ManageUserEventPayload(user)))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage user event")
		return response, InternalError
	}
	return response, nil
}

func (s *authService) Login(ctx context.Context, payload LoginPayload) (LoginResponse, error) {
	var response LoginResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	if user.DeletedAt != nil {
		logger.Error().Err(err).Msg("user deleted")
		return response, UserNotFoundErr
	}

	if user.EmailVerified == false {
		logger.Error().Err(err).Msg("user email not verified")
		return response, UserEmailNotVerifiedErr
	}

	if user.PhoneVerified == false {
		logger.Error().Err(err).Msg("user phone not verified")
		return response, UserPhoneNotVerifiedErr
	}

	if user.DeactivatedAt != nil {
		logger.Error().Err(err).Msg("user deactivated")
		return response, UserDeactivatedErr
	}

	password, err := s.repository.FindPasswordByUserId(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find password")
		return response, InvalidLoginMethodErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(payload.Password)); err != nil {
		logger.Error().Err(err).Msg("invalid password")
		return response, InvalidCredentialsErr
	}

	business, err := s.repository.FindBusinessesByUserID(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find business")
		return response, InternalError
	}
	var businessID string
	// if user has exactly one business, issue token for that business
	if len(business) == 1 {
		businessID = business[0].Business.ID.String()
		response.BusinessFound = true
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

	s.repository.CreateSession(ctx, dao.CreateSessionParams{
		HumanID:      cryptoutil.HumanID("session"),
		UserID:       user.ID,
		UserIp:       payload.UserIP,
		UserAgent:    payload.UserAgent,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(SessionExpiry),
		CreatedBy:    user.ID,
	})

	response.AccessToken = accessToken.Token
	response.AccessTokenLifetime = accessToken.Lifetime
	response.RefreshToken = refreshToken
	response.RefreshTokenLifetime = time.Now().Add(SessionExpiry)

	return response, nil
}

func (s *authService) SendEmailVerificationRequest(ctx context.Context, payload SendEmailVerificationRequestPayload) (SendEmailVerificationRequestResponse, error) {
	var response SendEmailVerificationRequestResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	if user.EmailVerified {
		logger.Error().Err(err).Msg("user email already verified")
		return response, UserEmailVerifiedErr
	}

	verificationRequest, err := s.repository.FindVerificationRequestByUserIdAndType(ctx, dao.FindVerificationRequestByUserIdAndTypeParams{
		UserID: user.ID,
		Type:   dao.VerificationTypeEmail,
	})
	if err == nil {
		if verificationRequest.CreatedAt.Add(VerificationRequestWaitTime).Before(time.Now()) {
			logger.Error().Err(err).Msg("too many verification requests")
			return response, TooManyVerificationRequestsErr
		}
	}

	otp, err := cryptoutil.GenerateOTP()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate otp for email")
		return response, InternalError
	}

	err = s.repository.CreateVerificationRequest(ctx, dao.CreateVerificationRequestParams{
		UserID:    user.ID,
		Code:      otp,
		Type:      dao.VerificationTypeEmail,
		ExpiresAt: time.Now().Add(VerificationRequestExpiry),
		CreatedBy: user.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create verification request")
		return response, InternalError
	}
	err = s.eventManager.EmitManageNotificationEvent(ctx, events.NewNotificationManageEvent(events.ManageNotificationEventPayload{
		Event: notification.EMAIL_VERIFICATION,
		Kind:  notification.P2P,
		Payload: []events.NotificationChannelPayload{
			{
				Channel: notification.EMAIL,
				Data:    notification.NewEmail(user.Email),
			},
		},
		Tokens: map[string]string{
			"Name":       user.Name,
			"Email":      user.Email,
			"Otp":        otp,
			"Expires_In": VerificationRequestExpiry.String(),
		},
	}))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage notification event")
		return response, err
	}

	return response, nil
}

func (s *authService) SendPhoneVerificationRequest(ctx context.Context, payload SendPhoneVerificationRequestPayload) (SendPhoneVerificationRequestResponse, error) {
	var response SendPhoneVerificationRequestResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	if user.PhoneVerified {
		logger.Error().Err(err).Msg("user phone already verified")
		return response, UserPhoneVerifiedErr
	}

	verificationRequest, err := s.repository.FindVerificationRequestByUserIdAndType(ctx, dao.FindVerificationRequestByUserIdAndTypeParams{
		UserID: user.ID,
		Type:   dao.VerificationTypePhone,
	})
	if err == nil {
		if verificationRequest.CreatedAt.Add(VerificationRequestWaitTime).Before(time.Now()) {
			logger.Error().Err(err).Msg("too many verification requests")
			return response, TooManyVerificationRequestsErr
		}
	}

	otp, err := cryptoutil.GenerateOTP()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate otp for phone")
		return response, InternalError
	}

	err = s.repository.CreateVerificationRequest(ctx, dao.CreateVerificationRequestParams{
		UserID:    user.ID,
		Code:      otp,
		Type:      dao.VerificationTypePhone,
		ExpiresAt: time.Now().Add(VerificationRequestExpiry),
		CreatedBy: user.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create verification request")
		return response, InternalError
	}
	err = s.eventManager.EmitManageNotificationEvent(ctx, events.NewNotificationManageEvent(events.ManageNotificationEventPayload{
		Event: notification.PHONE_VERIFICATION,
		Kind:  notification.P2P,
		Payload: []events.NotificationChannelPayload{
			{
				Channel: notification.SMS,
				Data:    notification.NewSMS(user.Phone),
			},
		},
		Tokens: map[string]string{
			"Name":       user.Name,
			"Phone":      user.Phone,
			"Otp":        otp,
			"Expires_In": VerificationRequestExpiry.String(),
		},
	}))
	if err != nil {
		logger.Error().Err(err).Msg("failed to emit manage notification event")
		return response, err
	}

	return response, nil
}

func (s *authService) ForgotPassword(ctx context.Context, payload ForgotPasswordPayload) (ForgotPasswordResponse, error) {
	var response ForgotPasswordResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	otp, err := cryptoutil.GenerateOTP()
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate otp for email")
		return response, InternalError
	}

	err = s.repository.CreateVerificationRequest(ctx, dao.CreateVerificationRequestParams{
		UserID:    user.ID,
		Code:      otp,
		Type:      dao.VerificationTypeResetPassword,
		ExpiresAt: time.Now().Add(VerificationRequestExpiry),
		CreatedBy: user.ID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to create verification request")
		return response, InternalError
	}

	return response, nil
}

// TODO: reset password email
func (s *authService) ResetPassword(ctx context.Context, payload ResetPasswordPayload) (ResetPasswordResponse, error) {
	var response ResetPasswordResponse
	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	errs = validation.ValidatePassword(payload.Password)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user")
		return response, UserNotFoundErr
	}

	verificationRequest, err := s.repository.FindVerificationRequestByUserIdAndType(ctx, dao.FindVerificationRequestByUserIdAndTypeParams{
		UserID: user.ID,
		Type:   dao.VerificationTypeResetPassword,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find verification request")
		return response, VerificationRequestExpiredErr
	}

	if verificationRequest.Code != payload.Code {
		logger.Error().Err(err).Msg("invalid verification code")
		return response, InvalidVerificationCodeErr
	}

	lastFourPasswords, err := s.repository.FindLastFourPasswordsByUserId(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find last four passwords")
		return response, InternalError
	}
	match := false
	for _, password := range lastFourPasswords {
		if err := bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(payload.Password)); err == nil {
			match = true
			break
		}
	}
	if match {
		logger.Error().Err(err).Msg("password already used")
		return response, PasswordAlreadyUsedErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate password hash")
		return response, InternalError
	}

	_, err = s.repository.DeletePassword(ctx, dao.DeletePasswordParams{
		UserID:    user.ID,
		DeletedBy: &user.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete password")
		return response, InternalError
	}

	err = s.repository.CreatePassword(ctx, dao.CreatePasswordParams{
		UserID:    user.ID,
		Password:  string(hashedPassword),
		CreatedBy: user.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create password")
		return response, InternalError
	}

	err = s.repository.SetVerificationRequestConsumedAt(ctx, verificationRequest.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to set verification request consumed at")
		return response, InternalError
	}
	return response, nil
}

func (s *authService) ChangePassword(ctx context.Context, initiator string, payload ChangePasswordPayload) (ChangePasswordResponse, error) {
	var response ChangePasswordResponse

	errs := validation.Validate(payload)
	if errs != nil {
		logger.Error().Err(errs).Msg("validation failed")
		return response, errs
	}

	errs = validation.ValidatePassword(payload.NewPassword)
	if errs != nil {
		logger.Error().Err(errs).Msg("password validation failed")
		return response, errs
	}

	user, err := s.repository.FindUserByEmail(ctx, payload.Email)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find user by email")
		return response, UserNotFoundErr
	}

	password, err := s.repository.FindPasswordByUserId(ctx, user.ID)

	if err != nil {
		logger.Error().Err(err).Msg("failed to find password by user id")
		return response, PasswordMismatchErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(payload.CurrentPassword)); err != nil {
		logger.Error().Err(err).Msg("password mismatch")
		return response, PasswordMismatchErr
	}

	lastFourPasswords, err := s.repository.FindLastFourPasswordsByUserId(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to find last four passwords")
		return response, InternalError
	}
	match := false
	for _, password := range lastFourPasswords {
		if err := bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(payload.NewPassword)); err == nil {
			match = true
			break
		}
	}
	if match {
		logger.Error().Err(err).Msg("password already used")
		return response, PasswordAlreadyUsedErr
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to generate password hash")
		return response, InternalError
	}

	_, err = s.repository.DeletePassword(ctx, dao.DeletePasswordParams{
		UserID:    user.ID,
		DeletedBy: &user.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete password")
		return response, InternalError
	}

	err = s.repository.CreatePassword(ctx, dao.CreatePasswordParams{
		UserID:    user.ID,
		Password:  string(hashedPassword),
		CreatedBy: user.ID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("failed to create password")
		return response, InternalError
	}

	return response, nil
}

// func (s *authService) Profile(ctx context.Context, initiator string, payload ProfilePayload) (ProfileResponse, error) {
// 	var response ProfileResponse
// 	user, err := s.repository.FindUserById(ctx, uuid.MustParse(initiator))
// 	if err != nil {
// 		logger.Error().Err(err).Msg("failed to find user by id")
// 		return response, UserNotFoundErr
// 	}
// 	response = ProfileResponse{
// 		HumanID: user.HumanID,
// 		Email:   user.Email,
// 		Name:    user.Name,
// 		Dp:      user.Dp,
// 		Phone:   user.Phone,
// 	}
// 	return response, nil
// }
