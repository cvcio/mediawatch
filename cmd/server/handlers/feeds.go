package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Feeds is the handler struct for org related enbpoints
type Feeds struct {
	DB  *db.MongoDB
	log *logrus.Logger
}

// List returns all the existing user in the system.
func (u *Feeds) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.List")
		defer span.End()

		// Get claims from request
		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// Create Options for the list query
		opts := feed.NewListOpts()

		l, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil {
			opts = append(opts, feed.Limit(l))
		}

		o, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err == nil {
			opts = append(opts, feed.Offset(o))
		}

		s := r.URL.Query().Get("status")
		if s != "" {
			opts = append(opts, feed.Status(s))
		}

		q := r.URL.Query().Get("q")
		if q != "" {
			opts = append(opts, feed.Q(q))
		}

		data, err := feed.List(ctx, u.DB, opts...)
		if err == feed.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.JSON(w, r, &data)
	}
}

func (u *Feeds) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Get")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		data, err := feed.Get(ctx, u.DB, id)
		if err == feed.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrExists)
			return
		}

		render.JSON(w, r, web.NewFeedResponse(data))
	}
}

func (u *Feeds) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Create")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		var newFeed feed.Feed
		if err := web.Unmarshal(r.Body, &newFeed); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		created, err := feed.Create(ctx, u.DB, &newFeed, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.NewFeedResponse(created))
	}
}

// Delete an existing organization from an id.
func (u *Feeds) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Delete")
		defer span.End()

		// Get id
		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			err := render.Render(w, r, web.ErrInvalidID)
			if err != nil {
				u.log.Error(err)
			}
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// Delete id
		err := feed.Delete(ctx, u.DB, id)
		if err == feed.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.Render(w, r, web.Deleted)
	}
}

func (u *Feeds) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Update")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		var updFeed feed.UpdateFeed
		if err := web.Unmarshal(r.Body, &updFeed); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		err := feed.Update(ctx, u.DB, id, &updFeed, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.Updated)
	}
}

func (u *Feeds) Export() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Update")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
	}
}

func (u *Feeds) Import() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Feeds.Update")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an org or admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && !claims.HasRole(auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
	}
}
