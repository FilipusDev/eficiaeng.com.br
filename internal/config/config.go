package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/lmittmann/tint"
)

// ServerConfig holds the configuration for the HTTP server.
type ServerConfig struct {
	Port int
}

// DatabaseConfig holds the configuration for database.
type DatabaseConfig struct {
	DSN string // Data Source Name
}

// Config is the top-level configuration struct for the entire application.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func BootStrapLooger() *slog.Logger {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO" // Default to INFO level if not set
	}

	var logLevel slog.Level
	// Create the final, level-aware logger
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.Kitchen,
		AddSource:  true,
	})

	return slog.New(logHandler)
}

func New(logger *slog.Logger) (*Config, error) {
	logger.Debug("\t !!! DEBUG !!! config.New function called")

	logger.Debug("\t !!! DEBUG !!! getting environment variable", "env_var", "APP_PORT")
	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		portStr = "8888" // Default to port 8888 if not set
		logger.Warn("environment variable NOT FOUND, fallbacking...", "portStr", portStr)
	}
	logger.Debug("\t !!! DEBUG !!! environment variable found", "portStr", portStr)

	logger.Debug("\t !!! DEBUG !!! converting portStr(string) to port(int)", "portStr", portStr)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Error("failed to convert portStr to integer", "portStr", portStr)
		return nil, err
	}
	logger.Debug("\t !!! DEBUG !!! successfully convertion portStr(string) to port(int)",
		"portStr",
		portStr,
		"port",
		port)

	// TODO: solve this crap at some point!!!
	logger.Debug("\t !!! DEBUG !!! getting environment variable", "env_var", "DATABASE_URL")
	dsn := os.Getenv("DATABASE_URL")
	// if dsn == "" {
	// 	logger.Error("environment variable NOT FOUND, erroring...", "dsn", dsn)
	// 	return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	// }
	logger.Debug("\t !!! DEBUG !!! environment variable found", "dsn", dsn)

	logger.Debug("\t !!! DEBUG !!! creating configuration struct", "cfg", "")
	cfg := &Config{
		Server:   ServerConfig{Port: port},
		Database: DatabaseConfig{DSN: dsn},
	}
	logger.Debug("\t !!! DEBUG !!! config.New ran successfully", "cfg", cfg)
	return cfg, nil
}
