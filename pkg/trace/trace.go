package trace

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// TraceIDHeader is the header added to outgoing requests which adds the traceID to it.
const TraceIDHeader = "X-Trace-ID"

// Key represents the type of value for the context key
type ctxKey int

// KeyValues is how request values are stored/retrieved
const KeyValues ctxKey = 1

// Values represent state for each request
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// Trace is a middleware that add a tracing id to a request
func Trace() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(TraceIDHeader) != nil {
				next.ServeHTTP(w, r)
			}

			v := Values{
				TraceID:    uuid.New().String(),
				Now:        time.Now(),
				StatusCode: http.StatusOK,
			}

			r.Header.Set(TraceIDHeader, v.TraceID)

			ctx := context.WithValue(r.Context(), KeyValues, v)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
