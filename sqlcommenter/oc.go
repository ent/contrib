package sqlcommenter

import (
	"context"

	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

const (
	traceparentHeader = "traceparent"
	tracestateHeader  = "tracestate"
)

// OCTagger is a Tagger that adds `traceparent` and `tracestate` tags to the SQL comment.
type OCTagger struct {
	format *tracecontext.HTTPFormat
}

func NewOCTagger() OCTagger {
	return OCTagger{&tracecontext.HTTPFormat{}}
}

// Tag finds trace information on the given context and returns SQL tags with trace information.
func (ot OCTagger) Tag(ctx context.Context) Tags {
	spanCtx := trace.FromContext(ctx).SpanContext()
	tp, ts := ot.format.SpanContextToHeaders(spanCtx)
	tags := Tags{
		traceparentHeader: tp,
	}
	if ts != "" {
		tags[tracestateHeader] = ts
	}
	return tags
}
