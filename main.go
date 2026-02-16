package main

import (
	"gofiber-starterkit/app/api/controllers"
	"gofiber-starterkit/app/api/services"
	"gofiber-starterkit/app/routes"
	"gofiber-starterkit/app/shared"
	"gofiber-starterkit/pkg/client/db"
	"gofiber-starterkit/pkg/client/redis"
	"gofiber-starterkit/pkg/client/s3"
	"gofiber-starterkit/pkg/config"
	"gofiber-starterkit/pkg/middlewares"
	"gofiber-starterkit/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()

	// Provide infrastructure
	c.Provide(db.New)
	c.Provide(redis.New)
	c.Provide(s3.New)

	// Provide application layers
	c.Provide(services.NewUserService)
	c.Provide(controllers.NewUserController)
	c.Provide(middlewares.NewAuthMiddleware)

	c.Provide(func() *fiber.App {
		cfg := config.FiberConfig()
		cfg.ErrorHandler = shared.RespondError

		app := fiber.New(cfg)

		// Global Middlewares
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))
		middlewares.FiberMiddleware(app)

		// Health Check
		app.Get(healthcheck.LivenessEndpoint, healthcheck.New())

		return app
	})

	// Start Application
	c.Invoke(func(
		app *fiber.App,
		userController *controllers.UserController,
		authMiddleware *middlewares.AuthMiddleware,
		dbClient *bun.DB,
		redisClient *redis.RedisClient,
	) {
		// Register Routes
		routes.RegisterRoutes(app, userController, authMiddleware)

		// Lifecycle management
		defer dbClient.Close()
		defer redisClient.Client.Close()

		if err := utils.StartServerWithGracefulShutdown(app); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	})
}
