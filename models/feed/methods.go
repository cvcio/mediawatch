package feed

import (
	"context"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const feedsCollection = "v2_feeds"

// EnsureIndex in mongodb.
func EnsureIndex(ctx context.Context, dbConn *db.MongoDB) error {
	index := []mongo.IndexModel{
		{
			Keys: bson.M{
				"user_name": 1,
			},
			Options: options.Index().SetUnique(true), // {Unique: true},
		},
		{
			Keys: bson.M{
				"twitter_id": 1,
			},
			Options: options.Index().SetUnique(true), // {Unique: true},
		},
		{
			Keys: bsonx.Doc{
				{Key: "user_name", Value: bsonx.String("text")},
				{Key: "url", Value: bsonx.String("text")},
				{Key: "name", Value: bsonx.String("text")},
			},
			Options: options.Index().SetDefaultLanguage("en").SetLanguageOverride("el"),
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		return errors.Wrap(err, "db.feeds.ensureIndex()")
	}
	return nil
}
