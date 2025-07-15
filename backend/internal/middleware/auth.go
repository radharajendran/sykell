package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Default secret for development
	}
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}