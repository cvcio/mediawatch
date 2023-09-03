package feed

import "strings"

// ListOpts implements feed's list options stuct.
type ListOpts struct {
	Limit        int
	Offset       int
	Id           string
	Q            string
	UserName     string
	Hostname     string
	StreamStatus int
	StreamType   int
	Lang         string
	Country      string
	SortKey      string
	SortOrder    int
}

func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.StreamStatus = 0
	l.StreamType = 0
	l.Lang = ""
	l.Country = ""
	l.SortKey = "_id"
	l.SortOrder = -1
	l.Id = ""
	l.Q = ""
	l.UserName = ""
	l.Hostname = ""
	return l
}

func NewListOpts() []func(*ListOpts) {
	return make([]func(*ListOpts), 0)
}

func Limit(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Limit = i
	}
}

func Offset(i int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Offset = i
	}
}

func Id(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Id = i
	}
}
func Q(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Q = i
	}
}

func UserName(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.UserName = i
	}
}

func Hostname(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Hostname = i
	}
}

func StreamStatus(s int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.StreamStatus = s
	}
}

func StreamType(s int) func(*ListOpts) {
	return func(l *ListOpts) {
		l.StreamType = s
	}
}

func Lang(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Lang = strings.ToUpper(s)
	}
}

func Country(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Country = s
	}
}

func SortKey(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		if i == "" {
			l.SortKey = "_id"
		} else {
			l.SortKey = i
		}
	}
}

func SortOrder(s int) func(*ListOpts) {
	return func(l *ListOpts) {
		values := map[int]bool{-1: true, 1: true}
		if _, ok := values[s]; !ok {
			l.SortOrder = 1
		} else {
			l.SortOrder = s
		}
	}
}
