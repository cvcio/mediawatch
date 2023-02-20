package relationships

type NodeArticle struct {
	Uid         string         `json:"uid"`
	DocId       string         `json:"doc_id"`
	Title       string         `json:"title"`
	Lang        string         `json:"lang"`
	CrawledAt   string         `json:"crawled_at"`
	PublishedAt string         `json:"published_at"`
	Url         string         `json:"url"`
	ScreenName  string         `json:"screen_name"`
	Type        string         `json:"type"`
	RelCount    int64          `json:"relCount,omitempty"`
	Similar     []*NodeArticle `json:"similar,omitempty"`
	Score       float64        `json:"score,omitempty"`
}

type NodeFeed struct {
	Uid        string `json:"uid"`
	FeedId     string `json:"feed_id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Url        string `json:"url"`
	Type       string `json:"type"`
}

type NodeEntity struct {
	Uid   string `json:"uid"`
	Label string `json:"label"`
	Type  string `json:"type"`
}
