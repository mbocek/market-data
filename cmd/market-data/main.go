package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	data "github.com/market-data/db"
	"github.com/market-data/internal/domain/market"
	"github.com/market-data/internal/providers/yahoo"

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
	cfg := initConfig()
	initLogger(&cfg.Logging)

	db := initDatabase(&cfg.Database)
	defer db.Close()

	runMigrations(cfg.Migrations.Enabled, cfg.Database.GetSchemaConnectionString())

	router := initRouter()
	// Create market repository and service
	marketRepo := market.NewMarketRepository(db)
	yahooClient := createYahooProvider(&cfg.YahooFinance)
	marketSvc := market.NewMarketService(marketRepo, yahooClient)

	// Configure auto-update settings
	marketSvc.SetAutoUpdateSettings(
		cfg.YahooFinance.GetUpdateInterval(),
		cfg.YahooFinance.EnableAutoUpdate,
	)

	// Start auto-update if enabled
	//if cfg.YahooFinance.EnableAutoUpdate {
	//	marketSvc.StartAutoUpdate()
	//}

	registerControllers(router, marketSvc)

	startServer(router, cfg.Server.Host, cfg.Server.Port)
}

func initConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	return cfg
}

func initLogger(lc *config.LoggingConfig) {
	lc.ConfigureLogging()
	if lc.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func initDatabase(dbCfg *config.DatabaseConfig) *database.DB {
	db, err := database.NewWithConfig(dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection is required")
	}

	log.Info().Msg("Database connection established")
	return db
}

func runMigrations(enabled bool, connString string) {
	migrator := migration.NewMigrator(enabled, connString, data.Migrations)
	if err := migrator.RunMigrations(); err != nil {
		log.Fatal().Err(err).Msg("Database migrations are required")
	}
}

func initRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLoggingMiddleware)
	return router
}

func requestLoggingMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	c.Next()
	latency := time.Since(start)
	log.Info().
		Str("method", c.Request.Method).
		Str("path", path).
		Int("status", c.Writer.Status()).
		Dur("latency", latency).
		Msg("Request processed")
}

func createYahooProvider(cfg *config.YahooFinanceConfig) *yahoo.Client {
	// Create and configure Yahoo Finance client
	return yahoo.NewClient(cfg)
}

func registerControllers(router *gin.Engine, marketSvc *market.MarketService) {
	// Register controllers
	healthController := api.NewHealthController()
	marketController := api.NewMarketController(marketSvc)
	healthController.RegisterRoutes(router)
	marketController.RegisterRoutes(router)
}

func startServer(router *gin.Engine, host, port string) {
	serverAddr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{
		Addr:              serverAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("address", serverAddr).Msg("Market Data Service starting")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	<-quit
	log.Info().Msg("Shutdown signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server exited gracefully")
	}
}
