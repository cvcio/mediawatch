package account

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"
)

const accountsCollection = "accounts"

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

	// ErrInvalid is a generic error
	ErrInvalid = errors.New("Invalid request")
)

// EnsureIndex fix the indexes in the account collections
func EnsureIndex(ctx context.Context, dbConn *db.MongoDB) error {
	index := []mongo.IndexModel{
		{
			Options: options.Index().SetUnique(true), // {Unique: true},
			Keys: bson.M{
				"email": 1,
			},
		},
		{
			Keys: bsonx.Doc{
				{Key: "firstName", Value: bsonx.String("text")},
				{Key: "lastName", Value: bsonx.String("text")},
				{Key: "screenName", Value: bsonx.String("text")},
				{Key: "email", Value: bsonx.String("text")},
			},
			Options: options.Index().SetDefaultLanguage("en").SetLanguageOverride("el"),
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	f := func(collection *mongo.Collection) error {
		_, err := collection.Indexes().CreateMany(ctx, index, opts) //EnsureIndex(index)
		return err
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		return errors.Wrap(err, "db.accounts.ensureIndex()")
	}
	return nil
}

// List retrieves a list of existing users from the database.
func List(ctx context.Context, dbConn *db.MongoDB, optionsList ...func(*db.ListOpts)) ([]*Account, error) {
	ctx, span := trace.StartSpan(ctx, "models.account.List")
	defer span.End()

	opts := db.DefaultOpts()
	for _, o := range optionsList {
		o(&opts)
	}

	u := []*Account{}

	var q bson.M
	q = bson.M{}
	q["deleted"] = opts.Deleted

	if opts.Org != "" {
		q["organization"] = opts.Org
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(opts.Limit))
	findOptions.SetSkip(int64(opts.Offset))
	findOptions.SetSort(bson.D{{"_id", -1}})

	f := func(collection *mongo.Collection) error {
		c, err := collection.Find(ctx, q, findOptions)
		if err != nil {
			return err
		}

		defer c.Close(ctx)
		for c.Next(ctx) {
			var a Account
			err := c.Decode(&a)
			if err != nil {
				return err
			}
			u = append(u, &a)
		}
		return nil
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		return nil, errors.Wrap(err, "db.account.find()")
	}

	return u, nil
}

// Get gets the specified user from the database.
func Get(ctx context.Context, dbConn *db.MongoDB, id string) (*Account, error) {
	ctx, span := trace.StartSpan(ctx, "models.account.Get")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	filter := bson.M{"_id": oid, "deleted": false}
	var u *Account
	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // Find(q).One(&u)
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.accounts.find(%s)", db.Query(filter)))
	}

	return u, nil
}

// Create inserts a new user into the database.
func Create(ctx context.Context, dbConn *db.MongoDB, nu *NewAccount, now time.Time) (*Account, error) {
	ctx, span := trace.StartSpan(ctx, "models.account.Create")
	defer span.End()

	// Mongo truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	pw, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := Account{
		ID:           primitive.NewObjectID(),
		Email:        nu.Email,
		PasswordHash: pw,
		Roles:        []string{auth.RoleUser},
		CreatedAt:    now,
		UpdatedAt:    now,
		Status:       "pending",
		AcceptTerms:  true,
		Pin:          OTP(4),
		Nonce:        OTP(24),
	}

	f := func(collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, &u) // (&u)

		// Normally we would return ErrDuplicateKey in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err != nil {
			fmt.Println(err)
			we, _ := err.(mongo.WriteException)
			if we.WriteErrors[0].Code == 11000 {
				return ErrInvalid
			}
		}

		return err
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.accounts.insert(%s)", db.Query(&u)))
	}

	return &u, nil
}

// Update replaces a user document in the database.
func Update(ctx context.Context, dbConn *db.MongoDB, id string, upd *UpdateAccount, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "models.account.Update")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	fields := make(bson.M)

	if upd.LastLoginAt != nil {
		fields["lastLoginAt"] = *upd.LastLoginAt
	}

	if upd.FirstName != nil {
		fields["firstName"] = *upd.FirstName
	}
	if upd.LastName != nil {
		fields["lastName"] = *upd.LastName
	}
	if upd.ScreenName != nil {
		fields["screenName"] = *upd.ScreenName
	}
	if upd.Country != nil {
		fields["country"] = *upd.Country
	}
	if upd.Language != nil {
		fields["language"] = *upd.Language
	}
	if upd.Mobile != nil {
		fields["mobile"] = *upd.Mobile
	}
	if upd.Industry != nil {
		fields["industry"] = *upd.Industry
	}
	if upd.Occupation != nil {
		fields["occupation"] = *upd.Occupation
	}

	if upd.Avatar != nil {
		fields["avatar"] = *upd.Avatar
	}
	if upd.FA2 != nil {
		fields["fa2"] = *upd.FA2
	}

	if upd.Email != nil {
		fields["email"] = *upd.Email
	}

	if upd.Organization != nil {
		fields["organization"] = upd.Organization
	}

	if upd.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*upd.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "generating password hash")
		}
		fields["password_hash"] = pw
	}

	if upd.Status != nil {
		fields["status"] = *upd.Status
	}
	if upd.Nonce != nil {
		fields["nonce"] = *upd.Nonce
	}
	if upd.Pin != nil {
		fields["pin"] = *upd.Pin
	}

	if upd.Roles != nil {
		fields["roles"] = *upd.Roles
	}

	// If there's nothing to update we can quit early.
	if len(fields) == 0 {
		return nil
	}

	fields["updatedAt"] = now

	update := bson.M{"$set": fields}
	filter := bson.M{"_id": oid, "deleted": false} // bson.ObjectIdHex(id)}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { //ErrNotFound {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.account.update(%s, %s)", db.Query(filter), db.Query(update)))
	}

	return nil
}

// Delete removes a user from the database.
func Delete(ctx context.Context, dbConn *db.MongoDB, id string) error {
	ctx, span := trace.StartSpan(ctx, "models.account.Delete")
	defer span.End()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	fields := make(bson.M)
	fields["deleted"] = true
	fields["updatedAt"] = time.Now()

	filter := bson.M{"_id": oid, "deleted": false}
	update := bson.M{"$set": fields}

	f := func(collection *mongo.Collection) error {
		_, err := collection.UpdateOne(ctx, filter, update)
		return err
	}
	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		if err == mongo.ErrNoDocuments { // .ErrNotFound {
			return ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.accounts.remove(%s)", db.Query(filter)))
	}

	return nil
}

// TokenGenerator is the behavior we need in our Authenticate to generate
// tokens for authenticated users.
type TokenGenerator interface {
	GenerateToken(auth.Claims) (string, error)
	ParseClaims(string) (auth.Claims, error)
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Token that can be used to authenticate in the future.
//
// The key, keyID, and alg are required for generating the token.
func Authenticate(ctx context.Context, tknGen TokenGenerator, now time.Time, u *Account) (Token, error) {
	_, span := trace.StartSpan(ctx, "models.account.Authenticate")
	defer span.End()

	// q := bson.M{"email": email}

	// var u *Account
	// f := func(collection *.Collection) error {
	// 	return collection.Find(q).One(&u)
	// }

	// if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {

	// 	// Normally we would return ErrNotFound in this scenario but we do not want
	// 	// to leak to an unauthenticated user which emails are in the system.
	// 	if err == .ErrNotFound {
	// 		return Token{}, ErrAuthenticationFailure
	// 	}
	// 	return Token{}, errors.Wrap(err, fmt.Sprintf("db.accounts.find(%s)", db.Query(q)))
	// }

	// // Compare the provided password with the saved hash. Use the bcrypt
	// // comparison function so it is cryptographically secure.
	// if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
	// 	return Token{}, ErrAuthenticationFailure
	// }

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	searchLimit := 7
	for _, role := range u.Roles {
		if role == "POWERUSER" {
			searchLimit = 7
		}
		if role == "ORGADMIN" {
			searchLimit = 30
		}
		if role == "ADMIN" {
			searchLimit = 30
		}
	}
	accessClaims := auth.NewClaims(u.ID.Hex(), u.Email, u.Roles, u.Organization, strings.ToUpper(u.Language), searchLimit, now, 30*time.Minute)
	tkn, err := tknGen.GenerateToken(accessClaims)
	if err != nil {
		return Token{}, errors.Wrap(err, "generating token")
	}

	refreshClaims := auth.NewClaims(u.ID.Hex(), u.Email, []string{}, "", strings.ToUpper(u.Language), searchLimit, now, 72*time.Hour)
	rfstkn, err := tknGen.GenerateToken(refreshClaims)
	if err != nil {
		return Token{}, errors.Wrap(err, "generating token")
	}

	return Token{AccessToken: tkn, RefreshToken: rfstkn}, nil
}

// ByEmail retrieves a user account by email
func ByEmail(ctx context.Context, dbConn *db.MongoDB, email string) (*Account, error) {
	ctx, span := trace.StartSpan(ctx, "models.account.ByEmail")
	defer span.End()

	filter := bson.M{"email": email, "deleted": false}

	var u *Account

	f := func(collection *mongo.Collection) error {
		return collection.FindOne(ctx, filter).Decode(&u) // (q).One(&u)
	}

	if err := dbConn.Execute(ctx, accountsCollection, f); err != nil {
		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.accounts.find(%s)", db.Query(filter)))
	}
	return u, nil
}

// PasswordOK compares the provided password with the saved hash. Use the bcrypt
// comparison function so it is cryptographically secure.
func PasswordOK(ctx context.Context, u *Account, password string) error {
	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return ErrAuthenticationFailure
	}
	return nil
}

// OTP One Time Password used for 2 factor authentication, generates a n
// digit password as a string (pin).
func OTP(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func ValidMailAddress(address string) (string, error) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}
