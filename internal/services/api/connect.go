package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/cvcio/mediawatch/internal/handlers"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2/articlesv2connect"
	"github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2/feedsv2connect"
	"github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2/passagesv2connect"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/redis"
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
		return err
	}
	log.Debugf("[SERVER] MongoDB connected on: %s", cfg.GetMongoURL())

	// Close mongo connection on exit
	defer func() { _ = mongo.Close() }()

	// ============================================================
	// Elasticsearch
	// ============================================================
	elastic, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		return err
	}
	log.Debugf("[SERVER] Elasticsearch connected on: %s", cfg.GetElasticsearchURL())

	// ============================================================
	// Neo4J
	// ============================================================
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		return err
	}
	log.Debugf("[SERVER] Neo4J connected on: %s", cfg.Neo.BOLT)
	defer func() { _ = neoClient.Client.Close(ctx) }()

	// ============================================================
	// Redis
	// ============================================================
	rdb, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		return err
	}
	log.Debugf("[SERVER] Redis connected on: %s", cfg.GetRedisURL())
	defer func() { _ = rdb.Close() }()

	// Create authenticator
	authenticator, err := auth.NewJWTAuthenticator(cfg.Auth.PrivateKeyFile, cfg.Auth.KeyID, cfg.Auth.Algorithm, cfg.Auth.Authorizer)
	if err != nil {
		return err
	}
	log.Debugf("[SERVER] Authenticator for %s created with key %s", cfg.Auth.Authorizer, cfg.Auth.KeyID)

	// ============================================================
	// HTTP Middleware
	// ============================================================
	// Create cors middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:  []string{"*.localhost", "*.cvcio.org", "*.mediawatch.io"},
		AllowOriginFunc: func(origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodPost,
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		AllowCredentials: true,
		MaxAge:           int(2 * time.Hour / time.Second), // Maximum value not ignored by any of major browsers
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
	errSignals := make(chan error, 1)

	// Listen for os signals
	osSignals := make(chan os.Signal, 1)

	// ============================================================
	// HTTP Muxer
	// ============================================================
	// ...
	mux := http.NewServeMux()

	// feeds
	feedsHandler, err := handlers.NewFeedsHandler(cfg, log, mongo, elastic, authenticator, rdb)
	if err != nil {
		return err
	}
	muxFeedsPath, muxFeedsHandler := feedsv2connect.NewFeedServiceHandler(feedsHandler, connect.WithCompressMinBytes(1024*100))
	mux.Handle(muxFeedsPath, muxFeedsHandler)

	// articles
	articlesHandler := handlers.NewArticlesHandler(cfg, log, mongo, elastic, neoClient, authenticator, rdb)
	muxArticlesPath, muxArticlesHandler := articlesv2connect.NewArticlesServiceHandler(articlesHandler, connect.WithCompressMinBytes(1024*100))
	mux.Handle(muxArticlesPath, muxArticlesHandler)

	//passages
	passagesHandler, err := handlers.NewPassagesHandler(cfg, log, mongo, elastic, authenticator, rdb)
	if err != nil {
		log.Errorf("[SERVER] Error while creating passages collection: %s", err.Error())
		return err
	}
	muxPassagesPath, muxPassagesHandler := passagesv2connect.NewPassageServiceHandler(passagesHandler, connect.WithCompressMinBytes(1024*100))
	mux.Handle(muxPassagesPath, muxPassagesHandler)
	// ============================================================
	// H2C Server
	// ============================================================
	// Use WebsocketProxy to expose the underlying handler as a bidi
	// websocket stream with newline-delimited JSON as the content encoding.
	server := &http.Server{
		Addr:              cfg.GetServiceURL(),
		Handler:           h2c.NewHandler(corsMiddleware.Handler(mux), &http2.Server{}),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       cfg.Service.ReadTimeout,
		WriteTimeout:      0, // set to 0 in order to stream to clients forever cfg.Service.WriteTimeout, //
		MaxHeaderBytes:    1 << 20,
	}

	// Start the service listening for requests.
	go func() {
		log.Debugf("[SERVER] Starting Connect/gRPC server on: %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errSignals <- err
		}
	}()

	// ============================================================
	// Termination
	// ============================================================
	// Listen for manual termination
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-errSignals:
		return err

	case s := <-osSignals:
		log.Infof("[SERVER] Server shutdown signal %s", s)
		// Asking listener shutdown and load shed.
		if err := server.Shutdown(ctx); err != nil {
			log.Debugf("[SERVER] Graceful shutdown did not complete in %s, exiting with error: %s", cfg.Service.ShutdownTimeout, err.Error())
			if err := server.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
