package sqlcommenter

import (
	"context"
)

type Commenter interface {
	Comments(context.Context) SQLComments
}

type (
	CommentsHandler func(context.Context) SQLComments
	Option          func(*options)
	options         struct {
		commenters     []Commenter
		globalComments SQLComments
	}
)

// WithCommenter overrides the default comments generator handler.
// default comments added via WithComments will still be applied.
func WithCommenter(commenters ...Commenter) Option {
	return func(opts *options) {
		opts.commenters = append(opts.commenters, commenters...)
	}
}

// WithComments appends the given comments to every sql query.
func WithComments(comments SQLComments) Option {
	return func(opts *options) {
		opts.commenters = append(opts.commenters, NewStaticCommenter(comments))
	}
}

func buildOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
