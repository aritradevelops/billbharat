package handlers

import (
	"net/http"

	"github.com/aritradevelops/billbharat/backend/product/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/product/internal/ports/httpd/authn"
	"github.com/aritradevelops/billbharat/backend/shared/translation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductCategoryHandler struct {
	service service.ProductCategoryService
}

func NewProductCategoryHandler(service service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		service: service,
	}
}

type CreateProductCategoryPayload struct {
	Name string `json:"name"`
}

type UpdateProductCategoryPayload struct {
	Name string `json:"name"`
}

type ListProductCategoriesQuery struct {
	Limit int `query:"limit"`
	Page  int `query:"page"`
}

func (h *ProductCategoryHandler) CreateProductCategory(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var payload CreateProductCategoryPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	category, err := h.service.CreateProductCategory(c.Context(), service.CreateProductCategoryPayload{
		Name:       payload.Name,
		BusinessID: uuid.MustParse(user.BusinessID),
		Initiator:  uuid.MustParse(user.UserID),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusCreated)
	return c.JSON(NewResponse(translation.Localize(c, "controller.create", map[string]string{
		"Entity": "Category",
	}), category, nil))
}

func (h *ProductCategoryHandler) UpdateProductCategory(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var payload UpdateProductCategoryPayload
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	category, err := h.service.UpdateProductCategory(c.Context(), service.UpdateProductCategoryPayload{
		ID:         uuid.MustParse(c.Params("id")),
		Name:       payload.Name,
		BusinessID: uuid.MustParse(user.BusinessID),
		Initiator:  uuid.MustParse(user.UserID),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusOK)
	return c.JSON(NewResponse(translation.Localize(c, "controller.update", map[string]string{
		"Entity": "Category",
	}), category, nil))
}

func (h *ProductCategoryHandler) ListProductCategories(c *fiber.Ctx) error {
	user, err := authn.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var query ListProductCategoriesQuery
	if err := c.QueryParser(&query); err != nil {
		return err
	}

	categories, err := h.service.ListProductCategories(c.Context(), service.ListProductCategoryPayload{
		BusinessID: uuid.MustParse(user.BusinessID),
		Page:       1,
		Limit:      10,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusOK)
	return c.JSON(NewResponse(translation.Localize(c, "controller.list", map[string]string{
		"Entity": "Categories",
	}), categories, nil))
}
