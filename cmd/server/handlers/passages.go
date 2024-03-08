package handlers

import (
	"net/http"
	"time"

	"github.com/cvcio/mediawatch/models/passage"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	scrape_pb "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Passages is the handler struct for org related enbpoints
type Passages struct {
	DB     *db.MongoDB
	log    *logrus.Logger
	scrape scrape_pb.ScrapeServiceClient
}

// Create new passage
func (u *Passages) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Passages.Create")
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

		var newPassage *passagesv2.Passage
		if err := web.Unmarshal(r.Body, &newPassage); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		created, err := passage.Create(ctx, u.DB, newPassage, time.Now()) // Remove the address-of operator (&) from newPassage
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.JSON(w, r, created)
	}
}

// Reload passages
func (u *Passages) Reload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Passages.Reload")
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

		var newPassage passage.Passage
		if err := web.Unmarshal(r.Body, &newPassage); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		var emptyReq scrape_pb.Empty
		if err := web.Unmarshal(r.Body, &emptyReq); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}
		reloadPassagesRes, err := u.scrape.ReloadPassages(ctx, &emptyReq)
		if err != nil {
			u.log.Println(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}
		render.JSON(w, r, reloadPassagesRes)

		render.JSON(w, r, nil)
	}
}
