package feed

import (
	"context"
	"fmt"
	"time"

	feedsv2 "github.com/cvcio/mediawatch/internal/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
				"hostname": 1,
			},
			Options: options.Index().SetUnique(true), // {Unique: true},
		},
		{
			Keys: bsonx.Doc{
				{Key: "user_name", Value: bsonx.String("text")},
				{Key: "hostname", Value: bsonx.String("text")},
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

func Get()           {}
func GetById()       {}
func GetByUserName() {}
func GetTargets()    {}
func List()          {}

// Create creates a new feed.
func Create(ctx context.Context, mg *db.MongoDB, feed *feedsv2.Feed) (*feedsv2.Feed, error) {
	feed.CreatedAt = time.Now().Format(time.RFC3339)

	f := func(collection *mongo.Collection) error {
		inserted, err := collection.InsertOne(ctx, &feed) // (&u)
		// Normally we would return ErrDuplicateKey in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err != nil {
			we, _ := err.(mongo.WriteException)
			if we.WriteErrors[0].Code == 11000 {
				return db.ErrInvalid
			}

			return err
		}

		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			feed.Id = oid.Hex()
		} else {
			return db.ErrInvalid
		}

		return nil
	}

	if err := mg.Execute(ctx, feedsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.insert(%s)", feed.Name))
	}

	return feed, nil
}

// Update updates a feed.
func Update(ctx context.Context, mg *db.MongoDB, feed *feedsv2.Feed) error {
	// convert id string to ObjectId
	oid, err := primitive.ObjectIDFromHex(feed.Id)
	if err != nil {
		return db.ErrInvalidID
	}

	// create the fields to update
	fields := make(bson.M)

	return nil
}

func Delete() {}
