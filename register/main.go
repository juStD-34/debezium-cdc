package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"register/config"
	"register/handler"
	"register/pkg/http"
	"register/pkg/logger"
	"register/service"
)

var (
	configFile string
	port       string
	kafkaURL   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cdc-registration",
		Short: "CDC Registration Service for Debezium",
		Long:  "A service to register CDC connectors for existing databases with Kafka Connect",
		Run:   runServer,
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "server port")
	rootCmd.PersistentFlags().StringVarP(&kafkaURL, "kafka-connect-url", "k", "http://localhost:8083", "Kafka Connect URL")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	// Initialize logger
	log := logger.NewZapLogger()
	defer log.Sync()
	cfg := config.Load()
	c := http.NewRestyClient(log)
	svc := service.NewCDCRegistrationService(cfg, log, c)
	log.Info("Starting CDC Registration Service")

	h := handler.NewCDCHandler(svc, log)

	r := http.NewGinServer(log)
	api := r.Group("/api")
	{
		api.POST("/connector", h.RegisterConnector)
		api.GET("connectors", h.ListConnectors)
		api.GET("/connectors/:name/status", h.GetConnectorStatus)
		api.DELETE("/connectors/:name", h.DeleteConnector)
	}

	log.Info("Starting CDC Registration Service")
	for _, route := range r.Routes() {
		log.Info("Registered route",
			logger.String("method", route.Method),
			logger.String("path", route.Path),
			logger.String("handler", route.Handler),
		)
	}
	err := r.Run(":" + port)
	if err != nil {
		return
	}
}
