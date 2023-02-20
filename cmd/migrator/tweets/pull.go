package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	enrich_pb "github.com/cvcio/mediawatch/internal/mediawatch/enrich/v2"
	scrape_pb "github.com/cvcio/mediawatch/internal/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/models/deprecated/nodes"
	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/kafka"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	kaf "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// kubectl port-forward elasticsearch-master-0 9200:9200 9300:9300
// kubectl port-forward neo4j-neo4j-core-0 7474:7474 7687:7687 7473:7473 6362:6362
// kubectl port-forward mongo-6459898f8c-44knt 27017:27017
// kubectl port-forward kafka-0 9092:9092 9093:9093

var api *anaconda.TwitterApi

// WorkerGroup struct
type WorkerGroup struct {
	ctx         context.Context
	log         *logrus.Logger
	kafkaClient *kafka.KafkaGoClient
	errChan     chan error

	dbConn    *db.MongoDB
	esClient  *es.ES
	esIndex   string
	neoClient *neo.Neo
	scrape    scrape_pb.ScrapeServiceClient
	enrich    enrich_pb.EnrichServiceClient
}

// Close closes the kafka client
func (worker *WorkerGroup) Close() {
	worker.kafkaClient.Close()
}

func (worker *WorkerGroup) Produce(msg kaf.Message) {
	err := worker.kafkaClient.Producer.WriteMessages(worker.ctx, msg)
	if err != nil {
		worker.errChan <- errors.Wrap(err, "failed to write messages to kafka")
	}
}

func NewWorkerGroup(
	log *logrus.Logger,
	kafkaClient *kafka.KafkaGoClient,
	errChan chan error,
	dbConn *db.MongoDB,
	esClient *es.ES,
	esIndex string,
	neoClient *neo.Neo,
	scrape scrape_pb.ScrapeServiceClient,
	enrich enrich_pb.EnrichServiceClient,
) *WorkerGroup {
	ctx := context.Background()
	return &WorkerGroup{
		ctx,
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
	var screenName string
	flag.StringVar(&screenName, "screenName", "", "screenName")

	flag.Parse()

	if screenName == "" {
		fmt.Println("screenName Required")
		os.Exit(0)
	}

	// ========================================
	// Configure
	cfg := config.NewConfig()
	log := logrus.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Error loading config: %s", err.Error())
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
	if cfg.Log.Dev {
		log.SetFormatter(&logrus.TextFormatter{})
	}

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
	esClient, err := es.NewElastic(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("[SVC-WORKER] Register Elasticsearch: %v", err)
	}

	log.Info("[SVC-WORKER] Connected to Elasticsearch")
	log.Info("[SVC-WORKER] Check for elasticsearch indexes")

	// now := time.Now()
	// cfg.Elasticsearch.Index + now.Format("2006-01"),
	// cfg.Elasticsearch.Index + now.AddDate(0, 1, 0).Format("2006-01")},
	err = es.CreateElasticIndexArticles(esClient, []string{cfg.Elasticsearch.Index})
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
	// Create the gRPC Service
	// Parse Server Options
	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithInsecure())

	// Create gRPC Scrape Connection
	scrapeGRPC, err := grpc.Dial(cfg.Scrape.Host, grpcOptions...)
	if err != nil {
		log.Fatalf("[SVC-WORKER] GRPC Scrape did not connect: %v", err)
	}
	defer scrapeGRPC.Close()
	scrape := scrape_pb.NewScrapeServiceClient(scrapeGRPC)

	// Create gRPC Enrich Connection
	enrichGRPC, err := grpc.Dial(cfg.Enrich.Host, grpcOptions...)
	if err != nil {
		log.Fatalf("[SVC-WORKER] GRPC Enrich did not connect: %v", err)
	}
	defer enrichGRPC.Close()
	enrich := enrich_pb.NewEnrichServiceClient(enrichGRPC)

	// Create a new twitter client
	api = anaconda.NewTwitterApiWithCredentials(
		"",
		"",
		"",
		"")

	// =========================================================================
	// Create kafka consumer/producer worker

	// create an error channel to forward errors
	kafkaChan := make(chan error, 1)

	// create a reader/writer kafka connection
	kafkaGoClient := kafka.NewGoClient(
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
	defer worker.Close()

	fmt.Printf("Start Crawling, Please Wait...")
	tweetChan := make(chan link.CatchedURL)

	go func() {
		getTweets(screenName, 0, "", tweetChan)
		defer close(tweetChan)
	}()
	for {
		select {
		case err := <-kafkaChan:
			log.Errorf("[SVC-WORKER] error from kafka: %s", err.Error())
		// Got Url from twitter
		case u := <-tweetChan:
			log.Printf("id: %v | user: %s | time: %v | url: %s", u.TweetID, u.ScreenName, u.CreatedAt, u.URL)
			// check if article exists before processing it
			// on nil error the article exists
			id_str := fmt.Sprintf("%d", u.TweetID)
			if exists := worker.articleNodeExtist(id_str); !exists {
				// process the article
				worker.articleProcess(u)
			}
			// urlTweet, err := json.Marshal(&u)
			// if err != nil {
			// 	log.Infof("error marshal tweet data: %s", err.Error())
			// 	return
			// }
			// fmt.Println(u.TwitterUserIDStr)

			// fssd, err := feed.ByID(context.Background(), dbConn, u.TwitterUserIDStr)
			// if err != nil {
			// 	log.Fatalf("main: error getting feeds list: %v", err)
			// 	return
			// }

			// log.Debug(fssd)
			// if u.URL != "" {
			// 	err = pubsub.Publish(cfg.Listen.Pub, string(urlTweet))
			// 	if err != nil {
			// 		log.Infof("publish error: %s", err.Error())
			// 	}
			// }
		}
	}
}

func getTweets(s string, n int, id string, tweetChan chan link.CatchedURL) {
	v := url.Values{}
	v.Set("count", "200")
	v.Set("exclude_replies", "true")
	v.Set("include_rts", "false")
	if id != "" {
		v.Set("max_id", id)
	}
	v.Set("screen_name", s)

	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		fmt.Println(err)
		time.Sleep(15 * time.Minute)
		getTweets(s, n, id, tweetChan)
	}
	for _, t := range tweets {
		if isTweet(t) {
			if len(t.Entities.Urls) == 0 {
				continue
			}

			for u := range t.Entities.Urls {
				l, err := link.Parse(t.Entities.Urls[u].Expanded_url)
				if err != nil {
					continue
				}

				// 	s.log.Debugf("Try URL: (%s) %s", t.User.ScreenName, l)
				createdTime, _ := t.CreatedAtTime()
				tweetChan <- link.CatchedURL{
					ID:               uuid.New().String(),
					URL:              l,
					TweetID:          t.Id,
					TwitterUserID:    t.User.Id,
					TwitterUserIDStr: t.User.IdStr,
					ScreenName:       t.User.ScreenName,
					CreatedAt:        createdTime.Format(time.RFC3339),
				}
			}
		}
	}

	n++
	if n*200 < 3200 {
		fmt.Println(s, n, tweets[len(tweets)-1].IdStr)
		getTweets(s, n, tweets[len(tweets)-1].IdStr, tweetChan)
	} else {
		os.Exit(0)
	}
}

func isTweet(t anaconda.Tweet) bool {
	if t.InReplyToStatusIdStr == "" && t.InReplyToUserIdStr == "" &&
		t.RetweetedStatus == nil && t.QuotedStatus == nil {
		return true
	}
	return false
}

func (worker *WorkerGroup) articleProcess(in link.CatchedURL) error {
	// get the feed from mongo and create nodeFeed if not exist
	feed, err := feed.ByID(context.Background(), worker.dbConn, in.TwitterUserIDStr)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] feed not found: %s", err.Error())
		return errors.Wrap(err, "feed not found")
	}

	fUID := "d0de4efc-c19e-4c48-94d8-8e841025dd33" // left
	// "da6fd44a-07e2-4165-81fa-152e35c8298e" // avgi

	a := new(article.Document)
	a.DocID = in.ID
	a.URL = in.URL
	a.TweetID = in.TweetID
	a.TweetIDStr = fmt.Sprintf("%d", in.TweetID)
	a.ScreenName = feed.ScreenName
	a.Lang = feed.Lang
	// Twitter CreatedAt
	a.CrawledAt = in.CreatedAt

	feedMeta := map[string]string{
		"title":       feed.MetaClasses.Title,
		"excerpt":     feed.MetaClasses.Excerpt,
		"body":        feed.MetaClasses.Body,
		"authors":     feed.MetaClasses.Authors,
		"sources":     feed.MetaClasses.Sources,
		"tags":        feed.MetaClasses.Tags,
		"categories":  feed.MetaClasses.Categories,
		"publishedAt": feed.MetaClasses.PublishedAt,
		"editedAt":    feed.MetaClasses.EditedAt,
	}
	if feed.MetaClasses.FeedType == "js" {
		feedMeta["api"] = feed.MetaClasses.API
	}

	f, err := json.Marshal(feedMeta)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] failed to marshal feed's meta: %s", err.Error())
		return errors.Wrap(err, "feed marchal error")
	}

	scrapeReq := scrape_pb.ScrapeRequest{
		Feed:       string(f),
		Url:        in.URL,
		Lang:       feed.Lang,
		TweetId:    "",
		ScreenName: strings.ToLower(feed.ScreenName),
		CrawledAt:  in.CreatedAtStr,
	}

	scrapeResp, err := worker.scrape.Scrape(context.Background(), &scrapeReq)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] scrape error: %s", err.Error())
		return errors.Wrap(err, "scrape error")
	}

	worker.log.Debugf("[SVC-WORKER] SCRAPED: %s - %s", in.ScreenName, in.URL)

	a.Content.Body = scrapeResp.Data.Content.Body
	a.Content.Authors = scrapeResp.Data.Content.Authors
	a.Content.Tags = scrapeResp.Data.Content.Tags
	a.Content.PublishedAt, err = time.Parse(time.RFC3339, scrapeResp.Data.Content.PublishedAt)
	if err != nil {
		a.Content.PublishedAt = a.CrawledAt
	}

	a.Content.Title = scrapeResp.Data.Content.Title

	// Send to enrich and populate article
	enrichReq := enrich_pb.EnrichRequest{
		Body: a.Content.Body,
		Lang: feed.Lang,
	}
	enrichResp, err := worker.enrich.NLP(context.Background(), &enrichReq)
	if err != nil {
		// TODO: check what to do here
		worker.log.Debugf("[SVC-WORKER] enrich error: %s", err.Error())
		return errors.Wrap(err, "enrich error")
	}

	worker.log.Debugf("[SVC-WORKER] ENRICH:  %s - %s", in.ScreenName, in.URL)

	a.NLP.Summary = enrichResp.Data.Nlp.Summary
	a.NLP.Keywords = enrichResp.Data.Nlp.Keywords
	a.NLP.StopWords = enrichResp.Data.Nlp.Stopwords
	// a.NLP.Topics = enrichResp.Data.Nlp.Topics
	// a.NLP.Quotes = enrichResp.Data.Nlp.Quotes
	// a.NLP.Claims = enrichResp.Data.Nlp.Claims
	// for _, entity := range enrichResp.Data.Nlp.Entities {
	// 	a.NLP.Entities = append(a.NLP.Entities, &article.Entity{
	// 		EntityText: entity.EntityText,
	// 		EntityType: entity.EntityType,
	// 	})
	// }

	// =========================================================================
	// Save the Document to Elasticsearch
	_, err = worker.esClient.Client.Index().
		Index(worker.esIndex).
		Id(a.DocID).
		BodyJson(a).
		Do(context.Background())
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] article indexing error: %s", err.Error())
		return errors.Wrap(err, "index error")
	}

	worker.log.Debugf("[SVC-WORKER] INDEXED: %s - %s", in.ScreenName, in.URL)

	// =========================================================================
	// Create a new nodeAuthor if not exist for each Atuhor extracted
	// by svc-scraper service and return the uid
	var authors []*nodes.NodeAuthor
	for _, author := range a.Content.Authors {
		uid, err := worker.mergeNodeAuthor(author)
		if err != nil {
			worker.log.Debugf("[SVC-WORKER] merge author error: %s", err.Error())
			continue
		}
		authors = append(authors, &nodes.NodeAuthor{
			UID: uid,
		})
	}

	// =========================================================================
	// Create a new nodeArticle
	nArticle := &nodes.NodeArticle{}
	nArticle.DocID = a.DocID
	nArticle.Lang = a.Lang
	nArticle.CrawledAt = a.CrawledAt
	nArticle.URL = a.URL
	nArticle.TweetID = a.TweetID
	nArticle.TweetIDStr = a.TweetIDStr
	nArticle.Title = a.Content.Title
	nArticle.Body = a.Content.Body
	nArticle.Summary = a.NLP.Summary
	nArticle.Tags = a.Content.Tags
	nArticle.Categories = a.Content.Categories
	nArticle.PublishedAt = a.Content.PublishedAt
	nArticle.EditedAt = a.Content.EditedAt
	nArticle.Keywords = a.NLP.Keywords
	nArticle.Topics = a.NLP.Topics

	resNeo, err := worker.createArticle(nArticle)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] merge article error: %s", err.Error())
		if !strings.Contains(err.Error(), "already exists with label `Article`") {
			// since the article exists we don't need to save it again
			// return nil to mark the message as read
			return errors.Wrap(err, "neo4j error")
		}
	}

	// relations
	if fUID != "" {
		go worker.mergeRel(nArticle.DocID, fUID, "PUBLISHED_AT")
	}

	for _, author := range authors {
		go worker.mergeRel(author.UID, nArticle.DocID, "AUTHOR_OF")
		if fUID != "" {
			go worker.mergeRel(author.UID, fUID, "WRITES_FOR")
		}
	}

	b, err := json.Marshal(nArticle)
	if err != nil {
		worker.log.Debugf("[SVC-WORKER] marshal article to message error: %s", err.Error())
		return err
	}
	go worker.Produce(kaf.Message{Value: []byte(b)})

	worker.log.Debugf("[SVC-WORKER] SAVED:   %s - %s", in.ScreenName, resNeo)

	return nil
}

func (worker *WorkerGroup) createArticle(article *nodes.NodeArticle) (string, error) {
	session := worker.neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(nodes.CreateNodeArticleTxFunc(article))

	if err != nil {
		return "", err
	}

	return res.(string), nil
}

func (worker *WorkerGroup) mergeNodeFeed(f *feed.Feed) (string, error) {
	session := worker.neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(nodes.MergeNodeFeedTxFunc(f))

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (worker *WorkerGroup) articleNodeExtist(docId string) bool {
	session := worker.neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	res, err := session.ReadTransaction(nodes.ExistsTxFunc(docId))
	if err != nil {
		return false
	}
	return res != nil
}

func (worker *WorkerGroup) mergeNodeAuthor(author string) (string, error) {
	session := worker.neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(nodes.MergeNodeAuthorTxFunc(author))

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func (worker *WorkerGroup) mergeRel(source, dest, rel string) error {
	session := worker.neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	var f neo4j.TransactionWork

	switch rel {
	case "PUBLISHED_AT":
		f = nodes.CreatePublishedAtTxFunc(source, dest)
	case "AUTHOR_OF":
		f = nodes.CreateAuthorOfTxFunc(source, dest)
	case "WRITES_FOR":
		f = nodes.CreateWritesForTxFunc(source, dest)
	case "HAS_ENTITY":
		f = nodes.CreateHasEntityTxFunc(source, dest)
	}

	if _, err := session.WriteTransaction(f); err != nil {
		return err
	}
	return nil
}
