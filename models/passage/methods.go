package passage

import (
	"context"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opencensus.io/trace"
)

const passagesCollection = "passages"

// Create inserts a new passage into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, passage *passagesv2.Passage, now time.Time) (*passagesv2.Passage, error) {

	ctx, span := trace.StartSpan(ctx, "model.passage.Create")
	defer span.End()

	filter := bson.M{"type": passage.Type, "text": passage.Text}

	var p *passagesv2.Passage

	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&p)
	}
	passage.Id = primitive.NewObjectID().Hex()

	f = func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &passage)
		return err
	}
	if err := dbConn.Execute(ctx, passagesCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.passages.insert(%s)", db.Query(&passage)))
	}

	return passage, nil
}
