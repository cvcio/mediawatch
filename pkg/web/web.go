package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	validator "gopkg.in/go-playground/validator.v8"
)

// // A Handler is a type that handles an http request within our own little mini framework
// type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*chi.Mux
	// mw []Middleware
}

// type Middleware func(http.Handler) http.Handler

// New creates an App value that handle a set of routes for the application.
// You can provide any number of middleware and they'll be used to wrap every
// request handler.
func New(mw ...func(http.Handler) http.Handler) *App { //mw ...Middleware) *App {
	mux := chi.NewMux()
	mux.Use(mw...)

	return &App{
		Mux: mux,
	}
}

func NewHealthCheck(addr string, healthCheckfunc func() http.HandlerFunc, mw ...func(http.Handler) http.Handler) *http.Server {
	mux := chi.NewMux()
	mux.Use(mw...)
	app := &App{
		Mux: mux,
	}
	app.HealthCheck(healthCheckfunc)
	// create the http.Server
	srv := http.Server{
		Addr:           addr,
		Handler:        app,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &srv
}

func (a *App) HealthCheck(healthCheckfunc func() http.HandlerFunc) {
	a.Handle("GET", "/healthcheck", healthCheckfunc())
}

// func (a *App) Walk() error {
// 	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
// 		log.Printf("%s %s\n", method, route)
// 		return nil
// 	}
// 	if err := chi.Walk(a, walkFunc); err != nil {
// 		return err
// 	}
// 	return nil
// }

// Handle is our mechanism from mounting Handlers for a given HTTP verb and path
// pair, this make for really easy, convenient routing.
// func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {
func (a *App) Handle(verb, pattern string, handler http.Handler, mw ...func(http.Handler) http.Handler) {
	a.With(mw...).Method(verb, pattern, handler)
}

// validate provides a validator for checking models.
var validate = validator.New(&validator.Config{
	TagName:      "validate",
	FieldNameTag: "json",
})

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

// InvalidError is a custom error type for invalid fields.
type InvalidError []Invalid

// Error implements the error interface for InvalidError.
func (err InvalidError) Error() string {
	var str string
	for _, v := range err {
		str = fmt.Sprintf("%s,{%s:%s}", str, v.Fld, v.Err)
	}
	return str
}

// Unmarshal decodes the input to the struct type and checks the
// fields to verify the value is in a proper state.
func Unmarshal(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	var inv InvalidError
	if fve := validate.Struct(v); fve != nil {
		for _, fe := range fve.(validator.ValidationErrors) {
			inv = append(inv, Invalid{Fld: fe.Field, Err: fe.Tag})
		}
		return inv
	}

	return nil
}
