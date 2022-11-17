package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/twillio"
	"github.com/cvcio/mediawatch/pkg/web"
)

// Account is the our handler struct for user specific requests
type Account struct {
	DB             *db.MongoDB
	TokenGenerator account.TokenGenerator
	// ADD other state like the logger and config here
	log     *logrus.Logger
	auth    map[string]*oauth2.Config
	mail    *mailer.Mailer
	twillio *twillio.Twillio
}

// List returns all the existing user in the system.
func (u *Account) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.List")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		opts := db.NewListOpts()

		l, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err == nil {
			opts = append(opts, db.Limit(l))
		}

		o, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err == nil {
			opts = append(opts, db.Offset(o))
		}

		dbConn := u.DB //.Copy()
		// defer dbConn.Close()

		var data []*account.Account

		// if 'org' param exist list users for org
		org := r.URL.Query().Get("org")
		if org != "" {
			// only if user is org member
			if claims.Organization == org {
				opts = append(opts, db.Org(org))

			}
			data, err = account.List(ctx, dbConn, opts...)
			if err == account.ErrNotFound {
				render.Render(w, r, web.ErrNotFound)
				return
			}
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}
			render.RenderList(w, r, web.NewAccountListResponse(data))
			return
		}

		if !claims.HasRole(auth.RoleAdmin, auth.RolePowerUser) {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		data, err = account.List(ctx, dbConn, opts...)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.RenderList(w, r, web.NewAccountListResponse(data))
	}
}

// Get returns an existing user from an id.
func (u *Account) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Get")
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
		if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		data, err := account.Get(ctx, u.DB, id)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.Render(w, r, web.NewAccountResponse(data))
	}
}

// Delete delete an existing user from an id.
func (u *Account) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Delete")
		defer span.End()

		// Get id
		id := chi.URLParam(r, "id")

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		// If you are not an admin and looking to retrieve someone else then you are rejected.
		if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		acc, err := account.ByEmail(ctx, u.DB, claims.User)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		// Delete id
		err = account.Delete(ctx, u.DB, id)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		err = mailer.SendAccountDeletion(ctx, u.mail, acc.Email, acc.FirstName)
		if err != nil {
			u.log.Debug(err)
			return
		}

		render.Render(w, r, web.Deleted)
	}
}

// Create a new account.
func (u *Account) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Accounts.Create")
		defer span.End()

		var newAccount account.NewAccount
		if err := web.Unmarshal(r.Body, &newAccount); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		newAccount.Email = strings.ToLower(newAccount.Email)

		_, err := account.ByEmail(ctx, u.DB, newAccount.Email)
		if err != account.ErrNotFound {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		created, err := account.Create(ctx, u.DB, &newAccount, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.NewAccountResponse(created))
	}
}

// Update an existing account.
func (u *Account) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Accounts.Update")
		defer span.End()
		id := chi.URLParam(r, "id")

		var updAccount account.UpdateAccount
		if err := web.Unmarshal(r.Body, &updAccount); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		// if the update proccess is not part of the registration proccess
		if *updAccount.Nonce == "" {
			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			// If you are not an admin and looking to retrieve someone else then you are rejected.
			if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}
		} else {
			data, err := account.Get(ctx, u.DB, id)
			if err == account.ErrNotFound {
				u.log.Debug(err)
				render.Render(w, r, web.ErrNotFound)
				return
			}
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}

			if data.Nonce != *updAccount.Nonce {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			if err := mailer.SendPin(ctx, u.mail, data.Email, "", data.Pin, data.ID.Hex()); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrAuthenticationFailure)
				return
			}

			// u.log.Debugf("PIN CODE %s || %s", *updAccount.Mobile, data.Pin)

			// err = twillio.SendPin(ctx, u.twillio, *updAccount.Mobile, data.Pin)
			// if err != nil {
			// 	u.log.Debug(err)
			// 	render.Render(w, r, web.ErrInvalidRequest(err))
			// 	return
			// }
		}

		err := account.Update(ctx, u.DB, id, &updAccount, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.Updated)
	}
}

// Login handles a request to authenticate a account. It expects a request using
// Basic Auth with a user's email and password. It responds with the user model.
func (u *Account) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Login")
		defer span.End()

		var login account.Login
		if err := web.Unmarshal(r.Body, &login); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		if login.Email == "" || login.Password == "" {
			render.Render(w, r, web.ErrInvalidCredentials)
			return
		}

		login.Email = strings.ToLower(login.Email)
		acc, err := account.ByEmail(ctx, u.DB, login.Email)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		if err := account.PasswordOK(ctx, acc, login.Password); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		// If user has 2FA enabled, update nonce and pin, and redirect to verification page
		if acc.FA2 {
			// create new pin and nonce
			pin := account.OTP(4)
			nonce := account.OTP(24)

			// update account
			if err := account.Update(ctx, u.DB, acc.ID.Hex(), &account.UpdateAccount{Pin: &pin, Nonce: &nonce}, time.Now()); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInvalidRequest(err))
				return
			}

			// send pin with twillio
			if err := twillio.SendPin(ctx, u.twillio, acc.Mobile, pin); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInvalidRequest(err))
				return
			}

			// return redirect to verify page
			render.Render(w, r, web.NewVerifyResponse(account.Verify{ID: acc.ID, Nonce: nonce}, acc.ID.Hex()))
			return
		}

		u.log.Info(acc.Status)

		if acc.Status == "suspended" {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		// if acc.Status == "pending" {
		// }

		// Update && Return
		now := time.Now()
		if err := account.Update(ctx, u.DB, acc.ID.Hex(), &account.UpdateAccount{LastLoginAt: &now}, time.Now()); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		type response struct {
			Account      *account.Account           `json:"account"`
			Organization *organization.Organization `json:"organization"`
			Token        account.Token              `json:"token"`
		}

		res := response{
			Account: acc,
		}

		if acc.Organization != "" {
			org, err := organization.Get(ctx, u.DB, acc.Organization)
			if err == organization.ErrNotFound {
				render.Render(w, r, web.ErrNotFound)
				return
			}
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}

			res.Organization = org
		}

		tkn, err := account.Authenticate(ctx, u.TokenGenerator, time.Now(), acc)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, translate(err))
			return
		}

		res.Token = tkn

		// If we are this far, the request is valid and the token successfully generated.

		// Write token to Authorization Header
		w.Header().Set("Authorization", "Bearer "+tkn.AccessToken)
		render.JSON(w, r, res)
	}
}

// Verify OTP
func (u *Account) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Token")
		defer span.End()

		// Get id
		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		var verify account.Verify
		if err := web.Unmarshal(r.Body, &verify); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		acc, err := account.Get(ctx, u.DB, id)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		if acc.Nonce != verify.Nonce || acc.Pin != verify.Pin {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		v := "active"
		var updAccount account.UpdateAccount
		updAccount.Nonce = new(string)
		updAccount.Pin = new(string)
		updAccount.Status = &v
		err = account.Update(ctx, u.DB, id, &updAccount, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		type response struct {
			Account      *account.Account           `json:"account"`
			Organization *organization.Organization `json:"organization"`
			Token        account.Token              `json:"token"`
		}

		res := response{
			Account: acc,
		}

		if acc.Organization != "" {
			org, err := organization.Get(ctx, u.DB, acc.Organization)
			if err == organization.ErrNotFound {
				render.Render(w, r, web.ErrNotFound)
				return
			}
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}

			res.Organization = org
		}

		tkn, err := account.Authenticate(ctx, u.TokenGenerator, time.Now(), acc)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, translate(err))
			return
		}

		res.Token = tkn

		// If we are this far, the request is valid and the token successfully generated.

		// Write token to Authorization Header
		w.Header().Set("Authorization", "Bearer "+tkn.AccessToken)

		render.JSON(w, r, res)
	}
}

// Token handles a request to authenticate a account. It expects a request using
// Basic Auth with a user's email and password. It responds with a JWT.
func (u *Account) Token() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Token")
		defer span.End()

		var token account.Token
		if err := web.Unmarshal(r.Body, &token); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		claims, err := u.TokenGenerator.ParseClaims(token.RefreshToken)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		acc, err := account.ByEmail(ctx, u.DB, claims.User)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		type response struct {
			Account      *account.Account           `json:"account"`
			Organization *organization.Organization `json:"organization"`
			Token        account.Token              `json:"token"`
		}

		res := response{
			Account: acc,
		}

		if acc.Organization != "" {
			org, err := organization.Get(ctx, u.DB, acc.Organization)
			if err == organization.ErrNotFound {
				render.Render(w, r, web.ErrNotFound)
				return
			}
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}

			res.Organization = org
		}

		tkn, err := account.Authenticate(ctx, u.TokenGenerator, time.Now(), acc)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, translate(err))
			return
		}

		res.Token = tkn

		// Write token to Authorization Header
		w.Header().Set("Authorization", "Bearer "+tkn.AccessToken)

		render.JSON(w, r, res)
	}
}

// OAuthGoogle Handler
func (u *Account) OAuthGoogle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := trace.StartSpan(r.Context(), "handlers.Google.Auth")
		defer span.End()

		url := u.auth["google"].AuthCodeURL(string(""))

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// OAuthGoogleCB Callback
func (u *Account) OAuthGoogleCB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Google.AuthCallback")
		defer span.End()

		// session := sessions.Default(c)
		// session.Clear()

		code := r.URL.Query().Get("code")

		accessToken, err := u.auth["google"].Exchange(context.Background(), code)
		if err != nil {
			render.Render(w, r, web.ErrInvalidCredentials)
			return
		}

		if !accessToken.Valid() {
			render.Render(w, r, web.ErrInvalidToken)
			return
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken.AccessToken)
		if err != nil {
			render.Render(w, r, web.ErrInternalError)
			return
		}
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		type googleUser struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			VerifiedEmail bool   `json:"verified_email"`
			Name          string `json:"name"`
			GivenName     string `json:"given_name"`
			FamilyName    string `json:"family_name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			Gender        string `json:"gender"`
			Locale        string `json:"locale"`
		}

		var user *googleUser
		err = json.Unmarshal(contents, &user)
		if err != nil {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		status := "active"
		var acc *account.Account
		acc, err = account.ByEmail(ctx, u.DB, user.Email)
		if err == account.ErrNotFound {
			nu := &account.NewAccount{}
			nu.Email = user.Email
			pass := RandomPassword(12)
			nu.Password = pass
			nu.PasswordConfirm = pass

			acc, err = account.Create(ctx, u.DB, nu, time.Now())
			if err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrAuthenticationFailure)
				return
			}

			// pin := account.OTP(4)
			// nonce := account.OTP(24)

			// enrich the newly added account with more fields
			// from Google Authentication proccess
			uu := &account.UpdateAccount{}
			uu.FirstName = &user.GivenName
			uu.LastName = &user.FamilyName
			uu.Avatar = &user.Picture
			// uu.Pin = &pin
			// uu.Nonce = &nonce
			uu.Status = &status

			if err := account.Update(ctx, u.DB, acc.ID.Hex(), uu, time.Now()); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrAuthenticationFailure)
				return
			}
		} else {
			uu := &account.UpdateAccount{}
			uu.Avatar = &user.Picture
			uu.Status = &status

			if err := account.Update(ctx, u.DB, acc.ID.Hex(), uu, time.Now()); err != nil {
				u.log.Debug(err)
				render.Render(w, r, web.ErrAuthenticationFailure)
				return
			}
		}

		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		tkn, err := account.Authenticate(ctx, u.TokenGenerator, time.Now(), acc)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		// We are good to redirect the user to the application
		// this is only temporary and we should find a better way to
		// get the enviroment variables (ex like cfg)

		successCallBackURL := os.Getenv("CLIENT_AUTH_CB_URL")
		if successCallBackURL != "" {
			http.Redirect(w, r, fmt.Sprintf(
				"%s?access_token=%s&refresh_token=%s",
				successCallBackURL,
				tkn.AccessToken,
				tkn.RefreshToken,
			), http.StatusTemporaryRedirect)
			return
		}
		render.Render(w, r, web.NewTokenResponse(tkn))
	}
}

// Reset passowrd request.
func (u *Account) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Token")
		defer span.End()

		var login struct {
			Email string `json:"email"`
		}
		if err := web.Unmarshal(r.Body, &login); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		if login.Email == "" {
			render.Render(w, r, web.ErrInvalidCredentials)
			return
		}

		acc, err := account.ByEmail(ctx, u.DB, login.Email)
		if err != nil {
			u.log.Debug(err)
			return
		}

		pin := account.OTP(4)
		nonce := account.OTP(24)

		err = account.Update(ctx, u.DB, acc.ID.Hex(), &account.UpdateAccount{Pin: &pin, Nonce: &nonce}, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		err = mailer.SendReset(ctx, u.mail, acc.Email, acc.FirstName, pin, acc.ID.Hex())
		if err != nil {
			u.log.Debug(err)
			return
		}

		type ret struct {
			ID    string `json:"id"`
			Nonce string `json:"nonce"`
		}
		render.JSON(w, r, &ret{ID: acc.ID.Hex(), Nonce: nonce})
	}
}

// ResetVerify password request.
func (u *Account) ResetVerify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Account.Token")
		defer span.End()

		id := chi.URLParam(r, "id")
		if !db.Valid(id) {
			render.Render(w, r, web.ErrInvalidID)
			return
		}

		var verify account.Verify
		if err := web.Unmarshal(r.Body, &verify); err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		acc, err := account.Get(ctx, u.DB, id)
		if err == account.ErrNotFound {
			render.Render(w, r, web.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, web.ErrInternalError)
			return
		}

		if acc.Nonce != verify.Nonce || acc.Pin != verify.Pin {
			render.Render(w, r, web.ErrAuthenticationFailure)
			return
		}

		pass := RandomPassword(12)

		var updAccount account.UpdateAccount
		updAccount.Nonce = new(string)
		updAccount.Pin = new(string)
		updAccount.Password = &pass

		err = account.Update(ctx, u.DB, id, &updAccount, time.Now())
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		err = mailer.SendNewPass(ctx, u.mail, acc.Email, acc.FirstName, pass)
		if err != nil {
			u.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}

		render.Render(w, r, web.Updated)
	}
}

// RandomPassword Generates a Random Password
func RandomPassword(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
