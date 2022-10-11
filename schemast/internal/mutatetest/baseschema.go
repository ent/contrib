package mutatetest

import (
	"entgo.io/contrib/schemast/internal/mutatetest/mixin"
	"entgo.io/contrib/schemast/internal/mutatetest/policy"
	"entgo.io/ent"
)

// TestBaseSchema is a test base schema
type TestBaseSchema struct {
	ent.Schema
}

func (TestBaseSchema) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.UUID{},
	}
}

func (TestBaseSchema) Policy() ent.Policy {
	return policy.IDFilterRule()
}
