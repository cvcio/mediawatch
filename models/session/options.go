package session

import (
	sessionsv2 "github.com/cvcio/mediawatch/internal/mediawatch/sessions/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListOptions model for account endpoint queries
type ListOptions struct {
	Limit   int64
	Id      string
	Key     string
	Value   string
	Message string
}

// DefaultOptions returns the defaults
func DefaultOptions() ListOptions {
	l := ListOptions{}
	l.Limit = 24
	return l
}

// NewListOpts returns list options
func NewListOpts() []func(*ListOptions) {
	return make([]func(*ListOptions), 0)
}

// ParseOptions parse query filters from interface
func ParseOptions(o *sessionsv2.Session) []func(*ListOptions) {
	opts := NewListOpts()

	if o.GetKey() != "" {
		opts = append(opts, Key(o.Id))
	}
	if o.GetValue() != "" {
		opts = append(opts, Value(o.Id))
	}
	if o.GetMessage() != "" {
		opts = append(opts, Message(o.Id))
	}

	return opts
}

// Key Option
func Key(i string) func(*ListOptions) {
	return func(l *ListOptions) {
		l.Key = i
	}
}

// Value Option
func Value(i string) func(*ListOptions) {
	return func(l *ListOptions) {
		l.Value = i
	}
}

// Message Option
func Message(i string) func(*ListOptions) {
	return func(l *ListOptions) {
		l.Value = i
	}
}

// Filter parses an option list to a bson.M query
func Filter(optionsList ...func(*ListOptions)) (ListOptions, bson.M) {
	options := DefaultOptions()
	for _, o := range optionsList {
		o(&options)
	}

	filter := bson.M{}
	if options.Id != "" {
		if oid, err := primitive.ObjectIDFromHex(options.Id); err == nil {
			filter["_id"] = bson.M{"$eq": oid}
		}
	}
	if options.Key != "" {
		filter["key"] = bson.M{"$eq": options.Key}
	}
	if options.Value != "" {
		filter["value"] = bson.M{"$eq": options.Value}
	}
	if options.Message != "" {
		filter["message"] = bson.M{"$eq": options.Message}
	}

	return options, filter
}
