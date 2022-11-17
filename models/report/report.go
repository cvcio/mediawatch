package report

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/pkg/db"
)

const reportCollection = "reports"

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

func EnsureIndex(ctx context.Context, dbConn *db.MongoDB) error {
	index := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateOne(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		return errors.Wrap(err, "db.reports.ensureIndex()")
	}
	return nil
}

type ListOpts struct {
	Limit   int
	Offset  int
	Deleted bool
	OrgID   primitive.ObjectID
	UserID  primitive.ObjectID
	Status  string
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

func OrgID(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.OrgID, _ = primitive.ObjectIDFromHex(s)
	}
}

func UserID(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.UserID, _ = primitive.ObjectIDFromHex(s)
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

// List retrieves a list of existing report from the database.
func List(ctx context.Context, dbConn *db.MongoDB, optionsList ...func(*db.ListOpts)) (*ReportsList, error) {
	ctx, span := trace.StartSpan(ctx, "models.report.List")
	defer span.End()

	opts := db.DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	d := new(ReportsList)

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))

	filter := bson.M{"deleted": false}
	if opts.Deleted {
		filter["deleted"] = true
	}

	pagination, err := GetPagination(ctx, dbConn, filter, opts.Limit, reportCollection)
	if err != nil {
		return nil, err
	}
	d.Pagination = pagination

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}
		defer c.Close(ctx)
		for c.Next(ctx) {
			var a Report
			err := c.Decode(&a)
			if err != nil {
				return err
			}
			d.Data = append(d.Data, &a)
		}
		return nil
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.report.find()")
	}
	return d, nil
}

// Get gets the specified user from the database.
func Get(ctx context.Context, dbConn *db.MongoDB, id string) (*Report, error) {
	ctx, span := trace.StartSpan(ctx, "models.report.Get")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.M{"_id": oid, "deleted": false}

	var r *Report
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&r) // Find(q).One(&u)
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.report.find(%s)", db.Query(filter)))
	}

	return r, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, nr *Report, now time.Time) (*Report, error) {
	ctx, span := trace.StartSpan(ctx, "models.report.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	nr.ID = primitive.NewObjectID()
	nr.CreatedAt = now

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &nr)
		return err
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.report.insert(%s)", db.Query(&nr)))
	}
	return nr, nil
}

// Update replaces a report document in the database.
func Update(ctx context.Context, dbConn *db.MongoDB, id string, upd *Report, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "models.report.Update")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	upd.UpdatedAt = now

	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, upd)
		return err
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.report.update(%s, %s)", db.Query(filter), db.Query(upd)))
	}

	return nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, dbConn *db.MongoDB, id string) error {
	ctx, span := trace.StartSpan(ctx, "models.report.Delete")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)
	fields["deleted"] = true
	fields["updatedAt"] = time.Now()

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": fields}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, reportCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.report.remove(%s)", db.Query(filter)))
	}
	return nil
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
