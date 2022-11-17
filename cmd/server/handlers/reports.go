package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/models/report"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/web"
)

// Reports is the our handler struct for reports api
type Reports struct {
	DB  *db.MongoDB
	log *logrus.Logger
}

// List returns all the existing user in the system.
func (u *Reports) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Reports.List")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// if sufficient claims
		// If you are not an org or admin or not having a valid id, then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) || claims.Id != "" {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// parse options
		// execute qurery
		// return the results

		// opts := db.NewListOpts()

		// l, err := strconv.Atoi(r.URL.Query().Get("limit"))
		// if err == nil {
		// 	opts = append(opts, db.Limit(l))
		// }

		// o, err := strconv.Atoi(r.URL.Query().Get("offset"))
		// if err == nil {
		// 	opts = append(opts, db.Offset(o))
		// }

		// dbConn := u.DB //.Copy()

		// data := []*report.Report{}

		// // if 'org' param exist list users for org
		// org := r.URL.Query().Get("org")
		// if org != "" {
		// 	// only if user is org member
		// 	if claims.Organization == org {
		// 		opts = append(opts, db.Org(org))

		// 	}
		// 	data, err = report.List(ctx, dbConn, opts...)
		// 	if err == report.ErrNotFound {
		// 		render.Render(w, r, web.ErrNotFound)
		// 		return
		// 	}
		// 	if err != nil {
		// 		u.log.Debug(err)
		// 		render.Render(w, r, web.ErrInternalError)
		// 		return
		// 	}
		// 	render.RenderList(w, r, web.NewReportListResponse(data))
		// 	return
		// }

		// if !claims.HasRole(auth.RoleAdmin, auth.RolePowerUser) {
		// 	render.RenderList(w, r, web.NewReportListResponse(data))
		// 	return
		// }

		// data, err = report.List(ctx, dbConn, opts...)
		// if err == account.ErrNotFound {
		// 	render.Render(w, r, web.ErrNotFound)
		// 	return
		// }

		// if err != nil {
		// 	u.log.Debug(err)
		// 	render.Render(w, r, web.ErrInternalError)
		// 	return
		// }

		// render.JSON(w, r, &data)
	}
}

func (u *Reports) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Reports.Get")
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

		// If you are not an org or admin or not having a valid id, then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) || claims.Id != "" {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		data, err := report.Get(ctx, u.DB, id)
		if err == report.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		if !claims.HasRole(auth.RoleOrgAdmin) || (data.UserID != claims.Id) {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.Render(w, r, web.NewReportResponse(data))
	}
}

func (u *Reports) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Reports.Get")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an org or admin or not having a valid id, then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) || claims.Id != "" {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		var newReport report.Report
		if err := web.Unmarshal(r.Body, &newReport); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		created, err := report.Create(ctx, u.DB, &newReport, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.NewReportResponse(created))
	}
}
