package handlers

import (
	"context"
	"encoding/json"

	"connectrpc.com/connect"
	"github.com/cvcio/mediawatch/models/article"
	"github.com/cvcio/mediawatch/models/relationships"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	articlesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2"
	"github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2/articlesv2connect"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/zap"
)

// ArticlesHandler handles operations related to articles, leveraging MongoDB, Elastic, Neo4j, JWT authentication, and Redis for caching.
type ArticlesHandler struct {
	log             *zap.SugaredLogger
	mg              *db.MongoDB
	elastic         *es.Elastic
	neo             *neo.Neo
	authenticator   *auth.JWTAuthenticator
	rdb             *redis.RedisClient
	articlesCh      chan []byte
	relationshipsCh chan []byte

	// Embed the unimplemented server
	articlesv2connect.UnimplementedArticlesServiceHandler
}

// NewArticlesHandler initializes and returns a new ArticlesHandler instance.
// It takes a configuration, logger, MongoDB, ElasticSearch, Neo4j, JWT authenticator, and Redis client as arguments.
func NewArticlesHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, neo *neo.Neo, authenticator *auth.JWTAuthenticator, rdb *redis.RedisClient) *ArticlesHandler {
	// articlesCh, _ := rdb.Subscribe("mediawatch_articles")
	// relationshipsCh, _ := rdb.Subscribe("mediawatch_relationships")
	return &ArticlesHandler{log: log, mg: mg, elastic: elastic, neo: neo, authenticator: authenticator, rdb: rdb, articlesCh: nil, relationshipsCh: nil}
}

// Consume processes messages from the articles channel and handles them accordingly.
func (h *ArticlesHandler) Consume() {}

// GetArticle retrieves an article by its ID from Elasticsearch and optionally counts similar articles if requested.
func (h *ArticlesHandler) GetArticle(ctx context.Context, req *connect.Request[articlesv2.QueryArticle]) (*connect.Response[articlesv2.Article], error) {
	data, err := article.GetById(ctx, h.elastic, "mediawatch_articles_el", req.Msg.DocId)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if req.Msg.CountCases {
		// return count per article
		ses := h.neo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer func() { _ = ses.Close(ctx) }()
		res, err := ses.Run(ctx, relationships.CountSimilarTpl, map[string]interface{}{
			"doc_id": data.DocId,
		})

		if err != nil {
			h.log.Error(err)
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		if res.Next(ctx) {
			record := res.Record()
			data.RelCount = 0
			if value, ok := record.Get("count"); ok {
				if value.(int64) > 0 {
					data.RelCount = value.(int64)
				}
			}
		}
	}

	data.Nlp.Stopwords = nil
	return connect.NewResponse(data), nil
}

// GetArticles retrieves a list of articles matching the specified query criteria.
func (h *ArticlesHandler) GetArticles(ctx context.Context, req *connect.Request[articlesv2.QueryArticle]) (*connect.Response[articlesv2.ArticleList], error) {
	j, _ := json.Marshal(req.Msg)
	opts := article.NewOptsForm(j)

	data, err := article.Search(ctx, h.elastic, opts)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if req.Msg.CountCases {
		// return count per article
		for _, v := range data.Data {
			func() {
				ses := h.neo.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
				defer func() { _ = ses.Close(ctx) }()

				res, err := ses.Run(ctx, relationships.CountSimilarTpl, map[string]interface{}{
					"doc_id": v.DocId,
				})

				if err != nil {
					h.log.Error(err)
					v.RelCount = 0
					return
				}

				if res.Next(ctx) {
					record := res.Record()
					v.RelCount = 0
					if value, ok := record.Get("count"); ok {
						if value.(int64) > 0 {
							v.RelCount = value.(int64)
						}
					}
				}
			}()
		}
	}
	for _, v := range data.Data {
		v.Nlp.Stopwords = nil
	}

	return connect.NewResponse(data), nil
}

// StreamArticles streams articles in real-time to the client.
func (h *ArticlesHandler) StreamArticles(ctx context.Context, req *connect.Request[articlesv2.QueryArticle], stream *connect.ServerStream[articlesv2.ArticleList]) error {
	h.log.Infof("[ARTICLES] Streaming articles: %s", req.Msg.String())
	pubsub, articlesCh, err := h.rdb.Subscribe(ctx, "mediawatch_articles_*")
	if err != nil {
		h.log.Errorf("[ARTICLES] Pub/Sub subscriprion error: %s", err.Error())
		return connect.NewError(connect.CodeInternal, err)
	}
	defer func() { _ = pubsub.Close() }()

	for m := range articlesCh {
		var a *articlesv2.Article
		if err := json.Unmarshal(m, &a); err != nil {
			h.log.Errorf("[ARTICLES] Unmarshal error: %s", err.Error())
			continue
		}

		data := &articlesv2.ArticleList{}
		data.Data = append(data.Data, a)

		h.log.Debugf("[ARTICLES] Streaming article: %s", article.DocId)

		if err := stream.Send(data); err != nil {
			h.log.Errorf("[ARTICLES] Internal error: %s", err.Error())
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}

// StreamRelatedArticles streams related articles based on the given query parameters.
func (h *ArticlesHandler) StreamRelatedArticles(ctx context.Context, req *connect.Request[articlesv2.QueryArticle], stream *connect.ServerStream[articlesv2.ArticleList]) error {
	return nil
}
