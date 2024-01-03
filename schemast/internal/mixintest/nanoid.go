package mixintest

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type mockPolicy struct {
}

func (m mockPolicy) EvalMutation(ctx context.Context, mutation ent.Mutation) error {
	return nil
}

func (m mockPolicy) EvalQuery(ctx context.Context, query ent.Query) error {
	return nil
}

// UUID defines the id field as a UUID string.
type UUID struct {
	mixin.Schema
}

// Fields of the UUID.
func (UUID) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(uuid.NewString).
			NotEmpty().
			Immutable().
			Unique(),
	}
}

// Policy of the UUID.
func (UUID) Policy() ent.Policy {
	return mockPolicy{}
}
