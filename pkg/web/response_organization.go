package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/organization"
	"github.com/go-chi/render"
)

type OrganizationResponse struct {
	*organization.Organization
}

func NewOrganizationResponse(org *organization.Organization) *OrganizationResponse {
	return &OrganizationResponse{Organization: org}
}

func (rd *OrganizationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type OrganizationListResponse []*OrganizationResponse

func NewOrganizationListResponse(orgs []*organization.Organization) []render.Renderer {
	list := []render.Renderer{}
	for _, o := range orgs {
		list = append(list, NewOrganizationResponse(o))
	}
	return list
}
