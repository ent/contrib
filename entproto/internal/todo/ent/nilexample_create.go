// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/contrib/entproto/internal/todo/ent/nilexample"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// NilExampleCreate is the builder for creating a NilExample entity.
type NilExampleCreate struct {
	config
	mutation *NilExampleMutation
	hooks    []Hook
}

// SetStrNil sets the "str_nil" field.
func (nec *NilExampleCreate) SetStrNil(s string) *NilExampleCreate {
	nec.mutation.SetStrNil(s)
	return nec
}

// SetNillableStrNil sets the "str_nil" field if the given value is not nil.
func (nec *NilExampleCreate) SetNillableStrNil(s *string) *NilExampleCreate {
	if s != nil {
		nec.SetStrNil(*s)
	}
	return nec
}

// SetTimeNil sets the "time_nil" field.
func (nec *NilExampleCreate) SetTimeNil(t time.Time) *NilExampleCreate {
	nec.mutation.SetTimeNil(t)
	return nec
}

// SetNillableTimeNil sets the "time_nil" field if the given value is not nil.
func (nec *NilExampleCreate) SetNillableTimeNil(t *time.Time) *NilExampleCreate {
	if t != nil {
		nec.SetTimeNil(*t)
	}
	return nec
}

// Mutation returns the NilExampleMutation object of the builder.
func (nec *NilExampleCreate) Mutation() *NilExampleMutation {
	return nec.mutation
}

// Save creates the NilExample in the database.
func (nec *NilExampleCreate) Save(ctx context.Context) (*NilExample, error) {
	var (
		err  error
		node *NilExample
	)
	if len(nec.hooks) == 0 {
		if err = nec.check(); err != nil {
			return nil, err
		}
		node, err = nec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NilExampleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nec.check(); err != nil {
				return nil, err
			}
			nec.mutation = mutation
			if node, err = nec.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(nec.hooks) - 1; i >= 0; i-- {
			if nec.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = nec.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nec.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (nec *NilExampleCreate) SaveX(ctx context.Context) *NilExample {
	v, err := nec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nec *NilExampleCreate) Exec(ctx context.Context) error {
	_, err := nec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nec *NilExampleCreate) ExecX(ctx context.Context) {
	if err := nec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nec *NilExampleCreate) check() error {
	return nil
}

func (nec *NilExampleCreate) sqlSave(ctx context.Context) (*NilExample, error) {
	_node, _spec := nec.createSpec()
	if err := sqlgraph.CreateNode(ctx, nec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (nec *NilExampleCreate) createSpec() (*NilExample, *sqlgraph.CreateSpec) {
	var (
		_node = &NilExample{config: nec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: nilexample.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: nilexample.FieldID,
			},
		}
	)
	if value, ok := nec.mutation.StrNil(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: nilexample.FieldStrNil,
		})
		_node.StrNil = &value
	}
	if value, ok := nec.mutation.TimeNil(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: nilexample.FieldTimeNil,
		})
		_node.TimeNil = &value
	}
	return _node, _spec
}

// NilExampleCreateBulk is the builder for creating many NilExample entities in bulk.
type NilExampleCreateBulk struct {
	config
	builders []*NilExampleCreate
}

// Save creates the NilExample entities in the database.
func (necb *NilExampleCreateBulk) Save(ctx context.Context) ([]*NilExample, error) {
	specs := make([]*sqlgraph.CreateSpec, len(necb.builders))
	nodes := make([]*NilExample, len(necb.builders))
	mutators := make([]Mutator, len(necb.builders))
	for i := range necb.builders {
		func(i int, root context.Context) {
			builder := necb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NilExampleMutation)
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
					_, err = mutators[i+1].Mutate(root, necb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, necb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, necb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (necb *NilExampleCreateBulk) SaveX(ctx context.Context) []*NilExample {
	v, err := necb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (necb *NilExampleCreateBulk) Exec(ctx context.Context) error {
	_, err := necb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (necb *NilExampleCreateBulk) ExecX(ctx context.Context) {
	if err := necb.Exec(ctx); err != nil {
		panic(err)
	}
}
