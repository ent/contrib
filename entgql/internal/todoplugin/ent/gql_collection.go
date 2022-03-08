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
	"database/sql/driver"
	"fmt"

	"entgo.io/contrib/entgql/internal/todoplugin/ent/todo"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (c *CategoryQuery) CollectFields(ctx context.Context, satisfies ...string) (*CategoryQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return c, nil
	}
	if err := c.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CategoryQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	return nil
}

type categoryPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []CategoryPaginateOption
}

func newCategoryPaginateArgs(rv map[string]interface{}) *categoryPaginateArgs {
	args := &categoryPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			var (
				err1, err2 error
				order      = &CategoryOrder{Field: &CategoryOrderField{}}
			)
			if d, ok := v[directionField]; ok {
				err1 = order.Direction.UnmarshalGQL(d)
			}
			if f, ok := v[fieldField]; ok {
				err2 = order.Field.UnmarshalGQL(f)
			}
			if err1 == nil && err2 == nil {
				args.opts = append(args.opts, WithCategoryOrder(order))
			}
		case *CategoryOrder:
			if v != nil {
				args.opts = append(args.opts, WithCategoryOrder(v))
			}
		}
	}
	if v := rv[whereField]; v != nil && v != (*CategoryWhereInput)(nil) {
		args.opts = append(args.opts, WithCategoryFilter(v.(*CategoryWhereInput).Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (t *TodoQuery) CollectFields(ctx context.Context, satisfies ...string) (*TodoQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return t, nil
	}
	if err := t.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TodoQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	for _, field := range graphql.CollectFields(op, field.Selections, satisfies) {
		switch field.Name {
		case "parent":
			var (
				path  = append(path, field.Name)
				query = &TodoQuery{config: t.config}
			)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			t.withParent = query
		case "children":
			var (
				path  = append(path, field.Name)
				query = &TodoQuery{config: t.config}
			)
			args := newTodoPaginateArgs(fieldArgs(ctx, path...))
			if err := validateFirstLast(args.first, args.last); err != nil {
				return fmt.Errorf("validate first and last in path %q: %w", path, err)
			}
			pager, err := newTodoPager(args.opts)
			if err != nil {
				return fmt.Errorf("create new pager in path %q: %w", path, err)
			}
			if query, err = pager.applyFilter(query); err != nil {
				return err
			}
			if f, l := args.first, args.last; !hasCollectedField(ctx, append(path, edgesField)...) || f != nil && *f == 0 || l != nil && *l == 0 {
				if hasCollectedField(ctx, append(path, totalCountField)...) || hasCollectedField(ctx, append(path, pageInfoField)...) {
					query := query.Clone()
					t.loadTotal = append(t.loadTotal, func(ctx context.Context, nodes []*Todo) error {
						ids := make([]driver.Value, len(nodes))
						for i := range nodes {
							ids[i] = nodes[i].ID
						}
						var v []struct {
							NodeID int `sql:"todo_children"`
							Count  int `sql:"count"`
						}
						query.Where(func(s *sql.Selector) {
							s.Where(sql.InValues(todo.ChildrenColumn, ids...))
						})
						if err := query.GroupBy(todo.ChildrenColumn).Aggregate(Count()).Scan(ctx, &v); err != nil {
							return err
						}
						m := make(map[int]int, len(v))
						for i := range v {
							m[v[i].NodeID] = v[i].Count
						}
						for i := range nodes {
							n := m[nodes[i].ID]
							nodes[i].Edges.totalCount[1] = &n
						}
						return nil
					})
				}
				continue
			}
			if (args.after != nil || args.first != nil || args.before != nil || args.last != nil) && hasCollectedField(ctx, append(path, totalCountField)...) {
				query := query.Clone()
				t.loadTotal = append(t.loadTotal, func(ctx context.Context, nodes []*Todo) error {
					ids := make([]driver.Value, len(nodes))
					for i := range nodes {
						ids[i] = nodes[i].ID
					}
					var v []struct {
						NodeID int `sql:"todo_children"`
						Count  int `sql:"count"`
					}
					query.Where(func(s *sql.Selector) {
						s.Where(sql.InValues(todo.ChildrenColumn, ids...))
					})
					if err := query.GroupBy(todo.ChildrenColumn).Aggregate(Count()).Scan(ctx, &v); err != nil {
						return err
					}
					m := make(map[int]int, len(v))
					for i := range v {
						m[v[i].NodeID] = v[i].Count
					}
					for i := range nodes {
						n := m[nodes[i].ID]
						nodes[i].Edges.totalCount[1] = &n
					}
					return nil
				})
			} else {
				t.loadTotal = append(t.loadTotal, func(_ context.Context, nodes []*Todo) error {
					for i := range nodes {
						n := len(nodes[i].Edges.Children)
						nodes[i].Edges.totalCount[1] = &n
					}
					return nil
				})
			}
			query = pager.applyCursors(query, args.after, args.before)
			if limit := paginateLimit(args.first, args.last); limit > 0 {
				modify := limitRows(todo.ChildrenColumn, limit, pager.orderExpr(args.last != nil))
				query.modifiers = append(query.modifiers, modify)
			} else {
				query = pager.applyOrder(query, args.last != nil)
			}
			path = append(path, edgesField, nodeField)
			if err := query.collectField(ctx, op, field, path, satisfies...); err != nil {
				return err
			}
			t.withChildren = query
		}
	}
	return nil
}

type todoPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []TodoPaginateOption
}

func newTodoPaginateArgs(rv map[string]interface{}) *todoPaginateArgs {
	args := &todoPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			var (
				err1, err2 error
				order      = &TodoOrder{Field: &TodoOrderField{}}
			)
			if d, ok := v[directionField]; ok {
				err1 = order.Direction.UnmarshalGQL(d)
			}
			if f, ok := v[fieldField]; ok {
				err2 = order.Field.UnmarshalGQL(f)
			}
			if err1 == nil && err2 == nil {
				args.opts = append(args.opts, WithTodoOrder(order))
			}
		case *TodoOrder:
			if v != nil {
				args.opts = append(args.opts, WithTodoOrder(v))
			}
		}
	}
	if v := rv[whereField]; v != nil && v != (*TodoWhereInput)(nil) {
		args.opts = append(args.opts, WithTodoFilter(v.(*TodoWhereInput).Filter))
	}
	return args
}

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (u *UserQuery) CollectFields(ctx context.Context, satisfies ...string) (*UserQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return u, nil
	}
	if err := u.collectField(ctx, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserQuery) collectField(ctx context.Context, op *graphql.OperationContext, field graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	return nil
}

type masteruserPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []MasterUserPaginateOption
}

func newMasterUserPaginateArgs(rv map[string]interface{}) *masteruserPaginateArgs {
	args := &masteruserPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v := rv[whereField]; v != nil && v != (*MasterUserWhereInput)(nil) {
		args.opts = append(args.opts, WithMasterUserFilter(v.(*MasterUserWhereInput).Filter))
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, path ...string) map[string]interface{} {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	for _, name := range path {
		var field *graphql.CollectedField
		for _, f := range graphql.CollectFields(oc, fc.Field.Selections, nil) {
			if f.Name == name {
				field = &f
				break
			}
		}
		if field == nil {
			return nil
		}
		cf, err := fc.Child(ctx, *field)
		if err != nil {
			args := field.ArgumentMap(oc.Variables)
			return unmarshalArgs(args)
		}
		fc = cf
	}
	return fc.Args
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
// Note that "where" (WhereInput) is not supported.
func unmarshalArgs(args map[string]interface{}) map[string]interface{} {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = &c
		}
	}
	return args
}

func limitRows(partitionBy string, limit int, orderBy ...sql.Querier) func(s *sql.Selector) {
	return func(s *sql.Selector) {
		d := sql.Dialect(s.Dialect())
		s.SetDistinct(false)
		with := d.With("src_query").
			As(s.Clone()).
			With("limited_query").
			As(
				d.Select("*").
					AppendSelectExprAs(
						sql.RowNumber().PartitionBy(partitionBy).OrderExpr(orderBy...),
						"row_number",
					).
					From(d.Table("src_query")),
			)
		t := d.Table("limited_query").As(s.TableName())
		*s = *d.Select(s.UnqualifiedColumns()...).
			From(t).
			Where(sql.LTE(t.C("row_number"), limit)).
			Prefix(with)
	}
}
