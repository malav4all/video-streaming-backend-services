package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/yourorg/go-user-service/internal/auth"
	"github.com/yourorg/go-user-service/internal/config"
	"github.com/yourorg/go-user-service/internal/delivery/http/handler"
	appmiddleware "github.com/yourorg/go-user-service/internal/delivery/http/middleware"
	"github.com/yourorg/go-user-service/internal/delivery/http/router"
	"github.com/yourorg/go-user-service/internal/domain/tenant"
	"github.com/yourorg/go-user-service/internal/domain/user"
	"github.com/yourorg/go-user-service/internal/repository/postgres"
	redisrepo "github.com/yourorg/go-user-service/internal/repository/redis"
	tenantUsecase "github.com/yourorg/go-user-service/internal/usecase/tenant"
	userUsecase "github.com/yourorg/go-user-service/internal/usecase/user"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// AutoMigrate is convenient for development. For production, prefer
	// versioned SQL migrations (see /migrations) run via a migration tool.
	if err := db.AutoMigrate(&user.User{}, &tenant.Tenant{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	redisClient := redisrepo.NewRedisClient(cfg)

	keycloakClient, err := auth.NewKeycloakClient(ctx, &auth.Config{
		BaseURL:      cfg.KeycloakBaseURL,
		Realm:        cfg.KeycloakRealm,
		ClientID:     cfg.KeycloakClientID,
		ClientSecret: cfg.KeycloakClientSecret,
		RedirectURL:  cfg.KeycloakRedirectURL,
	})
	if err != nil {
		log.Fatalf("failed to initialize keycloak client: %v", err)
	}

	// Wiring: repository (adapter) -> usecase (business logic) -> handler (delivery)
	userRepo := postgres.NewUserRepository(db)
	userSvc := userUsecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	tenantRepo := postgres.NewTenantRepository(db)
	tenantSvc := tenantUsecase.NewTenantUsecase(tenantRepo)
	tenantHandler := handler.NewTenantHandler(tenantSvc)

	stateStore := redisrepo.NewStateStore(redisClient)
	sessionStore := redisrepo.NewSessionStore(redisClient, cfg.SessionTTL)
	authHandler := handler.NewAuthHandler(keycloakClient, stateStore, sessionStore, cfg.SessionTTL)
	authMiddleware := appmiddleware.NewAuthMiddleware(keycloakClient, sessionStore, tenantSvc)

	app := fiber.New(fiber.Config{
		AppName:      "Video Streaming AI Platform Service",
		BodyLimit:    int(cfg.MaxRequestSize),
		ReadTimeout:  cfg.RequestTimeout,
		WriteTimeout: cfg.RequestTimeout,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.AllowedOrigins, ","),
		AllowCredentials: cfg.CORSAllowCredentials,
	}))

	router.SetupRoutes(app, userHandler, authHandler, tenantHandler, authMiddleware)

	log.Printf("server starting on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
