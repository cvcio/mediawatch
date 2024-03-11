package passage

import "strings"

// ListOpts implements passage's list options struct.
type ListOpts struct {
	Limit  int
	Offset int
	Id     string
	Lang   string
}

func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.Lang = ""
	l.Id = ""
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

func Lang(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Lang = strings.ToUpper(s)
	}
}
