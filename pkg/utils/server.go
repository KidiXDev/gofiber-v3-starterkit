package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func StartServer(app *fiber.App) error {
	addr := ConnectionString()

	enablePrefork := false
	if os.Getenv("APP_ENV") == "production" && os.Getenv("FIBER_PREFORK") == "true" {
		enablePrefork = true
	}

	return app.Listen(addr, fiber.ListenConfig{EnablePrefork: enablePrefork})
}

func StartServerWithGracefulShutdown(app *fiber.App) error {
	log.Info().Msg("Starting server...")
	addr := ConnectionString()

	enablePrefork := false
	if os.Getenv("APP_ENV") == "production" && os.Getenv("FIBER_PREFORK") == "true" {
		enablePrefork = true
	}

	go func() {
		if err := app.Listen(addr, fiber.ListenConfig{EnablePrefork: enablePrefork}); err != nil {
			log.Error().Err(err).Msg("server closed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	<-ctx.Done()
	log.Info().Msg("Server stopped")

	return nil
}
