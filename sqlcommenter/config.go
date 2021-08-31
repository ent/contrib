package sqlcommenter

import (
	"context"
)

// A Tagger is used by the driver to add tags to SQL queries.
type Tagger interface {
	Tag(context.Context) Tags
}

type (
	Option  func(*options)
	options struct {
		taggers []Tagger
	}
)

// WithTagger sets the taggers to be used to populate the SQL comment.
func WithTagger(taggers ...Tagger) Option {
	return func(opts *options) {
		opts.taggers = append(opts.taggers, taggers...)
	}
}

// WithTags appends the given tags to every SQL query.
func WithTags(tags Tags) Option {
	return func(opts *options) {
		opts.taggers = append(opts.taggers, NewStaticTagger(tags))
	}
}

// WithDriverVersion adds `db_driver` tag with the current version of ent.
func WithDriverVersion() Option {
	return func(opts *options) {
		opts.taggers = append(opts.taggers, NewDriverVersionTagger())
	}
}

func buildOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
