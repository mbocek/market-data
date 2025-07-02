package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/market-data/internal/config"
	"github.com/market-data/internal/controller"
	"github.com/market-data/internal/database"
	"github.com/market-data/pkg/marketservice"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Configure logger
	cfg.Logging.ConfigureLogging()

	// Set Gin mode based on environment
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	var db *database.DB
	var serviceOpts []marketservice.ServiceOption

	db, err = database.NewWithConfig(&cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		log.Fatal().Msg("Database connection is required")
	}
	defer db.Close()
	log.Info().Msg("Database connection established")
	serviceOpts = append(serviceOpts, marketservice.WithDatabase(db))

	// Initialize market data service
	service := marketservice.NewService(serviceOpts...)

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add zerolog middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Msg("Request processed")
	})

	// Initialize controllers
	healthController := controller.NewHealthController()
	marketController := controller.NewMarketController(service)

	// Register routes
	healthController.RegisterRoutes(router)
	marketController.RegisterRoutes(router)

	// Start HTTP server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	log.Info().Str("address", serverAddr).Msg("Market Data Service starting")
	if err := router.Run(serverAddr); err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}
