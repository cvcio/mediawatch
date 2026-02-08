package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cvcio/mediawatch/cmd/server/handlers"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
)

func main() {
	// ========================================
	// Configure
	cfg := config.NewConfigFromEnv()
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	defer log.Sync()

	log.Info("main: Starting")

	// ========================================
	// Start Mongo
	log.Info("main: Initialize Mongo")
	mongoClient, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Sugar().Fatalf("main: Register DB: %v", err)
	}
	log.Info("main: Connected to Mongo")
	defer mongoClient.Close()

	// ========================================
	// Start Elasticsearch
	log.Info("main: Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Sugar().Fatalf("[SVC-COMPARE] Register Elasticsearch: %v", err)
	}
	err = esClient.CreateElasticIndexWithLanguages(cfg.Elasticsearch.Index, cfg.Langs)
	if err != nil {
		log.Sugar().Fatalf("[SVC-WORKER] Index in Elasticsearch: %v", err)
	}
	log.Info("main: Connected to Elasticsearch")

	// ========================================
	// Start Neo4j client
	log.Info("main: Initialize Neo4J")
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Sugar().Fatalf("main: Register Neo4J: %v", err)
	}
	log.Info("main: Connected to Neo4J")
	defer neoClient.Client.Close(context.Background())

	// ========================================
	// Start redis client
	redisClient, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		log.Sugar().Fatal("Redis connection error", err)
	}
	defer func() { _ = redisClient.Close() }()

	// ========================================
	// Blocking main listening for requests
	serverErrors := make(chan error, 1)

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()

	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// create the http.Server
	api := http.Server{
		Addr:           cfg.GetApiURL(),
		Handler:        handlers.RegisterRoutes(cfg, log, mongoClient, esClient, neoClient, redisClient),
		ReadTimeout:    cfg.Api.ReadTimeout,
		WriteTimeout:   cfg.Api.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// ========================================
	// Start the prometheus http handler
	go func() {
		log.Sugar().Infof("main: Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		serverErrors <- promHandler.ListenAndServe()
	}()

	// Start the service listening for requests.
	log.Info("main: Ready to start")
	go func() {
		log.Sugar().Infof("main: Starting api Listening %s", cfg.GetApiURL())
		serverErrors <- api.ListenAndServe()
	}()

	// ========================================
	// Shutdown
	//
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Sugar().Fatalf("main: Error starting server: %v", err)

	case s := <-osSignals:
		log.Sugar().Debugf("[SVC-API] gRPC Server shutdown signal: %s", s)

		// Asking api to shutdown and load shed.
		if err := api.Shutdown(context.Background()); err != nil {
			log.Sugar().Errorf("[SVC-API] Graceful shutdown did not complete in %v: %v", cfg.Api.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Sugar().Fatalf("[SVC-API] Could not stop http server: %v", err)
			}
		}

		// Asking prometheus to shutdown and load shed.
		if err := promHandler.Shutdown(context.Background()); err != nil {
			log.Sugar().Errorf("[SVC-WORKER] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
			if err := promHandler.Close(); err != nil {
				log.Sugar().Fatalf("[SVC-WORKER] Could not stop http server: %v", err)
			}
		}
	}
}
