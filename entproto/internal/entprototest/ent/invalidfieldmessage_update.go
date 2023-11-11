// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent/invalidfieldmessage"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/contrib/entproto/internal/entprototest/ent/schema"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// InvalidFieldMessageUpdate is the builder for updating InvalidFieldMessage entities.
type InvalidFieldMessageUpdate struct {
	config
	hooks    []Hook
	mutation *InvalidFieldMessageMutation
}

// Where appends a list predicates to the InvalidFieldMessageUpdate builder.
func (ifmu *InvalidFieldMessageUpdate) Where(ps ...predicate.InvalidFieldMessage) *InvalidFieldMessageUpdate {
	ifmu.mutation.Where(ps...)
	return ifmu
}

// SetJSON sets the "json" field.
func (ifmu *InvalidFieldMessageUpdate) SetJSON(sj *schema.SomeJSON) *InvalidFieldMessageUpdate {
	ifmu.mutation.SetJSON(sj)
	return ifmu
}

// Mutation returns the InvalidFieldMessageMutation object of the builder.
func (ifmu *InvalidFieldMessageUpdate) Mutation() *InvalidFieldMessageMutation {
	return ifmu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ifmu *InvalidFieldMessageUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ifmu.sqlSave, ifmu.mutation, ifmu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ifmu *InvalidFieldMessageUpdate) SaveX(ctx context.Context) int {
	affected, err := ifmu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ifmu *InvalidFieldMessageUpdate) Exec(ctx context.Context) error {
	_, err := ifmu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ifmu *InvalidFieldMessageUpdate) ExecX(ctx context.Context) {
	if err := ifmu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ifmu *InvalidFieldMessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(invalidfieldmessage.Table, invalidfieldmessage.Columns, sqlgraph.NewFieldSpec(invalidfieldmessage.FieldID, field.TypeInt))
	if ps := ifmu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ifmu.mutation.JSON(); ok {
		_spec.SetField(invalidfieldmessage.FieldJSON, field.TypeJSON, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ifmu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{invalidfieldmessage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ifmu.mutation.done = true
	return n, nil
}

// InvalidFieldMessageUpdateOne is the builder for updating a single InvalidFieldMessage entity.
type InvalidFieldMessageUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *InvalidFieldMessageMutation
}

// SetJSON sets the "json" field.
func (ifmuo *InvalidFieldMessageUpdateOne) SetJSON(sj *schema.SomeJSON) *InvalidFieldMessageUpdateOne {
	ifmuo.mutation.SetJSON(sj)
	return ifmuo
}

// Mutation returns the InvalidFieldMessageMutation object of the builder.
func (ifmuo *InvalidFieldMessageUpdateOne) Mutation() *InvalidFieldMessageMutation {
	return ifmuo.mutation
}

// Where appends a list predicates to the InvalidFieldMessageUpdate builder.
func (ifmuo *InvalidFieldMessageUpdateOne) Where(ps ...predicate.InvalidFieldMessage) *InvalidFieldMessageUpdateOne {
	ifmuo.mutation.Where(ps...)
	return ifmuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ifmuo *InvalidFieldMessageUpdateOne) Select(field string, fields ...string) *InvalidFieldMessageUpdateOne {
	ifmuo.fields = append([]string{field}, fields...)
	return ifmuo
}

// Save executes the query and returns the updated InvalidFieldMessage entity.
func (ifmuo *InvalidFieldMessageUpdateOne) Save(ctx context.Context) (*InvalidFieldMessage, error) {
	return withHooks(ctx, ifmuo.sqlSave, ifmuo.mutation, ifmuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ifmuo *InvalidFieldMessageUpdateOne) SaveX(ctx context.Context) *InvalidFieldMessage {
	node, err := ifmuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ifmuo *InvalidFieldMessageUpdateOne) Exec(ctx context.Context) error {
	_, err := ifmuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ifmuo *InvalidFieldMessageUpdateOne) ExecX(ctx context.Context) {
	if err := ifmuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ifmuo *InvalidFieldMessageUpdateOne) sqlSave(ctx context.Context) (_node *InvalidFieldMessage, err error) {
	_spec := sqlgraph.NewUpdateSpec(invalidfieldmessage.Table, invalidfieldmessage.Columns, sqlgraph.NewFieldSpec(invalidfieldmessage.FieldID, field.TypeInt))
	id, ok := ifmuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "InvalidFieldMessage.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ifmuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, invalidfieldmessage.FieldID)
		for _, f := range fields {
			if !invalidfieldmessage.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != invalidfieldmessage.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ifmuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ifmuo.mutation.JSON(); ok {
		_spec.SetField(invalidfieldmessage.FieldJSON, field.TypeJSON, value)
	}
	_node = &InvalidFieldMessage{config: ifmuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ifmuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{invalidfieldmessage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ifmuo.mutation.done = true
	return _node, nil
}
