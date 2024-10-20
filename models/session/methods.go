package session

import (
	"context"
	"fmt"
	"time"

	sessionsv2 "github.com/cvcio/mediawatch/internal/mediawatch/sessions/v2"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var sessionsCollection = "sessions"

// EnsureIndex on mongodb
func EnsureIndex(ctx context.Context, mg *db.MongoDB) error {
	index := []mongo.IndexModel{
		{
			Keys: bson.M{
				"expires_at": 1,
			},
			Options: options.Index().SetExpireAfterSeconds(60),
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts)
		return err
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		return fmt.Errorf("db.sessions.EnsureIndex: %w", err)
	}

	return nil
}

// GetByAccountID returns a session object from sessions collection by key/value account_id if exists, otherwise an error
func GetByAccountID(ctx context.Context, mg *db.MongoDB, id string) (*sessionsv2.Session, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, err
	}

	filter := bson.M{"key": "account_id", "value": id}

	var res sessionsv2.Session
	f := func(accountsCollection *mongo.Collection) error {
		return accountsCollection.FindOne(ctx, filter).Decode(&res) // (q).One(&u)
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == mongo.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.sessions.GetByAccountID(%s)", filter))
	}

	return &res, nil
}

// GetByEmail returns a session object from sessions collection by key/value email if exists, otherwise an error
func GetByEmail(ctx context.Context, mg *db.MongoDB, email string) (*sessionsv2.Session, error) {
	filter := bson.M{"key": "email", "value": email}

	var res sessionsv2.Session
	f := func(accountsCollection *mongo.Collection) error {
		return accountsCollection.FindOne(ctx, filter).Decode(&res) // (q).One(&u)
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == mongo.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.sessions.GetByAccountID(%s)", filter))
	}

	return &res, nil
}

// GetByID returns a session object from sessions collection if exists, otherwise an error
func GetByID(ctx context.Context, mg *db.MongoDB, id string) (*sessionsv2.Session, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}

	var res sessionsv2.Session
	f := func(accountsCollection *mongo.Collection) error {
		return accountsCollection.FindOne(ctx, filter).Decode(&res) // (q).One(&u)
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == mongo.ErrNoDocuments {
			return nil, db.ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.sessions.GetById(%s)", filter))
	}

	return &res, nil
}

// List session with query
func List(ctx context.Context, mg *db.MongoDB, q *sessionsv2.Session) ([]*sessionsv2.Session, error) {
	opts, filter := Filter(ParseOptions(q)...)

	find := options.Find()
	find.SetLimit(opts.Limit)

	// TODO: set sort from query
	find.SetSort(bson.M{"created_at": -1})

	var res []*sessionsv2.Session
	// find documents with query
	f := func(collection *mongo.Collection) error {
		cursor, err := collection.Find(ctx, filter, find)
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var r sessionsv2.Session

			err := cursor.Decode(&r)
			if err != nil {
				return err
			}

			res = append(res, &r)
		}

		return nil
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.sessions.List(%v)", q))
	}

	return res, nil
}

// Create returns a new session object on successful insert
func Create(ctx context.Context, mg *db.MongoDB, res *sessionsv2.Session, expire bool) (*sessionsv2.Session, error) {
	now := time.Now()

	res.CreatedAt = timestamppb.New(now)

	if !expire {
		res.ExpiresAt = nil
	} else {
		res.ExpiresAt = timestamppb.New(now.Add(12 * time.Hour))
	}

	f := func(collection *mongo.Collection) error {
		inserted, err := collection.InsertOne(ctx, &res) // (&u)
		// Normally we would return ErrDuplicateKey in this scenario, but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err != nil {
			we, _ := err.(mongo.WriteException)
			if we.WriteErrors[0].Code == 11000 {
				return db.ErrInvalid
			}

			return err
		}

		if oid, ok := inserted.InsertedID.(primitive.ObjectID); ok {
			res.Id = oid.Hex()
		} else {
			return db.ErrInvalid
		}

		return nil
	}

	if err := mg.Execute(ctx, sessionsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.sessions.Create(%s)", res.Value))
	}

	return res, nil
}
