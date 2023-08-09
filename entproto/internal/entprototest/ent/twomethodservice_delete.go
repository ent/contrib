// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/contrib/entproto/internal/entprototest/ent/twomethodservice"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TwoMethodServiceDelete is the builder for deleting a TwoMethodService entity.
type TwoMethodServiceDelete struct {
	config
	hooks    []Hook
	mutation *TwoMethodServiceMutation
}

// Where appends a list predicates to the TwoMethodServiceDelete builder.
func (tmsd *TwoMethodServiceDelete) Where(ps ...predicate.TwoMethodService) *TwoMethodServiceDelete {
	tmsd.mutation.Where(ps...)
	return tmsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tmsd *TwoMethodServiceDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, TwoMethodServiceMutation](ctx, tmsd.sqlExec, tmsd.mutation, tmsd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tmsd *TwoMethodServiceDelete) ExecX(ctx context.Context) int {
	n, err := tmsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tmsd *TwoMethodServiceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(twomethodservice.Table, sqlgraph.NewFieldSpec(twomethodservice.FieldID, field.TypeInt))
	if ps := tmsd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tmsd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tmsd.mutation.done = true
	return affected, err
}

// TwoMethodServiceDeleteOne is the builder for deleting a single TwoMethodService entity.
type TwoMethodServiceDeleteOne struct {
	tmsd *TwoMethodServiceDelete
}

// Where appends a list predicates to the TwoMethodServiceDelete builder.
func (tmsdo *TwoMethodServiceDeleteOne) Where(ps ...predicate.TwoMethodService) *TwoMethodServiceDeleteOne {
	tmsdo.tmsd.mutation.Where(ps...)
	return tmsdo
}

// Exec executes the deletion query.
func (tmsdo *TwoMethodServiceDeleteOne) Exec(ctx context.Context) error {
	n, err := tmsdo.tmsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{twomethodservice.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tmsdo *TwoMethodServiceDeleteOne) ExecX(ctx context.Context) {
	if err := tmsdo.Exec(ctx); err != nil {
		panic(err)
	}
}
