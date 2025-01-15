// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"math"

	"entgo.io/contrib/entgql/internal/todo/ent/predicate"
	"entgo.io/contrib/entgql/internal/todo/ent/workspace"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// WorkspaceQuery is the builder for querying Workspace entities.
type WorkspaceQuery struct {
	config
	ctx        *QueryContext
	order      []workspace.OrderOption
	inters     []Interceptor
	predicates []predicate.Workspace
	loadTotal  []func(context.Context, []*Workspace) error
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the WorkspaceQuery builder.
func (wq *WorkspaceQuery) Where(ps ...predicate.Workspace) *WorkspaceQuery {
	wq.predicates = append(wq.predicates, ps...)
	return wq
}

// Limit the number of records to be returned by this query.
func (wq *WorkspaceQuery) Limit(limit int) *WorkspaceQuery {
	wq.ctx.Limit = &limit
	return wq
}

// Offset to start from.
func (wq *WorkspaceQuery) Offset(offset int) *WorkspaceQuery {
	wq.ctx.Offset = &offset
	return wq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (wq *WorkspaceQuery) Unique(unique bool) *WorkspaceQuery {
	wq.ctx.Unique = &unique
	return wq
}

// Order specifies how the records should be ordered.
func (wq *WorkspaceQuery) Order(o ...workspace.OrderOption) *WorkspaceQuery {
	wq.order = append(wq.order, o...)
	return wq
}

// First returns the first Workspace entity from the query.
// Returns a *NotFoundError when no Workspace was found.
func (wq *WorkspaceQuery) First(ctx context.Context) (*Workspace, error) {
	nodes, err := wq.Limit(1).All(setContextOp(ctx, wq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{workspace.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wq *WorkspaceQuery) FirstX(ctx context.Context) *Workspace {
	node, err := wq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Workspace ID from the query.
// Returns a *NotFoundError when no Workspace ID was found.
func (wq *WorkspaceQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wq.Limit(1).IDs(setContextOp(ctx, wq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{workspace.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (wq *WorkspaceQuery) FirstIDX(ctx context.Context) int {
	id, err := wq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Workspace entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Workspace entity is found.
// Returns a *NotFoundError when no Workspace entities are found.
func (wq *WorkspaceQuery) Only(ctx context.Context) (*Workspace, error) {
	nodes, err := wq.Limit(2).All(setContextOp(ctx, wq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{workspace.Label}
	default:
		return nil, &NotSingularError{workspace.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wq *WorkspaceQuery) OnlyX(ctx context.Context) *Workspace {
	node, err := wq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Workspace ID in the query.
// Returns a *NotSingularError when more than one Workspace ID is found.
// Returns a *NotFoundError when no entities are found.
func (wq *WorkspaceQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = wq.Limit(2).IDs(setContextOp(ctx, wq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{workspace.Label}
	default:
		err = &NotSingularError{workspace.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (wq *WorkspaceQuery) OnlyIDX(ctx context.Context) int {
	id, err := wq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Workspaces.
func (wq *WorkspaceQuery) All(ctx context.Context) ([]*Workspace, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryAll)
	if err := wq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Workspace, *WorkspaceQuery]()
	return withInterceptors[[]*Workspace](ctx, wq, qr, wq.inters)
}

// AllX is like All, but panics if an error occurs.
func (wq *WorkspaceQuery) AllX(ctx context.Context) []*Workspace {
	nodes, err := wq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Workspace IDs.
func (wq *WorkspaceQuery) IDs(ctx context.Context) (ids []int, err error) {
	if wq.ctx.Unique == nil && wq.path != nil {
		wq.Unique(true)
	}
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryIDs)
	if err = wq.Select(workspace.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wq *WorkspaceQuery) IDsX(ctx context.Context) []int {
	ids, err := wq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wq *WorkspaceQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryCount)
	if err := wq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, wq, querierCount[*WorkspaceQuery](), wq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (wq *WorkspaceQuery) CountX(ctx context.Context) int {
	count, err := wq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wq *WorkspaceQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, wq.ctx, ent.OpQueryExist)
	switch _, err := wq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (wq *WorkspaceQuery) ExistX(ctx context.Context) bool {
	exist, err := wq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the WorkspaceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wq *WorkspaceQuery) Clone() *WorkspaceQuery {
	if wq == nil {
		return nil
	}
	return &WorkspaceQuery{
		config:     wq.config,
		ctx:        wq.ctx.Clone(),
		order:      append([]workspace.OrderOption{}, wq.order...),
		inters:     append([]Interceptor{}, wq.inters...),
		predicates: append([]predicate.Workspace{}, wq.predicates...),
		// clone intermediate query.
		sql:       wq.sql.Clone(),
		path:      wq.path,
		modifiers: append([]func(*sql.Selector){}, wq.modifiers...),
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Workspace.Query().
//		GroupBy(workspace.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (wq *WorkspaceQuery) GroupBy(field string, fields ...string) *WorkspaceGroupBy {
	wq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &WorkspaceGroupBy{build: wq}
	grbuild.flds = &wq.ctx.Fields
	grbuild.label = workspace.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.Workspace.Query().
//		Select(workspace.FieldName).
//		Scan(ctx, &v)
func (wq *WorkspaceQuery) Select(fields ...string) *WorkspaceSelect {
	wq.ctx.Fields = append(wq.ctx.Fields, fields...)
	sbuild := &WorkspaceSelect{WorkspaceQuery: wq}
	sbuild.label = workspace.Label
	sbuild.flds, sbuild.scan = &wq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a WorkspaceSelect configured with the given aggregations.
func (wq *WorkspaceQuery) Aggregate(fns ...AggregateFunc) *WorkspaceSelect {
	return wq.Select().Aggregate(fns...)
}

func (wq *WorkspaceQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range wq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, wq); err != nil {
				return err
			}
		}
	}
	for _, f := range wq.ctx.Fields {
		if !workspace.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if wq.path != nil {
		prev, err := wq.path(ctx)
		if err != nil {
			return err
		}
		wq.sql = prev
	}
	return nil
}

func (wq *WorkspaceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Workspace, error) {
	var (
		nodes = []*Workspace{}
		_spec = wq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Workspace).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Workspace{config: wq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(wq.modifiers) > 0 {
		_spec.Modifiers = wq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, wq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	for i := range wq.loadTotal {
		if err := wq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (wq *WorkspaceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wq.querySpec()
	if len(wq.modifiers) > 0 {
		_spec.Modifiers = wq.modifiers
	}
	_spec.Node.Columns = wq.ctx.Fields
	if len(wq.ctx.Fields) > 0 {
		_spec.Unique = wq.ctx.Unique != nil && *wq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, wq.driver, _spec)
}

func (wq *WorkspaceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(workspace.Table, workspace.Columns, sqlgraph.NewFieldSpec(workspace.FieldID, field.TypeInt))
	_spec.From = wq.sql
	if unique := wq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if wq.path != nil {
		_spec.Unique = true
	}
	if fields := wq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, workspace.FieldID)
		for i := range fields {
			if fields[i] != workspace.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := wq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wq *WorkspaceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(wq.driver.Dialect())
	t1 := builder.Table(workspace.Table)
	columns := wq.ctx.Fields
	if len(columns) == 0 {
		columns = workspace.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if wq.sql != nil {
		selector = wq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if wq.ctx.Unique != nil && *wq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range wq.modifiers {
		m(selector)
	}
	for _, p := range wq.predicates {
		p(selector)
	}
	for _, p := range wq.order {
		p(selector)
	}
	if offset := wq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (wq *WorkspaceQuery) Modify(modifiers ...func(s *sql.Selector)) *WorkspaceSelect {
	wq.modifiers = append(wq.modifiers, modifiers...)
	return wq.Select()
}

// WorkspaceGroupBy is the group-by builder for Workspace entities.
type WorkspaceGroupBy struct {
	selector
	build *WorkspaceQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wgb *WorkspaceGroupBy) Aggregate(fns ...AggregateFunc) *WorkspaceGroupBy {
	wgb.fns = append(wgb.fns, fns...)
	return wgb
}

// Scan applies the selector query and scans the result into the given value.
func (wgb *WorkspaceGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, wgb.build.ctx, ent.OpQueryGroupBy)
	if err := wgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WorkspaceQuery, *WorkspaceGroupBy](ctx, wgb.build, wgb, wgb.build.inters, v)
}

func (wgb *WorkspaceGroupBy) sqlScan(ctx context.Context, root *WorkspaceQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(wgb.fns))
	for _, fn := range wgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*wgb.flds)+len(wgb.fns))
		for _, f := range *wgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*wgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// WorkspaceSelect is the builder for selecting fields of Workspace entities.
type WorkspaceSelect struct {
	*WorkspaceQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ws *WorkspaceSelect) Aggregate(fns ...AggregateFunc) *WorkspaceSelect {
	ws.fns = append(ws.fns, fns...)
	return ws
}

// Scan applies the selector query and scans the result into the given value.
func (ws *WorkspaceSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ws.ctx, ent.OpQuerySelect)
	if err := ws.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*WorkspaceQuery, *WorkspaceSelect](ctx, ws.WorkspaceQuery, ws, ws.inters, v)
}

func (ws *WorkspaceSelect) sqlScan(ctx context.Context, root *WorkspaceQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ws.fns))
	for _, fn := range ws.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ws.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ws.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (ws *WorkspaceSelect) Modify(modifiers ...func(s *sql.Selector)) *WorkspaceSelect {
	ws.modifiers = append(ws.modifiers, modifiers...)
	return ws
}
