package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cvcio/mediawatch/internal/services/streamer"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// ============================================================
	// Read Config
	// ============================================================
	// Create a new Config struct to store environment variables
	cfg := config.NewConfig()

	// Read config from env variables and panic on error,
	// as we can't continue
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Printf("API failed to parse environment variables, exiting with error: %s\n", err.Error())
		os.Exit(1)
	}

	// ============================================================
	// Set Logger
	// ============================================================
	sugar := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	defer sugar.Sync()
	log := sugar.Sugar()

	// ============================================================
	// Start the Service
	// ============================================================
	// Create the context to cancel at exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the gRPC Server
	if err := streamer.RunGRPC(ctx, cfg, log); err != nil {
		log.Fatalf("[STREAMER] Fatal failure, exiting with error: %s\n", err.Error())
	}
}
