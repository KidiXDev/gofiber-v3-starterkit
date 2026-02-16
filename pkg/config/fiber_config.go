package config

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func jsonEncoderNoEscape(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

func init() {
	if os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local" || os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Error().Err(err).Msg("Error loading .env file")
		}
	}
}

func FiberConfig() fiber.Config {
	log.Info().Msg("Loading fiber configuration...")

	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if readTimeoutSecondsCount <= 0 {
		readTimeoutSecondsCount = 5
	}

	return fiber.Config{
		ReadTimeout:   time.Second * time.Duration(readTimeoutSecondsCount),
		CaseSensitive: false,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "GoFiber Starterkit v1.0.0",
		ProxyHeader:   fiber.HeaderXForwardedFor,
		BodyLimit:     20 * 1024 * 1024,
		JSONEncoder:   jsonEncoderNoEscape,
		JSONDecoder:   json.Unmarshal,
	}
}

func CorsConfig() fiber.Handler {
	parse := func(v string) []string {
		if v == "" {
			return nil
		}
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	return cors.New(cors.Config{
		AllowOrigins:     parse(os.Getenv("CORS_ALLOW_ORIGINS")),
		AllowHeaders:     parse(os.Getenv("CORS_ALLOW_HEADERS")),
		AllowMethods:     parse(os.Getenv("CORS_ALLOW_METHODS")),
		AllowCredentials: true,
	})
}
