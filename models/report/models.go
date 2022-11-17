package report

import (
	"time"

	"github.com/cvcio/mediawatch/models/article"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Deleted   bool      `json:"-" bson:"deleted"`

	UserID     string             `json:"-"`
	OrgID      string             `json:"-"`
	Title      string             `json:"title" bson:"title"`
	Style      string             `json:"style"`
	Articles   []article.Document `json:"articles"`
	Components []Component        `json:"components" bson:"components"`
}

// ReportsList of Report
type ReportsList struct {
	Data       []*Report   `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}

type ComponentFactory func() (Component, error)

var Components = make(map[string]ComponentFactory)

func init() {
	Components["articles"] = NewCmpArticles
	Components["article-view"] = NewCmpArticleView
	Components["timeline"] = NewCmpTimeline
	Components["tags"] = NewCmpTags
	Components["text"] = NewCmpText
	Components["html"] = NewCmpHTML
	Components["markdown"] = NewCmpMarkdown
}

type Component interface {
	Kind() string
	Payload() interface{}
}

type CmpArticles struct {
	Type string             `json:"type" default:"articles"`
	Body []article.Document `json:"body"`
}

func (c CmpArticles) Kind() string {
	return c.Type
}
func (c CmpArticles) Payload() interface{} {
	return c.Body
}

func NewCmpArticles() (Component, error) {
	return CmpArticles{}, nil
}

type CmpArticleView struct {
	Type string           `json:"type" default:"article-view"`
	Body article.Document `json:"body"`
}

func (c CmpArticleView) Kind() string {
	return c.Type
}
func (c CmpArticleView) Payload() interface{} {
	return c.Body
}
func NewCmpArticleView() (Component, error) {
	return CmpArticleView{}, nil
}

type CmpTimeline struct {
	Type string             `json:"type" default:"timeline"`
	Body []article.Document `json:"body"`
}

func (c CmpTimeline) Kind() string {
	return c.Type
}
func (c CmpTimeline) Payload() interface{} {
	return c.Body
}
func NewCmpTimeline() (Component, error) {
	return CmpTimeline{}, nil
}

type CmpTags struct {
	Type string   `json:"type" default:"tags"`
	Body []string `json:"body"`
}

func (c CmpTags) Kind() string {
	return c.Type
}
func (c CmpTags) Payload() interface{} {
	return c.Body
}
func NewCmpTags() (Component, error) {
	return CmpTags{}, nil
}

// type CmpNetworkView struct {
// 	Type string     `json:"type" default:"network"`
// 	Body []article.Document //`json:"articles"`
// 	Name string     `json:"Name"`
// }

type CmpText struct {
	Type string `json:"type" default:"text"`
	Body string `json:"body"`
}

func (c CmpText) Kind() string {
	return c.Type
}
func (c CmpText) Payload() interface{} {
	return c.Body
}
func NewCmpText() (Component, error) {
	return CmpText{}, nil
}

type CmpHTML struct {
	Type string `json:"type" default:"html"`
	Body string `json:"body"`
}

func (c CmpHTML) Kind() string {
	return c.Type
}
func (c CmpHTML) Payload() interface{} {
	return c.Body
}
func NewCmpHTML() (Component, error) {
	return CmpHTML{}, nil
}

type CmpMarkdown struct {
	Type string `json:"type" default:"markdown"`
	Body string `json:"body"`
}

func (c CmpMarkdown) Kind() string {
	return c.Type
}
func (c CmpMarkdown) Payload() interface{} {
	return c.Body
}
func NewCmpMarkdown() (Component, error) {
	return CmpMarkdown{}, nil
}
