package routes

import (
	"gofiber-starterkit/app/api/controllers"
	"gofiber-starterkit/app/api/services"
	"gofiber-starterkit/pkg/client/redis"
	"gofiber-starterkit/pkg/client/s3"
	"gofiber-starterkit/pkg/middlewares"

	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"
)

func RegisterRoutes(app *fiber.App, db *bun.DB, redisClient *redis.RedisClient, s3Client *s3.S3Client) {
	const ApiVersion = "/api/v1"

	userService := services.NewUserService(db, redisClient, s3Client)
	authMiddleware := middlewares.NewAuthMiddleware(userService, redisClient)

	userController := controllers.NewUserController(userService)

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
