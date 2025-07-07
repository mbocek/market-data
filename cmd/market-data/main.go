package main

import (
	"fmt"
	"github.com/market-data/internal/domain/market"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/market-data/internal/config"
	"github.com/market-data/internal/database"
	"github.com/market-data/internal/database/migration"
	"github.com/market-data/internal/interfaces/api"
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
	db, err := database.NewWithConfig(&cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		log.Fatal().Msg("Database connection is required")
	}
	defer db.Close()
	log.Info().Msg("Database connection established")

	// Run database migrations
	migrator := migration.NewMigrator(&cfg.Migrations, cfg.Database.GetSchemaConnectionString())
	if err := migrator.RunMigrations(); err != nil {
		log.Error().Err(err).Msg("Failed to run database migrations")
		log.Fatal().Msg("Database migrations are required")
	}

	// Initialize repository
	marketRepo := market.NewMarketRepository(db)

	// Initialize market data service
	marketSvc := market.NewMarketService(marketRepo)

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
	healthController := api.NewHealthController()
	marketController := api.NewMarketController(marketSvc)

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
