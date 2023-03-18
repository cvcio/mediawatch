package feed

import "strings"

// ListOpts implements feed's list options stuct.
type ListOpts struct {
	Limit      int
	Offset     int
	Q          string
	Deleted    bool
	Status     string
	StreamType string
	Lang       string
}

func DefaultOpts() ListOpts {
	l := ListOpts{}
	l.Offset = 0
	l.Limit = 24
	l.Deleted = false
	l.Status = ""
	l.StreamType = ""
	l.Lang = "EL"
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

func Q(i string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Q = i
	}
}

func Deleted() func(*ListOpts) {
	return func(l *ListOpts) {
		l.Deleted = true
	}
}

func Status(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Status = s
	}
}

func StreamType(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.StreamType = s
	}
}

func Lang(s string) func(*ListOpts) {
	return func(l *ListOpts) {
		l.Lang = strings.ToUpper(s)
	}
}
