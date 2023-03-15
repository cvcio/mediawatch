package article

import (
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
)

// Entity extracted entities from enrich service
type Entity struct {
	EntityText string `json:"entity_text"`
	EntityType string `json:"entity_type"`
}

// Document Raw document stored on elasticsearch
type Document struct {
	DocID      string    `json:"docId"`
	Lang       string    `json:"lang"`
	CrawledAt  time.Time `json:"crawledAt"`
	ScreenName string    `json:"screen_name"`
	URL        string    `json:"url"`
	TweetID    int64     `json:"tweet_id"`
	TweetIDStr string    `json:"tweet_id_str"`
	Content    struct {
		Title       string    `json:"title,omitempty"`
		Excerpt     string    `json:"excerpt,omitempty"`
		Image       string    `json:"image,omitempty"`
		Body        string    `json:"body,omitempty"`
		Authors     []string  `json:"authors,omitempty"`
		Sources     []string  `json:"sources,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
		Categories  []string  `json:"categories,omitempty"`
		PublishedAt time.Time `json:"publishedAt,omitempty"`
		EditedAt    time.Time `json:"editedAt,omitempty"`
	} `json:"content"`
	NLP struct {
		Keywords  []string  `json:"keywords,omitempty"`
		StopWords []string  `json:"stopWords,omitempty"`
		Entities  []*Entity `json:"entities,omitempty"`
		Topics    []string  `json:"topics,omitempty"`
		Quotes    []string  `json:"quotes,omitempty"`
		Claims    []string  `json:"claims,omitempty"`
		Summary   string    `json:"summary,omitempty"`
	} `json:"nlp"`

	Feed     *feed.Feed `json:"feed,omitempty"`
	RelCount int64      `json:"relCount,omitempty"`
}

type Pagination struct {
	Total int `json:"total"`
}
