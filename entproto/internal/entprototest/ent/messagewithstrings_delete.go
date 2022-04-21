// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithstrings"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithStringsDelete is the builder for deleting a MessageWithStrings entity.
type MessageWithStringsDelete struct {
	config
	hooks    []Hook
	mutation *MessageWithStringsMutation
}

// Where appends a list predicates to the MessageWithStringsDelete builder.
func (mwsd *MessageWithStringsDelete) Where(ps ...predicate.MessageWithStrings) *MessageWithStringsDelete {
	mwsd.mutation.Where(ps...)
	return mwsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mwsd *MessageWithStringsDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(mwsd.hooks) == 0 {
		affected, err = mwsd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageWithStringsMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mwsd.mutation = mutation
			affected, err = mwsd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(mwsd.hooks) - 1; i >= 0; i-- {
			if mwsd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mwsd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mwsd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwsd *MessageWithStringsDelete) ExecX(ctx context.Context) int {
	n, err := mwsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mwsd *MessageWithStringsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: messagewithstrings.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithstrings.FieldID,
			},
		},
	}
	if ps := mwsd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, mwsd.driver, _spec)
}

// MessageWithStringsDeleteOne is the builder for deleting a single MessageWithStrings entity.
type MessageWithStringsDeleteOne struct {
	mwsd *MessageWithStringsDelete
}

// Exec executes the deletion query.
func (mwsdo *MessageWithStringsDeleteOne) Exec(ctx context.Context) error {
	n, err := mwsdo.mwsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{messagewithstrings.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mwsdo *MessageWithStringsDeleteOne) ExecX(ctx context.Context) {
	mwsdo.mwsd.ExecX(ctx)
}
