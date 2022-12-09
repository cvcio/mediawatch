package main

import (
	"sort"
	"time"

	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
)

type Ticker struct {
	log     *zap.SugaredLogger
	worker  *ListenGroup
	rdb     *redis.RedisClient
	targets []string
	ticker  time.Ticker
	done    chan bool
}

func NewTicker(log *zap.SugaredLogger, done chan bool, targets []string, worker *ListenGroup, rdb *redis.RedisClient) *Ticker {
	return &Ticker{
		log:     log,
		worker:  worker,
		rdb:     rdb,
		targets: targets,
		ticker:  *time.NewTicker(time.Second * 60),
		done:    done,
	}
}

func (ticker *Ticker) Fetch() {
	for _, v := range ticker.targets {
		parser := gofeed.NewParser()
		parser.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36" // "MediaWatch Bot/3.0 (mediawatch.io)"

		// parse feed
		data, err := parser.ParseURL(v)
		if err != nil {
			// TODO: investigate how often this happens
			ticker.log.Errorf("[SVC-FEEDS] Error parsing RSS feed for: %s - %s", v, err.Error())
			continue
		}
		// create a new slice with the feed items and sort by time published
		slice := data.Items
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].PublishedParsed.Before(*slice[j].PublishedParsed)
		})

		// iter over the items and check if the article is already processed.
		// we assume that the article is processed if the publish time of an item
		// is before or equal to the time stored in redis key/value store.
		for _, v := range slice {
			// get the last saved time in redis key/value store
			if last, _ := ticker.rdb.Get(data.Link); last != "" {
				lastDate, _ := time.Parse(time.RFC3339, last)
				if lastDate.After(*v.PublishedParsed) || lastDate.Equal(*v.PublishedParsed) {
					continue
				}
			}

			ticker.log.Debugf("[SVC-FEEDS] %s (%s) %s", v.PublishedParsed.Format(time.RFC3339), data.Title, v.Title)

			// write message to kafka
			// go ticker.worker.Produce(kafka.Message{})

			// update last time published per target in redis key/value store
			ticker.rdb.Set(data.Link, v.PublishedParsed.Format(time.RFC3339))
		}
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
