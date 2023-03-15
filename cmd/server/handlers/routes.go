package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/mediawatch/pkg/config"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/twillio"
	"github.com/prometheus/client_golang/prometheus"

	scrape_pb "github.com/cvcio/mediawatch/internal/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/mid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/web"

	"github.com/pkg/errors"
)

// API returns a handler for a set of routes
func API(
	log *logrus.Logger,
	registry *prometheus.Registry,
	db *db.MongoDB,
	es *es.Elastic,
	neoClient *neo.Neo,
	mw []func(http.Handler) http.Handler,
	authenticator auth.Authenticator,
	externalAuth map[string]*oauth2.Config,
	mail *mailer.Mailer,
	twillio *twillio.Twillio,
	twtt *anaconda.TwitterApi,
	scrape scrape_pb.ScrapeServiceClient,
	cfg *config.Config) http.Handler {

	// Create the application
	app := web.New(mw...)
	// authmw is used for authentication/authorization middleware.
	authmw := mid.Auth{
		Authenticator: authenticator,
	}

	// Bind all the user handlers.
	u := Account{
		DB:             db,
		TokenGenerator: authenticator,
		log:            log,
		auth:           externalAuth,
		mail:           mail,
		twillio:        twillio,
	}

	err := account.EnsureIndex(context.Background(), u.DB)
	if err != nil {
		log.Fatal(err)
	}

	o := Organization{
		DB:   db,
		log:  log,
		mail: mail,
	}

	err = organization.EnsureIndex(context.Background(), o.DB)
	if err != nil {
		log.Fatal(err)
	}

	f := Feeds{
		DB:  db,
		log: log,
	}

	err = feed.EnsureIndex(context.Background(), f.DB)
	if err != nil {
		log.Fatal(err)
	}

	a := NewArticlesHandler(log, db, es, neoClient, scrape)
	c := NewCasesHandler(log, db, es, neoClient)

	r := Reports{
		DB:  db,
		log: log,
	}

	// err = report.EnsureIndex(context.Background(), f.DB)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	t := Twitter{
		log:  log,
		twtt: twtt,
	}

	s := Subscription{
		DB:        db,
		log:       log,
		StripeKey: cfg.Stripe.Key,
	}

	p := Passages{
		DB:     db,
		log:    log,
		scrape: scrape,
	}

	// Extrernal Auth Service API Endpoints
	if len(u.auth) > 0 {
		for k := range u.auth {
			if k == "google" {
				app.Handle("GET", fmt.Sprintf("/v2/auth/%s", k), u.OAuthGoogle())
				app.Handle("GET", fmt.Sprintf("/v2/auth/%s/callback", k), u.OAuthGoogleCB())
			}
		}
	}

	// Internal Auth Service API Endpoints
	// TODO: We need to define a new auth middleware logging information about login, register (...) attempts.
	app.Handle("POST", "/v2/auth/login", u.Login())
	app.Handle("POST", "/v2/auth/register", u.Create())
	app.Handle("PUT", "/v2/auth/register/{id}", u.Update())
	app.Handle("PUT", "/v2/auth/verify/{id}", u.Verify())
	app.Handle("POST", "/v2/auth/token", u.Token())
	app.Handle("POST", "/v2/auth/reset", u.Reset())
	app.Handle("POST", "/v2/auth/reset/verify/{id}", u.ResetVerify())

	app.Handle("GET", "/v2/accounts", u.List(), authmw.Authenticate())
	app.Handle("POST", "/v2/accounts", u.Create())
	app.Handle("GET", "/v2/accounts/{id}", u.Get(), authmw.Authenticate())
	app.Handle("PUT", "/v2/accounts/{id}", u.Update(), authmw.Authenticate())
	app.Handle("DELETE", "/v2/accounts/{id}", u.Delete(), authmw.Authenticate())

	app.Handle("GET", "/v2/orgs", o.List(), authmw.Authenticate())
	app.Handle("POST", "/v2/orgs", o.Create(), authmw.Authenticate())
	app.Handle("GET", "/v2/orgs/{id}", o.Get(), authmw.Authenticate())
	app.Handle("PUT", "/v2/orgs/{id}", o.Update(), authmw.Authenticate())
	app.Handle("DELETE", "/v2/orgs/{id}", o.Delete(), authmw.Authenticate())
	app.Handle("PUT", "/v2/orgs/{id}/members", o.UpsertMember(), authmw.Authenticate())
	app.Handle("PUT", "/v2/orgs/{id}/members/{memberId}", o.UpsertMember(), authmw.Authenticate())
	app.Handle("DELETE", "/v2/orgs/{id}/members/{memberId}", o.RemoveMember(), authmw.Authenticate())
	app.Handle("PUT", "/v2/orgs/{id}/invitations/{memberId}/{method}", o.ProccessInvitation())

	app.Handle("POST", "/v2/subscription/session", s.CreateSession(), authmw.Authenticate())
	app.Handle("POST", "/v2/subscription/customer", s.CreateCustomer(), authmw.Authenticate())

	app.Handle("GET", "/v2/feeds", f.List(), authmw.Authenticate())
	app.Handle("POST", "/v2/feeds", f.Create(), authmw.Authenticate())
	app.Handle("GET", "/v2/feeds/{id}", f.Get(), authmw.Authenticate())
	app.Handle("PUT", "/v2/feeds/{id}", f.Update(), authmw.Authenticate())
	app.Handle("DELETE", "/v2/feeds/{id}", f.Delete(), authmw.Authenticate())

	app.Handle("POST", "/v2/passages", p.Create(), authmw.Authenticate())

	app.Handle("GET", "/v2/articles", a.List(), authmw.Authenticate())
	app.Handle("GET", "/v2/articles/{id}", a.Get(), authmw.Authenticate())

	app.Handle("GET", "/v2/cases", c.List(), authmw.Authenticate())
	app.Handle("GET", "/v2/cases/{id}", c.Get(), authmw.Authenticate())
	app.Handle("GET", "/v2/cases/{id}/count", c.Count(), authmw.Authenticate())

	app.Handle("GET", "/v2/reports", r.List(), authmw.Authenticate())
	app.Handle("POST", "/v2/reports", r.Create(), authmw.Authenticate())
	app.Handle("GET", "/v2/reports/{id}", r.Get(), authmw.Authenticate())

	app.Handle("POST", "/v2/external/parse/twitter-profile", t.Profile(), authmw.Authenticate())
	app.Handle("POST", "/v2/external/parse/article", a.ParseArticle(), authmw.Authenticate())
	app.Handle("GET", "/v2/external/export/feeds", f.Export(), authmw.Authenticate())
	app.Handle("POST", "/v2/external/import/feeds", f.Import(), authmw.Authenticate())

	return app
}

func translate(err error) *web.ErrResponse {
	switch errors.Cause(err) {
	case account.ErrAuthenticationFailure:
		return web.ErrUnauthorized
	case account.ErrForbidden:
		return web.ErrForbidden
	}
	return &web.ErrResponse{Err: err, ErrorText: err.Error()}
}
