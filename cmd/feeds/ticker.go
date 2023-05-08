package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/cvcio/mediawatch/models/link"
	"github.com/cvcio/mediawatch/pkg/helper"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Ticker struct {
	log      *zap.SugaredLogger
	worker   *ListenGroup
	rdb      *redis.RedisClient
	proxy    *http.Client
	ticker   time.Ticker
	done     chan bool
	targets  []*feedsv2.Feed
	init     bool
	interval time.Duration
}

type CacheLast struct {
	Id              string    `json:"feed_id"`
	Hostname        string    `json:"hostname"`
	LastArticleAt   time.Time `json:"last_article_at"`
	LastArticleLink string    `json:"last_article_link"`
}

func NewTicker(log *zap.SugaredLogger, worker *ListenGroup, rdb *redis.RedisClient, proxy *http.Client, done chan bool, targets []*feedsv2.Feed, init bool, interval time.Duration) *Ticker {
	return &Ticker{
		log:      log,
		worker:   worker,
		rdb:      rdb,
		proxy:    proxy,
		ticker:   *time.NewTicker(interval),
		done:     done,
		targets:  targets,
		init:     init,
		interval: interval,
	}
}

func (ticker *Ticker) Fetch() {
	// delay := time.Duration((ticker.interval / time.Duration(math.Ceil(float64(len(ticker.targets))/100))) / 100)
	for _, v := range ticker.targets {
		if _, err := url.Parse(v.Stream.StreamTarget); err != nil {
			ticker.rdb.Set("feed:status:"+v.Id, "offline", time.Hour*3)
			ticker.log.Errorf("Unable to validate URL: %s", v.Stream.StreamTarget)
			continue
		}
		if status, _ := ticker.rdb.Get("feed:status:" + v.Id); status == "offline" {
			continue
		}
		parser := gofeed.NewParser()
		// TODO: Find a way to use a proxy for the reqursts, without getting back too many 403s. Using Tor works, but with too many errors.
		if v.Stream.RequiresProxy {
			parser.Client = ticker.proxy
		}
		parser.UserAgent = helper.RandomUserAgent() // "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36" // "MediaWatch Bot/3.0 (mediawatch.io)"

		// parse feed
		data, err := parser.ParseURL(v.Stream.StreamTarget)
		if err != nil {
			// TODO: Investigate how often this happens
			// TODO: Add prometheus metrics with error codes per feed
			ticker.rdb.Set("feed:status:"+v.Id, "offline", time.Hour*3)
			ticker.log.Errorf("Error parsing RSS feed for: (%s) - %s", v.Hostname, err.Error())
			continue
		}

		if &data.Items == nil || len(data.Items) == 0 {
			continue
		}

		// create a new slice with the feed items and sort by time published
		slice := data.Items
		sort.Slice(slice, func(i, j int) bool {
			if slice[i].PublishedParsed == nil || slice[j].PublishedParsed == nil {
				return false
			}
			return slice[i].PublishedParsed.Before(*slice[j].PublishedParsed)
		})
		// populate all urls in a list and send in a go routine
		var urls []link.CatchedURL
		// iter over the items and check if the article is already processed.
		// we assume that the article is processed if the publish time of an item
		// is before or equal to the time stored in redis key/value store.
		for _, l := range slice {
			if l.PublishedParsed == nil {
				// Probably the item is empty, skip it.
				continue
			}

			timePublished := l.PublishedParsed.Truncate(time.Millisecond)

			// get the last saved time in redis key/value store
			if last, _ := ticker.rdb.Get("feed:last:" + v.Id); last != "" {
				var lastCache CacheLast
				if err := json.Unmarshal([]byte(last), &lastCache); err != nil {
					ticker.log.Errorf("Unable to unmarshal cache: %v", last)
					continue
				}

				if lastCache.LastArticleAt.Before(timePublished) == false {
					continue
				}

				if lastCache.LastArticleLink == l.Link {
					continue
				}
			}

			ticker.log.Debugf("%s (%s) %s", timePublished.Format(time.RFC3339), v.Hostname, l.Title)

			catchedURL := link.CatchedURL{
				DocId:     uuid.New().String(),
				Type:      "rss",
				Url:       l.Link,
				CreatedAt: timePublished.Format(time.RFC3339),
				Title:     l.Title,
				UserName:  v.UserName,
				Hostname:  v.Hostname,
			}

			newCache, err := json.Marshal(&CacheLast{
				Id:              v.Id,
				Hostname:        v.Hostname,
				LastArticleAt:   timePublished,
				LastArticleLink: catchedURL.Url,
			})
			if err != nil {
				ticker.log.Errorf("Unable to marshal cache: %s", err.Error())
				continue
			}

			// update last time published per target in redis key/value store
			ticker.rdb.Set("feed:last:"+v.Id, string(newCache), 0)

			if ticker.init {
				continue
			}

			urls = append(urls, catchedURL)
		}

		if len(urls) > 0 {
			go ticker.Produce(urls)
			// time.Sleep(delay)
		}
	}
}

func (ticker *Ticker) Produce(urls []link.CatchedURL) {
	for _, v := range urls {
		// write message to kafka
		message, err := json.Marshal(&v)
		if err != nil {
			ticker.log.Errorf("Unable to marshal message: %s", err.Error())
			continue
		}
		go ticker.worker.Produce(kafka.Message{
			Value: []byte(message),
		})

		time.Sleep(1 * time.Second)
	}
}

func (ticker *Ticker) Tick() {
	for {
		select {
		case <-ticker.done:
			return
		case <-ticker.ticker.C:
			go ticker.Fetch()
		}
	}
}
