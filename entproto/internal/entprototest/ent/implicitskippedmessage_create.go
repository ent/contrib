// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/bionicstork/bionicstork/pkg/entproto/internal/entprototest/ent/implicitskippedmessage"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ImplicitSkippedMessageCreate is the builder for creating a ImplicitSkippedMessage entity.
type ImplicitSkippedMessageCreate struct {
	config
	mutation *ImplicitSkippedMessageMutation
	hooks    []Hook
}

// Mutation returns the ImplicitSkippedMessageMutation object of the builder.
func (ismc *ImplicitSkippedMessageCreate) Mutation() *ImplicitSkippedMessageMutation {
	return ismc.mutation
}

// Save creates the ImplicitSkippedMessage in the database.
func (ismc *ImplicitSkippedMessageCreate) Save(ctx context.Context) (*ImplicitSkippedMessage, error) {
	var (
		err  error
		node *ImplicitSkippedMessage
	)
	if len(ismc.hooks) == 0 {
		if err = ismc.check(); err != nil {
			return nil, err
		}
		node, err = ismc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ImplicitSkippedMessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = ismc.check(); err != nil {
				return nil, err
			}
			ismc.mutation = mutation
			if node, err = ismc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(ismc.hooks) - 1; i >= 0; i-- {
			if ismc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = ismc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ismc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (ismc *ImplicitSkippedMessageCreate) SaveX(ctx context.Context) *ImplicitSkippedMessage {
	v, err := ismc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ismc *ImplicitSkippedMessageCreate) Exec(ctx context.Context) error {
	_, err := ismc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ismc *ImplicitSkippedMessageCreate) ExecX(ctx context.Context) {
	if err := ismc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ismc *ImplicitSkippedMessageCreate) check() error {
	return nil
}

func (ismc *ImplicitSkippedMessageCreate) sqlSave(ctx context.Context) (*ImplicitSkippedMessage, error) {
	_node, _spec := ismc.createSpec()
	if err := sqlgraph.CreateNode(ctx, ismc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (ismc *ImplicitSkippedMessageCreate) createSpec() (*ImplicitSkippedMessage, *sqlgraph.CreateSpec) {
	var (
		_node = &ImplicitSkippedMessage{config: ismc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: implicitskippedmessage.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: implicitskippedmessage.FieldID,
			},
		}
	)
	return _node, _spec
}

// ImplicitSkippedMessageCreateBulk is the builder for creating many ImplicitSkippedMessage entities in bulk.
type ImplicitSkippedMessageCreateBulk struct {
	config
	builders []*ImplicitSkippedMessageCreate
}

// Save creates the ImplicitSkippedMessage entities in the database.
func (ismcb *ImplicitSkippedMessageCreateBulk) Save(ctx context.Context) ([]*ImplicitSkippedMessage, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ismcb.builders))
	nodes := make([]*ImplicitSkippedMessage, len(ismcb.builders))
	mutators := make([]Mutator, len(ismcb.builders))
	for i := range ismcb.builders {
		func(i int, root context.Context) {
			builder := ismcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ImplicitSkippedMessageMutation)
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
					_, err = mutators[i+1].Mutate(root, ismcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ismcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ismcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ismcb *ImplicitSkippedMessageCreateBulk) SaveX(ctx context.Context) []*ImplicitSkippedMessage {
	v, err := ismcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ismcb *ImplicitSkippedMessageCreateBulk) Exec(ctx context.Context) error {
	_, err := ismcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ismcb *ImplicitSkippedMessageCreateBulk) ExecX(ctx context.Context) {
	if err := ismcb.Exec(ctx); err != nil {
		panic(err)
	}
}
