package sqlcommenter

import (
	"context"

	"go.opentelemetry.io/otel"
)

type (
	OtelTrace      struct{}
	CommentCarrier SQLComments
)

func NewOtelTrace() OtelTrace {
	return OtelTrace{}
}

func (hc OtelTrace) Comments(ctx context.Context) SQLComments {
	comments := NewCommentCarrier()
	otel.GetTextMapPropagator().Inject(ctx, comments)

	return SQLComments(comments)
}

func NewCommentCarrier() CommentCarrier {
	return make(CommentCarrier)
}

// Get returns the value associated with the passed key.
func (c CommentCarrier) Get(key string) string {
	return string(c[key])
}

// Set stores the key-value pair.
func (c CommentCarrier) Set(key string, value string) {
	c[key] = value
}

// Keys lists the keys stored in this carrier.
func (c CommentCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, string(k))
	}
	return keys
}
