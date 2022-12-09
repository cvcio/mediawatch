package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cvcio/mediawatch/internal/handlers"
	"github.com/cvcio/mediawatch/internal/mediawatch/articles/v2/articlesv2connect"
	"github.com/cvcio/mediawatch/internal/mediawatch/feeds/v2/feedsv2connect"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// RunConnect creates a new connect/h2c server to handle gRPC endpoints.
func RunConnect(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) error {
	// ============================================================
	// Mongo
	// ============================================================
	mongo, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Errorf("[SERVER] MongoDB connection failed with error: %s", err.Error())
		return err
	}
	log.Debugf("[SERVER] MongoDB connected on: %s", cfg.GetMongoURL())

	// Close mongo connection on exit
	defer mongo.Close()

	// ============================================================
	// Elasticsearch
	// ============================================================
	elastic, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Errorf("[SERVER] Elasticsearch connection failed with error: %", err.Error())
		return err
	}
	log.Debugf("[SERVER] Elasticsearch connected on: %s", cfg.GetElasticsearchURL())

	// Create authenticator
	authenticator, err := auth.NewJWTAuthenticator(cfg.Auth.PrivateKeyFile, cfg.Auth.KeyID, cfg.Auth.Algorithm, cfg.Auth.Authorizer)
	if err != nil {
		log.Errorf("[SERVER] Error while creating the authenticator: %s", err.Error())
		return err
	}
	log.Debugf("[SERVER] Authenticator for %s created with key %s", cfg.Auth.Authorizer, cfg.Auth.KeyID)

	// ============================================================
	// HTTP Middleware
	// ============================================================
	// Create cors middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// ============================================================
	// gRPC Interceptors
	// ============================================================
	// ...

	// ============================================================
	// Set Channels
	// ============================================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	errSingals := make(chan error, 1)

	// Listen for os signals
	osSignals := make(chan os.Signal, 1)

	// ============================================================
	// HTTP Muxer
	// ============================================================
	// ...
	mux := http.NewServeMux()
	// feeds
	feedsHandler := handlers.NewFeedsHandler(cfg, log, mongo, elastic, authenticator)
	muxFeedsPath, muxFeedsHandler := feedsv2connect.NewFeedServiceHandler(feedsHandler)
	mux.Handle(muxFeedsPath, muxFeedsHandler)

	// articles
	articlesHandler := handlers.NewArticlesHandler(cfg, log, mongo, elastic, authenticator)
	muxArticlesPath, muxArticlesHandler := articlesv2connect.NewArticlesServiceHandler(articlesHandler)
	mux.Handle(muxArticlesPath, muxArticlesHandler)

	// ============================================================
	// H2C Server
	// ============================================================
	// Use WebsocketProxy to expose the underlying handler as a bidi
	// websocket stream with newline-delimited JSON as the content encoding.
	server := &http.Server{
		Addr:           cfg.GetServiceURL(),
		Handler:        h2c.NewHandler(cors.Handler(mux), &http2.Server{}),
		ReadTimeout:    cfg.Service.ReadTimeout,
		WriteTimeout:   cfg.Service.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Start the service listening for requests.
	go func() {
		log.Debugf("[SERVER] Starting Connect/gRPC server on: %s", server.Addr)
		errSingals <- server.ListenAndServe()
	}()

	// ============================================================
	// Termination
	// ============================================================
	// Listen for manual termination
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-errSingals:
		log.Errorf("[SERVER] Server error: %s", err.Error())
		return err

	case signal := <-osSignals:
		log.Infof("[SERVER] Server shutdown signal %s", signal)
		// Asking listener to shutdown and load shed.
		if err := server.Shutdown(ctx); err != nil {
			log.Debugf("[SERVER] Graceful shutdown did not complete in %s, exiting with error: %s", cfg.Service.ShutdownTimeout, err.Error())
			if err := server.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
