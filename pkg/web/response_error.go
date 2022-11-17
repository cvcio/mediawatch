package web

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	Created = &CreatedResponse{HTTPStatusCode: http.StatusCreated, StatusText: "created"}
	Deleted = &DeletedReposnse{HTTPStatusCode: http.StatusGone, StatusText: "deleted"}
	Updated = &UpdatedResponse{HTTPStatusCode: http.StatusOK, StatusText: "updated"}

	ErrInvalidCredentials    = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: "invalid credentials"}
	ErrAuthenticationFailure = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: "Authentication failed"}

	ErrNotFound      = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "not found"}
	ErrInvalidID     = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: "bad Request"}
	ErrInternalError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: "Internal Server Error"}
	ErrUnauthorized  = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, StatusText: "Unauthorized"}
	ErrForbidden     = &ErrResponse{HTTPStatusCode: http.StatusForbidden, StatusText: "Forbidden"}
	ErrInvalidToken  = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, StatusText: "Invalid Token"}

	ErrNotImplemented = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "not implemented"}

	ErrExists = &ErrResponse{HTTPStatusCode: http.StatusForbidden, StatusText: "Entry already exists."}
)

type CreatedResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	StatusText string `json:"status"`
}

func (*CreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type DeletedReposnse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	StatusText string `json:"status"`
}

func (*DeletedReposnse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type UpdatedResponse struct {
	HTTPStatusCode int `json:"-"` // http response status code

	StatusText string `json:"status"`
}

func (*UpdatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render is for implementing the render.Renderer interface for ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// 	// ErrDBNotConfigured occurs when the DB is not initialized.
// 	ErrDBNotConfigured = errors.New("DB not initialized")
// )
