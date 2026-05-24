package main

import (
	"go-fiber-postgres/database/seeders"
	"go-fiber-postgres/internal/app/exception"
	"go-fiber-postgres/internal/config"
	"go-fiber-postgres/internal/pkg/validator"
	"go-fiber-postgres/internal/router"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg := config.NewAppConfig()

	// set logger

	// init database
	db := config.NewDatabase(cfg)

	// seeder
	args := os.Args
	if len(args) > 1 && args[1] == "--seed" {
		if err := seeders.Seed(db); err != nil {
			panic(err)
		}
		return
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       cfg.AppName,
		BodyLimit:     10 * 1024 * 1024, // 10 MB
		ErrorHandler:  exception.ErrorHandler,
	})

	allowedOrigins := make(map[string]bool)
	for _, origin := range cfg.CorsOrigins {
		allowedOrigins[origin] = true
	}
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// init validator
	validator := validator.NewValidator(db)

	app.Get("/", func(c fiber.Ctx) error {
		// c.SendString("Hello, World!")
		return fiber.NewError(404, "Custom error message")
	})

	router.Register(app, &router.RouterDependencies{
		AppConfig: cfg,
		DB:        db,
		Validator: validator,
	})

	app.Listen(":3000", fiber.ListenConfig{
		EnablePrefork:         cfg.AppPrefork,
		DisableStartupMessage: cfg.AppEnv == "production",
	})
}
