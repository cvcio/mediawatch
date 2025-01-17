package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	articlesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	enrichv2 "github.com/cvcio/mediawatch/pkg/mediawatch/enrich/v2"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kaf "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var (
	workerProcessDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "consumer_process_duration_seconds",
		Help:       "Duration of consumer processing requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})

	workerProcessErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "consumer_process_errors_total",
		Help: "Total number of consumer processing errors.",
	}, []string{"service"})
)

// WorkerGroup struct.
type WorkerGroup struct {
	ctx         context.Context
	log         *zap.Logger
	kafkaClient *kafka.Client
	errChan     chan error

	dbConn    *db.MongoDB
	esClient  *es.Elastic
	esIndex   string
	neoClient *neo.Neo

	ackBefore string

	rdb *redis.RedisClient
}

// Close closes the kafka client.
func (worker *WorkerGroup) Close() {
	worker.kafkaClient.Close()
}

// ArticleExists checks if an article exists in elasticsearch.
func (worker *WorkerGroup) ArticleExists(url string) bool {
	opts := article.NewOpts()
	opts.Index = worker.esIndex + "_*"
	opts.Url = url

	exists := article.Exists(context.Background(), worker.esClient, opts)
	return exists
}

// TimeTrack logs the duration of a process.
func (worker *WorkerGroup) TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	worker.log.Debug("Duration", zap.String("process", name), zap.Duration("elapsed", elapsed))
}

// Consume consumes kafka topics inside an infinite loop. In our logic we need
// to fetch a message from a topic (FetchMessage), parse the json (Unmarshal)
// and process the content (articleProcess) if it doesn't already exist.
// If, for any reason, any step fails with an error we will commit this message
// to kafka as we don't want to process this particular message again (it failed
// for some reason).
func (worker *WorkerGroup) Consume() {
	for {
		func() {
			timer := prometheus.NewTimer(workerProcessDuration.WithLabelValues("worker"))
			defer timer.ObserveDuration()
			//
			// fetch the message from kafka topic
			m, err := worker.kafkaClient.Consumer.FetchMessage(worker.ctx)
			if err != nil {
				// at this point we don't have a message, as such we don't commit
				// send the error to channel
				worker.errChan <- errors.Wrap(err, "failed to fetch messages from kafka")
				return
			}

			switch m.Topic {
			case "worker":
				// process the article
				// Unmarshal incoming json message
				var msg link.Link
				if err := json.Unmarshal(m.Value, &msg); err != nil {
					// mark message as read (commit)
					worker.Commit(m)
					worker.log.Error("Failed to unmarshal message from kafka", zap.Error(err))
					return
				}

				// Commit messages if the AckBefore environment variable is present and valid
				if worker.ackBefore != "" {
					if s, err := time.Parse(time.DateOnly, worker.ackBefore); err == nil {
						if e, err := time.Parse(time.RFC3339, msg.CreatedAt); err == nil {
							if e.Before(s) {
								worker.Commit(m)
								worker.log.Warn("Skip message (AckBefore)", zap.String("createdAt", msg.CreatedAt), zap.String("url", msg.Url))
								return
							}
						}
					}
				}

				// re-validate link and make sure it is a valid url
				if _, err := link.Validate(msg.Url); err != nil {
					// mark message as read (commit)
					worker.Commit(m)
					worker.log.Warn("Skip message (Invalid)", zap.String("createdAt", msg.CreatedAt), zap.String("url", msg.Url))
					return
				}

				name := msg.Hostname
				if msg.Type == "twitter" {
					name = msg.UserName
				}
				worker.log.Debug("Consume message", zap.String("hostname", name), zap.String("createdAt", msg.CreatedAt), zap.String("url", msg.Url))

				// check if article exists before processing it
				// on nil error the article exists
				if exists := worker.ArticleExists(msg.Url); !exists {
					// process the article
					if err := worker.ProcessArticle(msg); err != nil {
						worker.log.Error("Error while processing message", zap.String("hostname", name), zap.String("createdAt", msg.CreatedAt), zap.String("url", msg.Url), zap.Error(err))

						// send the error to channel
						worker.errChan <- errors.Wrap(err, "failed process article")
						return
					}
				}

			case "store":
				var msg articlesv2.Article
				if err := json.Unmarshal(m.Value, &msg); err != nil {
					// mark message as read (commit)
					worker.Commit(m)
					worker.log.Error("Failed to unmarshal message from kafka", zap.Error(err))
					return
				}
				if err := worker.StoreArticle(&msg); err != nil {
					worker.log.Error("Error while storing article", zap.String("hostname", msg.Hostname), zap.String("createdAt", msg.CrawledAt), zap.String("url", msg.Url), zap.Error(err))

					// send the error to channel
					worker.errChan <- errors.Wrap(err, "failed to store article")
					return
				}
				worker.log.Info("Store message", zap.String("doc_id", string(msg.DocId)))
			}

			// mark message as read (commit)
			worker.Commit(m)
		}()
	}
}

// Commit commits a message to the kafka topic.
func (worker *WorkerGroup) Commit(m kaf.Message) {
	if err := worker.kafkaClient.Consumer.CommitMessages(worker.ctx, m); err != nil {
		worker.errChan <- errors.Wrap(err, "failed to commit messages to kafka")
	}
}

// Produce writes a new message to the kafka topic.
func (worker *WorkerGroup) Produce(msg kaf.Message) {
	err := worker.kafkaClient.Producer.WriteMessages(worker.ctx, msg)
	if err != nil {
		worker.errChan <- errors.Wrap(err, "failed to write messages to kafka")
	}
}

// Publish writes a new message to the corresponding redis pub/sub channel.
func (worker *WorkerGroup) Publish(channel string, msg string) {
	err := worker.rdb.Publish(channel, msg)
	if err != nil {
		worker.errChan <- errors.Wrap(err, "failed to publish message to redis")
	}
}

// NewWorkerGroup implements a new WorkerGroup struct.
func NewWorkerGroup(
	log *zap.Logger,
	kafkaClient *kafka.Client,
	errChan chan error,
	dbConn *db.MongoDB,
	esClient *es.Elastic,
	esIndex string,
	neoClient *neo.Neo,
	ackBefore string,
	rdb *redis.RedisClient,
) *WorkerGroup {
	return &WorkerGroup{
		context.Background(),
		log,
		kafkaClient,
		errChan,
		dbConn,
		esClient,
		esIndex,
		neoClient,
		ackBefore,
		rdb,
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

	// ========================================
	// Create a reusable zap logger
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	// Sync the logger on exit
	defer func() {
		// ignore the error as Sync() is always returning nil
		// sync: error syncing '/dev/stdout': Invalid argument
		_ = log.Sync()
	}()

	// =========================================================================
	// Create mongo client
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatal("MongoDB connection error", zap.Error(err))
	}
	defer func() { _ = dbConn.Close() }()

	// =========================================================================
	// Start elasticsearch
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatal("Elasticsearch connection error", zap.Error(err))
	}

	log.Debug("Checking for elasticsearch indices")
	if err := esClient.CreateElasticIndexWithLanguages(cfg.Elasticsearch.Index, cfg.Langs); err != nil {
		log.Fatal("Elasticsearch indices error", zap.Error(err))
	}

	// =========================================================================
	// Start neo4j client
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Fatal("Neo4j connection error", zap.Error(err))
	}
	defer func() { _ = neoClient.Client.Close(context.Background()) }()

	// ============================================================
	// Redis
	// ============================================================
	rdb, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		log.Fatal("Redis connection error", zap.Error(err))
	}
	defer func() { _ = rdb.Close() }()

	// =========================================================================
	// Create kafka consumer/producer worker
	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewKafkaClient(
		true, true,
		cfg.GetKafkaBrokers(),
		cfg.Kafka.ConsumerTopic,
		cfg.Kafka.ConsumerGroup,
		cfg.Kafka.WorkerOffsetOldest,
	)

	// create a new worker
	worker := NewWorkerGroup(
		log, kafkaGoClient, kafkaChan,
		dbConn, esClient, cfg.Elasticsearch.Index, neoClient, cfg.Kafka.AckBefore, rdb,
	)

	// close connections on exit
	defer worker.Close()

	// run the worker
	go func() {
		defer close(kafkaChan)
		log.Debug("Starting kafka consumer")
		// consume messages from kafka
		worker.Consume()
	}()

	// ========================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	errSignals := make(chan error, 1)

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()
	registry.MustRegister(workerProcessErrors)
	registry.MustRegister(workerProcessDuration)

	// create a prometheus http.Server
	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// ========================================
	// Start the http service for prometheus
	go func() {
		errSignals <- promHandler.ListenAndServe()
	}()

	// ========================================
	// Shutdown
	//
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// ========================================
	// Stop API Service
	// Blocking main and waiting for shutdown.
	for {
		select {
		case err := <-kafkaChan:
			log.Error("Kafka Error", zap.Error(err))

		case err := <-errSignals:
			log.Error("Prometheus Error", zap.Error(err))
			os.Exit(1)

		case s := <-osSignals:
			log.Debug("Worker shutdown signal", zap.String("signal", s.String()))

			// Asking prometheus shutdown and load shed.
			if err := promHandler.Shutdown(context.Background()); err != nil {
				log.Error("Graceful shutdown did not complete", zap.Error(err))
				if err := promHandler.Close(); err != nil {
					log.Fatal("Could not stop http server", zap.Error(err))
				}
			}
		}
	}
}

// ProcessArticle processes a potential article from an incoming link.
//
// For each link we need to check if it contains a valid url, from a valid source,
// and save it into the next kafka topic for further processing.
func (worker *WorkerGroup) ProcessArticle(in link.Link) error {
	timer := prometheus.NewTimer(workerProcessDuration.WithLabelValues("process"))
	defer timer.ObserveDuration()

	// check if article exists before processing it
	var f *feedsv2.Feed
	if in.Type == "twitter" {
		tf, err := feed.GetByUserName(context.Background(), worker.dbConn, in.UserName)
		if err != nil {
			workerProcessErrors.WithLabelValues("feed not found").Inc()
			worker.log.Error("Feed not found", zap.String("username", in.UserName), zap.Error(err))
			return errors.Wrap(err, "feed not found")
		}
		f = tf
	} else if in.Type == "rss" {
		tf, err := feed.GetByHostname(context.Background(), worker.dbConn, in.Hostname)
		if err != nil {
			workerProcessErrors.WithLabelValues("feed not found").Inc()
			worker.log.Error("Feed not found", zap.String("hostname", in.Hostname), zap.Error(err))
			return errors.Wrap(err, "feed not found")
		}
		f = tf
	}

	if f == nil {
		workerProcessErrors.WithLabelValues("feed not found").Inc()
		worker.log.Error("Feed not found", zap.String("hostname", in.Hostname), zap.String("url", in.Url))
		return errors.New("feed not found")
	}

	// skip if offline
	if f.Stream != nil && f.Stream.StreamStatus == commonv2.Status_STATUS_OFFLINE {
		workerProcessErrors.WithLabelValues("feed offline").Inc()
		worker.log.Warn("Skip message, feed is offline", zap.String("hostname", f.Hostname), zap.String("url", in.Url))
		return nil
	}

	// create a new article document
	a := new(articlesv2.Article)
	a.Content = new(articlesv2.Content)
	a.Nlp = new(enrichv2.NLP)
	a.DocId = in.DocId
	a.CrawledAt = in.CreatedAt
	a.Content.PublishedAt = a.CrawledAt
	a.Hostname = f.Hostname
	a.ScreenName = f.UserName
	a.Lang = f.Localization.Lang
	a.Url = in.Url
	a.Feed = f
	a.FeedId = f.Id
	if in.Title != "" {
		a.Content.Title = in.Title
	}

	// marshal the article
	b, err := json.Marshal(a)
	if err != nil {
		workerProcessErrors.WithLabelValues("marshal article error").Inc()
		worker.log.Error("Marshal article to message error", zap.Error(err))
		return err
	}
	// write article as a new message in the scrape topic, so we can process it
	go worker.Produce(kaf.Message{Value: []byte(b), Topic: "scrape"})
	return nil
}

// StoreArticle stores an article into elasticsearch and neo4j.
func (worker *WorkerGroup) StoreArticle(a *articlesv2.Article) error {
	timer := prometheus.NewTimer(workerProcessDuration.WithLabelValues("store"))
	defer timer.ObserveDuration()

	// Save the Document to Elasticsearch
	c := a
	c.Feed = &feedsv2.Feed{}
	data, err := json.Marshal(c)
	if err != nil {
		workerProcessErrors.WithLabelValues("json marshal error").Inc()
		worker.log.Error("JSON marshal error", zap.Error(err))
		return errors.Wrap(err, "json marshal error")
	}

	index, err := worker.esClient.Client.Index(
		worker.esIndex+"_"+strings.ToLower(a.Lang),
		bytes.NewReader(data),
		worker.esClient.Client.Index.WithDocumentID(a.DocId),
	)
	if err != nil {
		workerProcessErrors.WithLabelValues("index error").Inc()
		worker.log.Error("Article indexing error", zap.Error(err))
		return errors.Wrap(err, "index error")
	}

	if index.IsError() {
		workerProcessErrors.WithLabelValues("index error").Inc()
		return errors.New(index.String())
	}
	_ = index.Body.Close()

	// Create a new nodeAuthor if not exist for each author extracted
	// by svc-scraper service and return the uid
	var entities []*relationships.NodeEntity
	for _, author := range a.Content.Authors {
		uid, err := relationships.MergeNodeEntity(worker.ctx, worker.neoClient, author, "author")
		if err != nil {
			workerProcessErrors.WithLabelValues("merge authors error").Inc()
			worker.log.Error("Merge author error", zap.Error(err))
			continue
		}
		entities = append(entities, &relationships.NodeEntity{
			Uid:   uid,
			Type:  "author",
			Label: author,
		})
	}

	for _, topic := range a.Nlp.Topics {
		uid, err := relationships.MergeNodeEntity(worker.ctx, worker.neoClient, topic.Text, "topic")
		if err != nil {
			workerProcessErrors.WithLabelValues("merge topics error").Inc()
			worker.log.Error("Merge topic error", zap.Error(err))
			continue
		}
		entities = append(entities, &relationships.NodeEntity{
			Uid:   uid,
			Type:  "topic",
			Label: topic.Text,
		})
	}

	// Create a new nodeArticle
	nArticle := &relationships.NodeArticle{}
	nArticle.Uid = a.DocId
	nArticle.DocId = a.DocId
	nArticle.Lang = a.Lang
	nArticle.CrawledAt = a.CrawledAt
	nArticle.Url = a.Url
	nArticle.Title = a.Content.Title
	nArticle.PublishedAt = a.Content.PublishedAt
	nArticle.Hostname = a.Hostname
	nArticle.Type = "article"

	if _, err := relationships.CreateNodeArticle(worker.ctx, worker.neoClient, nArticle); err != nil {
		workerProcessErrors.WithLabelValues("create node article error").Inc()
		worker.log.Error("Merge article error", zap.Error(err))
		if !strings.Contains(err.Error(), "already exists with label `Article`") {
			// since the article exists we don't need to save it again
			// return nil to mark the message as read
			return errors.Wrap(err, "neo4j error")
		}
	}

	// create a new node into the neo4j database, if not exists
	fUID, err := relationships.MergeNodeFeed(worker.ctx, worker.neoClient, a.Feed)
	if err != nil {
		workerProcessErrors.WithLabelValues("merge feed error").Inc()
		// if there is an error during feed node creation return an error
		// as will not be able to associate the article with a feed.
		worker.log.Error("Merge feed failed", zap.String("hostname", a.Hostname), zap.Error(err))
		return errors.Wrap(err, "merge feed error")
	}

	// =========================================================================
	// Save relationships
	// feed published at
	if fUID != "" {
		// disable lint for the error since we don't care about the response
		// nolint:errcheck
		go func() {
			_ = relationships.MergeRel(worker.ctx, worker.neoClient, nArticle.DocId, a.FeedId, "PUBLISHED_AT")
		}()
	}
	for _, entity := range entities {
		// if entity.Type == "author" {
		// 	go relationships.MergeRel(worker.ctx, worker.neoClient, entity.Uid, nArticle.DocId, "AUTHOR_OF")
		// 	if fUID != "" {
		// 		// writes for
		// 		go relationships.MergeRel(worker.ctx, worker.neoClient, entity.Uid, fUID, "WRITES_FOR")
		// 	}
		// } else
		if entity.Type == "topic" {
			// disable lint for the error since we don't care about the response
			// nolint:errcheck
			go func() {
				_ = relationships.MergeRel(worker.ctx, worker.neoClient, nArticle.DocId, entity.Uid, "IN_TOPIC")
			}()
		}
	}

	// marshal node article
	b, err := json.Marshal(nArticle)
	if err != nil {
		workerProcessErrors.WithLabelValues("marshal article error").Inc()
		worker.log.Error("Marshal article to message error", zap.Error(err))
		return err
	}
	// write node article as a new message in the compare topic, so we can process it
	// using the compare microservice.
	go worker.Produce(kaf.Message{Value: []byte(b), Topic: "compare"})

	// set values to stream to the frontend
	a.Nlp.Stopwords = nil

	// marshal article into string
	s, err := json.Marshal(a)
	if err != nil {
		workerProcessErrors.WithLabelValues("marshal article error").Inc()
		worker.log.Error("Marshal article to message error", zap.Error(err))
		return err
	}
	// publish the article to the corresponding redis channel
	go worker.Publish("mediawatch_articles_"+strings.ToLower(a.Lang), string(s))
	return nil
}
