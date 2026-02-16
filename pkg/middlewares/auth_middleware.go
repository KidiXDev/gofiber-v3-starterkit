package middlewares

import (
	"strings"

	"gofiber-starterkit/app/api/services"
	"gofiber-starterkit/app/shared"
	"gofiber-starterkit/pkg/client/redis"
	"gofiber-starterkit/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	userService *services.UserService
	redisClient *redis.RedisClient
}

func NewAuthMiddleware(userService *services.UserService, redisClient *redis.RedisClient) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		redisClient: redisClient,
	}
}

func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return shared.ErrUnauthorized("Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return shared.ErrUnauthorized("Invalid authorization header format")
		}

		token := parts[1]
		claims, err := utils.ValidateAccessToken(m.redisClient.Client, token)
		if err != nil {
			return shared.ErrUnauthorized("Invalid or expired token")
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return shared.ErrUnauthorized("Invalid user ID in token")
		}

		user, err := m.userService.GetByID(c.Context(), userID)
		if err != nil {
			return shared.ErrUnauthorized("User not found")
		}

		c.Locals("userID", user.ID)
		c.Locals("user", user)
		c.Locals("tokenID", claims.TokenID)

		return c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		token := parts[1]
		claims, err := utils.ValidateAccessToken(m.redisClient.Client, token)
		if err != nil {
			return c.Next()
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return c.Next()
		}

		user, err := m.userService.GetByID(c.Context(), userID)
		if err != nil {
			return c.Next()
		}

		c.Locals("userID", user.ID)
		c.Locals("user", user)
		c.Locals("tokenID", claims.TokenID)

		return c.Next()
	}
}
