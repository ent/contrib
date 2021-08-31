package sqlcommenter

import (
	"context"
)

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

func buildOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
