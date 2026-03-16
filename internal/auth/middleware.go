package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func JWTMiddleware(redis *redis.Client) fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized. Auth token required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := int(claims["user_id"].(float64))

		key := fmt.Sprintf("auth:session:%d", userID)

		val, err := redis.Get(c.Context(), key).Result()
		if err != nil || val != tokenString {
			return c.Status(401).JSON(fiber.Map{
				"error": "session expired",
			})
		}

		return c.Next()
	}
}
