package auth

import (
	"crypto/rsa"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// KeyFunc is used to map a JWT key id (kid) to the corresponding public key.
// It is a requirement for creating an Authenticator.
//
// * Private keys should be rotated. During the transition period, tokens
// signed with the old and new keys can coexist by looking up the correct
// public key by key id (kid).
//
// * Key-id-to-public-key resolution is usually accomplished via a public JWKS
// endpoint. See https://auth0.com/docs/jwks for more details.
type KeyFunc func(keyID string) (*rsa.PublicKey, error)

// NewSingleKeyFunc is a simple implementation of KeyFunc that only ever
// supports one key. This is easy for development but in production should be
// replaced with a caching layer that calls a JWKS endpoint.
func NewSingleKeyFunc(id string, key *rsa.PublicKey) KeyFunc {
	return func(kid string) (*rsa.PublicKey, error) {
		if id != kid {
			return nil, fmt.Errorf("Unrecognized kid %q", kid)
		}
		return key, nil
	}
}

// Authenticator is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Authenticator interface {
	GenerateToken(Claims) (string, error)
	ParseClaims(string) (Claims, error)
}

type DefaultAuthenticator struct {
	privateKey *rsa.PrivateKey
	keyID      string
	algorithm  string
	kf         KeyFunc
	parser     *jwt.Parser
}
type JWTAuthenticator struct {
	privateKey *rsa.PrivateKey
	keyID      string
	algorithm  string
	kf         KeyFunc
	parser     *jwt.Parser
	// expose
	Issuer string
}

// NewAuthenticator creates an *Authenticator for use. It will error if:
// - The private key is nil.
// - The public key func is nil.
// - The key ID is blank.
// - The specified algorithm is unsupported.
func NewDefaultAuthenticator(key *rsa.PrivateKey, keyID, algorithm string, publicKeyFunc KeyFunc) (*DefaultAuthenticator, error) {
	if key == nil {
		return nil, errors.New("private key cannot be nil")
	}
	if publicKeyFunc == nil {
		return nil, errors.New("public key function cannot be nil")
	}
	if keyID == "" {
		return nil, errors.New("keyID cannot be blank")
	}
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	// Create the token parser to use. The algorithm used to sign the JWT must be
	// validated to avoid a critical vulnerability:
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := DefaultAuthenticator{
		privateKey: key,
		keyID:      keyID,
		algorithm:  algorithm,
		kf:         publicKeyFunc,
		parser:     &parser,
	}

	return &a, nil
}

// NewJWTAuthenticator creates an *Authenticator for use. It will error if:
// - The private key is nil.
// - The public key func is nil.
// - The key ID is blank.
// - The specified algorithm is unsupported.
func NewJWTAuthenticator(keyCertificate, keyID, algorithm, issuer string) (*JWTAuthenticator, error) {
	// read private certificate file
	keyContents, err := os.ReadFile(keyCertificate)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	// parset rsa with jwt
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}
	// check for keyID
	if keyID == "" {
		return nil, errors.New("keyID cannot be blank")
	}

	// get jwt algorithm
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	// make lookup function
	publicKeyFunc := NewSingleKeyFunc(keyID, key.Public().(*rsa.PublicKey))

	// Create the token parser to use. The algorithm used to sign the JWT must be
	// validated to avoid a critical vulnerability:
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	return &JWTAuthenticator{
		privateKey: key,
		keyID:      keyID,
		algorithm:  algorithm,
		kf:         publicKeyFunc,
		parser:     &parser,
		Issuer:     issuer,
	}, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a DefaultAuthenticator) GenerateToken(claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)

	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.keyID

	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}

// ParseClaims recreates the Claims that were used to generate a token. It
// verifies that the token was signed using our key.
func (a DefaultAuthenticator) ParseClaims(tknStr string) (Claims, error) {

	// f is a function that returns the public key for validating a token. We use
	// the parsed (but unverified) token to find the key id. That ID is passed to
	// our KeyFunc to find the public key to use for verification.
	f := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("Missing key id (kid) in token header")
		}
		kidStr, ok := kid.(string)
		if !ok {
			return nil, errors.New("Token key id (kid) must be string")
		}

		return a.kf(kidStr)
	}

	var claims Claims
	tkn, err := a.parser.ParseWithClaims(tknStr, &claims, f)
	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	if !tkn.Valid {
		return Claims{}, errors.New("Invalid token")
	}

	return claims, nil
}

type DebugAuthenticator struct{}

func (a DebugAuthenticator) GenerateToken(claims Claims) (string, error) {

	return claims.Roles[0], nil
}
func (a DebugAuthenticator) ParseClaims(tkn string) (Claims, error) {
	switch tkn {
	case "user":
		return Claims{
			Roles: []string{
				RoleUser,
			},
		}, nil
	case "admin":
		return Claims{
			Roles: []string{
				RoleAdmin,
			},
		}, nil
	case "orgadmin":
		return Claims{
			Roles: []string{
				RoleOrgAdmin,
			},
		}, nil
	case "poweruser":
		return Claims{
			Roles: []string{
				RolePowerUser,
			},
		}, nil
	default:
		return Claims{
			Roles: []string{
				RoleUser,
			},
		}, nil
	}
}
