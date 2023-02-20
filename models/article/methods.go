package article

import (
	"context"

	articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	"github.com/cvcio/mediawatch/pkg/es"
	"go.opencensus.io/trace"
)

func GetById(ctx context.Context, es *es.Elastic) (*articlesv2.Article, error) {
	ctx, span := trace.StartSpan(ctx, "model.article.GetById")
	defer span.End()

	return nil, nil
}
