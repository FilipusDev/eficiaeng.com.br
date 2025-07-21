package main

import (
	"os"

	"github.com/FilipusDev/filipus.dev.br/internal/config"
	"github.com/FilipusDev/filipus.dev.br/internal/server"
)

// This variable will be set at build time by the -ldflags in the Dockerfile
var assetsVersion string

func main() {
	// 0. Bootstrap Looger
	logger := config.BootStrapLooger()
	logger.Debug("\t !!! DEBUG !!! step 0: app logger bootstrapped")

	// 1. Load configuration
	logger.Debug("\t !!! DEBUG !!! step 1: calling config.New(logger)")
	cfg, err := config.New(logger)
	if err != nil {
		logger.Error("could not load configuration", "error", err)
		os.Exit(1)
	}
	logger.Debug("\t !!! DEBUG !!! step 1: app config (cfg) successfully done")

	// TODO: 2. Connect to the database
	logger.Debug("\t !!! DEBUG !!! step 2: TBD calling database.NewConnection(ctx, dsn, logger)")

	// 3. Start the server
	// Pass all dependencies to the server
	logger.Debug("\t !!! DEBUG !!! step 3: calling server.start(cfg, dbpool, logger, assetsVersion)")
	server.Start(cfg, logger, assetsVersion)
}
