package handlers

import (
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// Twitter is the handler struct for org related enbpoints
type Twitter struct {
	log  *logrus.Logger
	twtt *anaconda.TwitterApi
}

// Profile retrieves a user profile from twitter
func (u *Twitter) Profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// setup tracing TODO
		ctx, span := trace.StartSpan(r.Context(), "handlers.Twitter.Profile")
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
		type Req struct {
			ScreenName string `json:"screen_name"`
		}
		var newReq Req
		if err := web.Unmarshal(r.Body, &newReq); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		user, err := u.twtt.GetUsersShow(newReq.ScreenName, url.Values{})
		if err != nil {
			render.Render(w, r, web.ErrNotFound)
		}

		render.JSON(w, r, &user)
	}
}
