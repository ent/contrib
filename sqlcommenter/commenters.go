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

func (dc DriverVersionCommenter) GetComments(ctx context.Context) SqlComments {
	return SqlComments{
		DbDriverCommentKey: CommentValue(dc.version),
	}
}

type ContextMapper struct {
	contextKey interface{}
	commentKey CommentKey
}

func NewContextMapper(commentKey CommentKey, contextKey interface{}) ContextMapper {
	return ContextMapper{
		commentKey: commentKey,
		contextKey: contextKey,
	}
}

func (cm ContextMapper) GetComments(ctx context.Context) SqlComments {
	switch v := ctx.Value(cm.contextKey).(type) {
	case string:
		return SqlComments{cm.commentKey: CommentValue(v)}
	case CommentValue:
		return SqlComments{cm.commentKey: v}
	default:
		return nil
	}
}

type StaticCommenter struct {
	comments SqlComments
}

func NewStaticCommenter(comments SqlComments) StaticCommenter {
	return StaticCommenter{comments}
}

func (sc StaticCommenter) GetComments(ctx context.Context) SqlComments {
	return sc.comments
}
