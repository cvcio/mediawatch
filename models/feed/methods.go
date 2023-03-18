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

// Get returns a single feed by feed query.
func Get(ctx context.Context, mg *db.MongoDB, optionsList ...func(*ListOpts)) (*feedsv2.Feed, error) {
	filter := bson.M{}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	if opts.Id != "" {
		oid, err := primitive.ObjectIDFromHex(opts.Id)
		if err != nil {
			return nil, db.ErrInvalid
		}

		filter["_id"] = oid
	}

	if opts.Hostname != "" {
		filter["hostname"] = opts.Hostname
	}

	if opts.UserName != "" {
		filter["username"] = opts.UserName
	}

	var data *feedsv2.Feed
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&data)
	}
	if err := mg.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, db.ErrNotFound
		}

		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.findOne(%s)", db.Query(filter)))
	}

	return data, nil
}

// GetFeedsStreamList returns a list of all active streams by feed query.
func GetFeedsStreamList(ctx context.Context, mg *db.MongoDB, optionsList ...func(*ListOpts)) ([]*feedsv2.Feed, error) {
	filter := bson.M{}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	if opts.Lang != "" {
		filter["localization.lang"] = opts.Lang
	}

	if opts.Country != "" {
		filter["localization.country"] = opts.Country
	}

	if opts.StreamType > 0 {
		filter["stream.streamtype"] = opts.StreamType
	}

	if opts.StreamStatus > 0 {
		filter["stream.streamstatus"] = opts.StreamStatus
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))
	findOptions.SetSort(bson.M{
		opts.SortKey: opts.SortOrder,
	})

	data := make([]*feedsv2.Feed, 0)

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}
		defer c.Close(ctx)
		for c.Next(ctx) {
			var f feedsv2.Feed
			err := c.Decode(&f)
			if err != nil {
				return err
			}
			data = append(data, &f)
		}
		return nil
	}

	if err := mg.Execute(ctx, feedsCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.feeds.find()")
	}

	return data, nil
}

func List() {}

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
				return db.ErrExists
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

func Update() {}

// Delete deletes a feed.
func Delete(ctx context.Context, mg *db.MongoDB, feed *feedsv2.Feed) error {
	oid, err := primitive.ObjectIDFromHex(feed.Id)
	if err != nil {
		return db.ErrInvalidID
	}

	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		c, err := collection.DeleteOne(ctx, filter)
		if c.DeletedCount == 0 {
			return db.ErrNotFound
		}
		return err
	}

	if err := mg.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return db.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.feeds.delete(%v)", db.Query(filter)))
	}

	return nil
}
