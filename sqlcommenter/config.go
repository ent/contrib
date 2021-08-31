package sqlcommenter

import (
	"context"
)

type Tagger interface {
	Tag(context.Context) SQLCommentTags
}

type (
	Option  func(*options)
	options struct {
		commenters     []Tagger
		globalComments SQLCommentTags
	}
)

// WithTagger sets the taggers to be used to populate the SQL comment
func WithTagger(taggers ...Tagger) Option {
	return func(opts *options) {
		opts.commenters = append(opts.commenters, taggers...)
	}
}

// WithTags appends the given tags to every SQL query.
func WithTags(tags SQLCommentTags) Option {
	return func(opts *options) {
		opts.commenters = append(opts.commenters, NewStaticTagger(tags))
	}
}

func buildOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
