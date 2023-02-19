package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cvcio/go-plagiarism"
	"github.com/kelseyhightower/envconfig"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"

	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/cvcio/mediawatch/models/deprecated/nodes"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	"github.com/cvcio/mediawatch/pkg/neo"
	kaf "github.com/segmentio/kafka-go"
)

var (
	compareProcessDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "consumer_process_duration_seconds",
		Help:       "Duration of consumer processing requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})
)

// Compare Service
type Compare struct {
	es        *es.ES
	index     string
	log       *zap.SugaredLogger
	cfg       *config.Config
	neoClient *neo.Neo
}

// CompareGroup struct.
type CompareGroup struct {
	ctx         context.Context
	log         *zap.SugaredLogger
	kafkaClient *kafka.KafkaGoClient
	errChan     chan error
	compare     *Compare
}

// Close closes the kafka client.
func (worker *CompareGroup) Close() {
	worker.kafkaClient.Close()
}

// Consume consumes kafka topics inside an infinite loop. In our logic we need
// to fetch a message from a topic (FetcMessage), parse the json (Unmarshal)
// and find similar articles. If, for any reason, any step fails with an error
// we will skip commiting this message to kafka as we want to re-process
// this particular message again.
func (worker *CompareGroup) Consume() {
	for {
		timer := prometheus.NewTimer(compareProcessDuration.WithLabelValues("compare"))
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
		var msg relationships.NodeArticle
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			// mark message as read (commit)
			worker.Commit(m)
			// send the error to channel
			worker.errChan <- errors.Wrap(err, "failed to unmarshall messages from kafka")
			// go to next
			continue
		}

		worker.log.Debugf("[SVC-COMPARE] COMPARE: %s (%s)", msg.DocId, msg.CrawledAt)

		// test the article for similar
		go worker.compare.FindAndCompare(msg.DocId)

		// mark message as read (commit)
		worker.Commit(m)
		timer.ObserveDuration()
	}
}

// Commit commits a message to the kafka topic.
func (worker *CompareGroup) Commit(m kaf.Message) {
	if err := worker.kafkaClient.Consumer.CommitMessages(worker.ctx, m); err != nil {
		worker.errChan <- errors.Wrap(err, "failed to commit messages to kafka")
	}
}

// NewCompareGroup implements a new CompareGroup struct.
func NewCompareGroup(
	log *zap.SugaredLogger,
	kafkaClient *kafka.KafkaGoClient,
	errChan chan error,
	compare *Compare,
) *CompareGroup {
	return &CompareGroup{
		context.Background(),
		log,
		kafkaClient,
		errChan,
		compare,
	}
}

func main() {
	// TODO: Change this to read from cfg
	runtime.GOMAXPROCS(2)

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
	log.Info("[SVC-COMPARE] Starting")

	// =========================================================================
	// Start elasticsearch
	log.Info("[SVC-COMPARE] Initialize Elasticsearch")
	esClient, err := es.NewElastic(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Register Elasticsearch: %v", err)
	}

	log.Info("[SVC-COMPARE] Connected to Elasticsearch")
	log.Info("[SVC-COMPARE] Check for elasticsearch indexes")
	err = es.CreateElasticIndexArticles(esClient, []string{cfg.Elasticsearch.Index})
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Index in elasticsearch: %v", err)
	}

	// =========================================================================
	// Start neo4j client
	log.Info("[SVC-COMPARE] Initialize Neo4J")
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Register Neo4J: %v", err)
	}
	log.Info("[SVC-COMPARE] Connected to Neo4J")
	defer neoClient.Client.Close()

	// Create compare struct
	compareCLient := new(Compare)
	compareCLient.es = esClient
	compareCLient.index = cfg.Elasticsearch.Index
	compareCLient.log = log
	compareCLient.cfg = cfg
	compareCLient.neoClient = neoClient

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()
	registry.MustRegister(compareProcessDuration)

	// =========================================================================
	// Create kafka consumer/producer worker

	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewGoClient(
		true, false,
		[]string{cfg.Kafka.Broker},
		cfg.Kafka.CompareTopic,
		cfg.Kafka.ConsumerGroupCompare,
		"",
		"",
		false,
	)

	// create a new worker
	worker := NewCompareGroup(
		log, kafkaGoClient, kafkaChan,
		compareCLient,
	)

	// close connections on exit
	defer worker.Close()

	// run the worker
	go func() {
		defer close(kafkaChan)
		log.Info("[SVC-COMPARE] Starting kafka consumer")
		// consume messages from kafka
		worker.Consume()
	}()

	// ========================================
	// Blocking main listening for requests
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
	// Start the http service
	go func() {
		log.Infof("[SVC-COMPARE] Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		errSingals <- promHandler.ListenAndServe()
	}()

	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Stop Service
	// Blocking main and waiting for shutdown.
	for {
		select {
		case err := <-kafkaChan:
			log.Errorf("[SVC-COMPARE] error from kafka: %s", err.Error())

		case err := <-errSingals:
			log.Errorf("[SVC-COMPARE] gRPC Server Error: %s", err.Error())
			os.Exit(1)

		case s := <-osSignals:
			log.Debugf("[SVC-COMPARE] gRPC Server shutdown signal: %s", s)

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

// FindAndCompare looks for articles stored in elasticsearch with similar
// keywords and tests (one-to-one) each article using the go-plagiarism
// algorithm.
// Read more about go-plagiarism -> https://github.com/cvcio/go-plagiarism
func (c Compare) FindAndCompare(id string) error {
	// retrive the source article we want to compare from elasticsearch
	source, err := article.Get(context.Background(), c.es, id)
	if err != nil {
		// in very rare occasions the document is missing
		c.log.Debugf("[SVC-COMPARE] failed to get document: %s", err.Error())
		return errors.Wrap(err, "failed to get document")
	}

	// in some occasions the article is too small or there was a problem while
	// extracting the keywords from the article using enrich microservice,
	// resulting to have only a few (<2) keywords.
	if len(source.NLP.Keywords) < 2 {
		c.log.Debugf("[SVC-COMPARE] article (%s) too small or could't extract keywords", id)
		return nil
	}

	now := source.CrawledAt
	// last 7 days
	from := now.AddDate(0, 0, -7)

	// create the elasticsearch query
	query := elastic.NewBoolQuery()
	queries := make([]elastic.Query, 0)
	// query only within same language
	queries = append(queries, elastic.NewQueryStringQuery(source.Lang).Field("lang"))
	// query articles with similar keywords
	queries = append(queries, elastic.NewQueryStringQuery(strings.Join(source.NLP.Keywords, " ")).Field("nlp.keywords"))
	// set query window from-now
	queries = append(queries, elastic.NewRangeQuery("content.publishedAt").Gte(from.Format(time.RFC3339)).Lte(now.Format(time.RFC3339)))
	query = query.Must(queries...)

	// select the index to query
	currentIndex := c.index

	// count total potential similar articles
	total, err := c.es.Client.Count(currentIndex).Query(query).Do(context.Background())
	if err != nil {
		c.log.Errorf("[SVC-COMPARE] Error counting total potential similar %s", err.Error())
		return errors.Wrap(err, "failed to get potential similar")
	}

	if total == 0 {
		// if there are no potential similar articles return
		c.log.Debugf("[SVC-COMPARE] No similar articles found for DocId: %s", source.DocID)
		return nil
	}

	// set size of the scroller
	SIZE := 64

	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		// Scroller
		scroll := c.es.Client.Scroll(currentIndex).Query(query).Size(SIZE)
		// Itterate
		for {
			results, err := scroll.Do(context.Background())
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				c.log.Errorf("Scroll Error: %s", err.Error())
				return err // something went wrong
			}
			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	})

	for i := 0; i < SIZE; i++ {
		g.Go(func() error {
			for hit := range hits {
				var dest article.Document
				err := json.Unmarshal(hit, &dest)
				if err != nil {
					c.log.Errorf("Unmarshal Error: %s", err.Error())
					continue
				}

				// do not compare same documents
				if source.DocID == dest.DocID {
					continue
				}
				// do not compare same documents
				if source.URL == dest.URL {
					continue
				}

				start := time.Now()
				// create the plagiarism detection interface
				detector, _ := plagiarism.NewDetector(plagiarism.SetLang(strings.ToLower(source.Lang)), plagiarism.SetN(8))
				// detect with extracted stopwords
				if err := detector.DetectWithStopWords(source.NLP.StopWords, dest.NLP.StopWords); err == nil {
					// save only if the score is higher than 0.25
					if detector.Score >= 0.25 {
						var a, b article.Document
						// set source and target article time to compare
						sourceTime := source.Content.PublishedAt
						targetTime := dest.Content.PublishedAt
						// set a,b as source and dest
						a = *source
						b = dest
						// swap direction if source time is after target time
						if sourceTime.Sub(targetTime).Minutes() >= 0 {
							a = dest
							b = *source

							// multiple with -1
							detector.Score = detector.Score * (-1)
						}

						c.log.Debugw(
							"DetectWithStopWords",
							"timeTaken", time.Since(start),
							"score", detector.Score,
							"similar", detector.Similar,
							"total", detector.Total,
							"stopwords", strings.Join(a.NLP.Keywords, ", "),
						)

						// save relation to neo4j database
						go nodes.CreateSimilar(context.Background(), c.neoClient, a, b, detector.Score)
					}
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		c.log.Errorf("WaitGroup Error: %s", err.Error())
	}
	return nil
}
