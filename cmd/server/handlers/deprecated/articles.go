package deprecated

import (
	"net/http"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/models/deprecated/nodes"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	scrape_pb "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

type Articles struct {
	ES     *es.ES
	DB     *db.MongoDB
	Neo    *neo.Neo
	log    *logrus.Logger
	scrape scrape_pb.ScrapeServiceClient
}

func (a *Articles) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Articles.List")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		opts := article.NewSearch(r.FormValue)
		opts.Langs = claims.Lang
		if opts.Langs == "" {
			opts.Langs = "EL"
		}

		maxfrom := time.Now().AddDate(0, 0, -claims.SearchLimit).UTC()
		maxto := time.Now().UTC()

		from, err := time.Parse(time.RFC3339, opts.Range.From)
		if err != nil {
			from = maxfrom
		}

		to, err := time.Parse(time.RFC3339, opts.Range.To)
		if err != nil {
			to = maxto
		}

		if maxfrom.Sub(from).Hours()/24 < 0 {
			from = maxfrom
		}

		opts.Range.From = from.Format(time.RFC3339)
		opts.Range.To = to.Format(time.RFC3339)

		docs, total, err := article.List(ctx, a.ES, opts)

		if err != nil {
			a.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
		}

		for _, v := range docs {
			v.Feed, _ = feed.ByScreenName(ctx, a.DB, v.ScreenName)

			ressession := a.Neo.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
			defer ressession.Close()

			resultsdata, err := ressession.Run(nodes.CountNodeTxFunc, map[string]interface{}{
				"docId": v.DocID,
			})
			if err != nil {
				a.log.Error(err)
				render.Render(w, r, web.ErrInternalError)
				return
			}

			if resultsdata.Next() {
				record := resultsdata.Record()
				v.RelCount = 0
				if value, ok := record.Get("count"); ok {
					if value.(int64) > 0 {
						v.RelCount = value.(int64)
					}
				}
			}
		}

		render.JSON(w, r, map[string]interface{}{
			"data": docs,
			"pagination": &nodes.Pagination{
				Total: total,
			},
		})
	}
}

func (a *Articles) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Articles.List")
		defer span.End()

		id := chi.URLParam(r, "id")

		_, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		data, err := article.Get(ctx, a.ES, id)
		if err != nil {
			a.log.Println(err)
			render.Render(w, r, web.ErrNotFound)
			return
		}
		render.Render(w, r, web.NewArticleResponse(data))
	}
}

// ParseArticle Uses Simple-Scrape gRPC endpoint to test the scraper
func (a *Articles) ParseArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Articles.List")
		defer span.End()

		_, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}
		var scrapeReq scrape_pb.SimpleScrapeRequest
		if err := web.Unmarshal(r.Body, &scrapeReq); err != nil {
			a.log.Debug(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}
		scrapeRes, err := a.scrape.SimpleScrape(ctx, &scrapeReq)
		if err != nil {
			a.log.Println(err)
			render.Render(w, r, web.ErrInvalidRequest(err))
			return
		}
		render.JSON(w, r, scrapeRes)
	}
}
