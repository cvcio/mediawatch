package mid

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cvcio/mediawatch/pkg/trace"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware is a http middleware for logging requests
func LoggerMiddleware(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				log = log.WithFields(logrus.Fields{"request_id": reqID}).Logger
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.WithTime(time.Now()).WithField("recover_info", rec).WithField("debug_stack", debug.Stack()).Error("error_request")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}

				// log end request
				log.WithTime(time.Now()).WithFields(logrus.Fields{
					"http_scheme":       scheme,
					"http_proto":        r.Proto,
					"remote_addr":       r.RemoteAddr,
					"host":              r.Host,
					"proto":             r.Proto,
					"http_method":       r.Method,
					"user_agent":        r.Header.Get("User-Agent"),
					"status":            ww.Status(),
					"took_ms":           float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
					"bytes_in":          r.Header.Get("Content-Length"),
					"bytes_out":         ww.BytesWritten(),
					"uri":               r.RequestURI,
					trace.TraceIDHeader: r.Header.Get(trace.TraceIDHeader),
				}).Info("handled_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
