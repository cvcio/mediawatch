package article

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Opts struct {
	Index string `json:"index,omitempty"`
	DocId string `json:"doc_id,omitempty"`

	Q         string `json:"q,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	Tags      string `json:"tags,omitempty"`
	Keywords  string `json:"keywords,omitempty"`
	Topics    string `json:"topics,omitempty"`
	Entities  string `json:"entities,omitempty"`
	Authors   string `json:"authors,omitempty"`
	Lang      string `json:"lang,omitempty"`
	Feeds     string `json:"feeds,omitempty"`
	Hostnames string `json:"hostnames,omitempty"`

	CountCases  bool `json:"count_cases,omitempty"`
	IncludeRels bool `json:"include_rels,omitempty"`

	Skip   int  `json:"skip,omitempty"`
	Limit  int  `json:"limit,omitempty"`
	Scroll bool `json:"-"`

	Sort struct {
		By  string `json:"by,omitempty"`
		Asc bool   `json:"asc,omitempty"`
	} `json:"sort,omitempty"`

	Range struct {
		By   string `json:"by,omitempty"`
		From string `json:"from,omitempty"`
		To   string `json:"to,omitempty"`
	} `json:"range,omitempty"`
}

func NewOpts() *Opts {
	opts := new(Opts)
	opts.Skip = 0
	opts.Limit = 24
	opts.Range.From = "now-7d"
	opts.Range.To = "now"
	opts.Range.By = "content.published_at"
	opts.Sort.By = "content.published_at"
	opts.Sort.Asc = false
	opts.Index = "mediawatch_articles_*"
	opts.Scroll = false
	return opts
}

func NewOptsForm(j []byte) *Opts {
	opts := NewOpts()
	json.Unmarshal(j, &opts)
	return opts
}

func (o *Opts) Query() map[string]interface{} {
	filter := []map[string]interface{}{}
	must := []map[string]interface{}{}

	if o.Q != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": o.Q,
			},
		})
	}
	if o.Title != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Title,
				"fields": []string{"content.title"},
			},
		})
	}
	if o.Body != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Body,
				"fields": []string{"content.body"},
			},
		})
	}
	if o.Tags != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Tags,
				"fields": []string{"content.tags"},
			},
		})
	}
	if o.Keywords != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Keywords,
				"fields": []string{"nlp.keywords"},
			},
		})
	}
	if o.Topics != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Topics,
				"fields": []string{"nlp.topics.text"},
			},
		})
	}
	if o.Entities != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Entities,
				"fields": []string{"nlp.entities.text"},
			},
		})
	}
	if o.Authors != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Authors,
				"fields": []string{"content.authors"},
			},
		})
	}
	if o.Lang != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"lang": o.Lang,
			},
		})
	}
	if o.Feeds != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Feeds,
				"fields": []string{"screen_name"},
			},
		})
	}
	if o.Hostnames != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query":  o.Hostnames,
				"fields": []string{"hostname"},
			},
		})
	}

	if o.Range.From != "" && o.Range.To != "" {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				o.Range.By: map[string]interface{}{
					"gte": o.Range.From,
					"lte": o.Range.To,
				},
			},
		})
	} else {
		if o.Range.From != "" {
			filter = append(filter, map[string]interface{}{
				"range": map[string]interface{}{
					o.Range.By: map[string]interface{}{
						"gte": o.Range.From,
					},
				},
			})
		}
		if o.Range.To != "" {
			filter = append(filter, map[string]interface{}{
				"range": map[string]interface{}{
					o.Range.By: map[string]interface{}{
						"lte": o.Range.To,
					},
				},
			})
		}
	}
	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filter,
				"must":   must,
			},
		},
	}
}

func Index(i string) func(*Opts) {
	return func(s *Opts) {
		s.Index = i
	}
}

func Sort(By string, Asc bool) func(*Opts) {
	return func(s *Opts) {
		s.Sort.By = By
		s.Sort.Asc = Asc
	}
}

func sortBy(i string) func(*Opts) {
	return func(s *Opts) {
		s.Sort.By = i
	}
}

func sortAsc(i bool) func(*Opts) {
	return func(s *Opts) {
		s.Sort.Asc = i
	}
}

func Skip(i int) func(*Opts) {
	return func(s *Opts) {
		s.Skip = i
	}
}

func Limit(i int) func(*Opts) {
	return func(s *Opts) {
		s.Limit = i
	}
}

func DocId(q string) func(*Opts) {
	return func(s *Opts) {
		s.DocId = q
	}
}

func Q(q string) func(*Opts) {
	return func(s *Opts) {
		s.Q = q
	}
}

func Title(f string) func(*Opts) {
	return func(s *Opts) {
		s.Title = f
	}
}

func Tags(f string) func(*Opts) {
	return func(s *Opts) {
		s.Tags = f
	}
}

func Body(f string) func(*Opts) {
	return func(s *Opts) {
		s.Body = f
	}
}

func Topics(f string) func(*Opts) {
	return func(s *Opts) {
		s.Topics = f
	}
}

func Entities(f string) func(*Opts) {
	return func(s *Opts) {
		s.Entities = f
	}
}

func Authors(f string) func(*Opts) {
	return func(s *Opts) {
		s.Authors = f
	}
}

func Lang(f string) func(*Opts) {
	return func(s *Opts) {
		s.Lang = f
	}
}

func Feeds(f string) func(*Opts) {
	return func(s *Opts) {
		s.Feeds = f
	}
}

func Hostnames(f string) func(*Opts) {
	return func(s *Opts) {
		s.Hostnames = f
	}
}

func Keywords(f string) func(*Opts) {
	return func(s *Opts) {
		s.Keywords = f
	}
}

func Range(By, From, To string) func(*Opts) {
	return func(s *Opts) {
		s.Range.By = By
		s.Range.From = From
		s.Range.To = To
	}
}

func RangeBy(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.By = i
	}
}

func RangeFrom(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.From = i
	}
}

func RangeTo(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.To = i
	}
}

func (o *Opts) NewArticlesSearchQuery(es esapi.Search) []func(*esapi.SearchRequest) {
	buf, _ := json.Marshal(o.Query())

	opts := make([]func(*esapi.SearchRequest), 0)
	opts = append(opts, esapi.Search.WithBody(es, strings.NewReader(string(buf))))
	opts = append(opts, esapi.Search.WithIndex(es, o.Index))
	opts = append(opts, esapi.Search.WithFrom(es, o.Skip))
	opts = append(opts, esapi.Search.WithSize(es, o.Limit))
	opts = append(opts, esapi.Search.WithTimeout(es, time.Second*10))
	opts = append(opts, esapi.Search.WithTrackTotalHits(es, true))
	if o.Scroll {
		opts = append(opts, esapi.Search.WithScroll(es, time.Second*10))
	}
	return opts
}

func (o *Opts) NewArticlesCountQuery(es esapi.Count) []func(*esapi.CountRequest) {
	buf, _ := json.Marshal(o.Query())

	opts := make([]func(*esapi.CountRequest), 0)
	opts = append(opts, esapi.Count.WithBody(es, strings.NewReader(string(buf))))
	opts = append(opts, esapi.Count.WithIndex(es, o.Index))
	return opts
}
