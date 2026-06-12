package router

import (
	"go-fiber-mysql/internal/app/controllers"
	"go-fiber-mysql/internal/app/enums"
	"go-fiber-mysql/internal/app/middlewares"
	"go-fiber-mysql/internal/app/repositories"
	"go-fiber-mysql/internal/app/services"

	"github.com/gofiber/fiber/v3"
)

func registerRouterApi(router fiber.Router, deps *RouterDependencies) {
	userRepository := repositories.NewUserRepository(deps.DB)
	userService := services.NewUserService(userRepository)

	authRouter(router, deps, userService)
	usersRouter(router, deps, userService)
	rolesRouter(router, deps, userService)
	permissionsRouter(router, deps, userService)
}

func authRouter(router fiber.Router, deps *RouterDependencies, userService services.UserService) {
	emailService := services.NewEmailService(deps.AppConfig.MailHost, deps.AppConfig.MailPort, deps.AppConfig.MailUser, deps.AppConfig.MailPass, deps.AppConfig.MailFrom)

	tokenRepository := repositories.NewTokenRepository(deps.DB)
	tokenService := services.NewTokenService(tokenRepository, deps.DB)

	userRepository := repositories.NewUserRepository(deps.DB)
	service := services.NewAuthService(deps.AppConfig, emailService, userRepository, tokenService)
	controller := controllers.NewAuthController(service, deps.Validator)

	r := router.Group("auth")
	r.Post("forgot-password", middlewares.AuthMiddleware(deps.AppConfig), controller.ForgotPassword)
	r.Get("me", middlewares.AuthMiddleware(deps.AppConfig), controller.Me)
	r.Post("refresh-token", middlewares.AuthMiddleware(deps.AppConfig), controller.RefreshToken)
	r.Post("reset-password", middlewares.AuthMiddleware(deps.AppConfig), controller.ResetPassword)
	r.Post("token", controller.Token)
}

func usersRouter(router fiber.Router, deps *RouterDependencies, service services.UserService) {
	controller := controllers.NewUserController(service, deps.Validator)

	r := router.Group("users", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", middlewares.PermissionMiddleware(service, enums.PermissionUserRead), controller.FindAll)
	r.Get("/:id", middlewares.PermissionMiddleware(service, enums.PermissionUserRead), controller.FindOne)
	r.Post("", middlewares.PermissionMiddleware(service, enums.PermissionUserCreate), controller.Create)
	r.Patch("/:id", middlewares.PermissionMiddleware(service, enums.PermissionUserUpdate), controller.Update)
	r.Delete("/:id", middlewares.PermissionMiddleware(service, enums.PermissionUserDelete), controller.Delete)
	r.Patch("/:id/restore", middlewares.PermissionMiddleware(service, enums.PermissionUserRestore), controller.Restore)
	r.Delete("/:id/force", middlewares.PermissionMiddleware(service, enums.PermissionUserForceDelete), controller.ForceDelete)
	r.Post("/:id/roles", middlewares.PermissionMiddleware(service, enums.PermissionUserSyncRoles), controller.SyncRoles)
}

func rolesRouter(router fiber.Router, deps *RouterDependencies, userService services.UserService) {
	repository := repositories.NewRoleRepository(deps.DB)
	service := services.NewRoleService(repository)
	controller := controllers.NewRoleController(service, deps.Validator)

	r := router.Group("roles", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", middlewares.PermissionMiddleware(userService, enums.PermissionRoleRead), controller.FindAll)
	r.Get("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionRoleRead), controller.FindOne)
	r.Post("", middlewares.PermissionMiddleware(userService, enums.PermissionRoleCreate), controller.Create)
	r.Patch("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionRoleUpdate), controller.Update)
	r.Delete("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionRoleDelete), controller.Delete)
	r.Patch("/:id/restore", middlewares.PermissionMiddleware(userService, enums.PermissionRoleRestore), controller.Restore)
	r.Delete("/:id/force", middlewares.PermissionMiddleware(userService, enums.PermissionRoleForceDelete), controller.ForceDelete)
	r.Post("/:id/permissions", middlewares.PermissionMiddleware(userService, enums.PermissionRoleSyncPermissions), controller.SyncPermissions)
}

func permissionsRouter(router fiber.Router, deps *RouterDependencies, userService services.UserService) {
	repository := repositories.NewPermissionRepository(deps.DB)
	service := services.NewPermissionService(repository)
	controller := controllers.NewPermissionController(service, deps.Validator)

	r := router.Group("permissions", middlewares.AuthMiddleware(deps.AppConfig))
	r.Get("", middlewares.PermissionMiddleware(userService, enums.PermissionPermissionRead), controller.FindAll)
	r.Get("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionPermissionRead), controller.FindOne)
	r.Post("", middlewares.PermissionMiddleware(userService, enums.PermissionRoleCreate), controller.Create)
	r.Patch("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionRoleUpdate), controller.Update)
	r.Delete("/:id", middlewares.PermissionMiddleware(userService, enums.PermissionRoleDelete), controller.Delete)
	r.Patch("/:id/restore", middlewares.PermissionMiddleware(userService, enums.PermissionRoleRestore), controller.Restore)
	r.Delete("/:id/force", middlewares.PermissionMiddleware(userService, enums.PermissionRoleForceDelete), controller.ForceDelete)
}
