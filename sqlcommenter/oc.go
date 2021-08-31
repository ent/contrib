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

type OCTrace struct {
	format *tracecontext.HTTPFormat
}

func NewOCTrace() OCTrace {
	return OCTrace{&tracecontext.HTTPFormat{}}
}

func (oc OCTrace) Tag(ctx context.Context) SQLCommentTags {
	spanCtx := trace.FromContext(ctx).SpanContext()
	tp, ts := oc.format.SpanContextToHeaders(spanCtx)
	cmts := SQLCommentTags{
		traceparentHeader: tp,
	}
	if ts != "" {
		cmts[tracestateHeader] = ts
	}
	return cmts
}