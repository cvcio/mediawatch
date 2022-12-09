package handlers

import (
	"context"

	"github.com/bufbuild/connect-go"
	articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	"github.com/cvcio/mediawatch/internal/mediawatch/articles/v2/articlesv2connect"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"go.uber.org/zap"
)

// ArticlesHandler implements feeds connect service
type ArticlesHandler struct {
	log           *zap.SugaredLogger
	mg            *db.MongoDB
	elastic       *es.Elastic
	authenticator *auth.JWTAuthenticator
	// Embed the unimplemented server
	articlesv2connect.UnimplementedArticlesServiceHandler
}

// NewArticlesHandler returns a new ArticlesHandler service.
func NewArticlesHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, authenticator *auth.JWTAuthenticator) *ArticlesHandler {
	return &ArticlesHandler{log: log, mg: mg, elastic: elastic, authenticator: authenticator}
}

// GetArticle return a single article.
func (h *ArticlesHandler) GetArticle(ctx context.Context, req *connect.Request[articlesv2.ArticleRequest]) (*connect.Response[articlesv2.ArticleResponse], error) {
	return connect.NewResponse(&articlesv2.ArticleResponse{}), nil
}

// GetArticles returns a list of aricles.
func (h *ArticlesHandler) GetArticles(ctx context.Context, req *connect.Request[articlesv2.ArticlesRequest]) (*connect.Response[articlesv2.ArticlesResponse], error) {
	return connect.NewResponse(&articlesv2.ArticlesResponse{}), nil
}

// Stream streams articles in real time.
func (h *ArticlesHandler) StreamArticles(ctx context.Context, req *connect.Request[articlesv2.ArticlesRequest], stream *connect.ServerStream[articlesv2.ArticleResponse]) error {
	return nil
}

// StreamRelatedArticles streams article relationships for a specific article in real time.
func (h *ArticlesHandler) StreamRelatedArticles(ctx context.Context, req *connect.Request[articlesv2.ArticlesRequest], stream *connect.ServerStream[articlesv2.ArticleResponse]) error {
	return nil
}
