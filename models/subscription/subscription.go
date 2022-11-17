package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

const subscriptionCollection = "subscriptions"

var (
	// ErrNotFound abstracts the  not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List retrieves a list of existing Subscription from the database.
func List(ctx context.Context, dbConn *db.MongoDB) ([]*Subscription, error) {
	ctx, span := trace.StartSpan(ctx, "models.Subscription.List")
	defer span.End()

	u := []*Subscription{}

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		defer c.Close(ctx)
		for c.Next(ctx) {
			var a Subscription
			err := c.Decode(&a)
			if err != nil {
				return err
			}
			u = append(u, &a)

		}
		return nil
	}
	if err := dbConn.Execute(ctx, subscriptionCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.Subscription.List()")
	}

	return u, nil
}

// Get gets the specified user from the database.
func Get(ctx context.Context, dbConn *db.MongoDB, id string) (*Subscription, error) {
	ctx, span := trace.StartSpan(ctx, "models.Subscription.Get")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.M{"_id": oid}

	var u *Subscription
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // Find(q).One(&u)
	}
	if err := dbConn.Execute(ctx, subscriptionCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.Subscription.find(%s)", db.Query(filter)))
	}

	return u, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, nu *Subscription) (*Subscription, error) {
	ctx, span := trace.StartSpan(ctx, "models.subscriptions.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now := time.Now()

	u := Subscription{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &u)
		return err
	}
	if err := dbConn.Execute(ctx, subscriptionCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.subscriptions.insert(%s)", db.Query(&u)))
	}
	return &u, nil
}
