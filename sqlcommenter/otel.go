package sqlcommenter

import (
	"context"

	"go.opentelemetry.io/otel"
)

type (
	OtelTrace      struct{}
	CommentCarrier SQLCommentTags
)

func NewOtelTrace() OtelTrace {
	return OtelTrace{}
}

func (hc OtelTrace) Tag(ctx context.Context) SQLCommentTags {
	c := NewCommentCarrier()
	otel.GetTextMapPropagator().Inject(ctx, c)

	return SQLCommentTags(c)
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
