package sqlcommenter

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const (
	KeyDBDriver    = "db_driver"
	KeyFramework   = "framework"
	KeyAppliaction = "application"
	KeyRoute       = "route"
	KeyController  = "controller"
	KeyAction      = "action"
)

type Tags map[string]string

func encodeValue(v string) string {
	urlEscape := strings.ReplaceAll(url.PathEscape(string(v)), "+", "%20")
	return fmt.Sprintf("'%s'", urlEscape)
}

func encodeKey(k string) string {
	return url.QueryEscape(string(k))
}

func (sc Tags) Marshal() string {
	kv := make([]struct{ k, v string }, 0, len(sc))
	for k := range sc {
		kv = append(kv, struct{ k, v string }{encodeKey(k), encodeValue(sc[k])})
	}
	sort.Slice(kv, func(i, j int) bool {
		return kv[i].k < kv[j].k
	})
	var b strings.Builder
	for i, p := range kv {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%s=%s", p.k, p.v)
	}
	return b.String()
}

func (sc Tags) Merge(tags ...Tags) Tags {
	for _, c := range tags {
		for k, v := range c {
			sc[k] = v
		}
	}
	return sc
}
