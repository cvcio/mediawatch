package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cvcio/go-plagiarism"
	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/logger"
	articlesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kaf "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	ScrollThreshold = 240
	ScoreThreshold  = 0.25
)

var (
	compareProcessDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "consumer_process_duration_seconds",
		Help:       "Duration of consumer processing requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})
)

// CompareService responsible for comparing articles.
//
// For each article coming through kafka, compares it with other articles in the same language
// and stores the similarity in a neo4j database.
type CompareService struct {
	ctx context.Context

	cfg *config.Config
	log *zap.SugaredLogger

	esClient    *es.Elastic   // elasticsearch client
	neoClient   *neo.Neo      // neo4j client
	kafkaClient *kafka.Client // kafka client

	errChan chan error
}

// NewCompareService returns a new CompareService.
func NewCompareService(cfg *config.Config, log *zap.SugaredLogger, esClient *es.Elastic, neoClient *neo.Neo, kafkaClient *kafka.Client) *CompareService {
	return &CompareService{
		ctx: context.Background(),

		cfg: cfg,
		log: log,

		esClient:    esClient,
		neoClient:   neoClient,
		kafkaClient: kafkaClient,

		errChan: make(chan error),
	}
}

// FindAndCompare looks for articles stored in elasticsearch with similar
// keywords and tests (one-to-one) each article using the go-plagiarism
// algorithm.
//
// Read more about go-plagiarism -> https://github.com/cvcio/go-plagiarism
func (c *CompareService) FindAndCompare(id string, lang string) error {
	// retrieve the source article we want to compare from elasticsearch
	source, err := article.GetById(c.ctx, c.esClient, c.cfg.Elasticsearch.Index+"_"+strings.ToLower(lang), id)
	if err != nil {
		// in very rare occasions the document is missing, but we should not restart the service here
		c.log.Errorf("[SVC-COMPARE] failed to get document: %s", err.Error())
		c.errChan <- err
		return err
	}

	// in some occasions the article is too small or there was a problem while
	// extracting the keywords from the article using enrich microservice,
	// resulting to have only a few (<2) keywords.
	if len(source.Nlp.Keywords) < 2 {
		c.log.Debugf("[SVC-COMPARE] article (%s) too small or could't extract keywords", id)
		return nil
	}

	now, _ := time.Parse(time.RFC3339, source.CrawledAt)
	// last 3 days
	from := now.AddDate(0, 0, -3)

	opts := article.NewOpts()
	opts.Index = c.cfg.Elasticsearch.Index + "_" + strings.ToLower(lang)
	opts.Lang = lang
	opts.Limit = ScrollThreshold

	opts.Range.From = from.Format(time.RFC3339)
	opts.Range.To = time.Now().Format(time.RFC3339)
	opts.Keywords = strings.Join(source.Nlp.Keywords, " ")

	total, err := article.Count(c.ctx, c.esClient, opts)
	if err != nil {
		c.log.Errorf("[SVC-COMPARE] Error counting total potential similar: %s", err.Error())
		c.errChan <- err
		return err
	}

	if total == 0 {
		// if there are no potential similar articles return
		c.log.Debugf("[SVC-COMPARE] No similar articles found for DocId: %s", source.DocId)
		return nil
	}

	if total > ScrollThreshold {
		opts.Scroll = true
		opts.Limit = ScrollThreshold
	}

	opts.Sort.By = "content.published_at"
	opts.Sort.Order = "desc"

	articles, err := article.Search(c.ctx, c.esClient, opts)
	if err != nil {
		c.log.Errorf("[SVC-COMPARE] Error retrieving articles: %s", err.Error())
		return err
	}

	c.log.Infof("[SVC-COMPARE] PSD: %s - %d/%d", id, len(articles.Data), total)

	similar := 0
	for _, dest := range articles.Data {
		// do not compare same documents
		if source.DocId == dest.DocId {
			continue
		}
		// do not compare same documents
		if source.Url == dest.Url {
			continue
		}

		// create the plagiarism detection interface
		detector, _ := plagiarism.NewDetector(plagiarism.SetLang(strings.ToLower(source.Lang)), plagiarism.SetN(8))
		// detect with extracted stopwords
		if err := detector.DetectWithStopWords(source.Nlp.Stopwords, dest.Nlp.Stopwords); err == nil {
			// save only if the score is higher than 0.25
			if detector.Score >= ScoreThreshold {
				var a, b *articlesv2.Article
				sourceTime, _ := time.Parse(time.RFC3339, source.Content.PublishedAt)
				targetTime, _ := time.Parse(time.RFC3339, dest.Content.PublishedAt)

				// set a,b as source and dest
				a = source
				b = dest

				// swap direction if source time is after target time
				if sourceTime.Sub(targetTime).Minutes() >= 0 {
					a = dest
					b = source

					// multiple with -1
					detector.Score = detector.Score * (-1)
				}

				similar++

				// save relation to neo4j database
				go func() {
					_ = relationships.CreateSimilar(c.ctx, c.neoClient, a.DocId, b.DocId, detector.Score)
				}()
			}
		}
	}
	c.log.Infof("[SVC-COMPARE] DWS: %s - %d/%d", id, similar, len(articles.Data))

	return nil
}

// Close closes the kafka client.
func (c *CompareService) Close() {
	c.kafkaClient.Close()
}

// Consume consumes kafka topics inside an infinite loop. In our logic we need
// to fetch a message from a topic (FetchMessage), parse the json (Unmarshal)
// and find similar articles. If, for any reason, any step fails with an error
// we will skip commiting this message to kafka as we want to re-process
// this particular message again.
func (c *CompareService) Consume() {
	for {
		timer := prometheus.NewTimer(compareProcessDuration.WithLabelValues("compare"))
		// read the message from kafka topic
		m, err := c.kafkaClient.Consumer.FetchMessage(c.ctx)
		if err != nil {
			// at this point we don't have a message, as such we don't commit
			// send the error to channel
			c.errChan <- errors.Wrap(err, "failed to fetch messages from kafka")
			// go to next
			continue
		}

		// Unmarshal incoming json message
		var msg relationships.NodeArticle
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			// mark message as read (commit)
			c.Commit(m)
			// send the error to channel
			c.errChan <- errors.Wrap(err, "failed to unmarshall messages from kafka")
			// go to next
			continue
		}

		c.log.Infof("[SVC-COMPARE] FAC: %s - %s", msg.DocId, msg.CrawledAt)

		// test the article for similar
		// go worker.compare.FindAndCompare(msg.DocId, msg.Lang)
		if err := c.FindAndCompare(msg.DocId, msg.Lang); err != nil {
			c.log.Errorf("[SCV-COMPARE] FAC: %s", err.Error())
			continue
		}

		// mark message as read (commit)
		c.Commit(m)
		timer.ObserveDuration()
	}
}

// Commit commits a message to the kafka topic.
func (c *CompareService) Commit(m kaf.Message) {
	if err := c.kafkaClient.Consumer.CommitMessages(c.ctx, m); err != nil {
		c.errChan <- errors.Wrap(err, "failed to commit messages to kafka")
	}
}

// ErrorChan returns the error channel.
func (c *CompareService) ErrorChan() chan error {
	return c.errChan
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
	// Logger
	sugar := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	defer func() { _ = sugar.Sync() }()
	log := sugar.Sugar()

	log.Info("[SVC-COMPARE] Starting")

	// ========================================
	// Start elasticsearch client
	log.Info("[SVC-COMPARE] Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Register Elasticsearch: %v", err)
	}

	// ========================================
	// Start neo4j client
	log.Info("[SVC-COMPARE] Initialize Neo4J")
	neoClient, err := neo.NewNeo(cfg.Neo.BOLT, cfg.Neo.User, cfg.Neo.Pass)
	if err != nil {
		log.Fatalf("[SVC-COMPARE] Register Neo4J: %v", err)
	}
	defer func() { _ = neoClient.Client.Close() }()

	// create a reader/writer kafka connection
	kafkaClient := kafka.NewKafkaClient(
		true, false,
		cfg.GetKafkaBrokers(),
		cfg.Kafka.CompareTopic,
		cfg.Kafka.ConsumerGroupCompare,
		cfg.Kafka.WorkerOffsetOldest,
	)

	// ========================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	errChan := make(chan error, 1)

	// ========================================
	// Create a new compare service
	svc := NewCompareService(cfg, log, esClient, neoClient, kafkaClient)
	// close connections on exit
	defer svc.Close()

	// ========================================
	// Create a registry and a web server for prometheus metrics
	registry := prometheus.NewRegistry()
	registry.MustRegister(compareProcessDuration)

	// run the worker
	go func() {
		log.Info("[SVC-COMPARE] Starting kafka consumer")
		// consume messages from kafka
		svc.Consume()
	}()

	// ========================================
	// Create a new http server for prometheus metrics
	promHandler := http.Server{
		Addr:           cfg.GetPrometheusURL(),
		Handler:        promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), // api(cfg.Log.Debug, registry),
		ReadTimeout:    cfg.Prometheus.ReadTimeout,
		WriteTimeout:   cfg.Prometheus.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// start the http service
	go func() {
		log.Infof("[SVC-COMPARE] Starting prometheus web server listening %s", cfg.GetPrometheusURL())
		errChan <- promHandler.ListenAndServe()
	}()

	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Stop Service
	// Blocking main and waiting for shutdown.
	for {
		select {
		case err := <-svc.ErrorChan():
			log.Errorf("[SVC-COMPARE] CompareService Error: %s", err.Error())
		case err := <-errChan:
			// most probably the service is dead, restart and initialize the service
			log.Errorf("[SVC-COMPARE] Service Error: %s", err.Error())
			os.Exit(1)
		case s := <-osSignals:
			log.Debugf("[SVC-COMPARE] gRPC Server shutdown signal: %s", s)

			// Asking prometheus shutdown and load shed.
			if err := promHandler.Shutdown(context.Background()); err != nil {
				log.Errorf("[SVC-COMPARE] Graceful shutdown did not complete in %v: %v", cfg.Prometheus.ShutdownTimeout, err)
				if err := promHandler.Close(); err != nil {
					log.Fatalf("[SVC-COMPARE] Could not stop http server: %v", err)
				}
			}
		}
	}
}
