package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/cvcio/mediawatch/pkg/twitter"
	"github.com/google/uuid"
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
	// ========================================
	// Configure
	cfg := config.NewConfig()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	// ** LOGGER
	// Create a reusable zap logger
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	log.Info("[SVC-LISTEN] Starting")

	// ============================================================
	// Start Mongo
	log.Info("[SVC-LISTEN] Initialize Mongo")
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("[SVC-LISTEN] Register DB: %v", err)
	}
	log.Info("[SVC-LISTEN] Connected to Mongo")
	defer dbConn.Close()

	// =========================================================================
	// Create kafka consumer/producer worker

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

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()

	// ============================================================
	// Get feeds list
	feeds, err := feed.List(context.Background(), dbConn, feed.Status("active"), feed.Limit(1000))
	if err != nil {
		log.Fatalf("[SVC-LISTEN] error getting feeds list: %v", err)
	}

	// ============================================================
	// Get tweeter ids from feeds
	fIDs := getScreenNames(feeds.Data)

	// ============================================================
	// Create a new twitter client
	twitterAPIClient, err := twitter.NewAPI(cfg.Twitter.TwitterConsumerKey,
		cfg.Twitter.TwitterConsumerSecret, cfg.Twitter.TwitterAccessToken,
		cfg.Twitter.TwitterAccessTokenSecret)
	if err != nil {
		log.Fatalf("[SVC-LISTEN] Error connecting to twitter: %s", err.Error())
	}

	// ============================================================
	// Create a new Listener service, with our twitter stream and the scrape service grpc conn
	log.Debugf("[SVC-LISTEN] Twitter ids to listen : %v", fIDs)

	svc, err := twitter.NewListener(
		twitterAPIClient, log,
		twitter.WithPublicStream(map[string][]string{"follow": fIDs}),
	)

	if err != nil {
		log.Fatalf("[SVC-LISTEN] error creating twitter listener: %s", err.Error())
	}

	// Create a channel to send catched urls from tweets
	tweetChan := make(chan link.CatchedURL, 1)

	// ========================================
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

	// ========================================
	// Here we start listening for tweets given a handler function
	go svc.TweetListen(handler(log, tweetChan))

	// Start the service listening for requests.
	log.Info("[SVC-LISTEN] Ready to start")
	go func() {
		log.Infof("[SVC-COMPARE] Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		errSingals <- promHandler.ListenAndServe()
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
	for {
		select {
		// Got Url from twitter
		case u := <-tweetChan:
			log.Infow(
				"catched",
				"tweetID", u.TweetID,
				"user", u.ScreenName,
				"timeCreated", u.CreatedAtStr,
				"url", u.URL,
			)

			urlTweet, err := json.Marshal(&u)
			if err != nil {
				log.Errorf("[SVC-LISTEN] error marshal tweet data: %s", err.Error())
				return
			}
			go worker.Produce(kaf.Message{Value: []byte(urlTweet)})

		case err := <-kafkaChan:
			log.Errorf("[SVC-LISTEN] error from kafka: %v", err)
			os.Exit(1)

		case err := <-errSingals:
			// Got Error from twitter stream
			log.Errorf("[SVC-LISTEN] Error while streaming tweets: %s", err.Error())
			os.Exit(1)

		case s := <-osSignals:
			log.Debugf("[SVC-LISTEN] Listen shutdown signal: %s", s)

			// Asking prometheus to shutdown and load shed.
			if err := promHandler.Shutdown(context.Background()); err != nil {
				log.Errorf("[SVC-COMPARE] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
				if err := promHandler.Close(); err != nil {
					log.Fatalf("[SVC-COMPARE] Could not stop http server: %v", err)
				}
			}
		}
	}
}

// getScreenNames list to listen
func getScreenNames(feeds []*feed.Feed) []string {
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

// handler handles incoming tweets
func handler(log *zap.SugaredLogger, tweetChan chan link.CatchedURL) func(t anaconda.Tweet) {
	return func(t anaconda.Tweet) {
		// Introduce a better logic here to classify Tweet and extract Url
		if !twitter.IsTweet(t) {
			return
		}
		if len(t.Entities.Urls) == 0 {
			return
		}
		for u := range t.Entities.Urls {
			// Introduce clean URL logic
			// Remove Twitter Share ID (i.e. /#.WpAW30E8tRc.twitter)
			l, err := link.Parse(t.Entities.Urls[u].Expanded_url)
			if err != nil {
				log.Error(err)
				continue
			}

			createdTime, _ := t.CreatedAtTime()
			tweetChan <- link.CatchedURL{
				ID:               uuid.New().String(),
				URL:              l,
				TweetID:          t.Id,
				TwitterUserID:    t.User.Id,
				TwitterUserIDStr: t.User.IdStr,
				ScreenName:       t.User.ScreenName,
				CreatedAt:        createdTime,
				CreatedAtStr:     t.CreatedAt,
			}
		}
	}
}
