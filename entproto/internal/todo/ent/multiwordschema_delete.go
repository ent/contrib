// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent/multiwordschema"
	"github.com/bionicstork/bionicstork/pkg/entproto/internal/todo/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MultiWordSchemaDelete is the builder for deleting a MultiWordSchema entity.
type MultiWordSchemaDelete struct {
	config
	hooks    []Hook
	mutation *MultiWordSchemaMutation
}

// Where appends a list predicates to the MultiWordSchemaDelete builder.
func (mwsd *MultiWordSchemaDelete) Where(ps ...predicate.MultiWordSchema) *MultiWordSchemaDelete {
	mwsd.mutation.Where(ps...)
	return mwsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mwsd *MultiWordSchemaDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(mwsd.hooks) == 0 {
		affected, err = mwsd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MultiWordSchemaMutation)
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
func (mwsd *MultiWordSchemaDelete) ExecX(ctx context.Context) int {
	n, err := mwsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mwsd *MultiWordSchemaDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: multiwordschema.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: multiwordschema.FieldID,
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

// MultiWordSchemaDeleteOne is the builder for deleting a single MultiWordSchema entity.
type MultiWordSchemaDeleteOne struct {
	mwsd *MultiWordSchemaDelete
}

// Exec executes the deletion query.
func (mwsdo *MultiWordSchemaDeleteOne) Exec(ctx context.Context) error {
	n, err := mwsdo.mwsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{multiwordschema.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mwsdo *MultiWordSchemaDeleteOne) ExecX(ctx context.Context) {
	mwsdo.mwsd.ExecX(ctx)
}
