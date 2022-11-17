package web

import (
	"net/http"
)

type HealthResponse struct {
	Status int
}

func NewHealthResponse(status int) *HealthResponse {
	return &HealthResponse{status}
}

func (rd HealthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}
