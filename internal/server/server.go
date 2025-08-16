package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/FilipusDev/filipus.dev.br/internal/config"
	"github.com/FilipusDev/filipus.dev.br/templates"
	"github.com/FilipusDev/filipus.dev.br/templates/layouts"
	"github.com/a-h/templ"
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
	if err := http.ListenAndServe(addr, srv.Routes(logger)); err != nil {
		logger.Error("server failed to start", "error", err)
	}
}

// ========================================================================
// Routing
// ========================================================================
// The Routes method is responsible for setting up the entire router,
// including all routes and middleware. It returns the final http.Handler
// that can be passed to ListenAndServe.
func (s *Server) Routes(logger *slog.Logger) http.Handler {
	logger.Debug("\t !!! DEBUG !!! srv.Routes function called")

	logger.Debug("\t !!! DEBUG !!! creating a http.ServeMux")
	mux := http.NewServeMux()

	// Register handler for static assets.
	logger.Debug("\t !!! DEBUG !!! registering handler: for static files' server")
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Register the main application handler.
	// We wrap our main logic handler with the logging middleware.
	logger.Debug("\t !!! DEBUG !!! registering handler: main web handler '/' (root)")
	mux.Handle("/", s.loggingMiddleware(http.HandlerFunc(s.handleWebSite)))

	logger.Debug("\t !!! DEBUG !!! http.serveMux created", "mux", mux)

	logger.Debug("\t !!! DEBUG !!! ----------------------------------")
	logger.Debug("\t !!! DEBUG !!! from here, every event will call  ")
	logger.Debug("\t !!! DEBUG !!! the http.ServeMux with the logger ")
	logger.Debug("\t !!! DEBUG !!! middleware and handle each 'route'")
	logger.Debug("\t !!! DEBUG !!! as described in srv.handleWebSite ")
	logger.Debug("\t !!! DEBUG !!! ----------------------------------")
	return mux
}

// ========================================================================
// Middleware
// ========================================================================
// This is our standard, reusable logging middleware, now implemented as a
// method on the Server struct. This gives it direct access to the server's
// configured logger without needing to pass it in every time.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		spy := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(spy, r)

		s.logger.Info("request completed",
			"status", spy.status,
			"method", r.Method,
			"path", r.URL.Path,
			"host", r.Host,
			"duration", time.Since(start),
		)
	})
}

// ========================================================================
// Helper for Middleware
// ========================================================================
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// ========================================================================
// Handlers
// ========================================================================
// This is the main application logic handler. It's now a method on the Server
// struct, which gives it access to dependencies like `s.webSiteRepo` and `s.logger`.
func (s *Server) handleWebSite(w http.ResponseWriter, r *http.Request) {
	var pageComponent templ.Component
	var pageTitle string

	switch r.URL.Path {
	case "/":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		pageComponent = templates.Hello("'new' Eficia Engenharia Website")
		pageTitle = "Eficia Engenharia"
	default:
		s.logger.Warn("path not found for website", "host", r.Host, "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	fullPage := layouts.Base(pageComponent, pageTitle, s.assetsVersion)

	// Render the final component.
	err := fullPage.Render(context.Background(), w)
	if err != nil {
		fmt.Println("error rendering contacts page:", err.Error())
	}
}
