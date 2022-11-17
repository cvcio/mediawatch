package api

import (
	"context"
	"fmt"
	"os"

	"github.com/cvcio/mediawatch/internal/services/api"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// Create a new Config struct to store environment variables
	cfg := config.NewConfig()

	// Read config from env variables and panic on error,
	// as we can't continue
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Printf("API failed to parse environment variables, exiting with error: %s\n", err.Error())
		os.Exit(1)
	}

	// Create the context to cancel at exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the gRPC Server
	if err := api.RunGRPC(ctx, cfg); err != nil {
		fmt.Printf("API failure, exiting with error: %s\n", err.Error())
		os.Exit(1)
	}
}
