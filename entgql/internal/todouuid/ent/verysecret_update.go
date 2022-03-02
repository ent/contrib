// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/predicate"
	"github.com/bionicstork/contrib/entgql/internal/todouuid/ent/verysecret"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VerySecretUpdate is the builder for updating VerySecret entities.
type VerySecretUpdate struct {
	config
	hooks    []Hook
	mutation *VerySecretMutation
}

// Where appends a list predicates to the VerySecretUpdate builder.
func (vsu *VerySecretUpdate) Where(ps ...predicate.VerySecret) *VerySecretUpdate {
	vsu.mutation.Where(ps...)
	return vsu
}

// SetPassword sets the "password" field.
func (vsu *VerySecretUpdate) SetPassword(s string) *VerySecretUpdate {
	vsu.mutation.SetPassword(s)
	return vsu
}

// Mutation returns the VerySecretMutation object of the builder.
func (vsu *VerySecretUpdate) Mutation() *VerySecretMutation {
	return vsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (vsu *VerySecretUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(vsu.hooks) == 0 {
		affected, err = vsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VerySecretMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			vsu.mutation = mutation
			affected, err = vsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(vsu.hooks) - 1; i >= 0; i-- {
			if vsu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = vsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (vsu *VerySecretUpdate) SaveX(ctx context.Context) int {
	affected, err := vsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (vsu *VerySecretUpdate) Exec(ctx context.Context) error {
	_, err := vsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsu *VerySecretUpdate) ExecX(ctx context.Context) {
	if err := vsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (vsu *VerySecretUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   verysecret.Table,
			Columns: verysecret.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: verysecret.FieldID,
			},
		},
	}
	if ps := vsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vsu.mutation.Password(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: verysecret.FieldPassword,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, vsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{verysecret.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// VerySecretUpdateOne is the builder for updating a single VerySecret entity.
type VerySecretUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *VerySecretMutation
}

// SetPassword sets the "password" field.
func (vsuo *VerySecretUpdateOne) SetPassword(s string) *VerySecretUpdateOne {
	vsuo.mutation.SetPassword(s)
	return vsuo
}

// Mutation returns the VerySecretMutation object of the builder.
func (vsuo *VerySecretUpdateOne) Mutation() *VerySecretMutation {
	return vsuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (vsuo *VerySecretUpdateOne) Select(field string, fields ...string) *VerySecretUpdateOne {
	vsuo.fields = append([]string{field}, fields...)
	return vsuo
}

// Save executes the query and returns the updated VerySecret entity.
func (vsuo *VerySecretUpdateOne) Save(ctx context.Context) (*VerySecret, error) {
	var (
		err  error
		node *VerySecret
	)
	if len(vsuo.hooks) == 0 {
		node, err = vsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VerySecretMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			vsuo.mutation = mutation
			node, err = vsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(vsuo.hooks) - 1; i >= 0; i-- {
			if vsuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = vsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (vsuo *VerySecretUpdateOne) SaveX(ctx context.Context) *VerySecret {
	node, err := vsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (vsuo *VerySecretUpdateOne) Exec(ctx context.Context) error {
	_, err := vsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vsuo *VerySecretUpdateOne) ExecX(ctx context.Context) {
	if err := vsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (vsuo *VerySecretUpdateOne) sqlSave(ctx context.Context) (_node *VerySecret, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   verysecret.Table,
			Columns: verysecret.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: verysecret.FieldID,
			},
		},
	}
	id, ok := vsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "VerySecret.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := vsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, verysecret.FieldID)
		for _, f := range fields {
			if !verysecret.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != verysecret.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := vsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := vsuo.mutation.Password(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: verysecret.FieldPassword,
		})
	}
	_node = &VerySecret{config: vsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, vsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{verysecret.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
