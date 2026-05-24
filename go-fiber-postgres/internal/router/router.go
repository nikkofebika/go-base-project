package router

import (
	"go-fiber-postgres/internal/config"
	"go-fiber-postgres/internal/pkg/validator"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type RouterDependencies struct {
	AppConfig *config.AppConfig
	DB        *gorm.DB
	Validator *validator.Validator
}

func Register(router fiber.Router, deps *RouterDependencies) {
	registerRouterApi(router, deps)
}
