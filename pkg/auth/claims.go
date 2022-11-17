package auth

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin        = "ADMIN"
	RolePowerUser    = "POWERUSER"
	RoleUser         = "USER"
	RoleOrgAdmin     = "ORGADMIN"
	RoleEditor       = "EDITOR"
	RoleInvestigator = "INVESTIGATOR"
	RoleAnnotator    = "ANNOTATOR"
	RoleReviewer     = "REVIEWER"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

// Claims represents the authorization claims transmitted via a JWT.
// TODO: add scope to claims
type Claims struct {
	Roles        []string `json:"roles"`
	User         string   `json:"user"`
	Organization string   `json:"organization"`
	Lang         string   `json:"lang"`
	SearchLimit  int      `json:"searchLimit"`
	jwt.StandardClaims
}

// NewClaims constructs a Claims value for the identified user. The Claims
// expire within a specified duration of the provided time. Additional fields
// of the Claims can be set after calling NewClaims is desired.
func NewClaims(subject string, user string, roles []string, organization string, lang string, searchLimit int, now time.Time, expires time.Duration) Claims {
	c := Claims{
		Roles:        roles,
		Organization: organization,
		User:         user,
		Lang:         lang,
		SearchLimit:  searchLimit,
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(expires).Unix(),
		},
	}

	return c
}

// Valid is called during the parsing of a token.
func (c Claims) Valid() error {
	for _, r := range c.Roles {
		switch r {
		case RoleAdmin, RoleUser, RolePowerUser, RoleOrgAdmin, RoleEditor, RoleInvestigator, RoleAnnotator, RoleReviewer: // Role is valid.
		default:
			return fmt.Errorf("invalid role %q", r)
		}
	}
	if err := c.StandardClaims.Valid(); err != nil {
		return errors.Wrap(err, "validating standard claims")
	}
	return nil
}

// HasRole returns true if the claims has at least one of the provided roles.
func (c Claims) HasRole(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}
