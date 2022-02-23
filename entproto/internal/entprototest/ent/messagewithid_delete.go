// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/messagewithid"
	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithIDDelete is the builder for deleting a MessageWithID entity.
type MessageWithIDDelete struct {
	config
	hooks    []Hook
	mutation *MessageWithIDMutation
}

// Where appends a list predicates to the MessageWithIDDelete builder.
func (mwid *MessageWithIDDelete) Where(ps ...predicate.MessageWithID) *MessageWithIDDelete {
	mwid.mutation.Where(ps...)
	return mwid
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mwid *MessageWithIDDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(mwid.hooks) == 0 {
		affected, err = mwid.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageWithIDMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mwid.mutation = mutation
			affected, err = mwid.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(mwid.hooks) - 1; i >= 0; i-- {
			if mwid.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mwid.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mwid.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwid *MessageWithIDDelete) ExecX(ctx context.Context) int {
	n, err := mwid.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mwid *MessageWithIDDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: messagewithid.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt32,
				Column: messagewithid.FieldID,
			},
		},
	}
	if ps := mwid.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, mwid.driver, _spec)
}

// MessageWithIDDeleteOne is the builder for deleting a single MessageWithID entity.
type MessageWithIDDeleteOne struct {
	mwid *MessageWithIDDelete
}

// Exec executes the deletion query.
func (mwido *MessageWithIDDeleteOne) Exec(ctx context.Context) error {
	n, err := mwido.mwid.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{messagewithid.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mwido *MessageWithIDDeleteOne) ExecX(ctx context.Context) {
	mwido.mwid.ExecX(ctx)
}
