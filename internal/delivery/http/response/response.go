package response

import "github.com/gofiber/fiber/v2"

// Meta carries pagination info for list endpoints.
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
}

// SuccessResponse is the standard envelope for successful responses.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorResponse is the standard envelope for error responses.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success writes a standard success response.
func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta writes a standard success response including pagination meta.
func SuccessWithMeta(c *fiber.Ctx, status int, message string, data interface{}, meta *Meta) error {
	return c.Status(status).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error writes a standard error response.
func Error(c *fiber.Ctx, status int, message string, errs interface{}) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}
