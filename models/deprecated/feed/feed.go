package feed

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const feedsCollection = "feeds"

var (
	// ErrNotFound abstracts the  not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrDuplicate occurs when a key alerady exists
	ErrDuplicate = errors.New("Entry already exists.")
)

func EnsureIndex(ctx context.Context, dbConn *db.MongoDB) error {
	index := []mongo.IndexModel{
		// {
		// 	Keys: bson.M{
		// 		"screen_name": 1,
		// 	},
		// 	Options: options.Index().SetUnique(true), // {Unique: true},
		// },
		// {
		// 	Keys: bson.M{
		// 		"twitter_id": 1,
		// 	},
		// 	Options: options.Index().SetUnique(true), // {Unique: true},
		// },
		{
			Keys: bsonx.Doc{
				{Key: "screen_name", Value: bsonx.String("text")},
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

type ListOpts struct {
	Limit      int
	Offset     int
	Q          string
	Deleted    bool
	Status     string
	StreamType string
	Lang       string
}

func Limit(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Limit = i
	}
}
func Offset(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Offset = i
	}
}
func Q(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Q = i
	}
}
func Deleted() func(*ListOpts) {
	return func(l *ListOpts) {
		l.Deleted = true
	}
}
func Status(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Status = s
	}
}
func StreamType(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.StreamType = s
	}
}
func Lang(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Lang = strings.ToUpper(s)
	}
}
func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.Deleted = false
	l.Status = ""
	l.Lang = "EL"
	return l
}
func NewListOpts() []func(*ListOpts) {
	return make([]func(*ListOpts), 0)
}

// List retrieves a list of existing feeds from the database.
func List(ctx context.Context, dbConn *db.MongoDB, optionsList ...func(*ListOpts)) (*FeedsList, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.List")
	defer span.End()

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}
	filter := bson.M{"deleted": false}
	if opts.Deleted {
		filter["deleted"] = true
	}
	if opts.Status != "" {
		filter["status"] = opts.Status
	}
	if opts.StreamType != "" {
		filter["stream_type"] = opts.StreamType
	}
	if opts.Lang != "" {
		filter["lang"] = opts.Lang
	}
	if opts.Q != "" {
		filter["$text"] = bson.M{"$search": opts.Q}
	}

	d := new(FeedsList)

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))
	findOptions.SetSort(bson.D{{"_id", -1}})

	pagination, err := GetPagination(ctx, dbConn, filter, opts.Limit, feedsCollection)
	if err != nil {
		return nil, err
	}
	d.Pagination = pagination

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, filter, findOptions) //.Decode(p) //(nil).Limit(opts.Limit).Skip(opts.Offset * opts.Limit).All(&p)
		if err != nil {
			return err
		}
		defer c.Close(ctx)
		for c.Next(ctx) {
			var f Feed
			err := c.Decode(&f)
			if err != nil {
				return err
			}
			d.Data = append(d.Data, &f)
		}
		return nil
	}

	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.feeds.find()")
	}

	return d, nil
}

// Get gets the specified feed from the database.
func Get(ctx context.Context, dbConn *db.MongoDB, id string, optionsList ...func(*db.ListOpts)) (*Feed, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.Get")
	defer span.End()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	filter := bson.M{"_id": oid, "deleted": false}

	opts := db.DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}
	if opts.Deleted {
		filter["deleted"] = true
	}

	var p *Feed
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&p)
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.find(%s)", db.Query(filter)))
	}

	return p, nil
}

// Create inserts a new feed into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, cf *Feed, now time.Time) (*Feed, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	cf.CreatedAt = now
	cf.UpdatedAt = now
	cf.ID = primitive.NewObjectID()

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &cf)
		return err
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.insert(%s)", db.Query(&cf)))
	}

	return cf, nil
}

// Update replaces a feed document in the database.
func Update(ctx context.Context, dbConn *db.MongoDB, id string, upd *UpdateFeed, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "model.feed.Update")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)
	if upd.Name != nil {
		fields["name"] = *upd.Name
	}
	if upd.ScreenName != nil {
		fields["screen_name"] = *upd.ScreenName
	}
	if upd.TwitterID != nil {
		fields["twitter_id"] = *upd.TwitterID
	}
	if upd.TwitterIDStr != nil {
		fields["twitter_id_str"] = *upd.TwitterIDStr
	}
	if upd.TwitterProfileImage != nil {
		fields["twitter_profile_image"] = *upd.TwitterProfileImage
	}
	if upd.Email != nil {
		fields["email"] = *upd.Email
	}
	if upd.BusinessType != nil {
		fields["business_type"] = *upd.BusinessType
	}
	if upd.Country != nil {
		fields["country"] = *upd.Country
	}
	if upd.Lang != nil {
		fields["lang"] = *upd.Lang
	}
	if upd.URL != nil {
		fields["url"] = *upd.URL
	}
	if upd.RSS != nil {
		fields["rss"] = *upd.RSS
	}
	if upd.StreamType != nil {
		fields["stream_type"] = *upd.StreamType
	}
	if upd.MetaClasses != nil {
		fields["meta_classes"] = *upd.MetaClasses
	}
	if upd.Status != nil {
		fields["status"] = *upd.Status
	}
	if upd.TestURL != nil {
		fields["testURL"] = *upd.TestURL
	}

	if upd.TestData != nil {
		fields["testData"] = *upd.TestData
	}

	if upd.ContentType != nil {
		fields["content_type"] = *upd.ContentType
	}
	if upd.Description != nil {
		fields["description"] = *upd.Description
	}
	if upd.Tier != nil {
		fields["tier"] = *upd.Tier
	}
	if upd.BusinessOwner != nil {
		fields["business_owner"] = *upd.BusinessOwner
	}
	if upd.Registered != nil {
		fields["registered"] = *upd.Registered
	}
	if upd.PoliticalStance != nil {
		fields["political_stance"] = *upd.PoliticalStance
	}
	if upd.PoliticalOrientation != nil {
		fields["political_orientation"] = *upd.PoliticalOrientation
	}
	if upd.Locality != nil {
		fields["locality"] = *upd.Locality
	}

	// If there's nothing to update we can quit early.
	if len(fields) == 0 {
		return nil
	}

	fields["updatedAt"] = now

	update := bson.M{"$set": fields}
	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}

		return errors.Wrap(err, fmt.Sprintf("db.feeds.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// Delete removes a feed from the database.
func Delete(ctx context.Context, dbConn *db.MongoDB, id string) error {
	ctx, span := trace.StartSpan(ctx, "model.feed.Delete")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.DeleteOne(ctx, filter)
		return err
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.feeds.remove(%v)", db.Query(filter)))
	}

	return nil
}

// ByScreenName gets the specified feed from the database.
func ByScreenName(ctx context.Context, dbConn *db.MongoDB, screenName string, optionsList ...func(*ListOpts)) (*Feed, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.ByScreenName")
	defer span.End()

	filter := bson.M{"screen_name": screenName, "deleted": false}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}
	if opts.Deleted {
		filter["deleted"] = true
	}

	var p *Feed
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&p)
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.BuScreenName(%s)", db.Query(filter)))
	}

	return p, nil
}

// ByID gets the specified feed from the database.
func ByID(ctx context.Context, dbConn *db.MongoDB, id string, optionsList ...func(*ListOpts)) (*Feed, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.ByID")
	defer span.End()

	filter := bson.M{"twitter_id_str": id, "deleted": false}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}
	if opts.Deleted {
		filter["deleted"] = true
	}

	var p *Feed
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&p)
	}
	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.feeds.ByID(%s)", db.Query(filter)))
	}

	return p, nil
}

// GetPagination paginate documents with query
func GetPagination(ctx context.Context, dbConn *db.MongoDB, filter bson.M, limit int, collection string) (*Pagination, error) {
	var res Pagination

	// count total documents with query
	f := func(collection *mongo.Collection) error {
		total, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		res.Total = total
		res.Pages = int64(math.Ceil(float64(total) / float64(limit)))

		return nil
	}

	if err := dbConn.Execute(ctx, collection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.%s.Count(%v)", collection, filter))
	}

	return &res, nil
}

func GetTargets(ctx context.Context, dbConn *db.MongoDB, optionsList ...func(*ListOpts)) (*FeedsList, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.GetTargets")
	defer span.End()

	filter := bson.M{
		"deleted": false,
		"status":  "active",
	}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	if opts.StreamType != "" {
		filter["stream_type"] = opts.StreamType
	}
	if opts.Lang != "" {
		filter["lang"] = opts.Lang
	}

	d := new(FeedsList)

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))
	findOptions.SetSort(bson.D{{"_id", -1}})

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}
		defer c.Close(ctx)
		for c.Next(ctx) {
			var f Feed
			err := c.Decode(&f)
			if err != nil {
				return err
			}
			d.Data = append(d.Data, &f)
		}
		return nil
	}

	if err := dbConn.Execute(ctx, feedsCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.feeds.find()")
	}

	return d, nil
}
