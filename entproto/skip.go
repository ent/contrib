package entproto

import (
	"entgo.io/ent/schema"
)

const SkipAnnotation = "ProtoSkip"

type skipped struct{}

// Skip annotates an ent.Schema to specify that this field will be skipped during .proto generation.
func Skip() schema.Annotation {
	return skipped{}
}

func (f skipped) Name() string {
	return SkipAnnotation
}
