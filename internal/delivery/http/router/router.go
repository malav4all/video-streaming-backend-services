package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourorg/go-user-service/internal/delivery/http/handler"
)

// SetupRoutes registers all API routes. As new modules are added
// (e.g. products, orders), give each its own group here, mirroring
// the user module's pattern.
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "service is healthy"})
	})

	users := api.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Future modules go here, e.g.:
	// products := api.Group("/products")
	// products.Post("/", productHandler.Create)
}
