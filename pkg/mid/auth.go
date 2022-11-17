package mid

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/web"
)

// Auth is used to authenticate and authorize HTTP requests.
type Auth struct {
	Authenticator auth.Authenticator
}

// Authenticate validates a JWT from the `Authorization` header.
func (a *Auth) Authenticate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHdr := r.Header.Get("Authorization")
			if authHdr == "" {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			tknStr, err := parseAuthHeader(authHdr)
			if err != nil {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			claims, err := a.Authenticator.ParseClaims(tknStr)
			if err != nil {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			// Add claims to the context so they can be retrieved later.
			ctx := context.WithValue(r.Context(), auth.Key, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// parseAuthHeader parses an authorization header. Expected header is of
// the format `Bearer <token>`.
func parseAuthHeader(bearerStr string) (string, error) {
	split := strings.Split(bearerStr, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("Expected Authorization header format: Bearer <token>")
	}

	return split[1], nil
}

// HasRole validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func (a *Auth) HasRole(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.Key).(auth.Claims)
			if !ok {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			if !claims.HasRole(roles...) {
				render.Render(w, r, web.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
