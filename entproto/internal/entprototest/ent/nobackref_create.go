// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/contrib/entproto/internal/entprototest/ent/image"
	"entgo.io/contrib/entproto/internal/entprototest/ent/nobackref"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// NoBackrefCreate is the builder for creating a NoBackref entity.
type NoBackrefCreate struct {
	config
	mutation *NoBackrefMutation
	hooks    []Hook
}

// AddImageIDs adds the "images" edge to the Image entity by IDs.
func (nbc *NoBackrefCreate) AddImageIDs(ids ...uuid.UUID) *NoBackrefCreate {
	nbc.mutation.AddImageIDs(ids...)
	return nbc
}

// AddImages adds the "images" edges to the Image entity.
func (nbc *NoBackrefCreate) AddImages(i ...*Image) *NoBackrefCreate {
	ids := make([]uuid.UUID, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return nbc.AddImageIDs(ids...)
}

// Mutation returns the NoBackrefMutation object of the builder.
func (nbc *NoBackrefCreate) Mutation() *NoBackrefMutation {
	return nbc.mutation
}

// Save creates the NoBackref in the database.
func (nbc *NoBackrefCreate) Save(ctx context.Context) (*NoBackref, error) {
	return withHooks[*NoBackref, NoBackrefMutation](ctx, nbc.sqlSave, nbc.mutation, nbc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (nbc *NoBackrefCreate) SaveX(ctx context.Context) *NoBackref {
	v, err := nbc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nbc *NoBackrefCreate) Exec(ctx context.Context) error {
	_, err := nbc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nbc *NoBackrefCreate) ExecX(ctx context.Context) {
	if err := nbc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nbc *NoBackrefCreate) check() error {
	return nil
}

func (nbc *NoBackrefCreate) sqlSave(ctx context.Context) (*NoBackref, error) {
	if err := nbc.check(); err != nil {
		return nil, err
	}
	_node, _spec := nbc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nbc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	nbc.mutation.id = &_node.ID
	nbc.mutation.done = true
	return _node, nil
}

func (nbc *NoBackrefCreate) createSpec() (*NoBackref, *sqlgraph.CreateSpec) {
	var (
		_node = &NoBackref{config: nbc.config}
		_spec = sqlgraph.NewCreateSpec(nobackref.Table, sqlgraph.NewFieldSpec(nobackref.FieldID, field.TypeInt))
	)
	if nodes := nbc.mutation.ImagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   nobackref.ImagesTable,
			Columns: []string{nobackref.ImagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(image.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// NoBackrefCreateBulk is the builder for creating many NoBackref entities in bulk.
type NoBackrefCreateBulk struct {
	config
	builders []*NoBackrefCreate
}

// Save creates the NoBackref entities in the database.
func (nbcb *NoBackrefCreateBulk) Save(ctx context.Context) ([]*NoBackref, error) {
	specs := make([]*sqlgraph.CreateSpec, len(nbcb.builders))
	nodes := make([]*NoBackref, len(nbcb.builders))
	mutators := make([]Mutator, len(nbcb.builders))
	for i := range nbcb.builders {
		func(i int, root context.Context) {
			builder := nbcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NoBackrefMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, nbcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, nbcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, nbcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (nbcb *NoBackrefCreateBulk) SaveX(ctx context.Context) []*NoBackref {
	v, err := nbcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (nbcb *NoBackrefCreateBulk) Exec(ctx context.Context) error {
	_, err := nbcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nbcb *NoBackrefCreateBulk) ExecX(ctx context.Context) {
	if err := nbcb.Exec(ctx); err != nil {
		panic(err)
	}
}
