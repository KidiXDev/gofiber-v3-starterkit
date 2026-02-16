package main

import (
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

	c.Provide(db.New)
	c.Provide(redis.New)
	c.Provide(s3.New)

	c.Provide(func(dbClient *bun.DB, redisClient *redis.RedisClient, s3Client *s3.S3Client) *fiber.App {
		cfg := config.FiberConfig()
		cfg.ErrorHandler = shared.RespondError

		app := fiber.New(cfg)
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))

		middlewares.FiberMiddleware(app)

		app.Get(healthcheck.LivenessEndpoint, healthcheck.New())

		routes.RegisterRoutes(app, dbClient, redisClient, s3Client)

		return app
	})

	c.Invoke(func(app *fiber.App, dbClient *bun.DB, redisClient *redis.RedisClient, s3Client *s3.S3Client) {
		defer dbClient.Close()
		defer redisClient.Client.Close()

		if err := utils.StartServerWithGracefulShutdown(app); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	})
}
