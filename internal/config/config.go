package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Migrations   MigrationsConfig   `mapstructure:"migrations"`
	YahooFinance YahooFinanceConfig `mapstructure:"yahoo_finance"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	User              string `mapstructure:"user"`
	Password          string `mapstructure:"password"`
	DBName            string `mapstructure:"dbname"`
	SSLMode           string `mapstructure:"sslmode"`
	MaxConnections    int    `mapstructure:"max_connections"`
	ConnectionTimeout int    `mapstructure:"connection_timeout"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// MigrationsConfig represents the database migrations configuration
type MigrationsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// YahooFinanceConfig represents the Yahoo Finance API configuration
type YahooFinanceConfig struct {
	BaseURL          string   `mapstructure:"base_url"`
	RequestTimeout   int      `mapstructure:"request_timeout"`
	RetryCount       int      `mapstructure:"retry_count"`
	RetryWaitTime    int      `mapstructure:"retry_wait_time"`
	DefaultSymbols   []string `mapstructure:"default_symbols"`
	UpdateInterval   int      `mapstructure:"update_interval"`
	EnableAutoUpdate bool     `mapstructure:"enable_auto_update"`
}

// GetRequestTimeout returns the request timeout as a time.Duration
func (yfc *YahooFinanceConfig) GetRequestTimeout() time.Duration {
	return time.Duration(yfc.RequestTimeout) * time.Second
}

// GetRetryWaitTime returns the retry wait time as a time.Duration
func (yfc *YahooFinanceConfig) GetRetryWaitTime() time.Duration {
	return time.Duration(yfc.RetryWaitTime) * time.Millisecond
}

// GetUpdateInterval returns the update interval as a time.Duration
func (yfc *YahooFinanceConfig) GetUpdateInterval() time.Duration {
	return time.Duration(yfc.UpdateInterval) * time.Minute
}

// GetConnectionString returns the database connection string
func (dc *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.DBName, dc.SSLMode,
	)
}

// GetSchemaConnectionString returns the database connection string including a schema parameter
func (dc *DatabaseConfig) GetSchemaConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.DBName, dc.SSLMode)
}

// GetConnectionTimeout returns the database connection timeout as a time.Duration
func (dc *DatabaseConfig) GetConnectionTimeout() time.Duration {
	return time.Duration(dc.ConnectionTimeout) * time.Second
}

// ConfigureLogging configures the zerolog logger based on the logging configuration
func (lc *LoggingConfig) ConfigureLogging() {
	// Set time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	// Set log level
	switch lc.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set output format
	if lc.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano})
	}
}

// Load loads the configuration from the config file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// Yahoo Finance defaults
	viper.SetDefault("yahoo_finance.base_url", "https://query1.finance.yahoo.com/v8/finance/chart/")
	viper.SetDefault("yahoo_finance.request_timeout", 10)
	viper.SetDefault("yahoo_finance.retry_count", 3)
	viper.SetDefault("yahoo_finance.retry_wait_time", 500)
	viper.SetDefault("yahoo_finance.default_symbols", []string{"AAPL", "MSFT", "GOOG"})
	viper.SetDefault("yahoo_finance.update_interval", 15)
	viper.SetDefault("yahoo_finance.enable_auto_update", true)

	// Read environment variables
	viper.AutomaticEnv()

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Could not read config file: %v\n", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %v", err)
	}

	return &config, nil
}
