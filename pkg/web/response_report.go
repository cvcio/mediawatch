package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/report"
	"github.com/go-chi/render"
)

type ReportResponse struct {
	*report.Report
}

func NewReportResponse(org *report.Report) *ReportResponse {
	return &ReportResponse{Report: org}
}

func (rd *ReportResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type ReportListResponse []*ReportResponse

func NewReportListResponse(orgs []*report.Report) []render.Renderer {
	list := []render.Renderer{}
	for _, o := range orgs {
		list = append(list, NewReportResponse(o))
	}
	return list
}
