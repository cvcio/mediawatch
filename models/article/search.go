package article

import (
	"strings"

	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Opts struct {
	Index string
	Sort  struct {
		By  string
		Asc bool
	}
	From int
	Size int

	DocId string

	Q        string
	Title    string
	Body     string
	Tags     string
	Topics   string
	Authors  string
	Langs    string
	Feeds    string
	Keywords string

	Range struct {
		By   string
		From string
		To   string
	}
}

func NewOpts() *Opts {
	opts := new(Opts)
	opts.From = 0
	opts.Size = 48
	opts.Langs = "EL"
	opts.Range.From = "now-7d"
	opts.Range.To = "now"
	opts.Range.By = "content.publishedAt"
	opts.Sort.By = "content.publishedAt"
	opts.Sort.Asc = false
	opts.Index = "mediawatch_articles"
	return opts
}

func (opts *Opts) Query() *strings.Reader {
	var b strings.Builder
	b.WriteString("")
	query := strings.NewReader(b.String())
	return query
}

func (opts *Opts) Do(elastic *es.Elastic) (*esapi.Response, error) {
	return elastic.Client.Search(
		elastic.Client.Search.WithIndex(opts.Index),
		elastic.Client.Search.WithBody(opts.Query()),
		elastic.Client.Search.WithSort(opts.Sort.By),
		elastic.Client.Search.WithFrom(opts.From),
		elastic.Client.Search.WithSize(opts.Size),
	)
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

func From(i int) func(*Opts) {
	return func(s *Opts) {
		s.From = i
	}
}

func Size(i int) func(*Opts) {
	return func(s *Opts) {
		s.Size = i
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

func Authors(f string) func(*Opts) {
	return func(s *Opts) {
		s.Authors = f
	}
}

func Langs(f string) func(*Opts) {
	return func(s *Opts) {
		s.Langs = f
	}
}

func Feeds(f string) func(*Opts) {
	return func(s *Opts) {
		s.Feeds = f
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
