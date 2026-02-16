package middleware

import (
	"golang-boilerplate-example/internal/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Generate request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("X-Request-ID", requestID)

		// Continue request
		err := c.Next()

		latency := time.Since(start)

		status := c.Response().StatusCode()

		logger.Log.Info("incoming request",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.String("ip", c.IP()),
			zap.Duration("latency", latency),
		)

		return err
	}
}
