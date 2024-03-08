package passage

import (
	"context"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/db"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const passagesCollection = "passages"

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
