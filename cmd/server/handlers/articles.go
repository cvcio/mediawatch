package handlers

import (
	"context"

	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"net/http"
	"time"

	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opencensus.io/trace"
)

// Articles represents the Articles handler
type Articles struct {
	log   *zap.Logger
	db    *db.MongoDB
	es    *es.Elastic
	neo4j *neo.Neo
}

// NewArticlesHandler returns a new instance of Articles
func NewArticlesHandler(log *zap.Logger, db *db.MongoDB, es *es.Elastic, neo4j *neo.Neo) (*Articles, error) {
	return &Articles{
		log:   log,
		db:    db,
		es:    es,
		neo4j: neo4j,
	}, nil
}

// List returns a list of articles
func (handler *Articles) List(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Articles.List")
	defer span.End()

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, web.ErrInvalidRequest(err))
		return
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	opts := article.NewOptsFromChiForm(c.Request.Form)
	if opts.Lang == "" {
		opts.Lang = claims.Lang
	}

	maxFrom := time.Now().AddDate(0, 0, -claims.SearchLimit).UTC()
	maxTo := time.Now().UTC()

	from, err := time.Parse(time.RFC3339, opts.Range.From)
	if err != nil {
		from = maxFrom
	}

	to, err := time.Parse(time.RFC3339, opts.Range.To)
	if err != nil {
		to = maxTo
	}

	if maxFrom.Sub(from).Hours()/24 < 0 {
		from = maxFrom
	}

	opts.Range.From = from.Format(time.RFC3339)
	opts.Range.To = to.Format(time.RFC3339)

	data, err := article.Search(context.Background(), handler.es, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	for _, v := range data.Data {
		f, err := feed.Get(ctx, handler.db, feed.Id(v.FeedId))
		if err != nil {
			continue
		}

		v.Feed = f
		v.Nlp.Stopwords = nil
		v.RelCount = 0

		if opts.CountCases {
			func() {
				ses := handler.neo4j.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
				defer func() { _ = ses.Close(ctx) }()

				res, err := ses.Run(ctx, relationships.CountSimilarTpl, map[string]interface{}{
					"doc_id": v.DocId,
				})

				if err != nil {
					return
				}

				if res.Next(ctx) {
					record := res.Record()
					if value, ok := record.Get("count"); ok {
						if value.(int64) > 0 {
							v.RelCount = value.(int64)
						}
					}
				}
			}()
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data":       data.Data,
		"pagination": data.Pagination,
	})
}

// Get returns a single article
func (handler *Articles) Get(c *gin.Context) {
	_, span := trace.StartSpan(c.Request.Context(), "handlers.Articles.Get")
	defer span.End()
	c.JSON(http.StatusInternalServerError, web.ErrInternalError)
}

// ParseArticle parses an article
func (handler *Articles) ParseArticle(c *gin.Context) {
	_, span := trace.StartSpan(c.Request.Context(), "handlers.Articles.ParseArticle")
	defer span.End()
	c.JSON(http.StatusInternalServerError, web.ErrInternalError)
}
