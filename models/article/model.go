package article

import (
	"encoding/json"

	articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	commonv2 "github.com/cvcio/mediawatch/internal/mediawatch/common/v2"
)

type ArticlesData struct {
	Data  []*articlesv2.Article `json:"data,omitempty"`
	Total int64                 `json:"total,omitempty"`
}

func ParseDocument(source map[string]interface{}) (*articlesv2.Article, error) {
	doc := source["_source"].(map[string]interface{})
	jsonString, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var data *articlesv2.Article
	json.Unmarshal([]byte(jsonString), &data)
	return data, nil
}

func ParseCount(source map[string]interface{}) (float64, error) {
	count := source["count"].(float64)
	return count, nil
}

func ParseDocuments(source map[string]interface{}) (*articlesv2.ArticlesResponse, error) {
	var data articlesv2.ArticlesResponse
	var docs []*articlesv2.Article

	for _, hit := range source["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc, err := ParseDocument(hit.(map[string]interface{}))
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}

	data.Data = docs
	data.Pagination = &commonv2.Pagination{
		Total: int64(source["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	}

	return &data, nil
}
