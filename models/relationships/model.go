package relationships

import "strings"

// NodeArticle struct.
type NodeArticle struct {
	Uid         string         `json:"uid"`
	DocId       string         `json:"doc_id"`
	Title       string         `json:"title,omitempty"`
	Lang        string         `json:"lang,omitempty"`
	CrawledAt   string         `json:"crawled_at,omitempty"`
	PublishedAt string         `json:"published_at,omitempty"`
	Url         string         `json:"url,omitempty"`
	ScreenName  string         `json:"screen_name,omitempty"`
	Type        string         `json:"type,omitempty"`
	RelCount    int64          `json:"relCount,omitempty"`
	Similar     []*NodeArticle `json:"similar,omitempty"`
	Score       float64        `json:"score,omitempty"`
}

// NodeFeed struct.
type NodeFeed struct {
	Uid        string `json:"uid"`
	FeedId     string `json:"feed_id,omitempty"`
	Name       string `json:"name,omitempty"`
	ScreenName string `json:"screen_name,omitempty"`
	Url        string `json:"url,omitempty"`
	Type       string `json:"type,omitempty"`
}

// NodeEntity struct.
type NodeEntity struct {
	Uid   string `json:"uid"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`
}

// getEntityType return the type of an entity.
func getEntityType(entityType string) string {
	switch strings.ToLower(entityType) {
	case "feed":
		return "Feed"
	case "gpe":
		return "GPE"
	case "org":
		return "Organization"
	case "person":
		return "Person"
	case "author":
		return "Author"
	case "topic":
		return "Topic"
	default:
		return "Article"
	}
}
