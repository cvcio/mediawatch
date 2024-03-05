package passage

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

const passagesCollection = "passages"

// Create inserts a new feed into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, cf *Passage, now time.Time) (*Passage, error) {
	ctx, span := trace.StartSpan(ctx, "model.feed.Create")
	defer span.End()

	filter := bson.M{"type": cf.Type, "text": cf.Text}

	var p *Passage
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&p)
	}
	if err := dbConn.Execute(ctx, passagesCollection, f); err == nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s db.passages.find(%s)", "Passage already exists", cf.Text))
	}

	cf.ID = primitive.NewObjectID()

	f = func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &cf)
		return err
	}
	if err := dbConn.Execute(ctx, passagesCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.passages.insert(%s)", db.Query(&cf)))
	}

	return cf, nil
}
