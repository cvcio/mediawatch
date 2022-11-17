package handlers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/customer"

	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/models/subscription"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.opencensus.io/trace"
)

// Subscription is the handler struct for org related enbpoints
type Subscription struct {
	DB        *db.MongoDB
	log       *logrus.Logger
	StripeKey string
}

// Get ...
func (u *Subscription) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Subscription.Get")
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

func (u *Subscription) CreateCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Subscription.Create")
		defer span.End()

		_, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		// user, err := account.Get(ctx, u.DB, claims.Subject)
		// if err == account.ErrNotFound {
		// 	render.Render(w, r, web.ErrNotFound)
		// 	return
		// }

		// if err != nil {
		// 	u.log.Debug(err)
		// 	render.Render(w, r, web.ErrInternalError)
		// 	return
		// }

		var customerParams stripe.CustomerParams
		if err := web.Unmarshal(r.Body, &customerParams); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		stripe.Key = u.StripeKey
		u.log.Infoln(&customerParams)
		// create a customer
		c, err := customer.New(&customerParams)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.JSON(w, r, c)
	}
}

func (u *Subscription) CreateSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Subscription.Create")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		user, err := account.Get(ctx, u.DB, claims.Subject)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		var customerParams stripe.CustomerParams
		if err := web.Unmarshal(r.Body, &customerParams); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		stripe.Key = "sk_test_mgwiMIv5fEgmym173jEpy4GL00bu6r7GN7"
		u.log.Infoln(&customerParams)
		// create a customer
		c, err := customer.New(&customerParams)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		// create subscription
		ns := subscription.Subscription{
			UserID:         user.ID.Hex(),
			OrganizationID: user.Organization,
		}

		subscription.Create(ctx, u.DB, &ns)
		// create stripe customer
		// c, err := customer.New(&customerParams)
		// if err != nil {
		// 	u.log.Debug(err)
		// 	render.Render(w, r, web.ErrInternalError)
		// 	return
		// }
		// create billing portal session
		billingPortalParams := &stripe.BillingPortalSessionParams{
			Customer:  stripe.String(c.ID),
			ReturnURL: stripe.String("http://localhost:8080/account"),
		}
		s, err := session.New(billingPortalParams)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}
		// if success create organization from customer
		// else delete customer, subscription

		// redirect to billing portal

		// if paid and has org create organization
		// else delete customer, subscription
		render.JSON(w, r, &subscription.SubsciptionSessionResponse{
			Status: 302,
			URL:    s.URL,
		})
		// http.Redirect(w, r, s.URL, 302)
	}
}

// Delete an existing organization from an id.
func (u *Subscription) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Subscription.Delete")
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

func (u *Subscription) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Subscription.Update")
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

		render.Render(w, r, web.Created)
	}
}
