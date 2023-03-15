package article

import (
	"context"
	"encoding/json"

	"go.opencensus.io/trace"

	"github.com/cvcio/mediawatch/pkg/es"
)

// Get Get Raw Article by ID from ealasticsearch
func Get(ctx context.Context, es *es.ES, id string) (*Document, error) {
	ctx, span := trace.StartSpan(ctx, "model.article.Get")
	defer span.End()

	// Retrive Source Article
	opts := newOpts()
	// opts.Index = "articles_new"
	article, err := es.Client.Get().Index(opts.Index).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}
	var a Document
	err = json.Unmarshal(article.Source, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// List All Articles from ealasticsearch
func List(ctx context.Context, es *es.ES, options *Opts) ([]*Document, int64, error) {
	_, span := trace.StartSpan(ctx, "model.article.List")
	defer span.End()

	// opts := newOpts()

	// for _, o := range options {
	// 	o(opts)
	// }

	results, err := options.Do(es)

	if err != nil {
		return nil, 0, err
	}

	total := results.TotalHits()

	var docs []*Document
	for _, doc := range results.Hits.Hits {
		var a Document
		err := json.Unmarshal(doc.Source, &a)
		if err != nil {
			return nil, 0, err
		}
		docs = append(docs, &a)
	}

	return docs, total, nil
}
