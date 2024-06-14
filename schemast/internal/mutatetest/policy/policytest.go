package policy

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
)

type (
	FilterQueryRule func(context.Context, ent.Query) error

	OnlyQuery interface {
		OnlyID(ctx context.Context) (id string, err error)
	}
)

func (f FilterQueryRule) EvalQuery(ctx context.Context, q ent.Query) error {
	return f(ctx, q)
}

func (FilterQueryRule) EvalMutation(context.Context, ent.Mutation) error {
	return nil
}

func IDFilterRule() ent.Policy {
	return FilterQueryRule(func(ctx context.Context, q ent.Query) error {
		IDQuery, ok := q.(OnlyQuery)
		if ok {
			_, err := IDQuery.OnlyID(ctx)
			if err != nil {
				return privacy.Deny
			}
		}

		return privacy.Skip
	})
}
