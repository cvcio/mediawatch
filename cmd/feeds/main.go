package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kaf "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// ListenGroup struct.
type ListenGroup struct {
	ctx         context.Context
	log         *zap.SugaredLogger
	kafkaClient *kafka.KafkaClient
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
	kafkaClient *kafka.KafkaClient,
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
	sugar := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	defer sugar.Sync()
	log := sugar.Sugar()

	// ============================================================
	// Mongo
	// ============================================================
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("Register DB: %v", err)
	}
	defer dbConn.Close()

	// ============================================================
	// Redis
	// ============================================================
	rdb, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		log.Fatalf("Error connecting to Redis: %s", err.Error())
	}
	defer rdb.Close()

	// ============================================================
	// Kafka
	// ============================================================
	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewKafkaClient(
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
	feeds, err := feed.GetFeedsStreamList(
		context.Background(),
		dbConn,
		feed.Limit(cfg.Streamer.Size),
		feed.Lang(strings.ToUpper(cfg.Streamer.Lang)),
		// feed.StreamType(int(commonv2.StreamType_STREAM_TYPE_RSS)),
	)
	if err != nil {
		log.Fatalf("error getting feeds list: %v", err)
	}

	feeds = filter(feeds)
	log.Debugf("Loaded feeds: %d", len(feeds))
	if len(feeds) == 0 {
		log.Infof("No feeds to listen, exiting.")
		os.Exit(0)
	}

	// ============================================================
	// Ticker
	// ============================================================
	done := make(chan bool, 1)
	defer close(done)

	// create chunks
	targets := chunks(feeds, cfg.Streamer.Chunks)

	// run the tickers
	go tick(cfg, log, worker, rdb, done, targets, cfg.Streamer.Init, cfg.Streamer.Interval)

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
	log.Info("Ready to start")
	go func() {
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
	for {
		select {
		case err := <-kafkaChan:
			// Ignore the Error
			log.Errorf("Error from kafka: %v", err)

		case err := <-errSingals:
			// Got Error from stream
			log.Errorf("Error while streaming tweets: %s", err.Error())
			os.Exit(1)

		case s := <-osSignals:
			log.Debugf("Listen shutdown signal: %s", s)

			// Asking prometheus to shutdown and load shed.
			if err := promHandler.Shutdown(context.Background()); err != nil {
				log.Errorf("Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
				if err := promHandler.Close(); err != nil {
					log.Fatalf("Could not stop http server: %v", err)
				}
			}
		}
	}
}

func tick(cfg *config.Config, log *zap.SugaredLogger, worker *ListenGroup, rdb *redis.RedisClient, done chan bool, targets [][]*feedsv2.Feed, init bool, interval time.Duration) {
	// delay := interval / time.Duration(math.Ceil(float64(len(targets))/100))
	// delay := time.Second * time.Duration(len(targets))
	for _, v := range targets {
		ticker := NewTicker(cfg, log, worker, rdb, done, v, init, interval)
		go ticker.Tick()
		time.Sleep(time.Second * time.Duration(150))
	}
}

func chunks(feeds []*feedsv2.Feed, size int) [][]*feedsv2.Feed {
	var chunks [][]*feedsv2.Feed
	for i := 0; i < len(feeds); i += size {
		d := i + size
		if d > len(feeds) {
			d = len(feeds)
		}
		chunks = append(chunks, feeds[i:d])
	}
	return chunks
}

func filter(feeds []*feedsv2.Feed) []*feedsv2.Feed {
	var f []*feedsv2.Feed
	for _, v := range feeds {
		if v.Stream.StreamStatus == commonv2.Status_STATUS_ACTIVE {
			if v.Stream.StreamType == commonv2.StreamType_STREAM_TYPE_OTHER || v.Stream.StreamType == commonv2.StreamType_STREAM_TYPE_RSS {
				f = append(f, v)
			}
		}
		// if v.Meta.ContentType != commonv2.ContentType_CONTENT_TYPE_AUTO &&
		// 	v.Meta.ContentType != commonv2.ContentType_CONTENT_TYPE_MUSIC &&
		// 	v.Meta.ContentType != commonv2.ContentType_CONTENT_TYPE_ENTERTAINMENT &&
		// 	v.Meta.ContentType != commonv2.ContentType_CONTENT_TYPE_SPORTS {
		// 	f = append(f, v)
		// }
	}
	return f
}
