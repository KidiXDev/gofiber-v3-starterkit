package routes

import (
	"gofiber-starterkit/app/api/controllers"
	"gofiber-starterkit/pkg/middlewares"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(
	app *fiber.App,
	userController *controllers.UserController,
	authMiddleware *middlewares.AuthMiddleware,
) {
	const ApiVersion = "/api/v1"

	api := app.Group(ApiVersion)

	auth := api.Group("/auth")
	auth.Post("/register", userController.Register)
	auth.Post("/login", userController.Login)
	auth.Post("/refresh", userController.Refresh)

	protected := api.Group("")
	protected.Use(authMiddleware.AuthRequired())

	protected.Get("/auth/me", userController.Me)
	protected.Put("/auth/me", userController.UpdateProfile)
	protected.Post("/auth/logout", userController.Logout)
	protected.Post("/auth/logout-all", userController.LogoutAll)

	users := protected.Group("/users")
	users.Get("", userController.List)
	users.Get("/:id", userController.Get)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)
}
