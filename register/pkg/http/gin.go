package http

import (
	"go.uber.org/zap"
	"net/http"
	"register/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func NewGinServer(logger logger.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(GinLogger(logger))

	// Health check
	r.GET("/health", healthCheck)

	return r
}

func GinLogger(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("duration", duration.String()),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "cdc-registration",
	})
}
