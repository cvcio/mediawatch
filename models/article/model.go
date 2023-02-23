package article

import articlesv2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"

type ArticlesData struct {
	Data  []*articlesv2.Article `json:"data,omitempty"`
	Total int64                 `json:"total,omitempty"`
}
