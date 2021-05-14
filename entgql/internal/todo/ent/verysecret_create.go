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

	"entgo.io/contrib/entgql/internal/todo/ent/verysecret"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VerySecretCreate is the builder for creating a VerySecret entity.
type VerySecretCreate struct {
	config
	mutation *VerySecretMutation
	hooks    []Hook
}

// SetPassword sets the "password" field.
func (vsc *VerySecretCreate) SetPassword(s string) *VerySecretCreate {
	vsc.mutation.SetPassword(s)
	return vsc
}

// Mutation returns the VerySecretMutation object of the builder.
func (vsc *VerySecretCreate) Mutation() *VerySecretMutation {
	return vsc.mutation
}

// Save creates the VerySecret in the database.
func (vsc *VerySecretCreate) Save(ctx context.Context) (*VerySecret, error) {
	var (
		err  error
		node *VerySecret
	)
	if len(vsc.hooks) == 0 {
		if err = vsc.check(); err != nil {
			return nil, err
		}
		node, err = vsc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*VerySecretMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = vsc.check(); err != nil {
				return nil, err
			}
			vsc.mutation = mutation
			node, err = vsc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(vsc.hooks) - 1; i >= 0; i-- {
			mut = vsc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, vsc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (vsc *VerySecretCreate) SaveX(ctx context.Context) *VerySecret {
	v, err := vsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// check runs all checks and user-defined validators on the builder.
func (vsc *VerySecretCreate) check() error {
	if _, ok := vsc.mutation.Password(); !ok {
		return &ValidationError{Name: "password", err: errors.New("ent: missing required field \"password\"")}
	}
	return nil
}

func (vsc *VerySecretCreate) sqlSave(ctx context.Context) (*VerySecret, error) {
	_node, _spec := vsc.createSpec()
	if err := sqlgraph.CreateNode(ctx, vsc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (vsc *VerySecretCreate) createSpec() (*VerySecret, *sqlgraph.CreateSpec) {
	var (
		_node = &VerySecret{config: vsc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: verysecret.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: verysecret.FieldID,
			},
		}
	)
	if value, ok := vsc.mutation.Password(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: verysecret.FieldPassword,
		})
		_node.Password = value
	}
	return _node, _spec
}

// VerySecretCreateBulk is the builder for creating many VerySecret entities in bulk.
type VerySecretCreateBulk struct {
	config
	builders []*VerySecretCreate
}

// Save creates the VerySecret entities in the database.
func (vscb *VerySecretCreateBulk) Save(ctx context.Context) ([]*VerySecret, error) {
	specs := make([]*sqlgraph.CreateSpec, len(vscb.builders))
	nodes := make([]*VerySecret, len(vscb.builders))
	mutators := make([]Mutator, len(vscb.builders))
	for i := range vscb.builders {
		func(i int, root context.Context) {
			builder := vscb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*VerySecretMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, vscb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, vscb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, vscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (vscb *VerySecretCreateBulk) SaveX(ctx context.Context) []*VerySecret {
	v, err := vscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
