// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/bionicstork/contrib/entproto/internal/entprototest/ent/onemethodservice"
	"github.com/bionicstork/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// OneMethodServiceUpdate is the builder for updating OneMethodService entities.
type OneMethodServiceUpdate struct {
	config
	hooks    []Hook
	mutation *OneMethodServiceMutation
}

// Where appends a list predicates to the OneMethodServiceUpdate builder.
func (omsu *OneMethodServiceUpdate) Where(ps ...predicate.OneMethodService) *OneMethodServiceUpdate {
	omsu.mutation.Where(ps...)
	return omsu
}

// Mutation returns the OneMethodServiceMutation object of the builder.
func (omsu *OneMethodServiceUpdate) Mutation() *OneMethodServiceMutation {
	return omsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (omsu *OneMethodServiceUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(omsu.hooks) == 0 {
		affected, err = omsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*OneMethodServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			omsu.mutation = mutation
			affected, err = omsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(omsu.hooks) - 1; i >= 0; i-- {
			if omsu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = omsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, omsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (omsu *OneMethodServiceUpdate) SaveX(ctx context.Context) int {
	affected, err := omsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (omsu *OneMethodServiceUpdate) Exec(ctx context.Context) error {
	_, err := omsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (omsu *OneMethodServiceUpdate) ExecX(ctx context.Context) {
	if err := omsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (omsu *OneMethodServiceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   onemethodservice.Table,
			Columns: onemethodservice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: onemethodservice.FieldID,
			},
		},
	}
	if ps := omsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if n, err = sqlgraph.UpdateNodes(ctx, omsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{onemethodservice.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// OneMethodServiceUpdateOne is the builder for updating a single OneMethodService entity.
type OneMethodServiceUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *OneMethodServiceMutation
}

// Mutation returns the OneMethodServiceMutation object of the builder.
func (omsuo *OneMethodServiceUpdateOne) Mutation() *OneMethodServiceMutation {
	return omsuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (omsuo *OneMethodServiceUpdateOne) Select(field string, fields ...string) *OneMethodServiceUpdateOne {
	omsuo.fields = append([]string{field}, fields...)
	return omsuo
}

// Save executes the query and returns the updated OneMethodService entity.
func (omsuo *OneMethodServiceUpdateOne) Save(ctx context.Context) (*OneMethodService, error) {
	var (
		err  error
		node *OneMethodService
	)
	if len(omsuo.hooks) == 0 {
		node, err = omsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*OneMethodServiceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			omsuo.mutation = mutation
			node, err = omsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(omsuo.hooks) - 1; i >= 0; i-- {
			if omsuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = omsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, omsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (omsuo *OneMethodServiceUpdateOne) SaveX(ctx context.Context) *OneMethodService {
	node, err := omsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (omsuo *OneMethodServiceUpdateOne) Exec(ctx context.Context) error {
	_, err := omsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (omsuo *OneMethodServiceUpdateOne) ExecX(ctx context.Context) {
	if err := omsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (omsuo *OneMethodServiceUpdateOne) sqlSave(ctx context.Context) (_node *OneMethodService, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   onemethodservice.Table,
			Columns: onemethodservice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: onemethodservice.FieldID,
			},
		},
	}
	id, ok := omsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "OneMethodService.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := omsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, onemethodservice.FieldID)
		for _, f := range fields {
			if !onemethodservice.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != onemethodservice.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := omsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	_node = &OneMethodService{config: omsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, omsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{onemethodservice.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
