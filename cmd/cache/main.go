package main

import (
	"context"

	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/logger"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// ============================================================
	// Configure
	cfg := config.NewConfig()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	// ============================================================
	// Create a reusable zap logger
	sugar := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	defer sugar.Sync()

	log := sugar.Sugar()

	log.Info("[SVC-CACHE] Starting")

	// =========================================================================
	// Create mongo client
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("Register DB: %v", err)
	}
	defer dbConn.Close()

	// ============================================================
	// Redis
	rdb, err := redis.NewRedisClient(context.Background(), cfg.GetRedisURL(), "")
	if err != nil {
		log.Fatalf("Error connecting to Redis: %s", err.Error())
	}
	defer rdb.Close()

	// ============================================================
	// Get all feeds
	if feeds, err := feed.List(context.Background(), dbConn, feed.Limit(int(5000))); err == nil {
		go cacheFeeds(rdb, feeds)
	} else {
		log.Errorf("Error getting feeds: %s", err.Error())
	}
}

func cacheFeeds(rdb *redis.RedisClient, feeds *feedsv2.FeedList) {
	for _, v := range feeds.Data {
		rdb.Client.Set(context.Background(), v.Id, v, 0)
	}
}

// func cacheTrims(rdb *redis.RedisClient, key string, trims []string) {
// }
