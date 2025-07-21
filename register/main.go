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

	log.Info("Starting CDC Registration Service")

	// Load configuration
	cfg := config.Load(configFile, port, kafkaURL)

	// Initialize database (if needed for storing connector metadata)
	//var db *db.QueryBuilder
	//if cfg.DatabaseURL != "" {
	//	db = db.NewQueryBuilder(cfg.DatabaseURL, log)
	//}
	//
	// Initialize HTTP client
	httpClient := http.NewRestyClient(log)
	//
	//// Initialize repository
	////repo := db.NewConnectorRepository(db, log)
	//
	//// Initialize service
	svc := service.NewCDCRegistrationService("kafka")

	// Initialize handler
	h := handler.NewCDCHandler(svc, log)

	// Setup HTTP server
	server := http.NewGinServer(h, log)

	log.Info("Server starting", logger.String("port", cfg.Port))
	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server", logger.Error(err))
	}
}
