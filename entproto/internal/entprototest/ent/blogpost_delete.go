// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/blogpost"
	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// BlogPostDelete is the builder for deleting a BlogPost entity.
type BlogPostDelete struct {
	config
	hooks    []Hook
	mutation *BlogPostMutation
}

// Where appends a list predicates to the BlogPostDelete builder.
func (bpd *BlogPostDelete) Where(ps ...predicate.BlogPost) *BlogPostDelete {
	bpd.mutation.Where(ps...)
	return bpd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (bpd *BlogPostDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(bpd.hooks) == 0 {
		affected, err = bpd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BlogPostMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			bpd.mutation = mutation
			affected, err = bpd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(bpd.hooks) - 1; i >= 0; i-- {
			if bpd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = bpd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, bpd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (bpd *BlogPostDelete) ExecX(ctx context.Context) int {
	n, err := bpd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (bpd *BlogPostDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: blogpost.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: blogpost.FieldID,
			},
		},
	}
	if ps := bpd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, bpd.driver, _spec)
}

// BlogPostDeleteOne is the builder for deleting a single BlogPost entity.
type BlogPostDeleteOne struct {
	bpd *BlogPostDelete
}

// Exec executes the deletion query.
func (bpdo *BlogPostDeleteOne) Exec(ctx context.Context) error {
	n, err := bpdo.bpd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{blogpost.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (bpdo *BlogPostDeleteOne) ExecX(ctx context.Context) {
	bpdo.bpd.ExecX(ctx)
}
