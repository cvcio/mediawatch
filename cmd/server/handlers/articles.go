package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	scrape_pb "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

type Articles struct {
	log    *logrus.Logger
	db     *db.MongoDB
	es     *es.Elastic
	neo4j  *neo.Neo
	scrape scrape_pb.ScrapeServiceClient
}

func NewArticlesHandler(log *logrus.Logger, db *db.MongoDB, es *es.Elastic, neo4j *neo.Neo, scrape scrape_pb.ScrapeServiceClient) *Articles {
	return &Articles{
		log:    log,
		db:     db,
		es:     es,
		neo4j:  neo4j,
		scrape: scrape,
	}
}

// q, from, to, title, body, feeds, topics, skip
func (handler *Articles) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "handlers.Articles.List")
		defer span.End()

		claims, ok := ctx.Value(auth.Key).(auth.Claims)
		if !ok {
			render.Render(w, r, web.ErrUnauthorized)
			return
		}

		opts := article.NewOptsForm(nil)
		opts.Lang = claims.Lang
		if opts.Lang == "" {
			opts.Lang = "EL"
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

		data, err := article.Search(context.Background(), handler.es, opts)
		if err != nil {
			handler.log.Error(err)
			render.Render(w, r, web.ErrInternalError)
		}

		render.JSON(w, r, map[string]interface{}{
			"data":       data.Data,
			"pagination": data.Pagination,
		})
	}
}
func (handler *Articles) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, web.ErrInternalError)
	}
}
func (handler *Articles) ParseArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, web.ErrInternalError)
	}
}
