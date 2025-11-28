package http

import (
	"recipe-processor/internal/shared/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests with structured logging
func LoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)

		log.Info("http request",
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", query),
			logger.Int("status", c.Writer.Status()),
			logger.Int("size", c.Writer.Size()),
			logger.Duration("duration", duration),
			logger.String("client_ip", c.ClientIP()),
		)
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
