package http

import (
	"go.uber.org/zap"
	"register/handler"
	"register/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func NewGinServer(handler *handler.CDCHandler, logger logger.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(GinLogger(logger))

	// Health check
	r.GET("/health", handler.HealthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/connectors", handler.RegisterConnector)
		api.GET("/connectors", handler.ListConnectors)
		api.GET("/connectors/:name/status", handler.GetConnectorStatus)
		api.DELETE("/connectors/:name", handler.DeleteConnector)
	}

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
