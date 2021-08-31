package sqlcommenter

import (
	"context"
	"fmt"
	"runtime/debug"
)

type DriverVersionTagger struct {
	version string
}

func NewDriverVersionTagger() DriverVersionTagger {
	info, ok := debug.ReadBuildInfo()
	ver := "ent"
	if ok {
		ver = fmt.Sprintf("ent:%s", info.Main.Version)
	}
	return DriverVersionTagger{ver}
}

func (dc DriverVersionTagger) Tag(ctx context.Context) SQLCommentTags {
	return SQLCommentTags{
		DbDriverTagKey: dc.version,
	}
}

type ContextMapper struct {
	contextKey interface{}
	tagKey     string
}

func NewContextMapper(tagKey string, contextKey interface{}) ContextMapper {
	return ContextMapper{
		tagKey:     tagKey,
		contextKey: contextKey,
	}
}

func (cm ContextMapper) Tag(ctx context.Context) SQLCommentTags {
	switch v := ctx.Value(cm.contextKey).(type) {
	case string:
		return SQLCommentTags{cm.tagKey: v}
	default:
		return nil
	}
}

type StaticTagger struct {
	tags SQLCommentTags
}

func NewStaticTagger(tags SQLCommentTags) StaticTagger {
	return StaticTagger{tags}
}

func (sc StaticTagger) Tag(ctx context.Context) SQLCommentTags {
	return sc.tags
}
