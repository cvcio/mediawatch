package article

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	"github.com/cvcio/mediawatch/pkg/es"
)

func GetById(ctx context.Context, es *es.Elastic, index string, id string) (*articlesv2.Article, error) {
	res, err := es.Client.Get(index, id, es.Client.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// retrun on response error
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

func Count(ctx context.Context, es *es.Elastic, index string, body map[string]interface{}) (int64, error) {
	if body == nil {
		body = map[string]interface{}{}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return 0, err
	}

	res, err := es.Client.Count(
		es.Client.Count.WithIndex(index),
		es.Client.Count.WithContext(ctx),
		es.Client.Count.WithBody(&buf),
	)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	// retrun on response error
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

func Search(ctx context.Context, es *es.Elastic, index string, body map[string]interface{}, size int) (*articlesv2.ArticlesResponse, error) {
	if body == nil {
		body = map[string]interface{}{}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	res, err := es.Client.Search(
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithContext(ctx),
		es.Client.Search.WithBody(&buf),
		es.Client.Search.WithSize(size),
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// retrun on response error
	if res.IsError() {
		return nil, errors.New(res.String())
	}

	// // map data to interface
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.New(res.String())
	}
	parsed, err := ParseDocuments(r)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
