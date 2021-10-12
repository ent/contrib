package sqlcommenter

import (
	"context"
	"fmt"
	"runtime/debug"
)

// driverVersionTagger adds `db_driver` tag with "ent:<version>"
type driverVersionTagger struct {
	version string
}

func NewDriverVersionTagger() driverVersionTagger {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return driverVersionTagger{"ent"}
	}
	for _, d := range info.Deps {
		if d.Path == "entgo.io/ent" {
			return driverVersionTagger{fmt.Sprintf("ent:%s", d.Version)}
		}
	}
	return driverVersionTagger{"ent"}
}

func (dv driverVersionTagger) Tag(ctx context.Context) Tags {
	return Tags{
		KeyDBDriver: dv.version,
	}
}

// ContextMapper is a Tagger that maps context values to tags
// for example, if you want to add 'route' tag to your SQL comment, put the url path on request context:
//  type routeKey struct{}
//  middleware := func(next http.Handler) http.Handler {
//  	fn := func(w http.ResponseWriter, r *http.Request) {
//  		c := context.WithValue(r.Context(), routeKey{}, r.URL.Path)
//  		next.ServeHTTP(w, r.WithContext(c))
//  	}
//  	return http.HandlerFunc(fn)
//  }
// and use ContextMapper to map that route to SQL tag, in your sqlcommenter init code:
//  sqc.NewDriver(drv, sqc.WithTagger(sqc.NewContextMapper("route", routeKey{})))
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

func (st StaticTagger) Tag(ctx context.Context) Tags {
	return st.tags
}
