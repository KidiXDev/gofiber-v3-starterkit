package middlewares

import (
	"gofiber-starterkit/pkg/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(helmet.New(),
		config.CorsConfig(),
		logger.New(),
		recover.New(),
	)
}
