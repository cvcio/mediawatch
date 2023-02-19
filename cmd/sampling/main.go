package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/kelseyhightower/envconfig"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	field := flag.String("field", "text", "field to sample")
	output := flag.String("output", "tmp/sample.jsonl", "outputfile to write")
	index := flag.String("index", "articles", "elastic search index")
	format := flag.String("format", "txt", "output format")
	flag.Parse()
	// ========================================
	// Configure
	cfg := config.NewConfig()
	log := logrus.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("main: Error loading config: %s", err.Error())
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
	log.SetFormatter(&logrus.TextFormatter{})
	if cfg.Log.Dev {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	log.Info("main: Starting")

	// =========================================================================
	// Start elasticsearch
	log.Info("main: Initialize Elasticsearch V7")
	esClient, err := es.NewElastic(cfg.Elasticsearch.Host, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("main: Register Elasticsearch V7: %v", err)
	}

	log.Info("main: Connected to Elasticsearch V7")

	if cfg.Log.Dev {
		log.Info("main: Check for elasticsearch indexes")
		err = es.CreateElasticIndex(esClient)
		if err != nil {
			log.Fatalf("main: Index in elasticsearch V7: %v", err)
		}
	}

	f, err := os.Create(*output)
	if err != nil {
		log.Fatalf("main: cant create file: %v", err)
	}
	defer f.Close()
	// =========================================================================
	// Query elasticsearch
	query := elastic.NewBoolQuery()
	// query = query.Must(
	// 	elastic.NewRangeQuery("crawledAt").Gte("2022-02-01"),
	// )

	C := 0
	total, _ := esClient.Client.Count(*index).Query(query).Do(context.Background())
	if total == 0 {
		return
	}

	log.Infof("[SAMPLING] Total Docs in Elasticsearch to migrate %d", total)
	log.Infoln("[SAMPLING] Start Scrolling Queries...")

	SIZE := 1000
	begin := time.Now()
	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		// Scroller
		log.Infof("[SAMPLING] Scroll %d remaining Docs %d", C, (int(total) - ((C + 1) * SIZE)))
		scroll := esClient.Client.Scroll(*index).Query(query).Size(SIZE)
		C++
		// Iterate
		for {
			results, err := scroll.Do(context.Background())
			if err == io.EOF {
				log.Infoln("[SAMPLING] Done...!")
				f.Sync()
				return nil // all results retrieved
			}
			if err != nil {
				log.Errorln("[ERROR] [SAMPLING] Error :", err)
				continue
				//return err // something went wrong
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

	type model struct {
		DocID     string   `json:"docId"`
		Text      string   `json:"text"`
		CreatedAt string   `json:"createdAt"`
		Source    string   `json:"source"`
		Link      string   `json:"link"`
		Tags      []string `json:"tags"`
		Topics    []string `json:"topics"`
	}

	var t uint64
	// Init Go Routines for Search Results
	for i := 0; i < 500; i++ {
		g.Go(func() error {
			for hit := range hits {
				if hit == nil {
					return nil
				}
				current := atomic.AddUint64(&t, 1)
				dur := time.Since(begin).Seconds()
				sec := int(dur)
				pps := int64(float64(current) / dur)
				log.Infof("%10d | %6d req/s | %02d:%02d\r", current, pps, sec/60, sec%60)

				// Do Migrate
				var doc article.Document
				err := json.Unmarshal(hit, &doc)
				if err != nil {
					continue
				}

				text := doc.Content.Body
				if *field == "summary" {
					text = doc.NLP.Summary
				}

				if *field == "title" {
					text = doc.Content.Title
				}

				m := &model{
					DocID:     doc.DocID,
					CreatedAt: doc.CrawledAt.String(),
					Source:    doc.ScreenName,
					Link:      doc.URL,
					Text:      text,
					Tags:      doc.Content.Tags,
					Topics:    doc.NLP.Topics,
				}

				if *format == "txt" {
					_, err = f.WriteString(m.Text + "\n")
					if err != nil {
						continue
					}
				}

				if *format == "json" {
					b, _ := json.Marshal(m)
					_, err = f.WriteString(string(b) + "\n")
					if err != nil {
						continue
					}
				}

				// Terminate
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		log.Print(err)
	}
}
