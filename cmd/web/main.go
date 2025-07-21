package main

import (
	"fmt"
	"os"

	"github.com/FilipusDev/filipus.dev.br/internal/config"
)

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

	fmt.Printf("%+v\n", cfg)
}
