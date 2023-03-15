package article

import (
	"context"
	"strconv"
	"strings"

	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/olivere/elastic/v7"
)

type Opts struct {
	Index string
	Type  string
	Sort  struct {
		By  string
		Asc bool
	}

	DataType int
	Cases    bool

	Skip  int
	Limit int

	DocID string

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

func newOpts() *Opts {
	s := new(Opts)
	s.Skip = 0
	s.Limit = 48
	s.Langs = "EL"
	s.Cases = false
	s.DataType = 1
	s.Range.From = "now-7d"
	s.Range.To = "now"
	s.Range.By = "content.publishedAt"
	s.Sort.By = "content.publishedAt"
	s.Sort.Asc = false
	s.Index = "mediawatch_articles"
	s.Type = "document"
	return s
}

func Index(i string) func(*Opts) {
	return func(s *Opts) {
		s.Index = i
	}
}

func Type(i string) func(*Opts) {
	return func(s *Opts) {
		s.Type = i
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

func Cases(searchCases bool) func(*Opts) {
	return func(s *Opts) {
		s.Cases = searchCases
	}
}
func DataType(d int) func(*Opts) {
	return func(s *Opts) {
		s.DataType = d
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

func DocID(q string) func(*Opts) {
	return func(s *Opts) {
		s.DocID = q
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

func NewSearch(urlQuery func(string) string) *Opts {
	s := newOpts()

	if i, err := strconv.Atoi(urlQuery("limit")); err == nil {
		if i < 0 {
			i = 48
		}
		s.Limit = i
	}

	if i, err := strconv.Atoi(urlQuery("skip")); err == nil {
		if i < 0 {
			i = 0
		}
		s.Skip = i * s.Limit
	}

	if i := urlQuery("q"); i != "" {
		s.Q = i
	}

	if i, err := strconv.Atoi(urlQuery("dataType")); err == nil {
		s.DataType = i
	}

	if i := urlQuery("docId"); i != "" {
		s.DocID = i
	}

	if i := urlQuery("title"); i != "" {
		s.Title = i
	}

	if i := urlQuery("body"); i != "" {
		s.Body = i
	}

	if i := urlQuery("tags"); i != "" {
		s.Tags = i
	}

	if i := urlQuery("topics"); i != "" {
		s.Topics = i
	}

	if i := urlQuery("authors"); i != "" {
		s.Authors = i
	}

	if i := urlQuery("feeds"); i != "" {
		s.Feeds = i
	}

	if i := urlQuery("keywords"); i != "" {
		s.Keywords = i
	}

	if i := urlQuery("from"); i != "" {
		s.Range.From = i
	}

	if i := urlQuery("to"); i != "" {
		s.Range.To = i
	}

	if i := urlQuery("rangeBy"); i != "" {
		s.Range.By = i
	}

	if i := urlQuery("sortBy"); i != "" {
		s.Sort.By = i
	}

	if i := urlQuery("asc"); i != "" {
		s.Sort.Asc = true
	}

	return s
}

func OptsFromURL(urlQuery func(string) string) []func(*Opts) {
	opts := make([]func(*Opts), 0)

	if i, err := strconv.Atoi(urlQuery("limit")); err == nil {
		if i <= 0 {
			i = 48
		}
		opts = append(opts, Limit(i))
	}

	if i, err := strconv.Atoi(urlQuery("skip")); err == nil {
		if i < 0 {
			i = 0
		}
		opts = append(opts, Skip(i))
	}

	if i, err := strconv.Atoi(urlQuery("dataType")); err == nil {
		opts = append(opts, DataType(i))
	}

	if i := urlQuery("q"); i != "" {
		opts = append(opts, Q(i))
	}

	if i := urlQuery("docId"); i != "" {
		opts = append(opts, DocID(i))
	}

	if i := urlQuery("title"); i != "" {
		opts = append(opts, Title(i))
	}

	if i := urlQuery("body"); i != "" {
		opts = append(opts, Body(i))
	}

	if i := urlQuery("tags"); i != "" {
		opts = append(opts, Tags(i))
	}

	if i := urlQuery("topics"); i != "" {
		opts = append(opts, Topics(i))
	}

	if i := urlQuery("authors"); i != "" {
		opts = append(opts, Authors(i))
	}

	if i := urlQuery("feeds"); i != "" {
		opts = append(opts, Feeds(i))
	}

	if i := urlQuery("keywords"); i != "" {
		opts = append(opts, Keywords(i))
	}

	if i := urlQuery("from"); i != "" {
		opts = append(opts, RangeFrom(i))
	}

	if i := urlQuery("to"); i != "" {
		opts = append(opts, RangeTo(i))
	}

	if i := urlQuery("rangeBy"); i != "" {
		opts = append(opts, RangeBy(i))
	}

	if i := urlQuery("sortBy"); i != "" {
		opts = append(opts, sortBy(i))
	}

	if i := urlQuery("asc"); i != "" {
		opts = append(opts, sortAsc(true))
	}

	return opts
}

func (s *Opts) parse() []elastic.Query {
	queries := make([]elastic.Query, 0)

	if s.DocID != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.DocID).Field("docId"))
	}
	if s.Q != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Q))
	}
	if s.Title != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Title).Field("content.title"))
	}
	if s.Body != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Body).Field("content.body"))
	}
	if s.Tags != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Tags).Field("content.tags"))
	}
	if s.Topics != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Topics).Field("nlp.topics"))
	}
	if s.Authors != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Authors).Field("content.authors"))
	}
	if s.Langs != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Langs).Field("lang"))
	}
	if s.Feeds != "" {
		f := strings.Split(s.Feeds, ",")
		for i, v := range f {
			if !strings.Contains(v, "*") {
				f[i] = "*" + v + "*"
			}
		}
		queries = append(queries, elastic.NewQueryStringQuery(strings.Join(f, " OR "))) // append(queries, elastic.NewTermsQuery("screen_name", strings.Split(s.Feeds, ",")))
	}
	if s.Keywords != "" {
		queries = append(queries, elastic.NewQueryStringQuery(s.Keywords).Field("nlp.keywords"))
	}

	if s.Range.From != "" && s.Range.To != "" {
		queries = append(queries, elastic.NewRangeQuery(s.Range.By).Gte(s.Range.From).Lte(s.Range.To))
	} else {
		if s.Range.From != "" {
			queries = append(queries, elastic.NewRangeQuery(s.Range.By).Gte(s.Range.From))
		}
		if s.Range.To != "" {
			queries = append(queries, elastic.NewRangeQuery(s.Range.By).Lte(s.Range.To))
		}
	}

	return queries
}

func (s *Opts) Do(es *es.ES) (*elastic.SearchResult, error) {
	// parse queries from search
	queries := s.parse()

	// create a bool query that MUST with the queries from url string
	q := elastic.NewBoolQuery()
	q = q.Filter(queries...)

	// fmt.Println(q.Source())
	// skip := 0
	// if s.Skip > 0 {
	// 	skip = s.Skip * s.Limit
	// }

	// currentIndexTime := time.Now()
	// currentIndex := s.Index + currentIndexTime.Format("2006-01") + "," + s.Index + currentIndexTime.AddDate(0, -1, 0).Format("2006-01")

	return es.Client.Search().
		Index(s.Index).
		Query(q).
		From(s.Skip).Size(s.Limit).Sort(s.Sort.By, s.Sort.Asc).
		Do(context.Background())
}
