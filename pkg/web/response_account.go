package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/deprecated/account"

	"github.com/go-chi/render"
)

// AccountResponse is the response payload for the Account data model.
// See NOTE above in AccountRequest as well.
//
// In the AccountResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type AccountResponse struct {
	*account.Account

	// Account *AccountPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	// Elapsed int64 `json:"elapsed"`
}

func NewAccountResponse(Account *account.Account) *AccountResponse {
	resp := &AccountResponse{Account: Account}

	return resp
}

func (rd *AccountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type AccountListResponse []*AccountResponse

func NewAccountListResponse(users []*account.Account) []render.Renderer {
	list := []render.Renderer{}
	for _, u := range users {
		list = append(list, NewAccountResponse(u))
	}
	return list
}

type TokenResponse struct {
	account.Token
	// Account *AccountPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	// Elapsed int64 `json:"elapsed"`
}

func NewTokenResponse(token account.Token) *TokenResponse {
	resp := &TokenResponse{token}

	return resp
}

func (rd *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type VerifyResponse struct {
	Redirect int `json:"redirect"` // http response status code
	account.Verify
}

func NewVerifyResponse(verify account.Verify, id string) *VerifyResponse {
	resp := &VerifyResponse{301, verify}

	return resp
}

func (rd *VerifyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}
