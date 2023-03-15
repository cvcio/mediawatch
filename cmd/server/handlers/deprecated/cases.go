package deprecated

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/models/deprecated/nodes"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
)

// Cases Struct
type Cases struct {
	ES  *es.ES
	Neo *neo.Neo
	log *logrus.Logger
}

// List Get Cases List from Neo4J
func (a *Cases) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Cases.List")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		opts := nodes.NewSearch(r.FormValue)
		opts.Langs = claims.Lang
		if opts.Langs == "" {
			opts.Langs = "EL"
		}

		// type listResponse struct {
		// 	Data       []*nodes.NodeArticle `json:"data"`
		// 	Pagination nodes.Pagination     `json:"pagination"`
		// }
		// var response listResponse

		maxfrom := time.Now().AddDate(0, 0, -claims.SearchLimit).UTC()
		to := time.Now().UTC()
		from, _ := time.Parse(time.RFC3339, opts.Range.From)

		if from.Sub(maxfrom).Hours()/24 < 0 {
			from = maxfrom
		}

		if opts.Q != "" || opts.Title != "" || opts.Body != "" {
			list := make([]string, 0)
			if opts.Q != "" {
				list = append(list, opts.Q)
			}
			if opts.Title != "" {
				list = append(list, "title:"+opts.Title)
			}
			if opts.Body != "" {
				list = append(list, "summary:"+opts.Body)
			}
			opts.Q = strings.Join(list, " ")
		}

		templateOptions := map[string]interface{}{
			"withFullText":        false,
			"withOptionalSimilar": false,
			"withSimilar":         false,
			"withNotSimilar":      false,
			"withFeeds":           false,
			"withURL":             false,
			"skip":                opts.Skip,
			"limit":               opts.Limit,
			"q":                   opts.Q,
			"dataType":            ``,
			"feeds":               strings.Join(opts.Feeds, " OR "),
			"url":                 opts.URL,
			"includeRels":         ``,
			"from":                from.Format(time.RFC3339),
			"to":                  to.Format(time.RFC3339),
		}

		// Build the query string
		if opts.Q != "" || opts.Title != "" || opts.Body != "" {
			list := make([]string, 0)
			if opts.Q != "" {
				list = append(list, opts.Q)
			}
			if opts.Title != "" {
				list = append(list, "title:'"+opts.Title+"'")
			}
			if opts.Body != "" {
				list = append(list, "body:'"+opts.Body+"'")
			}
			opts.Q = strings.Join(list, " ")
			templateOptions["withFullText"] = true
		}

		if opts.Range.To != "" {
			to, _ = time.Parse(time.RFC3339, opts.Range.To)
		}

		if len(opts.Feeds) > 0 {
			templateOptions["withFeeds"] = true
		}
		if opts.URL != "" {
			templateOptions["withURL"] = true
		}

		templateOptions["includeRels"] = `, relCount: count(r)`

		switch opts.DataType {
		case 1:
			templateOptions["withOptionalSimilar"] = true
		case 2:
			templateOptions["withSimilar"] = true
		case 3:
			templateOptions["withNotSimilar"] = true
			templateOptions["includeRels"] = ``
		}

		// Pagination Query
		totalstmpl, _ := template.New("totalsquery").Parse(nodes.ListCountNodesResultsTxFunc)
		var totalstpl bytes.Buffer
		if err := totalstmpl.Execute(&totalstpl, templateOptions); err != nil {
			a.log.Fatal(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		totals := totalstpl.String()
		// fmt.Printf("\nTOTALS\n%s\n", totals)

		totalssession := a.Neo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer totalssession.Close()

		totalsdata, err := totalssession.Run(totals, nil)
		if err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		if totalsdata.Err() != nil {
			a.log.Error(totalsdata.Err())
			render.Render(w, r, web.ErrInternalError)
			return
		}

		// Parse results
		var pagination *nodes.Pagination
		if totalsdata.Next() {
			record := totalsdata.Record()

			if value, ok := record.Get("total"); ok {
				pagination = &nodes.Pagination{
					Total: value.(int64),
				}

				// Break Request
				if value.(int64) <= 0 {
					render.JSON(w, r, map[string]interface{}{
						"data":       map[string]interface{}{},
						"pagination": pagination,
					})
					return
				}
			}
		}

		// Data Query
		resultstmpl, _ := template.New("resultsquery").Parse(nodes.ListNodesTxFunc)
		var resultstpl bytes.Buffer
		if err := resultstmpl.Execute(&resultstpl, templateOptions); err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		results := resultstpl.String()
		// fmt.Printf("\nRESULTS\n%s\n", results)

		resultssession := a.Neo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer resultssession.Close()

		resultsdata, err := resultssession.Run(results, nil)
		if err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		if resultsdata.Err() != nil {
			a.log.Error(resultsdata.Err())
			render.Render(w, r, web.ErrInternalError)
			return
		}

		// Parse results
		var articles []*nodes.NodeArticle
		for resultsdata.Next() {
			record := resultsdata.Record()
			if value, ok := record.Get("data"); ok {
				var article *nodes.NodeArticle

				bodyBytes, _ := json.Marshal(value)
				json.Unmarshal(bodyBytes, &article)

				articles = append(articles, article)
			}
		}

		render.JSON(w, r, map[string]interface{}{
			"data":       articles,
			"pagination": pagination,
		})
	}
}

// Get Single Case from Neo4J
func (a *Cases) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Cases.List")
		defer span.End()

		id := chi.URLParam(r, "id")

		_, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		ressession := a.Neo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer ressession.Close()

		resultsdata, err := ressession.Run(nodes.GetNodeTxFunc, map[string]interface{}{
			"docId": id,
		})
		if err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		var articles []*nodes.NodeArticle
		for resultsdata.Next() {
			record := resultsdata.Record()

			if articleData, nodeExists := record.Get("json"); nodeExists {
				var article *nodes.NodeArticle

				bodyBytes, _ := json.Marshal(articleData)
				json.Unmarshal(bodyBytes, &article)

				articles = append(articles, article)
			}
		}

		if err = resultsdata.Err(); err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"data":       articles,
			"pagination": nil,
		})
	}
}

// Count Relationships of a Single Case from Neo4J
func (a *Cases) Count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Cases.List")
		defer span.End()

		id := chi.URLParam(r, "id")

		_, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		ressession := a.Neo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer ressession.Close()

		resultsdata, err := ressession.Run(nodes.CountNodeTxFunc, map[string]interface{}{
			"docId": id,
		})
		if err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
			return
		}

		if resultsdata.Next() {
			record := resultsdata.Record()

			if value, ok := record.Get("count"); ok {
				// Break Request
				if value.(int64) > 0 {
					render.JSON(w, r, map[string]interface{}{
						"count": value.(int64),
					})
					return
				}
			}
		}

		render.JSON(w, r, map[string]interface{}{
			"count": 0,
		})
	}
}
