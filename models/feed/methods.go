package feed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		filter["user_name"] = opts.UserName
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
		filter["stream.stream_type"] = opts.StreamType
	}

	if opts.StreamStatus > 0 {
		filter["stream.stream_status"] = opts.StreamStatus
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

// GetByUserName returns a feed with a specific user_name.
func GetByUserName(ctx context.Context, mg *db.MongoDB, username string) (*feedsv2.Feed, error) {
	filter := bson.M{"user_name": username}

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

// GetByHostname returns a feed with a specific hostname.
func GetByHostname(ctx context.Context, mg *db.MongoDB, hostname string) (*feedsv2.Feed, error) {
	filter := bson.M{"hostname": hostname}

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

// List returns a list of feeds by feed query.
func List(ctx context.Context, mg *db.MongoDB, optionsList ...func(*ListOpts)) (*feedsv2.FeedList, error) {
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
		filter["stream.stream_type"] = opts.StreamType
	}

	if opts.StreamStatus > 0 {
		filter["stream.streams_tatus"] = opts.StreamStatus
	}

	if opts.Q != "" {
		filter["$text"] = bson.M{"$search": opts.Q}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))
	findOptions.SetSort(bson.M{
		opts.SortKey: opts.SortOrder,
	})

	data := make([]*feedsv2.Feed, 0)
	pagination := &commonv2.Pagination{}

	p, err := db.GetPagination(ctx, mg, filter, opts.Limit, feedsCollection)
	if err != nil {
		return nil, err
	}

	if _, ok := p["total"]; ok {
		pagination.Total = p["total"].(int64)
	}
	if _, ok := p["pages"]; ok {
		pagination.Pages = p["pages"].(int64)
	}

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

	return &feedsv2.FeedList{
		Data:       data,
		Pagination: pagination,
	}, nil
}

// Create creates a new feed.
func Create(ctx context.Context, mg *db.MongoDB, feed *feedsv2.Feed) (*feedsv2.Feed, error) {
	now := time.Now().Truncate(time.Millisecond)
	feed.CreatedAt = timestamppb.New(now)

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

// Update updates a feed.
func Update(ctx context.Context, mg *db.MongoDB, feed *feedsv2.Feed) error {
	oid, err := primitive.ObjectIDFromHex(feed.Id)
	if err != nil {
		return db.ErrInvalidID
	}

	now := time.Now().Truncate(time.Millisecond)

	// create the fields to update
	fields := make(bson.M)
	fields["updated_at"] = timestamppb.New(now)

	j, err := json.Marshal(feed)
	if err != nil {
		return errors.New("unable to marshal proto")
	}

	if err := bson.UnmarshalExtJSON(j, true, fields); err != nil {
		return errors.New("unable to unmarshal json to bson")
	}

	delete(fields, "id")

	update := bson.M{"$set": fields}
	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}

	if err := mg.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return db.ErrNotFound
		}

		return errors.Wrap(err, fmt.Sprintf("db.feeds.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

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
