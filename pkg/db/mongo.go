package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	codecs "github.com/amsokol/mongo-go-driver-protobuf"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opencensus.io/trace"
)

var (
	// ErrNotFound abstracts the  not found error.
	ErrNotFound = errors.New("document not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	// ErrAuthenticationFailure = errors.New("authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	// ErrForbidden = errors.New("attempted action is not allowed")

	// ErrInvalid is a generic error
	ErrInvalid = errors.New("invalid request")

	// ErrInvalidDBProvided wrong database
	// ErrInvalidDBProvided = errors.New("invalid DB provided")

	// ErrConnectionFailed connection error
	ErrConnectionFailed = errors.New("connection failed")

	// ErrExists document exists error
	ErrExists = errors.New("document exists")
)

// MongoDB client struct.
type MongoDB struct {
	Database string
	Client   *mongo.Client
	Context  context.Context
}

// NewMongoDB returns a new mongodb client.
func NewMongoDB(uri, database string, timeout time.Duration) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reg := codecs.Register(bson.NewRegistryBuilder()).Build()

	clientOpts := &options.ClientOptions{
		ConnectTimeout: &timeout,
		Registry:       reg,
	}

	client, err := mongo.Connect(ctx, clientOpts.ApplyURI(uri).SetRegistry(reg))
	if err != nil {
		return nil, errors.Wrap(ErrConnectionFailed, err.Error())
	}

	// Call Ping to verify that the deployment is up and the Client was configured successfully.
	// As mentioned in the Ping documentation, this reduces application resiliency as the server may be
	// temporarily unavailable when Ping is called.
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(ErrConnectionFailed, err.Error())
	}

	return &MongoDB{
		Client:   client,
		Database: database,
		Context:  ctx,
	}, nil
}

// Close closes the database connection.
func (db *MongoDB) Close() error {
	return db.Client.Disconnect(db.Context)
}

// Copy returns a copy of the database.
func (db *MongoDB) Copy() *mongo.Database {
	return db.Client.Database(db.Database)
}

// Execute executes database queries.
func (db *MongoDB) Execute(ctx context.Context, collName string, f func(*mongo.Collection) error) error {
	return f(db.Client.Database(db.Database).Collection(collName))
}

// Valid returns true if a given id is a valid mongo id
func (db *MongoDB) Valid(id string) bool {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return false
	}
	return true
}

// StatusCheck validates the DB status good.
func (db *MongoDB) StatusCheck(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "pkg.DB.StatusCheck")
	defer span.End()

	return nil
}

// IsDuplicateError validates whether the query error is a key duplicate error
func (db *MongoDB) IsDuplicateError(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}

// Query provides a string version of the value
func Query(value interface{}) string {
	j, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(j)
}

// Valid returns true if a given id is a valid mongo id
func Valid(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	return true
}

type ListOpts struct {
	Limit   int
	Offset  int
	Org     string
	Deleted bool
	Status  string
	// Q      ListFunc
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

func Org(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Org = i
	}
}

// func Deleted() func(*ListOpts) {
// 	return func(l *ListOpts) {
// 		l.Deleted = true
// 	}
// }

func Status(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Status = s
	}
}

func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.Deleted = false
	l.Status = ""
	return l
}

func NewListOpts() []func(*ListOpts) {
	return make([]func(*ListOpts), 0)
}

// GetPagination paginate documents with query
func GetPagination(ctx context.Context, mg *MongoDB, filter bson.M, limit int, collection string) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	res["total"] = int64(0)
	res["pages"] = int64(0)

	// count total documents with query
	f := func(collection *mongo.Collection) error {
		total, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		res["total"] = total
		if limit > 0 {
			res["pages"] = int64(math.Ceil(float64(total) / float64(limit)))
		}
		return nil
	}

	if err := mg.Execute(ctx, collection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.%s.Count(%v)", collection, filter))
	}

	return res, nil
}
