package main

import (
	"context"
	"net/url"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/twitter"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
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
	log.SetFormatter(&logrus.JSONFormatter{})
	if cfg.Log.Dev {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	log.Info("main: Starting")
	// ============================================== ==============
	// Start Mongo
	log.Info("main: Initialize Mongo")
	dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("main: Register DB: %v", err)
	}
	log.Info("main: Connected to Mongo")
	defer dbConn.Close()

	// Get feeds list
	feeds, err := feed.List(context.Background(), dbConn, feed.Status("active"), feed.Limit(0))
	if err != nil {
		log.Fatalf("main: error getting feeds list: %v", err)
	}

	// Create a new twitter client
	twtt, err := twitter.NewAPI("",
		"", "",
		"")
	if err != nil {
		log.Fatalf("Error connecting to twitter: %s", err.Error())
	}

	for _, f := range feeds.Data {
		if err != nil {
			log.Println(err)
			continue
		}

		user, err := twtt.GetUsersShow(f.ScreenName, url.Values{})
		if err != nil {
			log.Debug(err)
			continue
		}

		f.TwitterIDStr = user.IdStr
		f.TwitterProfileImage = user.ProfileImageUrlHttps

		nf := &feed.UpdateFeed{}
		nf.TwitterIDStr = &f.TwitterIDStr
		nf.TwitterProfileImage = &f.TwitterProfileImage

		log.Println(f.TwitterID, f.TwitterIDStr)
		time.Sleep(500 * time.Millisecond)
		ctx := context.TODO()
		err = feed.Update(ctx, dbConn, f.ID.Hex(), nf, time.Now())
		if err != nil {
			log.Debug(err)
			return
		}
	}
}
