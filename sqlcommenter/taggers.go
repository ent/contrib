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
	if !ok {
		return DriverVersionTagger{"ent"}
	}
	for _, d := range info.Deps {
		if d.Path == "entgo.io/ent" {
			return DriverVersionTagger{fmt.Sprintf("ent:%s", d.Version)}
		}
	}
	return DriverVersionTagger{"ent"}
}

func (dc DriverVersionTagger) Tag(ctx context.Context) Tags {
	return Tags{
		KeyDBDriver: dc.version,
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

func (cm ContextMapper) Tag(ctx context.Context) Tags {
	switch v := ctx.Value(cm.contextKey).(type) {
	case string:
		return Tags{cm.tagKey: v}
	default:
		return nil
	}
}

type StaticTagger struct {
	tags Tags
}

func NewStaticTagger(tags Tags) StaticTagger {
	return StaticTagger{tags}
}

func (sc StaticTagger) Tag(ctx context.Context) Tags {
	return sc.tags
}
