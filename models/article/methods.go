package article

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cvcio/mediawatch/pkg/es"
	articlesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/articles/v2"
)

func GetById(ctx context.Context, es *es.Elastic, index string, id string) (*articlesv2.Article, error) {
	res, err := es.Client.Get(index, id, es.Client.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	// map data to interface
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(res.Status())
	}

	parsed, err := ParseDocument(r)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func Count(ctx context.Context, es *es.Elastic, opts *Opts) (int64, error) {
	args := opts.NewArticlesCountQuery(es.Client.Count)
	args = append(args, es.Client.Count.WithContext(ctx))

	res, err := es.Client.Count(args...)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, errors.New(res.Status())
	}

	// map data to interface
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, errors.New(res.Status())
	}

	parsed, _ := ParseCount(r)
	return int64(parsed), nil
}

func Search(ctx context.Context, es *es.Elastic, opts *Opts) (*articlesv2.ArticleList, error) {
	args := opts.NewArticlesSearchQuery(es.Client.Search)
	args = append(args, es.Client.Search.WithContext(ctx))

	res, err := es.Client.Search(args...)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	// map data to interface
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(res.String())
	}

	parsed, err := ParseDocuments(r)
	if err != nil {
		return nil, err
	}

	if _, ok := r["_scroll_id"]; !ok {
		return parsed, nil
	}

	scrollId := r["_scroll_id"]
	data, err := scroll(ctx, es, fmt.Sprintf("%s", scrollId))
	if err != nil {
		return nil, err
	}
	parsed.Data = append(parsed.Data, data...)

	return parsed, nil
}

func scroll(ctx context.Context, es *es.Elastic, scrollId string) ([]*articlesv2.Article, error) {
	var data []*articlesv2.Article

	for {
		res, err := es.Client.Scroll(
			es.Client.Scroll.WithScrollID(scrollId),
			es.Client.Scroll.WithScroll(time.Second*10),
		)

		if err != nil {
			return nil, err
		}

		if res.IsError() {
			return nil, errors.New(res.String())
		}

		// map data to interface
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return nil, errors.New(res.String())
		}

		res.Body.Close()

		parsed, err := ParseDocuments(r)
		if err != nil {
			return nil, err
		}
		if _, ok := r["_scroll_id"]; !ok {
			break
		}
		scrollId = fmt.Sprintf("%s", r["_scroll_id"])

		if len(parsed.Data) < 1 {
			break
		}

		data = append(data, parsed.Data...)

		if len(data) > 3600 {
			break
		}
	}

	_ = clearScroll(ctx, es, scrollId)
	return data, nil
}

func clearScroll(ctx context.Context, es *es.Elastic, scrollId string) error {
	if scrollId == "" {
		return nil
	}

	_, err := es.Client.ClearScroll(
		es.Client.ClearScroll.WithScrollID(scrollId),
	)

	return err
}

func Exists(ctx context.Context, es *es.Elastic, opts *Opts) bool {
	args := opts.NewArticlesExistsQuery(es.Client.Count)
	args = append(args, es.Client.Count.WithContext(ctx))

	res, err := es.Client.Count(args...)
	if err != nil {
		return false
	}

	defer res.Body.Close()

	if res.IsError() {
		return false
	}

	// map data to interface
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return false
	}

	parsed, err := ParseCount(r)
	if err != nil {
		return false
	}

	return int64(parsed) > 0
}
