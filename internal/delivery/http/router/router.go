package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourorg/go-user-service/internal/delivery/http/handler"
	"github.com/yourorg/go-user-service/internal/delivery/http/middleware"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	tenantHandler *handler.TenantHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "service is healthy"})
	})

	// Keycloak OAuth2/OIDC flow — public routes.
	authRoutes := api.Group("/auth")
	authRoutes.Get("/login", authHandler.Login)
	authRoutes.Get("/callback", authHandler.Callback)
	authRoutes.Post("/logout", authHandler.Logout)
	authRoutes.Get("/me", authMiddleware.RequireAuth(), authHandler.Me)
	// User CRUD — now requires an authenticated Keycloak session.
	users := api.Group("/users", authMiddleware.RequireAuth())
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.GetAll)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
	// ...

	tenants := api.Group("/tenants", authMiddleware.RequireAuth())
	tenants.Post("/", tenantHandler.Create)
	tenants.Get("/", tenantHandler.GetAll)
	tenants.Get("/:id", tenantHandler.GetByID)
	tenants.Put("/:id", tenantHandler.Update)
	tenants.Delete("/:id", tenantHandler.Delete)
}
