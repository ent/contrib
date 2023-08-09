// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithoptionals"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithOptionalsQuery is the builder for querying MessageWithOptionals entities.
type MessageWithOptionalsQuery struct {
	config
	ctx        *QueryContext
	order      []messagewithoptionals.OrderOption
	inters     []Interceptor
	predicates []predicate.MessageWithOptionals
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MessageWithOptionalsQuery builder.
func (mwoq *MessageWithOptionalsQuery) Where(ps ...predicate.MessageWithOptionals) *MessageWithOptionalsQuery {
	mwoq.predicates = append(mwoq.predicates, ps...)
	return mwoq
}

// Limit the number of records to be returned by this query.
func (mwoq *MessageWithOptionalsQuery) Limit(limit int) *MessageWithOptionalsQuery {
	mwoq.ctx.Limit = &limit
	return mwoq
}

// Offset to start from.
func (mwoq *MessageWithOptionalsQuery) Offset(offset int) *MessageWithOptionalsQuery {
	mwoq.ctx.Offset = &offset
	return mwoq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mwoq *MessageWithOptionalsQuery) Unique(unique bool) *MessageWithOptionalsQuery {
	mwoq.ctx.Unique = &unique
	return mwoq
}

// Order specifies how the records should be ordered.
func (mwoq *MessageWithOptionalsQuery) Order(o ...messagewithoptionals.OrderOption) *MessageWithOptionalsQuery {
	mwoq.order = append(mwoq.order, o...)
	return mwoq
}

// First returns the first MessageWithOptionals entity from the query.
// Returns a *NotFoundError when no MessageWithOptionals was found.
func (mwoq *MessageWithOptionalsQuery) First(ctx context.Context) (*MessageWithOptionals, error) {
	nodes, err := mwoq.Limit(1).All(setContextOp(ctx, mwoq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{messagewithoptionals.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) FirstX(ctx context.Context) *MessageWithOptionals {
	node, err := mwoq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first MessageWithOptionals ID from the query.
// Returns a *NotFoundError when no MessageWithOptionals ID was found.
func (mwoq *MessageWithOptionalsQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mwoq.Limit(1).IDs(setContextOp(ctx, mwoq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{messagewithoptionals.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) FirstIDX(ctx context.Context) int {
	id, err := mwoq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single MessageWithOptionals entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one MessageWithOptionals entity is found.
// Returns a *NotFoundError when no MessageWithOptionals entities are found.
func (mwoq *MessageWithOptionalsQuery) Only(ctx context.Context) (*MessageWithOptionals, error) {
	nodes, err := mwoq.Limit(2).All(setContextOp(ctx, mwoq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{messagewithoptionals.Label}
	default:
		return nil, &NotSingularError{messagewithoptionals.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) OnlyX(ctx context.Context) *MessageWithOptionals {
	node, err := mwoq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only MessageWithOptionals ID in the query.
// Returns a *NotSingularError when more than one MessageWithOptionals ID is found.
// Returns a *NotFoundError when no entities are found.
func (mwoq *MessageWithOptionalsQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mwoq.Limit(2).IDs(setContextOp(ctx, mwoq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{messagewithoptionals.Label}
	default:
		err = &NotSingularError{messagewithoptionals.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) OnlyIDX(ctx context.Context) int {
	id, err := mwoq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of MessageWithOptionalsSlice.
func (mwoq *MessageWithOptionalsQuery) All(ctx context.Context) ([]*MessageWithOptionals, error) {
	ctx = setContextOp(ctx, mwoq.ctx, "All")
	if err := mwoq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*MessageWithOptionals, *MessageWithOptionalsQuery]()
	return withInterceptors[[]*MessageWithOptionals](ctx, mwoq, qr, mwoq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) AllX(ctx context.Context) []*MessageWithOptionals {
	nodes, err := mwoq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of MessageWithOptionals IDs.
func (mwoq *MessageWithOptionalsQuery) IDs(ctx context.Context) (ids []int, err error) {
	if mwoq.ctx.Unique == nil && mwoq.path != nil {
		mwoq.Unique(true)
	}
	ctx = setContextOp(ctx, mwoq.ctx, "IDs")
	if err = mwoq.Select(messagewithoptionals.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) IDsX(ctx context.Context) []int {
	ids, err := mwoq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mwoq *MessageWithOptionalsQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mwoq.ctx, "Count")
	if err := mwoq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mwoq, querierCount[*MessageWithOptionalsQuery](), mwoq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) CountX(ctx context.Context) int {
	count, err := mwoq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mwoq *MessageWithOptionalsQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mwoq.ctx, "Exist")
	switch _, err := mwoq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mwoq *MessageWithOptionalsQuery) ExistX(ctx context.Context) bool {
	exist, err := mwoq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MessageWithOptionalsQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mwoq *MessageWithOptionalsQuery) Clone() *MessageWithOptionalsQuery {
	if mwoq == nil {
		return nil
	}
	return &MessageWithOptionalsQuery{
		config:     mwoq.config,
		ctx:        mwoq.ctx.Clone(),
		order:      append([]messagewithoptionals.OrderOption{}, mwoq.order...),
		inters:     append([]Interceptor{}, mwoq.inters...),
		predicates: append([]predicate.MessageWithOptionals{}, mwoq.predicates...),
		// clone intermediate query.
		sql:  mwoq.sql.Clone(),
		path: mwoq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		StrOptional string `json:"str_optional,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.MessageWithOptionals.Query().
//		GroupBy(messagewithoptionals.FieldStrOptional).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mwoq *MessageWithOptionalsQuery) GroupBy(field string, fields ...string) *MessageWithOptionalsGroupBy {
	mwoq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &MessageWithOptionalsGroupBy{build: mwoq}
	grbuild.flds = &mwoq.ctx.Fields
	grbuild.label = messagewithoptionals.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		StrOptional string `json:"str_optional,omitempty"`
//	}
//
//	client.MessageWithOptionals.Query().
//		Select(messagewithoptionals.FieldStrOptional).
//		Scan(ctx, &v)
func (mwoq *MessageWithOptionalsQuery) Select(fields ...string) *MessageWithOptionalsSelect {
	mwoq.ctx.Fields = append(mwoq.ctx.Fields, fields...)
	sbuild := &MessageWithOptionalsSelect{MessageWithOptionalsQuery: mwoq}
	sbuild.label = messagewithoptionals.Label
	sbuild.flds, sbuild.scan = &mwoq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a MessageWithOptionalsSelect configured with the given aggregations.
func (mwoq *MessageWithOptionalsQuery) Aggregate(fns ...AggregateFunc) *MessageWithOptionalsSelect {
	return mwoq.Select().Aggregate(fns...)
}

func (mwoq *MessageWithOptionalsQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mwoq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mwoq); err != nil {
				return err
			}
		}
	}
	for _, f := range mwoq.ctx.Fields {
		if !messagewithoptionals.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mwoq.path != nil {
		prev, err := mwoq.path(ctx)
		if err != nil {
			return err
		}
		mwoq.sql = prev
	}
	return nil
}

func (mwoq *MessageWithOptionalsQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*MessageWithOptionals, error) {
	var (
		nodes = []*MessageWithOptionals{}
		_spec = mwoq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*MessageWithOptionals).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &MessageWithOptionals{config: mwoq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mwoq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (mwoq *MessageWithOptionalsQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mwoq.querySpec()
	_spec.Node.Columns = mwoq.ctx.Fields
	if len(mwoq.ctx.Fields) > 0 {
		_spec.Unique = mwoq.ctx.Unique != nil && *mwoq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mwoq.driver, _spec)
}

func (mwoq *MessageWithOptionalsQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(messagewithoptionals.Table, messagewithoptionals.Columns, sqlgraph.NewFieldSpec(messagewithoptionals.FieldID, field.TypeInt))
	_spec.From = mwoq.sql
	if unique := mwoq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mwoq.path != nil {
		_spec.Unique = true
	}
	if fields := mwoq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, messagewithoptionals.FieldID)
		for i := range fields {
			if fields[i] != messagewithoptionals.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mwoq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mwoq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mwoq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mwoq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mwoq *MessageWithOptionalsQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mwoq.driver.Dialect())
	t1 := builder.Table(messagewithoptionals.Table)
	columns := mwoq.ctx.Fields
	if len(columns) == 0 {
		columns = messagewithoptionals.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mwoq.sql != nil {
		selector = mwoq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mwoq.ctx.Unique != nil && *mwoq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range mwoq.predicates {
		p(selector)
	}
	for _, p := range mwoq.order {
		p(selector)
	}
	if offset := mwoq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mwoq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MessageWithOptionalsGroupBy is the group-by builder for MessageWithOptionals entities.
type MessageWithOptionalsGroupBy struct {
	selector
	build *MessageWithOptionalsQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mwogb *MessageWithOptionalsGroupBy) Aggregate(fns ...AggregateFunc) *MessageWithOptionalsGroupBy {
	mwogb.fns = append(mwogb.fns, fns...)
	return mwogb
}

// Scan applies the selector query and scans the result into the given value.
func (mwogb *MessageWithOptionalsGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mwogb.build.ctx, "GroupBy")
	if err := mwogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MessageWithOptionalsQuery, *MessageWithOptionalsGroupBy](ctx, mwogb.build, mwogb, mwogb.build.inters, v)
}

func (mwogb *MessageWithOptionalsGroupBy) sqlScan(ctx context.Context, root *MessageWithOptionalsQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mwogb.fns))
	for _, fn := range mwogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mwogb.flds)+len(mwogb.fns))
		for _, f := range *mwogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mwogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mwogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// MessageWithOptionalsSelect is the builder for selecting fields of MessageWithOptionals entities.
type MessageWithOptionalsSelect struct {
	*MessageWithOptionalsQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (mwos *MessageWithOptionalsSelect) Aggregate(fns ...AggregateFunc) *MessageWithOptionalsSelect {
	mwos.fns = append(mwos.fns, fns...)
	return mwos
}

// Scan applies the selector query and scans the result into the given value.
func (mwos *MessageWithOptionalsSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mwos.ctx, "Select")
	if err := mwos.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MessageWithOptionalsQuery, *MessageWithOptionalsSelect](ctx, mwos.MessageWithOptionalsQuery, mwos, mwos.inters, v)
}

func (mwos *MessageWithOptionalsSelect) sqlScan(ctx context.Context, root *MessageWithOptionalsQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(mwos.fns))
	for _, fn := range mwos.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*mwos.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mwos.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
