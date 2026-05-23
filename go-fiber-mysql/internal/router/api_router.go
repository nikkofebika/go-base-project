package router

import (
	"go-fiber-mysql/internal/app/controllers"
	"go-fiber-mysql/internal/app/middlewares"
	"go-fiber-mysql/internal/app/repositories"
	"go-fiber-mysql/internal/app/services"

	"github.com/gofiber/fiber/v3"
)

func registerRouterApi(router fiber.Router, deps *RouterDependencies) {
	authRouter(router, deps)
	usersRouter(router, deps)
	rolesRouter(router, deps)
	permissionsRouter(router, deps)
}

func authRouter(router fiber.Router, deps *RouterDependencies) {
	emailService := services.NewEmailService(deps.AppConfig.MailHost, deps.AppConfig.MailPort, deps.AppConfig.MailUser, deps.AppConfig.MailPass, deps.AppConfig.MailFrom)

	tokenRepository := repositories.NewTokenRepository(deps.DB)
	tokenService := services.NewTokenService(tokenRepository, deps.DB)

	userRepository := repositories.NewUserRepository(deps.DB)
	service := services.NewAuthService(deps.AppConfig, emailService, userRepository, tokenService)
	controller := controllers.NewAuthController(service, deps.Validator)

	r := router.Group("auth")
	r.Post("token", controller.Token)
	r.Post("refresh-token", middlewares.AuthMiddleware(deps.AppConfig), controller.RefreshToken)
}

func usersRouter(router fiber.Router, deps *RouterDependencies) {
	repository := repositories.NewUserRepository(deps.DB)
	service := services.NewUserService(repository)
	controller := controllers.NewUserController(service, deps.Validator)

	r := router.Group("users", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", controller.FindAll)
	r.Get("/:id", controller.FindOne)
	r.Post("", controller.Create)
	r.Patch("/:id", controller.Update)
	r.Delete("/:id", controller.Delete)
	r.Patch("/:id/restore", controller.Restore)
	r.Delete("/:id/force", controller.ForceDelete)
	r.Post("/:id/roles", controller.SyncRoles)
}

func rolesRouter(router fiber.Router, deps *RouterDependencies) {
	repository := repositories.NewRoleRepository(deps.DB)
	service := services.NewRoleService(repository)
	controller := controllers.NewRoleController(service, deps.Validator)

	r := router.Group("roles", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", controller.FindAll)
	r.Get("/:id", controller.FindOne)
	r.Post("", controller.Create)
	r.Patch("/:id", controller.Update)
	r.Delete("/:id", controller.Delete)
	r.Patch("/:id/restore", controller.Restore)
	r.Delete("/:id/force", controller.ForceDelete)
	r.Post("/:id/permissions", controller.SyncPermissions)
}

func permissionsRouter(router fiber.Router, deps *RouterDependencies) {
	repository := repositories.NewPermissionRepository(deps.DB)
	service := services.NewPermissionService(repository)
	controller := controllers.NewPermissionController(service, deps.Validator)

	r := router.Group("permissions", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", controller.FindAll)
	r.Get("/:id", controller.FindOne)
	r.Post("", controller.Create)
	r.Patch("/:id", controller.Update)
	r.Delete("/:id", controller.Delete)
	r.Patch("/:id/restore", controller.Restore)
	r.Delete("/:id/force", controller.ForceDelete)
}
