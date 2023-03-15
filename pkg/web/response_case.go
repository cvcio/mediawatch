package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/nodes"
	"github.com/go-chi/render"
)

type CaseResponse struct {
	*nodes.NodeArticle
}

func NewCaseResponse(f *nodes.NodeArticle) *CaseResponse {
	return &CaseResponse{f}
}

func (rd *CaseResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type CaseListResponse []*CaseResponse

func NewCaseListResponse(cases []*nodes.NodeArticle) []render.Renderer {
	list := []render.Renderer{}
	for _, o := range cases {
		list = append(list, NewCaseResponse(o))
	}
	return list
}
