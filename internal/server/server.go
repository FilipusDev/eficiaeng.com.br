package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FilipusDev/filipus.dev.br/internal/config"
)

// ========================================================================
// Server Struct and Constructor
// ========================================================================
// By creating a Server struct, we can attach our handlers as methods.
// This is a common and clean pattern that keeps related logic grouped together
// and gives all handlers easy access to shared dependencies like the logger and database.
type Server struct {
	logger        *slog.Logger
	assetsVersion string
}

// New creates a new instance of our Server with all its dependencies.
func newServer(logger *slog.Logger, assetsVersion string) *Server {
	if assetsVersion == "" {
		// Fallback for local dev if the version isn't injected
		assetsVersion = fmt.Sprintf("dev-%d", time.Now().Unix())
	}
	return &Server{
		logger:        logger,
		assetsVersion: assetsVersion,
	}
}

// ========================================================================
// Main Server Start Function
// ========================================================================
// This is the primary public function for our package. Its only job is to
// create a new server instance and start it. All the complex setup logic
// has been moved into the server's own methods.
func Start(cfg *config.Config, logger *slog.Logger, assertsVersion string) {
	logger.Debug("\t !!! DEBUG !!! server.Start function called")

	logger.Debug("\t !!! DEBUG !!! creating server", "srv", "")
	srv := newServer(logger, assertsVersion)
	logger.Debug("\t !!! DEBUG !!! server created", "srv", srv)

	logger.Debug("\t !!! DEBUG !!! formatting server addr with port", "addr", "", "port", cfg.Server.Port)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Debug("\t !!! DEBUG !!! server addr formatted", "addr", addr)
	logger.Info("server starting", "addr", addr)

	// We call a method on our struct to get the final, configured router.
	logger.Debug("\t !!! DEBUG !!! calling ListenAndServe func, passing the addr and calling srv.Routes func", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("server failed to start", "error", err)
	}
}
