package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/yourorg/go-user-service/internal/config"
	"github.com/yourorg/go-user-service/internal/delivery/http/handler"
	"github.com/yourorg/go-user-service/internal/delivery/http/router"
	"github.com/yourorg/go-user-service/internal/domain/user"
	"github.com/yourorg/go-user-service/internal/repository/postgres"
	userUsecase "github.com/yourorg/go-user-service/internal/usecase/user"
)

func main() {
	cfg := config.Load()
	log.Printf("DEBUG cfg: env=%s host=%s port=%s user=%s dbname=%s", cfg.AppEnv, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// AutoMigrate is convenient for development. For production, prefer
	// versioned SQL migrations (see /migrations) run via a migration tool.
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Wiring: repository (adapter) -> usecase (business logic) -> handler (delivery)
	userRepo := postgres.NewUserRepository(db)
	userSvc := userUsecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

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

	router.SetupRoutes(app, userHandler)

	log.Printf("server starting on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
