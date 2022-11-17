package feed

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feed is the fsource of articles
type Feed struct {
	ID  primitive.ObjectID `bson:"_id" json:"id"`
	UID string             `bson:"uid" json:"uid,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
	Deleted   bool      `json:"-" bson:"deleted"`

	// Name
	Name        string `json:"name,omitempty" bson:"name"`
	ScreenName  string `json:"screen_name,omitempty" bson:"screen_name"`
	Description string `json:"description,omitempty" bson:"description"`

	// Feed Info
	TwitterID           int64  `json:"twitter_id,omitempty" bson:"twitter_id"`
	TwitterIDStr        string `json:"twitter_id_str,omitempty" bson:"twitter_id_str"`
	TwitterProfileImage string `json:"twitter_profile_image,omitempty" bson:"twitter_profile_image"`
	RSS                 string `json:"rss,omitempty" bson:"rss"`

	// Info
	Email                string `json:"email,omitempty" bson:"email"`
	BusinessType         string `json:"business_type,omitempty" bson:"business_type"`
	ContentType          string `json:"content_type" bson:"content_type"`
	Country              string `json:"country,omitempty" bson:"country"`
	Lang                 string `json:"lang,omitempty" bson:"lang"`
	URL                  string `json:"url,omitempty" bson:"url"`
	Tier                 string `json:"tier,omitempty" bson:"tier"`
	BusinessOwner        string `json:"business_owner,omitempty" bson:"business_owner"`
	Registered           string `json:"registered,omitempty" bson:"registered"`
	PoliticalStance      string `json:"political_stance,omitempty" bson:"political_stance"`
	PoliticalOrientation string `json:"political_orientation,omitempty" bson:"political_orientation"`
	Locality             string `json:"locality,omitempty" bson:"locality"`

	// Scrape Info
	MetaClasses MetaClasses `json:"meta_classes,omitempty" bson:"meta_classes"`
	Status      string      `json:"status,omitempty" bson:"status"`
	TestURL     string      `json:"testURL,omitempty" bson:"testURL"`
	TestData    TestData    `json:"testData,omitempty" bson:"testData"`
}

// FeedsList of Feed
type FeedsList struct {
	Data       []*Feed     `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}

// UpdateFeed defines what information may be provided to modify an existing
// Account. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateFeed struct {
	Name                 *string      `json:"name" bson:"name"`
	ScreenName           *string      `json:"screen_name" bson:"screen_name"`
	TwitterID            *int64       `json:"twitter_id" bson:"twitter_id"`
	TwitterIDStr         *string      `json:"twitter_id_str" bson:"twitter_id_str,omitempty"`
	TwitterProfileImage  *string      `json:"twitter_profile_image" bson:"twitter_profile_image,omitempty"`
	Email                *string      `json:"email" bson:"email"`
	BusinessType         *string      `json:"business_type" bson:"business_type"`
	Country              *string      `json:"country" bson:"country"`
	Lang                 *string      `json:"lang" bson:"lang"`
	URL                  *string      `json:"url" bson:"url"`
	RSS                  *string      `json:"rss" bson:"rss"`
	MetaClasses          *MetaClasses `json:"meta_classes,omitempty" bson:"meta_classes,omitempty"`
	Status               *string      `json:"status" bson:"status"`
	TestURL              *string      `json:"testURL" bson:"testURL"`
	TestData             *TestData    `json:"testData,omitempty" bson:"testData,omitempty"`
	Tier                 *string      `json:"tier" bson:"tier"`
	BusinessOwner        *string      `json:"business_owner" bson:"business_owner"`
	Registered           *string      `json:"registered" bson:"registered"`
	PoliticalStance      *string      `json:"political_stance" bson:"political_stance"`
	PoliticalOrientation *string      `json:"political_orientation" bson:"political_orientation"`
	Locality             *string      `json:"locality" bson:"locality"`
	ContentType          *string      `json:"content_type" bson:"content_type"`
	Description          *string      `json:"description" bson:"description"`
}

// MetaClasses Model
type MetaClasses struct {
	API            string `json:"api,omitempty" bson:"api,omitempty"`
	FeedType       string `json:"feed_type,omitempty" bson:"feed_type,omitempty"`
	Title          string `json:"title,omitempty" bson:"title,omitempty"`
	Excerpt        string `json:"excerpt,omitempty" bson:"excerpt,omitempty"`
	Body           string `json:"body,omitempty" bson:"body,omitempty"`
	Authors        string `json:"authors,omitempty" bson:"authors,omitempty"`
	Sources        string `json:"sources,omitempty" bson:"sources,omitempty"`
	Tags           string `json:"tags,omitempty" bson:"tags,omitempty"`
	Categories     string `json:"categories,omitempty" bson:"categories,omitempty"`
	PublishedAt    string `json:"published_at,omitempty" bson:"published_at,omitempty"`
	EditedAt       string `json:"edited_at,omitempty" bson:"edited_at,omitempty"`
	TimezoneOffset string `json:"timezoneOffset,omitempty" bson:"timezoneOffset,omitempty"`
}

// TestData Model
type TestData struct {
	Title       string    `json:"title,omitempty" bson:"title,omitempty"`
	Body        string    `json:"body,omitempty" bson:"body,omitempty"`
	PublishedAt time.Time `json:"publishedAt,omitempty" bson:"publishedAt,omitempty"`
}
