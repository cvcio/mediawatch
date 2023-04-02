package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/twitter"
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

	// ============================================================
	// Start Mongo
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("Register DB: %v", err)
	}
	defer dbConn.Close()

	// =========================================================================
	// Create kafka consumer/producer worker

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

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()

	// ============================================================
	// Get feeds list
	feeds, err := feed.GetFeedsStreamList(
		context.Background(),
		dbConn,
		feed.Limit(cfg.Streamer.Size),
		feed.Lang(strings.ToUpper(cfg.Streamer.Lang)),
		feed.StreamType(int(commonv2.StreamType_STREAM_TYPE_TWITTER)),
	)
	if err != nil {
		log.Fatalf("Error getting feeds list: %v", err)
	}

	log.Infof("Loaded feeds: %d", len(feeds))
	if len(feeds) == 0 {
		log.Infof("No feeds to listen, exiting.")
		os.Exit(0)
	}

	// ============================================================
	// Get tweeter ids from feeds
	fUsernames := getUsernames(feeds)

	// ============================================================
	// Create a new twitter client
	api, err := twitter.NewTwitter(cfg.Twitter.TwitterConsumerKey, cfg.Twitter.TwitterConsumerSecret)
	if err != nil {
		log.Fatalf("Error connecting to twitter: %s", err.Error())
	}

	// ============================================================
	// Remove all active filter stream rules
	if _, err := removeRules(api); err != nil {
		log.Fatalf("Error while removing filter stream rules: %s", err.Error())
	}

	// ============================================================
	// Add new stream rules
	rules := splitFrom512(fUsernames)
	if _, err := addRules(api, rules); err != nil {
		log.Fatalf("Error while adding filter stream rules: %s", err.Error())
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

	// ============================================================
	// Create a new Listener service, with our twitter stream and the scrape service grpc conn
	log.Debugf("Twitter rules to listen : %v", rules)

	v := url.Values{}
	v.Add("expansions", "author_id,attachments.media_keys")
	v.Add("user.fields", "id,name,profile_image_url,url,username,verified")
	v.Add("tweet.fields", "created_at,id,author_id,lang,entities,in_reply_to_user_id")

	go func() {
		stream, err := api.GetFilterStream(v)
		if err != nil {
			errSingals <- err
			return
		}

		for t := range stream.C {
			f, _ := t.(twitter.StreamData)
			handler(log, f, tweetChan)
		}
	}()

	// Start the service listening for requests.
	log.Info("Ready to start")
	go func() {
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
			log.Debugf("New tweet from: %.16s %s", u.UserName, u.Url)

			urlTweet, err := json.Marshal(&u)
			if err != nil {
				log.Errorf("Error marshal tweet data: %s", err.Error())
				return
			}
			go worker.Produce(kaf.Message{Value: []byte(urlTweet)})

		case err := <-kafkaChan:
			log.Errorf("Error from kafka: %v", err)

		case err := <-errSingals:
			// Got Error from twitter stream
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

// handler handles incoming tweets
func handler(log *zap.SugaredLogger, t twitter.StreamData, tweetChan chan link.CatchedURL) {
	for _, v := range t.MatchingRules {
		if v.Tag != "mediawatch-listener" {
			return
		}
	}
	if t.Data.InReplyToUserID != "" {
		return
	}
	if len(t.Data.Entities.URLs) == 0 {
		return
	}
	for _, u := range t.Data.Entities.URLs {
		// Introduce clean URL logic
		// Remove Twitter Share ID (i.e. /#.WpAW30E8tRc.twitter)
		l, err := link.Parse(u.ExpandedURL)
		if err != nil {
			continue
		}
		createdAt, _ := t.Data.CreatedAtTime()
		messsage := link.CatchedURL{
			DocId:         uuid.New().String(),
			Type:          "twitter",
			Url:           l,
			TweetId:       t.Data.ID,
			TwitterUserId: t.Data.AuthorID,
			UserName:      getUserNameFromTweet(t.Data.AuthorID, t.Includes.Users),
			CreatedAt:     createdAt.Format(time.RFC3339),
		}
		tweetChan <- messsage
	}
}

// getUsernames list to listen.
func getUsernames(feeds []*feedsv2.Feed) []string {
	twitterUsernames := make([]string, 0)
	for _, f := range feeds {
		if f.UserName != "" {
			twitterUsernames = append(twitterUsernames, f.UserName)
		}
	}
	return twitterUsernames
}

// split screen names into string with max length 512 chars.
func splitFrom512(input []string) []string {
	var output []string
	current := ""
	for _, v := range input {
		s := "from:" + v
		if len(current) <= 512-(4+len(s)) {
			current += s + " OR "
		} else {
			if current[len(current)-4:] == " OR " {
				current = current[0 : len(current)-4]
			}
			output = append(output, current)
			current = ""
		}
	}
	return output
}

// retrieve screen name from tweet response.
func getUserNameFromTweet(authorId string, users []*twitter.User) string {
	for _, v := range users {
		if v.ID == authorId {
			return v.UserName
		}
	}
	return ""
}

// remove twitter api rules.
func removeRules(api *twitter.Twitter) (bool, error) {
	rules, err := api.GetFilterStreamRules(nil)
	if err != nil {
		return false, err
	}
	var ids []string
	for _, v := range rules.Data {
		if v.Tag == "mediawatch-listener" {
			ids = append(ids, v.ID)
		}
	}
	if len(ids) > 0 {
		rulesdel := new(twitter.Rules)
		rulesdel.Delete = &twitter.RulesDelete{
			Ids: ids,
		}
		deleted, err := api.PostFilterStreamRules(nil, rulesdel)
		if err != nil {
			return false, err
		}
		if deleted == nil {
			return false, errors.New("Rules not deleted")
		}
	}
	return true, nil
}

// add twitter api rules.
func addRules(api *twitter.Twitter, usernames []string) (bool, error) {
	rules := new(twitter.Rules)
	for _, v := range usernames {
		rules.Add = append(rules.Add, &twitter.RulesData{
			Value: v,
			Tag:   "mediawatch-listener",
		})
	}
	added, err := api.PostFilterStreamRules(nil, rules)
	if err != nil {
		return false, err
	}
	if added == nil {
		return false, errors.New("Rules not added")
	}
	return true, nil
}
