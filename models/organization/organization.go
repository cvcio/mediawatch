package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/pkg/db"
)

const orgCollection = "organizations"

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
	index := []mongo.IndexModel{
		{
			Keys: bson.M{
				"screenName": 1,
			},
			Options: options.Index().SetUnique(true), // {Unique: true},
		},
		{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true), // {Unique: true},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		return errors.Wrap(err, "db.organizations.ensureIndex()")
	}
	return nil
}

// List retrieves a list of existing organization from the database.
func List(ctx context.Context, dbConn *db.MongoDB, optionsList ...func(*db.ListOpts)) ([]*Organization, error) {
	ctx, span := trace.StartSpan(ctx, "models.organization.List")
	defer span.End()

	opts := db.DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	u := []*Organization{}

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))

	filter := bson.M{"deleted": false}
	if opts.Deleted {
		filter["deleted"] = true
	}

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, filter, findOptions)
		if err != nil {
			return err
		}

		defer c.Close(ctx)
		for c.Next(ctx) {
			var a Organization
			err := c.Decode(&a)
			if err != nil {
				return err
			}
			u = append(u, &a)

		}
		return nil
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.organization.find()")
	}

	return u, nil
}

// Get gets the specified user from the database.
func Get(ctx context.Context, dbConn *db.MongoDB, id string) (*Organization, error) {
	ctx, span := trace.StartSpan(ctx, "models.organization.Get")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.M{"_id": oid, "deleted": false}

	var u *Organization
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // Find(q).One(&u)
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.organization.find(%s)", db.Query(filter)))
	}

	return u, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, nu *Organization, now time.Time) (*Organization, error) {
	ctx, span := trace.StartSpan(ctx, "models.organization.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	u := Organization{
		ID:         primitive.NewObjectID(),
		Name:       nu.Name,
		Industry:   nu.Industry,
		Type:       nu.Type,
		Size:       nu.Size,
		Email:      nu.Email,
		Phone:      nu.Phone,
		ScreenName: nu.ScreenName,
		URL:        nu.URL,
		Country:    nu.Country,
		City:       nu.City,
		Language:   nu.Language,
		Timezone:   nu.Timezone,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &u)
		return err
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.organization.insert(%s)", db.Query(&u)))
	}
	return &u, nil
}

// Update replaces a organization document in the database.
func Update(ctx context.Context, dbConn *db.MongoDB, id string, upd *UpdOrg, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "models.organization.Update")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)

	if upd.Name != nil {
		fields["name"] = *upd.Name
	}
	if upd.Type != nil {
		fields["type"] = *upd.Type
	}
	if upd.Email != nil {
		fields["email"] = *upd.Email
	}
	if upd.ScreenName != nil {
		fields["screenName"] = *upd.ScreenName
	}
	if upd.URL != nil {
		fields["url"] = upd.URL
	}
	if upd.Country != nil {
		fields["country"] = upd.Country
	}
	if upd.City != nil {
		fields["city"] = upd.City
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
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.organization.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// InsertMember adds a member
func InsertMember(ctx context.Context, dbConn *db.MongoDB, id string, member *Member) error {
	ctx, span := trace.StartSpan(ctx, "models.organization.InsertMember")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)
	fields["updatedAt"] = time.Now().Truncate(time.Millisecond)

	update := bson.M{"$set": fields, "$push": bson.M{"members": member}}
	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.organization.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// UpdateMember updates a member
func UpdateMember(ctx context.Context, dbConn *db.MongoDB, id string, member *Member) error {
	ctx, span := trace.StartSpan(ctx, "models.organization.UpdateMember")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)
	fields["updatedAt"] = time.Now().Truncate(time.Millisecond)
	fields["members.$"] = member

	update := bson.M{"$set": fields}
	filter := bson.M{"_id": oid, "members.id": member.ID}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.organization.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// RemoveMember removes a member from an organization
func RemoveMember(ctx context.Context, dbConn *db.MongoDB, id string, memberId string) error {
	ctx, span := trace.StartSpan(ctx, "models.organization.RemoveMember")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)
	fields["updatedAt"] = time.Now().Truncate(time.Millisecond)

	update := bson.M{"$set": fields, "$pull": bson.M{"members": bson.M{"id": memberId}}}
	filter := bson.M{"_id": oid}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.organization.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, dbConn *db.MongoDB, id string) error {
	ctx, span := trace.StartSpan(ctx, "models.organization.Delete")
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
	if err := dbConn.Execute(ctx, orgCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.organization.remove(%s)", db.Query(filter)))
	}
	return nil
}
