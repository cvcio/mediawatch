package account

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Account represents someone with access to our system.
type Account struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
	LastLoginAt time.Time `json:"lastLoginAt" bson:"lastLoginAt"`
	Deleted     bool      `json:"-" bson:"deleted"`
	Nonce       string    `json:"nonce" bson:"nonce"`

	Status string `bson:"status" json:"status"`

	FirstName   string `bson:"firstName" json:"firstName"`
	LastName    string `bson:"lastName" json:"lastName"`
	ScreenName  string `bson:"screenName" json:"screenName"`
	Country     string `bson:"country" json:"country"`
	Language    string `bson:"language" json:"language"`
	Mobile      string `bson:"mobile" json:"mobile"`
	Industry    string `bson:"industry" json:"industry"`
	Occupation  string `bson:"occupation" json:"occupation"`
	Avatar      string `bson:"avatar" json:"avatar"`
	AcceptTerms bool   `bson:"acceptTerms" json:"acceptTerms"`

	Pin string `bson:"pin" json:"-"`
	FA2 bool   `bson:"fa2" json:"fa2"`

	Organization string `bson:"organization" json:"organization"`

	Roles        []string `bson:"roles" json:"roles"`
	Email        string   `bson:"email" json:"email"` // TODO(jlw) enforce uniqueness
	PasswordHash []byte   `bson:"password_hash" json:"-"`
}

// NewAccount contains information needed to create a new Account.
type NewAccount struct {
	Email           string `json:"email" validate:"required"` // TODO(jlw) enforce uniqueness.
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,eqfield=Password"`
}

// UpdateAccount defines what information may be provided to modify an existing
// Account. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateAccount struct {
	Status      *string    `json:"status"`
	Nonce       *string    `json:"nonce"`
	LastLoginAt *time.Time `json:"lastLoginAt" bson:"lastLoginAt"`

	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	ScreenName  *string `json:"screenName"`
	Country     *string `json:"country"`
	Language    *string `json:"language"`
	Mobile      *string `json:"mobile"`
	Industry    *string `json:"industry"`
	Occupation  *string `json:"occupation"`
	Avatar      *string `json:"avatar"`
	AcceptTerms *bool   `json:"acceptTerms"`

	Pin *string `json:"pin"`
	FA2 *bool   `json:"fa2"`

	Organization *string `json:"organization"`

	Roles           *[]string `json:"roles"`
	Email           *string   `json:"email"` // TODO(jlw) enforce uniqueness
	Password        *string   `json:"password"`
	PasswordConfirm *string   `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}

// Token is the payload we deliver to users when they authenticate.
// We return both, a short lived AccessToken for 15 minutes
// and a long lived RefreshToken for 72 hours
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login defines login request information
type Login struct {
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
}

// Verify defines verification request information
type Verify struct {
	ID    primitive.ObjectID `json:"id"`
	Pin   string             `json:"pin"`
	Nonce string             `json:"nonce"`
}
