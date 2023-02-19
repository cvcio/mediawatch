package nodes

import (
	"strconv"
	"strings"
	"time"
)

type Opts struct {
	Skip     int
	Limit    int
	DataType int
	Q        string
	Title    string
	Body     string
	Feeds    []string
	URL      string
	Langs    string
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
	s.DataType = 1
	s.Range.From = time.Now().AddDate(0, 0, -2).UTC().Format(time.RFC3339)
	s.Range.To = time.Now().UTC().Format(time.RFC3339)
	s.Range.By = "publishedAt"
	return s
}

// Skip Cursor
func Skip(i int) func(*Opts) {
	return func(s *Opts) {
		s.Skip = i
	}
}

// Limit Size
func Limit(i int) func(*Opts) {
	return func(s *Opts) {
		s.Limit = i
	}
}

// DataType query
func DataType(d int) func(*Opts) {
	return func(s *Opts) {
		s.DataType = d
	}
}

// Q Free-text search
func Q(q string) func(*Opts) {
	return func(s *Opts) {
		s.Q = q
	}
}

// Title match
func Title(f string) func(*Opts) {
	return func(s *Opts) {
		s.Title = f
	}
}

// Body match
func Body(f string) func(*Opts) {
	return func(s *Opts) {
		s.Body = f
	}
}

// Range query
func Range(By, From, To string) func(*Opts) {
	return func(s *Opts) {
		s.Range.By = By
		s.Range.From = From
		s.Range.To = To
	}
}

// RangeBy field
func RangeBy(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.By = i
	}
}

// RangeFrom query
func RangeFrom(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.From = i
	}
}

// RangeTo query
func RangeTo(i string) func(*Opts) {
	return func(s *Opts) {
		s.Range.To = i
	}
}

// Feeds query
func Feeds(i []string) func(*Opts) {
	return func(s *Opts) {
		s.Feeds = i
	}
}

// URL query
func URL(i string) func(*Opts) {
	return func(s *Opts) {
		s.URL = i
	}
}

// NewSearch parses query from url
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

	if i, err := strconv.Atoi(urlQuery("dataType")); err == nil {
		s.DataType = i
	}

	if i := urlQuery("q"); i != "" {
		s.Q = i
	}

	if i := urlQuery("title"); i != "" {
		s.Title = i
	}
	if i := urlQuery("body"); i != "" {
		s.Body = i
	}

	if i := urlQuery("feeds"); i != "" {
		s.Feeds = strings.Split(i, ",")
		for i, v := range s.Feeds {
			s.Feeds[i] = `f.screen_name =~ "(?i)` + v + `"`
		}
	}

	if i := urlQuery("url"); i != "" {
		s.URL = i
	}

	if i := urlQuery("from"); i != "" {
		s.Range.From = i
	}

	if i := urlQuery("to"); i != "" {
		s.Range.To = i
	}

	return s
}
