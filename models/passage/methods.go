package passage

import (
	"context"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const passagesCollection = "passages"

// EnsureIndex in mongodb.
func EnsureIndex(ctx context.Context, dbConn *db.MongoDB) error {
	index := []mongo.IndexModel{
		{
			Keys: bson.M{
				"type": 1,
				"text": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bsonx.Doc{
				{Key: "language", Value: bsonx.String("text")},
			},
			Options: options.Index().SetDefaultLanguage("en").SetLanguageOverride("el"),
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, passagesCollection, f); err != nil {
		return errors.Wrap(err, "db.passages.ensureIndex()")
	}
	return nil
}

// Create inserts a new passage into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, passage *passagesv2.Passage, now time.Time) (*passagesv2.Passage, error) {

	f := func(collection *mongo.Collection) error {
		inserted, err := collection.InsertOne(ctx, &passage) // (&u)
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
			passage.Id = oid.Hex()
		} else {
			return db.ErrInvalid
		}

		return nil
	}

	if err := dbConn.Execute(ctx, passagesCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.passages.insert(%s)", passage.Text))
	}

	return passage, nil
}

// List returns a list of passages by passage query.
func List(ctx context.Context, mg *db.MongoDB, optionsList ...func(*ListOpts)) (*passagesv2.PassageList, error) {
	filter := bson.M{}

	opts := DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	if opts.Lang != "" {
		filter["language"] = opts.Lang
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))

	data := make([]*passagesv2.Passage, 0)
	pagination := &commonv2.Pagination{}

	p, err := db.GetPagination(ctx, mg, filter, opts.Limit, passagesCollection)
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
			var p passagesv2.Passage
			err := c.Decode(&p)
			if err != nil {
				return err
			}
			data = append(data, &p)
		}
		return nil
	}

	if err := mg.Execute(ctx, passagesCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.passages.find()")
	}

	return &passagesv2.PassageList{
		Data:       data,
		Pagination: pagination,
	}, nil
}
