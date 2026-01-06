package httpd

import "github.com/aritradevelops/billbharat/backend/product/internal/ports/httpd/authn"

func (s *Server) SetupRoutes() {
	router := s.app
	authMiddleware := authn.Middleware(s.jwtManager)
	router.Get("/api/v1/product-srv/health", s.handlers.Health)
	router.Get("/api/v1/product-srv/product-categories/list", authMiddleware, s.handlers.Category.ListProductCategories)
	router.Post("/api/v1/product-srv/product-categories/create", authMiddleware, s.handlers.Category.CreateProductCategory)
	router.Put("/api/v1/product-srv/product-categories/update/:id", authMiddleware, s.handlers.Category.UpdateProductCategory)
}
