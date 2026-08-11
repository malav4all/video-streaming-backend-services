package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-user-service/internal/delivery/http/response"
	"github.com/yourorg/go-user-service/internal/domain/tenant"
	appvalidator "github.com/yourorg/go-user-service/pkg/validator"
)

type TenantHandler struct {
	service tenant.Service
}

func NewTenantHandler(service tenant.Service) *TenantHandler {
	return &TenantHandler{service: service}
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var req tenant.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if errs := appvalidator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation failed", errs)
	}

	res, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create tenant", err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "tenant created successfully", res)
}

func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	res, err := h.service.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "tenant not found", nil)
	}
	return response.Success(c, fiber.StatusOK, "tenant fetched successfully", res)
}

func (h *TenantHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	tenants, total, err := h.service.GetAll(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to fetch tenants", err.Error())
	}
	return response.SuccessWithMeta(c, fiber.StatusOK, "tenants fetched successfully", tenants, &response.Meta{
		Page: page, PageSize: pageSize, TotalCount: total,
	})
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	var req tenant.UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if errs := appvalidator.Validate(req); errs != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation failed", errs)
	}

	res, err := h.service.Update(c.Context(), c.Params("id"), &req)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "tenant not found", nil)
	}
	return response.Success(c, fiber.StatusOK, "tenant updated successfully", res)
}

func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id")); err != nil {
		return response.Error(c, fiber.StatusNotFound, "tenant not found", nil)
	}
	return response.Success(c, fiber.StatusOK, "tenant deleted successfully", nil)
}
