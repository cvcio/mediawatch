package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/auth"
	articlesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
)

type Cases struct {
	log   *zap.Logger
	db    *db.MongoDB
	es    *es.Elastic
	neo4j *neo.Neo
}

func NewCasesHandler(log *zap.Logger, db *db.MongoDB, es *es.Elastic, neo4j *neo.Neo) (*Cases, error) {
	return &Cases{
		log:   log,
		db:    db,
		es:    es,
		neo4j: neo4j,
	}, nil
}

func (handler *Cases) List(c *gin.Context) {
	_, span := trace.StartSpan(c.Request.Context(), "handlers.Cases.List")
	defer span.End()
	c.JSON(http.StatusInternalServerError, web.ErrInternalError)
}
func (handler *Cases) Get(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Cases.Get")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	docId := c.Param("id")
	lang := strings.ToLower(claims.Lang)

	var docs []*articlesv2.Article

	a, err := article.GetById(ctx, handler.es, "mediawatch_articles_"+lang, docId)
	if err != nil {
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}
	x, err := feed.Get(ctx, handler.db, feed.Id(a.FeedId))
	if err != nil {
		c.JSON(http.StatusNotFound, web.ErrNotFound)
		return
	}

	a.Feed = x
	a.Nlp.Stopwords = nil

	docs = append(docs, a)

	ses := handler.neo4j.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() { _ = ses.Close(ctx) }()

	res, err := ses.Run(ctx, relationships.GetSimilarTpl, map[string]interface{}{
		"doc_id": docId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	if res.Next(ctx) {
		record := res.Record()
		if value, ok := record.Get("json"); ok {
			bodyBytes, _ := json.Marshal(value)
			_ = json.Unmarshal(bodyBytes, &a.Similar)
		}
	}

	for i := range a.Similar {
		if a.Similar[i].DocId != docId {
			d, err := article.GetById(ctx, handler.es, "mediawatch_articles_"+lang, a.Similar[i].DocId)
			if err != nil {
				continue
			}

			f, err := feed.Get(ctx, handler.db, feed.Id(d.FeedId))
			if err != nil {
				continue
			}

			d.Feed = f
			d.Nlp.Stopwords = nil

			docs = append(docs, d)
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data":       docs,
		"pagination": nil,
	})
}

func (handler *Cases) Count(c *gin.Context) {
	ctx, span := trace.StartSpan(c.Request.Context(), "handlers.Cases.Count")
	defer span.End()

	_, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, web.ErrUnauthorized)
		return
	}

	docId := c.Param("id")

	ses := handler.neo4j.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() { _ = ses.Close(ctx) }()

	res, err := ses.Run(ctx, relationships.CountSimilarTpl, map[string]interface{}{
		"doc_id": docId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, web.ErrInternalError)
		return
	}

	var a articlesv2.Article

	if res.Next(ctx) {
		record := res.Record()
		if value, ok := record.Get("count"); ok {
			if value.(int64) > 0 {
				a.RelCount = value.(int64)
			}
		}
	}

	c.JSON(http.StatusOK, a)
}
