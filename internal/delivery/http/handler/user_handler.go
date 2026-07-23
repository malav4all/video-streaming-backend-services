package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourorg/go-user-service/internal/delivery/http/response"
	"github.com/yourorg/go-user-service/internal/domain/user"
	appvalidator "github.com/yourorg/go-user-service/pkg/validator"
)

// UserHandler handles HTTP requests for the user module and delegates
// all business logic to the injected user.Service.
type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req user.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if errs := appvalidator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation failed", errs)
	}

	res, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create user", err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "user created successfully", res)
}

// GetByID handles GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id", nil)
	}

	res, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "user not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "user fetched successfully", res)
}

// GetAll handles GET /api/v1/users?page=&page_size=
func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	users, total, err := h.service.GetAll(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to fetch users", err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "users fetched successfully", users, &response.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	})
}

// Update handles PUT /api/v1/users/:id
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id", nil)
	}

	var req user.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if errs := appvalidator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation failed", errs)
	}

	res, err := h.service.Update(c.Context(), uint(id), &req)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "user not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "user updated successfully", res)
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id", nil)
	}

	if err := h.service.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(c, fiber.StatusNotFound, "user not found", nil)
	}

	return response.Success(c, fiber.StatusOK, "user deleted successfully", nil)
}
