package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/go-chi/render"
)

type ArticleResponse struct {
	*article.Document
}

func NewArticleResponse(f *article.Document) *ArticleResponse {
	return &ArticleResponse{f}
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type ArticleListResponse []*ArticleResponse

func NewArticleListResponse(articles []*article.Document) []render.Renderer {
	list := []render.Renderer{}
	for _, o := range articles {
		list = append(list, NewArticleResponse(o))
	}
	return list
}
