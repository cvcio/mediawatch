package web

import (
	"net/http"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/go-chi/render"
)

type FeedResponse struct {
	*feed.Feed
}

func NewFeedResponse(f *feed.Feed) *FeedResponse {
	return &FeedResponse{Feed: f}
}

func (rd *FeedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// rd.Elapsed = 10
	return nil
}

type FeedListResponse []*FeedResponse

func NewFeedListResponse(feeds []*feed.Feed) []render.Renderer {
	list := []render.Renderer{}
	for _, o := range feeds {
		list = append(list, NewFeedResponse(o))
	}
	return list
}
