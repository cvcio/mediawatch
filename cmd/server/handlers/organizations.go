package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.opencensus.io/trace"
)

// Organization is the handler struct for org related enbpoints
type Organization struct {
	DB   *db.MongoDB
	log  *logrus.Logger
	mail *mailer.Mailer
}

// List returns all the existing user in the system.
func (u *Organization) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.List")
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
		// Copy db connection
		dbConn := u.DB //.Copy()
		// defer dbConn.Close()

		// Create Options for the list query
		opts := db.NewListOpts()
		l, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil {
			opts = append(opts, db.Limit(l))
		}

		o, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err == nil {
			opts = append(opts, db.Offset(o))
		}

		data, err := organization.List(ctx, dbConn, opts...)
		if err == organization.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.RenderList(w, r, web.NewOrganizationListResponse(data))
	}
}

func (u *Organization) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.Get")
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

		data, err := organization.Get(ctx, u.DB, id)
		if err == organization.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.Render(w, r, web.NewOrganizationResponse(data))
	}
}

func (u *Organization) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.organizations.Create")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		oldUser, err := account.Get(ctx, u.DB, claims.Subject)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}
		if oldUser.Organization != "" {
			u.log.Debug("Error creating org, user already have one")
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		var newOrganization organization.Organization
		if err := web.Unmarshal(r.Body, &newOrganization); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		created, err := organization.Create(ctx, u.DB, &newOrganization, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		orgID := created.ID.Hex()
		roles := append(oldUser.Roles, auth.RoleOrgAdmin)
		updUser := account.UpdateAccount{
			Roles:        &roles,
			Organization: &orgID,
		}

		err = account.Update(ctx, u.DB, oldUser.ID.Hex(), &updUser, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.NewOrganizationResponse(created))
	}
}

// Delete an existing organization from an id.
func (u *Organization) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.Delete")
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
		// If you are not an admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// Delete id

		err := organization.Delete(ctx, u.DB, id)
		if err == account.ErrNotFound {
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

func (u *Organization) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.Update")
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

		// If you are not an admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// if not org admin, reject
		if !claims.HasRole(auth.RoleOrgAdmin) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		var updOrg organization.UpdOrg
		if err := web.Unmarshal(r.Body, &updOrg); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		err := organization.Update(ctx, u.DB, id, &updOrg, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.Updated)
	}
}

func (u *Organization) UpsertMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.UpsertMember")
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

		// If you are not an admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// if not org admin, reject
		if !claims.HasRole(auth.RoleOrgAdmin) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// parse request body
		var member organization.Member
		if err := web.Unmarshal(r.Body, &member); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}
		member.Email = strings.ToLower(member.Email)

		// get the organization
		org, err := organization.Get(ctx, u.DB, id)
		if err != nil {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// get the inviter
		inviter, err := account.ByEmail(ctx, u.DB, claims.User)
		if err != nil {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		member.Nonce = account.OTP(4)
		member.Status = "pending"

		// check whether member exists
		exists, err := account.ByEmail(ctx, u.DB, member.Email)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrNotFound)
			return
		}

		// notify user and create accept request
		member.ID = exists.ID.Hex()
		err = mailer.SendInviteExistingUser(ctx, u.mail, member.Email, inviter.FirstName, inviter.LastName, inviter.Email, member.Nonce, org.Name, id, member.ID)
		if err != nil {
			u.log.Debug(err)
			return
		}

		if err := organization.InsertMember(ctx, u.DB, id, &member); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.JSON(w, r, member)
	}
}

func (u *Organization) RemoveMember() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.UpsertMember")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		memberId := chi.URLParam(r, "memberId")
		if !db.Valid(memberId) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// If you are not an admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && claims.Organization != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// if not org admin, reject
		if !claims.HasRole(auth.RoleOrgAdmin) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		if err := organization.RemoveMember(ctx, u.DB, id, memberId); err != nil {
			u.log.Error(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.Deleted)
	}
}

func (u *Organization) ProccessInvitation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Organization.ProccessInvitation")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		memberId := chi.URLParam(r, "memberId")
		if !db.Valid(memberId) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		method := chi.URLParam(r, "method")

		// parse request body
		var member organization.Member
		if err := web.Unmarshal(r.Body, &member); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		// get the organization
		org, err := organization.Get(ctx, u.DB, id)
		if err != nil {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		var found *organization.Member
		for _, v := range org.Members {
			if v.ID == memberId {
				found = v
			}
		}

		if found == nil {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if member.Nonce != found.Nonce {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		if found.Status == "accepted" {
			render.Render(w, r, web.Updated)
			return
		}

		// check whether member exists
		exists, err := account.Get(ctx, u.DB, memberId)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if method == "accept" {
			exists.Organization = id
			exists.Roles = append(exists.Roles, found.Role)
			if err := account.Update(ctx, u.DB, exists.ID.Hex(), &account.UpdateAccount{
				Organization: &id,
				Roles:        &exists.Roles,
			}, time.Now().Truncate(time.Millisecond)); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrRender(err))
				return
			}

			found.Nonce = ""
			found.Status = "accepted"
			if err := organization.UpdateMember(ctx, u.DB, id, found); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrRender(err))
				return
			}
		} else if method == "reject" {
			found.Status = "rejected"
			if err := organization.UpdateMember(ctx, u.DB, id, found); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrRender(err))
				return
			}
		}

		render.Render(w, r, web.Updated)
	}
}
