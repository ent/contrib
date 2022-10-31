// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithpackagename"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithPackageNameCreate is the builder for creating a MessageWithPackageName entity.
type MessageWithPackageNameCreate struct {
	config
	mutation *MessageWithPackageNameMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (mwpnc *MessageWithPackageNameCreate) SetName(s string) *MessageWithPackageNameCreate {
	mwpnc.mutation.SetName(s)
	return mwpnc
}

// Mutation returns the MessageWithPackageNameMutation object of the builder.
func (mwpnc *MessageWithPackageNameCreate) Mutation() *MessageWithPackageNameMutation {
	return mwpnc.mutation
}

// Save creates the MessageWithPackageName in the database.
func (mwpnc *MessageWithPackageNameCreate) Save(ctx context.Context) (*MessageWithPackageName, error) {
	var (
		err  error
		node *MessageWithPackageName
	)
	if len(mwpnc.hooks) == 0 {
		if err = mwpnc.check(); err != nil {
			return nil, err
		}
		node, err = mwpnc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MessageWithPackageNameMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = mwpnc.check(); err != nil {
				return nil, err
			}
			mwpnc.mutation = mutation
			if node, err = mwpnc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(mwpnc.hooks) - 1; i >= 0; i-- {
			if mwpnc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mwpnc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, mwpnc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*MessageWithPackageName)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from MessageWithPackageNameMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (mwpnc *MessageWithPackageNameCreate) SaveX(ctx context.Context) *MessageWithPackageName {
	v, err := mwpnc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mwpnc *MessageWithPackageNameCreate) Exec(ctx context.Context) error {
	_, err := mwpnc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwpnc *MessageWithPackageNameCreate) ExecX(ctx context.Context) {
	if err := mwpnc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mwpnc *MessageWithPackageNameCreate) check() error {
	if _, ok := mwpnc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "MessageWithPackageName.name"`)}
	}
	return nil
}

func (mwpnc *MessageWithPackageNameCreate) sqlSave(ctx context.Context) (*MessageWithPackageName, error) {
	_node, _spec := mwpnc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mwpnc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (mwpnc *MessageWithPackageNameCreate) createSpec() (*MessageWithPackageName, *sqlgraph.CreateSpec) {
	var (
		_node = &MessageWithPackageName{config: mwpnc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: messagewithpackagename.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithpackagename.FieldID,
			},
		}
	)
	if value, ok := mwpnc.mutation.Name(); ok {
		_spec.SetField(messagewithpackagename.FieldName, field.TypeString, value)
		_node.Name = value
	}
	return _node, _spec
}

// MessageWithPackageNameCreateBulk is the builder for creating many MessageWithPackageName entities in bulk.
type MessageWithPackageNameCreateBulk struct {
	config
	builders []*MessageWithPackageNameCreate
}

// Save creates the MessageWithPackageName entities in the database.
func (mwpncb *MessageWithPackageNameCreateBulk) Save(ctx context.Context) ([]*MessageWithPackageName, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mwpncb.builders))
	nodes := make([]*MessageWithPackageName, len(mwpncb.builders))
	mutators := make([]Mutator, len(mwpncb.builders))
	for i := range mwpncb.builders {
		func(i int, root context.Context) {
			builder := mwpncb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MessageWithPackageNameMutation)
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
					_, err = mutators[i+1].Mutate(root, mwpncb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mwpncb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mwpncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mwpncb *MessageWithPackageNameCreateBulk) SaveX(ctx context.Context) []*MessageWithPackageName {
	v, err := mwpncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mwpncb *MessageWithPackageNameCreateBulk) Exec(ctx context.Context) error {
	_, err := mwpncb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mwpncb *MessageWithPackageNameCreateBulk) ExecX(ctx context.Context) {
	if err := mwpncb.Exec(ctx); err != nil {
		panic(err)
	}
}
