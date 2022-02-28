// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/predicate"
	"github.com/bionicstork/contrib/schemast/internal/mutatetest/ent/withmodifiedfield"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// WithModifiedFieldDelete is the builder for deleting a WithModifiedField entity.
type WithModifiedFieldDelete struct {
	config
	hooks    []Hook
	mutation *WithModifiedFieldMutation
}

// Where appends a list predicates to the WithModifiedFieldDelete builder.
func (wmfd *WithModifiedFieldDelete) Where(ps ...predicate.WithModifiedField) *WithModifiedFieldDelete {
	wmfd.mutation.Where(ps...)
	return wmfd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (wmfd *WithModifiedFieldDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(wmfd.hooks) == 0 {
		affected, err = wmfd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WithModifiedFieldMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wmfd.mutation = mutation
			affected, err = wmfd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(wmfd.hooks) - 1; i >= 0; i-- {
			if wmfd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = wmfd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wmfd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (wmfd *WithModifiedFieldDelete) ExecX(ctx context.Context) int {
	n, err := wmfd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (wmfd *WithModifiedFieldDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: withmodifiedfield.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: withmodifiedfield.FieldID,
			},
		},
	}
	if ps := wmfd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, wmfd.driver, _spec)
}

// WithModifiedFieldDeleteOne is the builder for deleting a single WithModifiedField entity.
type WithModifiedFieldDeleteOne struct {
	wmfd *WithModifiedFieldDelete
}

// Exec executes the deletion query.
func (wmfdo *WithModifiedFieldDeleteOne) Exec(ctx context.Context) error {
	n, err := wmfdo.wmfd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{withmodifiedfield.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (wmfdo *WithModifiedFieldDeleteOne) ExecX(ctx context.Context) {
	wmfdo.wmfd.ExecX(ctx)
}
