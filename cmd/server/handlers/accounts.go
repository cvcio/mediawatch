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

	"github.com/cvcio/mediawatch/models/deprecated/account"
	"github.com/cvcio/mediawatch/models/organization"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/twillio"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// AccountsHandler provides a handler for managing accounts.
type AccountsHandler struct {
	log             *zap.Logger
	mongoDB         *db.MongoDB
	tokenGenerator  account.TokenGenerator
	oauth2Providers map[string]*oauth2.Config
	mailService     *mailer.Mailer
	smsService      *twillio.Twillio
}

// NewAccountsHandler returns a new instance of AccountsHandler.
func NewAccountsHandler(log *zap.Logger, mongoDB *db.MongoDB,
	tokenGenerator account.TokenGenerator, oauth2Providers map[string]*oauth2.Config,
	mailService *mailer.Mailer, smsService *twillio.Twillio) (*AccountsHandler, error) {
	ctx, span := trace.StartSpan(context.Background(), "account.EnsureIndex")
	defer span.End()

	if err := account.EnsureIndex(ctx, mongoDB); err != nil {
		return nil, err
	}

	return &AccountsHandler{
		log:             log,
		mongoDB:         mongoDB,
		tokenGenerator:  tokenGenerator,
		oauth2Providers: oauth2Providers,
		mailService:     mailService,
		smsService:      smsService,
	}, nil
}

// List returns all the existing users.
func (u *AccountsHandler) List(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.List")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		u.log.Sugar().Infof("ok: %v", ok)
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	opts := db.NewListOpts()
	l, err := strconv.Atoi(c.Query("limit"))
	if err == nil {
		opts = append(opts, db.Limit(l))
	}

	o, err := strconv.Atoi(c.Query("offset"))
	if err == nil {
		opts = append(opts, db.Offset(o))
	}

	var data []*account.Account

	// if 'org' param exists, list users for org
	org := c.Query("org")
	if org != "" {
		// only if user is org member
		if claims.Organization == org {
			opts = append(opts, db.Org(org))
		}
		data, err = account.List(ctx, u.mongoDB, opts...)
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrInternalError)
			return
		}
		c.JSON(http.StatusOK, web.NewAccountListResponse(data))
		return
	}

	if !claims.HasRole(auth.RoleAdmin, auth.RolePowerUser) {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	data, err = account.List(ctx, u.mongoDB, opts...)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, data)
}

// Get returns an existing user from an id.
func (u *AccountsHandler) Get(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Get")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	data, err := account.Get(ctx, u.mongoDB, id)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	c.JSON(http.StatusOK, web.NewAccountResponse(data))
}

// Delete delete an existing user from an id.
func (u *AccountsHandler) Delete(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Delete")
	defer span.End()

	// Get id
	id := c.Param("id")
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}
	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	acc, err := account.ByEmail(ctx, u.mongoDB, claims.User)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	// Delete id
	err = account.Delete(ctx, u.mongoDB, id)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	err = mailer.SendAccountDeletion(ctx, u.mailService, acc.Email, acc.FirstName)
	if err != nil {
		u.log.Sugar().Debug(err)
		return
	}

	c.JSON(http.StatusOK, web.Deleted)
}

// Create a new account.
func (u *AccountsHandler) Create(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Accounts.Create")
	defer span.End()

	var acc account.NewAccount
	if err := c.ShouldBindJSON(&acc); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	acc.Email = strings.ToLower(acc.Email)

	_, err := account.ByEmail(ctx, u.mongoDB, acc.Email)
	if err != account.ErrNotFound {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	created, err := account.Create(ctx, u.mongoDB, &acc, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.NewAccountResponse(created))
}

// Update an existing account.
func (u *AccountsHandler) Update(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Accounts.Update")
	defer span.End()

	id := c.Param("id")
	var updAccount account.UpdateAccount
	if err := c.ShouldBindJSON(&updAccount); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// If you are not an admin and looking to retrieve someone else then you are rejected.
	if !claims.HasRole(auth.RoleAdmin) && claims.Subject != id {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	if *updAccount.Nonce != "" {
		data, err := account.Get(ctx, u.mongoDB, id)
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrInternalError)
			return
		}

		if data.Nonce != *updAccount.Nonce {
			c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
			return
		}

		if err := mailer.SendPin(ctx, u.mailService, data.Email, "", data.Pin, data.ID.Hex()); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
			return
		}

		if err := u.smsService.SendPin(ctx, *updAccount.Mobile, data.Pin); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
			return
		}
	}

	err := account.Update(ctx, u.mongoDB, id, &updAccount, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.Updated)
}

// Login handles a request to authenticate an account. It expects a request using
// Basic Auth with a user's email and password. It responds with the user model.
func (u *AccountsHandler) Login(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Login")
	defer span.End()

	var login account.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	if login.Email == "" || login.Password == "" {
		c.JSON(http.StatusUnauthorized, web.ErrInvalidCredentials)
		return
	}

	login.Email = strings.ToLower(login.Email)
	acc, err := account.ByEmail(ctx, u.mongoDB, login.Email)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	if err := account.PasswordOK(ctx, acc, login.Password); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	// If user has 2FA enabled, update nonce and pin, and redirect to verification page
	if acc.FA2 {
		// create new pin and nonce
		pin := account.OTP(4)
		nonce := account.OTP(24)

		// update account
		if err := account.Update(ctx, u.mongoDB, acc.ID.Hex(), &account.UpdateAccount{Pin: &pin, Nonce: &nonce}, time.Now()); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
			return
		}

		// send pin with sms
		if err := u.smsService.SendPin(ctx, acc.Mobile, pin); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
			return
		}

		// return redirect to verify page
		c.JSON(http.StatusOK, web.NewVerifyResponse(account.Verify{ID: acc.ID, Nonce: nonce}, acc.ID.Hex()))
		return
	}

	if acc.Status == "suspended" {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	// Update && Return
	now := time.Now()
	if err := account.Update(ctx, u.mongoDB, acc.ID.Hex(), &account.UpdateAccount{LastLoginAt: &now}, time.Now()); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	type response struct {
		AccountsHandler *account.Account           `json:"account"`
		Organization    *organization.Organization `json:"organization"`
		Token           account.Token              `json:"token"`
	}

	res := response{
		AccountsHandler: acc,
	}

	if acc.Organization != "" {
		org, err := organization.Get(ctx, u.mongoDB, acc.Organization)
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrInternalError)
			return
		}

		res.Organization = org
	}

	tkn, err := account.Authenticate(ctx, u.tokenGenerator, time.Now(), acc)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	res.Token = tkn

	c.Header("Authorization", "Bearer "+tkn.AccessToken)
	c.JSON(http.StatusOK, res)
}

func (u *AccountsHandler) Verify(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Verify")
	defer span.End()

	// Get id
	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	var verify account.Verify
	if err := c.ShouldBindJSON(&verify); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	acc, err := account.Get(ctx, u.mongoDB, id)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	if acc.Nonce != verify.Nonce || acc.Pin != verify.Pin {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	v := "active"
	var updAccount account.UpdateAccount
	updAccount.Nonce = new(string)
	updAccount.Pin = new(string)
	updAccount.Status = &v
	err = account.Update(ctx, u.mongoDB, id, &updAccount, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	type response struct {
		AccountsHandler *account.Account           `json:"account"`
		Organization    *organization.Organization `json:"organization"`
		Token           account.Token              `json:"token"`
	}

	res := response{
		AccountsHandler: acc,
	}

	if acc.Organization != "" {
		org, err := organization.Get(ctx, u.mongoDB, acc.Organization)
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrInternalError)
			return
		}

		res.Organization = org
	}

	tkn, err := account.Authenticate(ctx, u.tokenGenerator, time.Now(), acc)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	res.Token = tkn

	// Write token to Authorization Header
	c.Header("Authorization", "Bearer "+tkn.AccessToken)

	c.JSON(http.StatusOK, res)
}

// Token handles a request to authenticate an account. It expects a request using
// Basic Auth with a user's email and password. It responds with a JWT.
func (u *AccountsHandler) Token(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Token")
	defer span.End()

	var token account.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	claims, err := u.tokenGenerator.ParseClaims(token.RefreshToken)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	acc, err := account.ByEmail(ctx, u.mongoDB, claims.User)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	type response struct {
		AccountsHandler *account.Account           `json:"account"`
		Organization    *organization.Organization `json:"organization"`
		Token           account.Token              `json:"token"`
	}

	res := response{
		AccountsHandler: acc,
	}

	if acc.Organization != "" {
		org, err := organization.Get(ctx, u.mongoDB, acc.Organization)
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusInternalServerError, web.ErrInternalError)
			return
		}

		res.Organization = org
	}

	tkn, err := account.Authenticate(ctx, u.tokenGenerator, time.Now(), acc)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	res.Token = tkn

	// Write token to Authorization Header
	c.Header("Authorization", "Bearer "+tkn.AccessToken)

	c.JSON(http.StatusOK, res)
}

// OAuthGoogle Handler
func (u *AccountsHandler) OAuthGoogle(c *gin.Context) {
	_, span := trace.StartSpan(c.Request.Context(), "handlers.Google.Auth")
	defer span.End()

	url := u.oauth2Providers["google"].AuthCodeURL(string(""))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (u *AccountsHandler) OAuthGoogleCB(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Google.AuthCallback")
	defer span.End()

	code := c.Query("code")

	accessToken, err := u.oauth2Providers["google"].Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, web.ErrInvalidCredentials)
		return
	}

	if !accessToken.Valid() {
		c.JSON(http.StatusUnauthorized, web.ErrInvalidToken)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
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
	if err := json.Unmarshal(contents, &user); err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	status := "active"
	var acc *account.Account
	acc, err = account.ByEmail(ctx, u.mongoDB, user.Email)
	if err != nil && err == account.ErrNotFound {
		nu := &account.NewAccount{}
		nu.Email = user.Email
		pass := RandomPassword(12)
		nu.Password = pass
		nu.PasswordConfirm = pass

		acc, err = account.Create(ctx, u.mongoDB, nu, time.Now())
		if err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
			return
		}

		uu := &account.UpdateAccount{}
		uu.FirstName = &user.GivenName
		uu.LastName = &user.FamilyName
		uu.Avatar = &user.Picture
		uu.Status = &status

		if err := account.Update(ctx, u.mongoDB, acc.ID.Hex(), uu, time.Now()); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
			return
		}
	} else {
		uu := &account.UpdateAccount{}
		uu.Avatar = &user.Picture
		uu.Status = &status

		if err := account.Update(ctx, u.mongoDB, acc.ID.Hex(), uu, time.Now()); err != nil {
			u.log.Sugar().Debug(err)
			c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
			return
		}
	}

	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	tkn, err := account.Authenticate(ctx, u.tokenGenerator, time.Now(), acc)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	successCallBackURL := os.Getenv("CLIENT_AUTH_CB_URL")
	if successCallBackURL != "" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(
			"%s?access_token=%s&refresh_token=%s",
			successCallBackURL,
			tkn.AccessToken,
			tkn.RefreshToken,
		))
		return
	}
	c.JSON(http.StatusOK, web.NewTokenResponse(tkn))
}

// Reset password request.
func (u *AccountsHandler) Reset(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Token")
	defer span.End()

	var login struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&login); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	if login.Email == "" {
		c.JSON(http.StatusUnauthorized, web.ErrInvalidCredentials)
		return
	}

	acc, err := account.ByEmail(ctx, u.mongoDB, login.Email)
	if err != nil {
		u.log.Sugar().Debug(err)
		return
	}

	pin := account.OTP(4)
	nonce := account.OTP(24)

	err = account.Update(ctx, u.mongoDB, acc.ID.Hex(), &account.UpdateAccount{Pin: &pin, Nonce: &nonce}, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	err = mailer.SendReset(ctx, u.mailService, acc.Email, acc.FirstName, pin, acc.ID.Hex())
	if err != nil {
		u.log.Sugar().Debug(err)
		return
	}

	type ret struct {
		ID    string `json:"id"`
		Nonce string `json:"nonce"`
	}
	c.JSON(http.StatusOK, &ret{ID: acc.ID.Hex(), Nonce: nonce})
}

// ResetVerify password request.
func (u *AccountsHandler) ResetVerify(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.AccountsHandler.Token")
	defer span.End()

	id := c.Param("id")
	if !db.Valid(id) {
		c.JSON(http.StatusBadRequest, web.ErrInvalidID)
		return
	}

	var verify account.Verify
	if err := c.ShouldBindJSON(&verify); err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	acc, err := account.Get(ctx, u.mongoDB, id)
	if err == account.ErrNotFound {
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	if acc.Nonce != verify.Nonce || acc.Pin != verify.Pin {
		c.JSON(http.StatusUnauthorized, web.ErrAuthenticationFailure)
		return
	}

	pass := RandomPassword(12)

	var updAccount account.UpdateAccount
	updAccount.Nonce = new(string)
	updAccount.Pin = new(string)
	updAccount.Password = &pass

	err = account.Update(ctx, u.mongoDB, id, &updAccount, time.Now())
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	err = mailer.SendNewPass(ctx, u.mailService, acc.Email, acc.FirstName, pass)
	if err != nil {
		u.log.Sugar().Debug(err)
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, web.Updated)
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
