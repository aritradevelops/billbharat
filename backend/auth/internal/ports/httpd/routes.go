package httpd

func (s *Server) SetupRoutes() {
	router := s.app
	router.Get("/api/v1/auth-srv/health", s.handlers.Health)
	router.Post("/api/v1/auth-srv/auth/register", s.handlers.Auth.Register)
}
