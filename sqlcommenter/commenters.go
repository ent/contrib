package sqlcommenter

import (
	"context"
	"fmt"
	"runtime/debug"
)

type DriverVersionCommenter struct {
	version string
}

func NewDriverVersionCommenter() DriverVersionCommenter {
	info, ok := debug.ReadBuildInfo()
	ver := "ent"
	if ok {
		ver = fmt.Sprintf("ent:%s", info.Main.Version)
	}
	return DriverVersionCommenter{ver}
}

func (dc DriverVersionCommenter) Comments(ctx context.Context) SQLComments {
	return SQLComments{
		DbDriverCommentKey: dc.version,
	}
}

type ContextMapper struct {
	contextKey interface{}
	commentKey string
}

func NewContextMapper(commentKey string, contextKey interface{}) ContextMapper {
	return ContextMapper{
		commentKey: commentKey,
		contextKey: contextKey,
	}
}

func (cm ContextMapper) Comments(ctx context.Context) SQLComments {
	switch v := ctx.Value(cm.contextKey).(type) {
	case string:
		return SQLComments{cm.commentKey: v}
	default:
		return nil
	}
}

type StaticCommenter struct {
	comments SQLComments
}

func NewStaticCommenter(comments SQLComments) StaticCommenter {
	return StaticCommenter{comments}
}

func (sc StaticCommenter) Comments(ctx context.Context) SQLComments {
	return sc.comments
}
