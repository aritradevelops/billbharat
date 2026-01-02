package httpd

import "github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/authn"

func (s *Server) SetupRoutes() {
	router := s.app
	authMiddleware := authn.Middleware(s.jwtManager)
	router.Get("/api/v1/auth-srv/health", s.handlers.Health)

	// Authentication routes
	router.Post("/api/v1/auth-srv/auth/register", s.handlers.Auth.Register)
	router.Post("/api/v1/auth-srv/auth/login", s.handlers.Auth.Login)
	router.Post("/api/v1/auth-srv/auth/forgot-password", s.handlers.Auth.ForgotPassword)
	router.Post("/api/v1/auth-srv/auth/reset-password", s.handlers.Auth.ResetPassword)
	router.Post("/api/v1/auth-srv/auth/verify-email", s.handlers.Auth.VerifyEmail)
	router.Post("/api/v1/auth-srv/auth/verify-phone", s.handlers.Auth.VerifyPhone)
	router.Post("/api/v1/auth-srv/auth/send-email-verification-request", s.handlers.Auth.SendEmailVerificationRequest)
	router.Post("/api/v1/auth-srv/auth/send-phone-verification-request", s.handlers.Auth.SendPhoneVerificationRequest)
	router.Post("/api/v1/auth-srv/auth/change-password", authMiddleware, s.handlers.Auth.ChangePassword)

	// User routes
	router.Get("/api/v1/auth-srv/users/profile/:id", authMiddleware, s.handlers.User.Profile)
	router.Post("/api/v1/auth-srv/users/change-profile-picture", authMiddleware, s.handlers.User.UpdateDP)
	router.Post("/api/v1/auth-srv/users/invite", authMiddleware, s.handlers.User.Invite)

	// Business routes
	router.Post("/api/v1/auth-srv/businesses/create", authMiddleware, s.handlers.Business.Create)
	router.Get("/api/v1/auth-srv/businesses/list", authMiddleware, s.handlers.Business.List)
	router.Post("/api/v1/auth-srv/businesses/select/:business_id", authMiddleware, s.handlers.Business.Select)
}
