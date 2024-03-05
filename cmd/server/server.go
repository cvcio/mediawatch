package main

import (
	"compress/flate"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cvcio/mediawatch/internal/startup"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/twillio"
	"github.com/cvcio/mediawatch/pkg/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	scrape_pb "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"

	"github.com/cvcio/mediawatch/cmd/server/handlers"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/mid"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"github.com/kelseyhightower/envconfig"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func main() {
	// ========================================
	// Configure
	cfg := config.NewConfig()
	log := logrus.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	// Configure logger
	// Default level for this example is info, unless debug flag is present
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = logrus.InfoLevel
		log.Error(err.Error())
	}
	log.SetLevel(level)

	// Adjust logging format
	log.SetFormatter(&logrus.JSONFormatter{})
	if cfg.Env == "development" {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	log.Info("main: Starting")

	// ============================================== ==============
	// Start Mongo
	log.Info("main: Initialize Mongo")
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("main: Register DB: %v", err)
	}
	log.Info("main: Connected to Mongo")
	defer dbConn.Close()

	// =========================================================================
	// Start elasticsearch
	log.Info("main: Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Register Elasticsearch: %v", err)
	}

	log.Info("main: Connected to Elasticsearch")
	log.Info("main: Check for elasticsearch indexes")
	err = esClient.CreateElasticIndexWithLanguages(cfg.Elasticsearch.Index, cfg.Langs)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Index in elasticsearch: %v", err)
	}

	// =========================================================================
	// Start neo4j client
	log.Info("main: Initialize Neo4J")
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Fatalf("main: Register Neo4J: %v", err)
	}
	log.Info("main: Connected to Neo4J")
	defer neoClient.Client.Close()

	// =========================================================================
	// Create authenticator
	authenticator, err := startup.GetAuthenticator(cfg)
	if err != nil {
		log.Fatalf("main: Authenticator: %v", err)
	}

	// =========================================================================
	// Create the gRPC Service
	// Parse Server Options
	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithInsecure())

	// Create gRPC Scrape Connection
	scrapeGRPC, err := grpc.Dial(cfg.Scrape.Host, grpcOptions...)
	if err != nil {
		log.Debugf("main: GRPC Scrape did not connect: %v", err)
	}
	defer scrapeGRPC.Close()

	scrape := scrape_pb.NewScrapeServiceClient(scrapeGRPC)

	// Create GoogleAuth
	// TODO: create map of auths
	externalAuths, err := cfg.ExternalAuths()
	if err != nil {
		log.Error(err)
	}

	log.Info("main: Created auth keys")
	// Create mail service
	mail := mailer.New(
		cfg.SMTP.Server,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
		cfg.SMTP.FromName,
		cfg.SMTP.Reply,
	)
	log.Info("main: Created mail service")
	// Create Twillio SMS Service
	twillio := twillio.New(
		cfg.Twillio.SID,
		cfg.Twillio.Token,
		"MediaWatch",
	)
	log.Info("main: Created twillio sms service")

	// Create Stripe Client

	// Create Twitter API Service
	twtt, err := twitter.NewAPI(cfg.Twitter.TwitterConsumerKey,
		cfg.Twitter.TwitterConsumerSecret, cfg.Twitter.TwitterAccessToken,
		cfg.Twitter.TwitterAccessTokenSecret)
	if err != nil {
		log.Fatalf("Error connecting to twitter: %s", err.Error())
	}
	log.Info("main: Created twitter api service")
	// ========================================
	// Create a server

	// Create cors middleware
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, //, "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		// Debug:            cfg.Web.Debug,
	})

	// create list of middlewares to enable on http server
	mw := []func(http.Handler) http.Handler{
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler,
		// trace.Trace(),
		mid.LoggerMiddleware(log),
		middleware.NewCompressor(flate.DefaultCompression).Handler,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		httprate.LimitByIP(980, 1*time.Minute),
	}

	// ========================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()

	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), // api(cfg.Log.Debug, registry),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// create the http.Server
	api := http.Server{
		Addr:           cfg.GetApiURL(),
		Handler:        handlers.API(log, registry, dbConn, esClient, neoClient, mw, authenticator, externalAuths, mail, twillio, twtt, scrape, cfg),
		ReadTimeout:    cfg.Api.ReadTimeout,
		WriteTimeout:   cfg.Api.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// ========================================
	// Start the prometheus http handler
	go func() {
		log.Infof("main: Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		serverErrors <- promHandler.ListenAndServe()
	}()

	// Start the service listening for requests.
	log.Info("main: Ready to start")
	go func() {
		log.Infof("main: Starting api Listening %s", cfg.GetApiURL())
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
		log.Fatalf("main: Error starting server: %v", err)

	case s := <-osSignals:
		log.Debugf("[SVC-API] gRPC Server shutdown signal: %s", s)

		// Asking api to shutdown and load shed.
		if err := api.Shutdown(context.Background()); err != nil {
			log.Errorf("[SVC-API] Graceful shutdown did not complete in %v: %v", cfg.Api.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("[SVC-API] Could not stop http server: %v", err)
			}
		}

		// Asking prometheus to shutdown and load shed.
		if err := promHandler.Shutdown(context.Background()); err != nil {
			log.Errorf("[SVC-WORKER] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
			if err := promHandler.Close(); err != nil {
				log.Fatalf("[SVC-WORKER] Could not stop http server: %v", err)
			}
		}
	}

}
