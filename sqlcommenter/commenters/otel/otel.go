package otel

import (
	"context"

	sqc "entgo.io/contrib/sqlcommenter"
	"go.opentelemetry.io/otel"
)

type (
	OtelTraceCommenter struct{}
	CommentCarrier     sqc.SqlComments
)

func NewOtelTraceCommenter() OtelTraceCommenter {
	return OtelTraceCommenter{}
}

func (hc OtelTraceCommenter) GetComments(ctx context.Context) sqc.SqlComments {
	comments := NewCommentCarrier()
	otel.GetTextMapPropagator().Inject(ctx, comments)

	return sqc.SqlComments(comments)
}

func NewCommentCarrier() CommentCarrier {
	return make(CommentCarrier)
}

// Get returns the value associated with the passed key.
func (c CommentCarrier) Get(key string) string {
	return string(c[sqc.CommentKey(key)])
}

// Set stores the key-value pair.
func (c CommentCarrier) Set(key string, value string) {
	c[sqc.CommentKey(key)] = sqc.CommentValue(value)
}

// Keys lists the keys stored in this carrier.
func (c CommentCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, string(k))
	}
	return keys
}
