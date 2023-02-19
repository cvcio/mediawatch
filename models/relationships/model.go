package relationships

import "time"

type NodeArticle struct {
	Uid         string         `json:"uid"`
	DocId       string         `json:"doc_id"`
	Title       string         `json:"title"`
	Lang        string         `json:"lang"`
	CrawledAt   time.Time      `json:"crawled_at"`
	PublishedAt time.Time      `json:"published_at"`
	Url         string         `json:"url"`
	ScreenName  string         `json:"screen_name"`
	Type        string         `json:"type"`
	RelCount    int64          `json:"relCount,omitempty"`
	Similar     []*NodeArticle `json:"similar,omitempty"`
	Score       float64        `json:"score,omitempty"`
}

type NodeFeed struct {
	Uid   string `json:"uid"`
	DocId string `json:"doc_id"`
	Label string `json:"label"`
	Url   string `json:"url"`
	Type  string `json:"type"`
}

type NodeGpe struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type NodeOrg struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type NodePerson struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type NodeAuthor struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type NodeTopic struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}
