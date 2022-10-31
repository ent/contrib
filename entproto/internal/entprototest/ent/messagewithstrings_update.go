// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithstrings"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
)

// MessageWithStringsUpdate is the builder for updating MessageWithStrings entities.
type MessageWithStringsUpdate struct {
	config
	hooks    []Hook
	mutation *MessageWithStringsMutation
}

// Where appends a list predicates to the MessageWithStringsUpdate builder.
func (mwsu *MessageWithStringsUpdate) Where(ps ...predicate.MessageWithStrings) *MessageWithStringsUpdate {
	mwsu.mutation.Where(ps...)
	return mwsu
}

// SetStrings sets the "strings" field.
func (mwsu *MessageWithStringsUpdate) SetStrings(s []string) *MessageWithStringsUpdate {
	mwsu.mutation.SetStrings(s)
	return mwsu
}

// AppendStrings appends s to the "strings" field.
func (mwsu *MessageWithStringsUpdate) AppendStrings(s []string) *MessageWithStringsUpdate {
	mwsu.mutation.AppendStrings(s)
	return mwsu
}

// Mutation returns the MessageWithStringsMutation object of the builder.
func (mwsu *MessageWithStringsUpdate) Mutation() *MessageWithStringsMutation {
	return mwsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mwsu *MessageWithStringsUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(mwsu.hooks) == 0 {
		affected, err = mwsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageWithStringsMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mwsu.mutation = mutation
			affected, err = mwsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(mwsu.hooks) - 1; i >= 0; i-- {
			if mwsu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mwsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, mwsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (mwsu *MessageWithStringsUpdate) SaveX(ctx context.Context) int {
	affected, err := mwsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mwsu *MessageWithStringsUpdate) Exec(ctx context.Context) error {
	_, err := mwsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwsu *MessageWithStringsUpdate) ExecX(ctx context.Context) {
	if err := mwsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mwsu *MessageWithStringsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   messagewithstrings.Table,
			Columns: messagewithstrings.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithstrings.FieldID,
			},
		},
	}
	if ps := mwsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mwsu.mutation.Strings(); ok {
		_spec.SetField(messagewithstrings.FieldStrings, field.TypeJSON, value)
	}
	if value, ok := mwsu.mutation.AppendedStrings(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, messagewithstrings.FieldStrings, value)
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mwsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{messagewithstrings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// MessageWithStringsUpdateOne is the builder for updating a single MessageWithStrings entity.
type MessageWithStringsUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MessageWithStringsMutation
}

// SetStrings sets the "strings" field.
func (mwsuo *MessageWithStringsUpdateOne) SetStrings(s []string) *MessageWithStringsUpdateOne {
	mwsuo.mutation.SetStrings(s)
	return mwsuo
}

// AppendStrings appends s to the "strings" field.
func (mwsuo *MessageWithStringsUpdateOne) AppendStrings(s []string) *MessageWithStringsUpdateOne {
	mwsuo.mutation.AppendStrings(s)
	return mwsuo
}

// Mutation returns the MessageWithStringsMutation object of the builder.
func (mwsuo *MessageWithStringsUpdateOne) Mutation() *MessageWithStringsMutation {
	return mwsuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (mwsuo *MessageWithStringsUpdateOne) Select(field string, fields ...string) *MessageWithStringsUpdateOne {
	mwsuo.fields = append([]string{field}, fields...)
	return mwsuo
}

// Save executes the query and returns the updated MessageWithStrings entity.
func (mwsuo *MessageWithStringsUpdateOne) Save(ctx context.Context) (*MessageWithStrings, error) {
	var (
		err  error
		node *MessageWithStrings
	)
	if len(mwsuo.hooks) == 0 {
		node, err = mwsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageWithStringsMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mwsuo.mutation = mutation
			node, err = mwsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(mwsuo.hooks) - 1; i >= 0; i-- {
			if mwsuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mwsuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, mwsuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*MessageWithStrings)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from MessageWithStringsMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (mwsuo *MessageWithStringsUpdateOne) SaveX(ctx context.Context) *MessageWithStrings {
	node, err := mwsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (mwsuo *MessageWithStringsUpdateOne) Exec(ctx context.Context) error {
	_, err := mwsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwsuo *MessageWithStringsUpdateOne) ExecX(ctx context.Context) {
	if err := mwsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (mwsuo *MessageWithStringsUpdateOne) sqlSave(ctx context.Context) (_node *MessageWithStrings, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   messagewithstrings.Table,
			Columns: messagewithstrings.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithstrings.FieldID,
			},
		},
	}
	id, ok := mwsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "MessageWithStrings.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := mwsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, messagewithstrings.FieldID)
		for _, f := range fields {
			if !messagewithstrings.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != messagewithstrings.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := mwsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := mwsuo.mutation.Strings(); ok {
		_spec.SetField(messagewithstrings.FieldStrings, field.TypeJSON, value)
	}
	if value, ok := mwsuo.mutation.AppendedStrings(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, messagewithstrings.FieldStrings, value)
		})
	}
	_node = &MessageWithStrings{config: mwsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, mwsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{messagewithstrings.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
