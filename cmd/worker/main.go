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

	articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	enrich_pb "github.com/cvcio/mediawatch/internal/mediawatch/enrich/v2"
	scrape_pb "github.com/cvcio/mediawatch/internal/mediawatch/scrape/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/cvcio/mediawatch/pkg/neo"
	kaf "github.com/segmentio/kafka-go"
)

var (
	grpcDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "grpc_response_time_seconds",
		Help:       "Duration of GRPC requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})

	workerProcessDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "consumer_process_duration_seconds",
		Help:       "Duration of consumer processing requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})
)

// WorkerGroup struct.
type WorkerGroup struct {
	ctx         context.Context
	log         *zap.SugaredLogger
	kafkaClient *kafka.KafkaClient
	errChan     chan error

	dbConn    *db.MongoDB
	esClient  *es.Elastic
	esIndex   string
	neoClient *neo.Neo
	scrape    scrape_pb.ScrapeServiceClient
	enrich    enrich_pb.EnrichServiceClient
}

// Close closes the kafka client.
func (worker *WorkerGroup) Close() {
	worker.kafkaClient.Close()
}

// Consume consumes kafka topics inside an infinite loop. In our logic we need
// to fetch a message from a topic (FetcMessage), parse the json (Unmarshal)
// and process the content (articleProcess) if it doesn't already exists.
// If, for any reason, any step fails with an error we will commit this message
// to kafka as we don't want to process this particular message again (it failed
// for some reason).
func (worker *WorkerGroup) Consume() {
	for {
		timer := prometheus.NewTimer(workerProcessDuration.WithLabelValues("worker"))
		//
		// fetch the message from kafka topic
		m, err := worker.kafkaClient.Consumer.FetchMessage(worker.ctx)
		if err != nil {
			// at this point we don't have a message, as such we don't commit
			// send the error to channel
			worker.errChan <- errors.Wrap(err, "failed to fetch messages from kafka")
			// go to next
			continue
		}

		// Unmarshal incoming json message
		var msg link.CatchedURL
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			// mark message as read (commit)
			worker.Commit(m)
			// send the error to channel
			worker.errChan <- errors.Wrap(err, "failed to unmarshall messages from kafka")
			// go to next
			continue
		}

		// re-validate link and make sure it is a valid url
		if _, err := link.Validate(msg.URL); err != nil {
			// mark message as read (commit)
			worker.Commit(m)
			// send the error to channel
			worker.errChan <- errors.Wrap(err, "url is not valid, skipping")
			// go to next
			continue
		}

		worker.log.Debugf("[SVC-WORKER] CONSUME: %s - %s - %s", msg.ScreenName, msg.CreatedAt, msg.ID)

		// check if article exists before processing it
		// on nil error the article exists
		// if exists := nodes.ArticleNodeExtist(worker.ctx, worker.neoClient, fmt.Sprintf("%d", msg.TweetID)); !exists {
		// process the article
		if err := worker.ProcessArticle(msg); err != nil {
			worker.log.Debugf("[SVC-WORKER] ERRORED: %s - %s", msg.ScreenName, err.Error())
			// send the error to channel
			worker.errChan <- errors.Wrap(err, "failed process article")
			continue
		}
		// }

		// mark message as read (commit)
		worker.Commit(m)
		timer.ObserveDuration()
	}
}

// Commit commits a message to the kafka topic.
func (worker *WorkerGroup) Commit(m kaf.Message) {
	if err := worker.kafkaClient.Consumer.CommitMessages(worker.ctx, m); err != nil {
		worker.errChan <- errors.Wrap(err, "failed to commit messages to kafka")
	}
}

// Procuce writes a new message to the kafka topic.
func (worker *WorkerGroup) Produce(msg kaf.Message) {
	err := worker.kafkaClient.Producer.WriteMessages(worker.ctx, msg)
	if err != nil {
		worker.errChan <- errors.Wrap(err, "failed to write messages to kafka")
	}
}

// NewWorkerGroup implements a new WorkerGroup struct.
func NewWorkerGroup(
	log *zap.SugaredLogger,
	kafkaClient *kafka.KafkaClient,
	errChan chan error,
	dbConn *db.MongoDB,
	esClient *es.Elastic,
	esIndex string,
	neoClient *neo.Neo,
	scrape scrape_pb.ScrapeServiceClient,
	enrich enrich_pb.EnrichServiceClient,
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
		scrape,
		enrich,
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
	log.Info("[SVC-WORKER] Starting")

	// =========================================================================
	// Create mongo client
	log.Info("[SVC-WORKER] Initialize Mongo")
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Register DB: %v", err)
	}
	log.Info("[SVC-WORKER] Connected to Mongo")
	defer dbConn.Close()

	// =========================================================================
	// Start elasticsearch
	log.Info("[SVC-WORKER] Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Register Elasticsearch: %v", err)
	}

	log.Info("[SVC-WORKER] Connected to Elasticsearch")
	log.Info("[SVC-WORKER] Check for elasticsearch indexes")

	err = esClient.CreateElasticIndexWithLanguages(cfg.Elasticsearch.Index, cfg.Langs)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Index in elasticsearch: %v", err)
	}

	// =========================================================================
	// Start neo4j client
	log.Info("[SVC-WORKER] Initialize Neo4J")
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Register Neo4J: %v", err)
	}
	log.Info("[SVC-WORKER] Connected to Neo4J")
	defer neoClient.Client.Close()

	// =========================================================================
	// Create the gRPC service clients
	// Parse Server Options
	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithInsecure())

	// Create gRPC Scrape Connection
	scrapeGRPC, err := grpc.Dial(cfg.Scrape.Host, grpcOptions...)
	if err != nil {
		log.Fatalf("[SVC-WORKER] GRPC Scrape did not connect: %v", err)
	}
	defer scrapeGRPC.Close()
	// Create gRPC Scrape client
	scrape := scrape_pb.NewScrapeServiceClient(scrapeGRPC)

	// Create gRPC Enrich Connection
	enrichGRPC, err := grpc.Dial(cfg.Enrich.Host, grpcOptions...)
	if err != nil {
		log.Fatalf("[SVC-WORKER] GRPC Enrich did not connect: %v", err)
	}
	defer enrichGRPC.Close()
	// Create gRPC Enrich Connection
	enrich := enrich_pb.NewEnrichServiceClient(enrichGRPC)

	// =========================================================================
	// Create kafka consumer/producer worker

	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewKafkaClient(
		true, true,
		[]string{cfg.Kafka.Broker},
		cfg.Kafka.WorkerTopic,
		cfg.Kafka.ConsumerGroupWorker,
		cfg.Kafka.CompareTopic,
		cfg.Kafka.ConsumerGroupCompare,
		cfg.Kafka.WorkerOffsetOldest,
	)

	// create a new worker
	worker := NewWorkerGroup(
		log, kafkaGoClient, kafkaChan,
		dbConn, esClient, cfg.Elasticsearch.Index, neoClient, scrape, enrich,
	)

	// close connections on exit
	// defer worker.Close()

	// run the worker
	go func() {
		defer close(kafkaChan)
		log.Info("[SVC-WORKER] Starting kafka consumer")
		// consume messages from kafka
		worker.Consume()
	}()

	// ========================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	errSingals := make(chan error, 1)

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()
	registry.MustRegister(grpcDuration)
	registry.MustRegister(workerProcessDuration)

	// create a prometheus http.Server
	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), // api(cfg.Log.Debug, registry),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// ========================================
	// Start the http service for prometheus
	go func() {
		log.Infof("[SVC-WORKER] Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		errSingals <- promHandler.ListenAndServe()
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
			log.Errorf("[SVC-WORKER] error from kafka: %s", err.Error())
			continue

		case err := <-errSingals:
			log.Errorf("[SVC-WORKER] Error: %s", err.Error())
			os.Exit(1)

		case s := <-osSignals:
			log.Debugf("[SVC-WORKER] gRPC Server shutdown signal: %s", s)

			// Asking prometheus to shutdown and load shed.
			if err := promHandler.Shutdown(context.Background()); err != nil {
				log.Errorf("[SVC-WORKER] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
				if err := promHandler.Close(); err != nil {
					log.Fatalf("[SVC-WORKER] Could not stop http server: %v", err)
				}
			}
		}
	}
}

// ProcessArticle processes an incoming link (potential article). At this step we process
// incoming links catched by the listen microservice. If there is an error at any point and
// for any reason, it will return an error. The message is always commited to the kafka topic
// even if it fails, otherwise we save the article in the elasticsearch index and any relations
// in the neo4j database. To save the final article object we first need to retrieve the
// corresponding feed from the mongo database, create a new article object to save within
// elasticsearch, scrape the article using the scraper microservice, enrich/extract contextual
// information using the enrich microservice, compare the body with other articles using the
// compare microservice and finally write to storage.
func (worker *WorkerGroup) ProcessArticle(in link.CatchedURL) error {
	// retrieve feed infoirmation (language, id, etc.)
	feed, err := feed.ByScreenName(context.Background(), worker.dbConn, in.ScreenName)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] feed not found: %s", err.Error())
		return errors.Wrap(err, "feed not found")
	}

	// create a new node into the neo4j database, if not exists
	fUID, err := relationships.MergeNodeFeed(worker.ctx, worker.neoClient, feed)
	if err != nil {
		// if there is an error during feed node creation return an error
		// as will not be able to associate the article with a feed.
		worker.log.Debugf("[SVC-WORKER] merge feed failed: %s", err.Error())
		return errors.Wrap(err, "merge feed error")
	}

	// create a new article document
	a := new(articlesv2.Article)
	a.Content = new(articlesv2.Content)
	a.Nlp = new(enrich_pb.NLP)

	a.DocId = in.ID
	a.Url = in.URL
	a.ScreenName = feed.ScreenName
	a.Hostname = feed.Hostname()
	a.Lang = feed.Lang
	a.CrawledAt = in.CreatedAt

	// =========================================================================
	timer := prometheus.NewTimer(grpcDuration.WithLabelValues("scraper"))

	// create the scrape request
	scrapeReq := scrape_pb.ScrapeRequest{
		Url:        in.URL,
		Lang:       feed.Lang,
		ScreenName: strings.ToLower(feed.ScreenName),
		CrawledAt:  in.CreatedAt,
	}

	// scrape the article
	scrapeResp, err := worker.scrape.Scrape(context.Background(), &scrapeReq)
	if err != nil {
		// if there is an error while scraping, return.
		worker.log.Debugf("[SVC-WORKER] scrape error: %s", err.Error())
		return errors.Wrap(err, "scrape error")
	}
	timer.ObserveDuration()

	worker.log.Debugf("[SVC-WORKER] SCRAPED: %s - %s", in.ScreenName, in.URL)

	// set scraped data to content
	a.Content.Body = scrapeResp.Data.Content.Body
	a.Content.Authors = scrapeResp.Data.Content.Authors
	a.Content.Tags = scrapeResp.Data.Content.Tags
	if in.Title != "" {
		a.Content.Title = in.Title
	} else {
		a.Content.Title = scrapeResp.Data.Content.Title
	}
	a.Content.Excerpt = scrapeResp.Data.Content.Description
	a.Content.Image = scrapeResp.Data.Content.Image

	// make sure to parse the datetime object
	if _, err := time.Parse(time.RFC3339, scrapeResp.Data.Content.PublishedAt); err == nil {
		a.Content.PublishedAt = scrapeResp.Data.Content.PublishedAt
	} else {
		// otherwise set published datetime to crawled time
		a.Content.PublishedAt = a.CrawledAt
	}

	// =========================================================================
	timer = prometheus.NewTimer(grpcDuration.WithLabelValues("enrich"))

	// create the enrich request
	enrichReq := enrich_pb.EnrichRequest{
		Body: a.Content.Body,
		Lang: feed.Lang,
	}

	// send to enrich
	enrichResp, err := worker.enrich.NLP(context.Background(), &enrichReq)
	if err != nil {
		// TODO: check what to do here
		worker.log.Debugf("[SVC-WORKER] enrich error: %s", err.Error())
		return errors.Wrap(err, "enrich error")
	}

	worker.log.Debugf("[SVC-WORKER] ENRICH:  %s - %s", in.ScreenName, in.URL)

	// set enriched data to nlp
	a.Nlp.Summary = enrichResp.Data.Nlp.Summary
	a.Nlp.Keywords = enrichResp.Data.Nlp.Keywords
	a.Nlp.Stopwords = enrichResp.Data.Nlp.Stopwords
	a.Nlp.Topics = enrichResp.Data.Nlp.Topics
	a.Nlp.Quotes = enrichResp.Data.Nlp.Quotes
	a.Nlp.Claims = enrichResp.Data.Nlp.Claims
	a.Nlp.Entities = enrichResp.Data.Nlp.Entities

	timer.ObserveDuration()

	// =========================================================================
	// Save the Document to Elasticsearch
	data, err := json.Marshal(a)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] json marshal error: %s", err.Error())
		return errors.Wrap(err, "json marshal error")
	}

	index, err := worker.esClient.Client.Index(
		worker.esIndex+"_"+strings.ToLower(a.Lang),
		bytes.NewReader(data),
		worker.esClient.Client.Index.WithDocumentID(a.DocId),
	)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] article indexing error: %s", err.Error())
		return errors.Wrap(err, "index error")
	}

	// retrun on response error
	if index.IsError() {
		return errors.New(index.String())
	}
	index.Body.Close()
	worker.log.Debugf("[SVC-WORKER] INDEXED: %s - %s", in.ScreenName, in.URL)

	// =========================================================================
	// Create a new nodeAuthor if not exist for each Atuhor extracted
	// by svc-scraper service and return the uid
	var entities []*relationships.NodeEntity
	for _, author := range a.Content.Authors {
		uid, err := relationships.MergeNodeEntity(worker.ctx, worker.neoClient, author, "author")
		if err != nil {
			worker.log.Debugf("[SVC-WORKER] merge author error: %s", err.Error())
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
			worker.log.Debugf("[SVC-WORKER] merge topic error: %s", err.Error())
			continue
		}
		entities = append(entities, &relationships.NodeEntity{
			Uid:   uid,
			Type:  "topic",
			Label: topic.Text,
		})
	}

	// =========================================================================
	// Create a new nodeArticle
	nArticle := &relationships.NodeArticle{}
	nArticle.Uid = a.DocId
	nArticle.DocId = a.DocId
	nArticle.Lang = a.Lang
	nArticle.CrawledAt = a.CrawledAt
	nArticle.Url = a.Url
	nArticle.Title = a.Content.Title
	nArticle.PublishedAt = a.Content.PublishedAt
	nArticle.ScreenName = feed.ScreenName
	nArticle.Type = "article"

	resNeo, err := relationships.CreateNodeArticle(worker.ctx, worker.neoClient, nArticle)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] merge article error: %s", err.Error())
		if !strings.Contains(err.Error(), "already exists with label `Article`") {
			// since the article exists we don't need to save it again
			// return nil to mark the message as read
			return errors.Wrap(err, "neo4j error")
		}
	}

	// =========================================================================
	// Save relationships
	// feed published at
	if fUID != "" {
		go relationships.MergeRel(worker.ctx, worker.neoClient, nArticle.DocId, fUID, "PUBLISHED_AT")
	}
	for _, entity := range entities {
		if entity.Type == "author" {
			go relationships.MergeRel(worker.ctx, worker.neoClient, entity.Uid, nArticle.DocId, "AUTHOR_OF")
			if fUID != "" {
				// writes for
				go relationships.MergeRel(worker.ctx, worker.neoClient, entity.Uid, fUID, "WRITES_FOR")
			}
		} else if entity.Type == "topic" {
			go relationships.MergeRel(worker.ctx, worker.neoClient, nArticle.DocId, entity.Uid, "IN_TOPIC")
		}
	}

	// marshal node article
	b, err := json.Marshal(nArticle)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] marshal article to message error: %s", err.Error())
		return err
	}

	// write node article as a new message in the compare topic, so we can process it
	// using the compare microservice.
	go worker.Produce(kaf.Message{Value: []byte(b)})

	worker.log.Debugf("[SVC-WORKER] SAVED: %s - %s", in.ScreenName, resNeo)

	return nil
}
