package nodes

import "time"

// NodeArticle Article Struct for Graph Database
type NodeArticle struct {
	UID         string         `json:"uid"`
	DocID       string         `json:"docId"`
	Lang        string         `json:"lang"`
	CrawledAt   time.Time      `json:"crawledAt"`
	URL         string         `json:"url"`
	TweetID     int64          `json:"tweet_id"`
	TweetIDStr  string         `json:"tweet_id_str"`
	Title       string         `json:"title"`
	Body        string         `json:"body,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Authors     []*NodeAuthor  `json:"authors,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Categories  []string       `json:"categories,omitempty"`
	PublishedAt time.Time      `json:"publishedAt,omitempty"`
	EditedAt    time.Time      `json:"editedAt,omitempty"`
	Keywords    []string       `json:"keywords,omitempty"`
	Topics      []string       `json:"topics,omitempty"`
	Entities    []string       `json:"entities,omitempty"`
	Feed        *NodeFeed      `json:"feed,omitempty"`
	RelCount    int64          `json:"relCount,omitempty"`
	Similar     []*NodeArticle `json:"similar,omitempty"`
	Score       float64        `json:"score,omitempty"`
	ScreenName  string         `json:"screen_name,omitempty"`
}

// NodeFeed Feed Struct for Graph Database
type NodeFeed struct {
	UID                 string `json:"uid,omitempty"`
	Name                string `json:"name,omitempty" bson:"name"`
	ScreenName          string `json:"screen_name,omitempty"`
	TwitterID           int64  `json:"twitter_id,omitempty"`
	TwitterIDStr        string `json:"twitter_id_str,omitempty"`
	TwitterProfileImage string `json:"twitter_profile_image,omitempty"`
	URL                 string `json:"url,omitempty"`
}

// NodeEntity Entity Struct for Graph Database
type NodeEntity struct {
	UID        string `json:"uid,omitempty"`
	EntityText string `json:"entity_text,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
}

// NodeAuthor Author Struct for Graph Database
type NodeAuthor struct {
	UID    string `json:"uid,omitempty"`
	Author string `json:"author,omitempty"`
}

type Pagination struct {
	Total int64 `json:"total"`
}
