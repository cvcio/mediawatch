package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/cvcio/twitter"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kaf "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var FEEDS []string = []string{
	"https://armynews.gr/?feed=rss2",
	"https://alterthess.gr/feed/",
	"https://www.anatropinews.gr/feed/",
	"https://antinews.gr/feed/",
	"https://www.thatslife.gr/feed/",
	"https://www.documentonews.gr/feed/",
	"https://www.ekklisiaonline.gr/feed/",
	"https://www.enikos.gr/feed/",
	"https://www.ert.gr/feed/",
	"https://www.avgi.gr/rss.xml",
}

// ListenGroup struct.
type ListenGroup struct {
	ctx         context.Context
	log         *zap.SugaredLogger
	kafkaClient *kafka.KafkaGoClient
	errChan     chan error
}

// Close closes the kafka client.
func (worker *ListenGroup) Close() {
	worker.kafkaClient.Close()
}

// Procuce writes a new message to the kafka topic.
func (worker *ListenGroup) Produce(msg kaf.Message) {
	err := worker.kafkaClient.Producer.WriteMessages(worker.ctx, msg)
	if err != nil {
		worker.errChan <- errors.Wrap(err, "failed to write messages to kafka")
	}
}

// NewListenGroup implements a new ListenGroup struct.
func NewListenGroup(
	log *zap.SugaredLogger,
	kafkaClient *kafka.KafkaGoClient,
	errChan chan error,
) *ListenGroup {
	return &ListenGroup{
		context.Background(),
		log,
		kafkaClient,
		errChan,
	}
}

func main() {
	// ============================================================
	// Read Config
	// ============================================================
	cfg := config.NewConfig()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	// ============================================================
	// Set Logger
	// ============================================================
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	log.Info("[SVC-FEEDS] Starting")

	// ============================================================
	// Mongo
	// ============================================================
	log.Info("[SVC-FEEDS] Initialize Mongo")
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("[SVC-FEEDS] Register DB: %v", err)
	}
	log.Info("[SVC-FEEDS] Connected to Mongo")
	defer dbConn.Close()

	// ============================================================
	// Redis
	// ============================================================
	rdb, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		log.Fatalf("[SVC-FEEDS] Error connecting to Redis: %s", err.Error())
	}
	log.Info("[SVC-FEEDS] Connected to Redis")
	defer rdb.Close()
	// ============================================================
	// Kafka
	// ============================================================
	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewGoClient(
		false, true,
		[]string{cfg.Kafka.Broker},
		"",
		"",
		cfg.Kafka.WorkerTopic,
		cfg.Kafka.ConsumerGroupWorker,
		false,
	)

	// create a new worker
	worker := NewListenGroup(log, kafkaGoClient, kafkaChan)
	// close connections on exit
	defer worker.Close()

	// ============================================================
	// Prometheus
	// ============================================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()

	// ============================================================
	// Get feeds list
	feeds, err := feed.List(context.Background(), dbConn, feed.Status("active"), feed.Limit(1000))
	if err != nil {
		log.Fatalf("[SVC-FEEDS] error getting feeds list: %v", err)
	}

	log.Infoln(feeds)

	// ============================================================
	// Ticker
	// ============================================================
	done := make(chan bool, 1)
	defer close(done)

	go func(log *zap.SugaredLogger, done chan bool, targets []string, worker *ListenGroup, rdb *redis.RedisClient) {
		c := chunks(targets, 4)
		for _, v := range c {
			ticker := NewTicker(log, done, v, worker, rdb)
			go ticker.Tick()
		}
	}(log, done, FEEDS, worker, rdb)

	// ============================================================
	// Set Channels
	// ============================================================
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	errSingals := make(chan error, 1)

	// api will be our http.Server
	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), // api(cfg.Log.Debug, registry),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Start the service listening for requests.
	log.Info("[SVC-FEEDS] Ready to start")
	go func() {
		log.Infof("[SVC-FEEDS] Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		errSingals <- promHandler.ListenAndServe()
	}()

	// ============================================================
	// Termination
	// ============================================================
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// Blocking main and waiting for shutdown.
	select {
	case err := <-kafkaChan:
		log.Errorf("[SVC-FEEDS] Error from kafka: %v", err)
		os.Exit(1)

	case err := <-errSingals:
		// Got Error from twitter stream
		log.Errorf("[SVC-FEEDS] Error while streaming tweets: %s", err.Error())
		os.Exit(1)

	case s := <-osSignals:
		log.Debugf("[SVC-FEEDS] Listen shutdown signal: %s", s)

		// Asking prometheus to shutdown and load shed.
		if err := promHandler.Shutdown(context.Background()); err != nil {
			log.Errorf("[SVC-COMPARE] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
			if err := promHandler.Close(); err != nil {
				log.Fatalf("[SVC-COMPARE] Could not stop http server: %v", err)
			}
		}
	}
}

func chunks(feeds []string, size int) [][]string {
	var chunks [][]string
	for i := 0; i < len(feeds); i += size {
		d := i + size
		if d > len(feeds) {
			d = len(feeds)
		}
		chunks = append(chunks, feeds[i:d])
	}
	return chunks
}

// getIds list to listen
func getIds(feeds []*feed.Feed) []string {
	twitterIDs := make([]string, 0)
	for _, f := range feeds {
		if f.TwitterIDStr != "" {
			twitterIDs = append(twitterIDs, f.TwitterIDStr)
		} else {
			twitterIDs = append(twitterIDs,
				strconv.FormatInt(f.TwitterID, 10))
		}
	}
	return twitterIDs
}

// getUsernames list to listen
func getUsernames(feeds []*feed.Feed) []string {
	twitterUsernames := make([]string, 0)
	for _, f := range feeds {
		if f.ScreenName != "" {
			twitterUsernames = append(twitterUsernames, f.ScreenName)
		}
	}
	return twitterUsernames
}

func getUserNameFromTweet(authorId string, users []*twitter.User) string {
	for _, v := range users {
		if v.ID == authorId {
			return v.UserName
		}
	}
	return ""
}
