// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/messagewithenum"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithEnumQuery is the builder for querying MessageWithEnum entities.
type MessageWithEnumQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.MessageWithEnum
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MessageWithEnumQuery builder.
func (mweq *MessageWithEnumQuery) Where(ps ...predicate.MessageWithEnum) *MessageWithEnumQuery {
	mweq.predicates = append(mweq.predicates, ps...)
	return mweq
}

// Limit adds a limit step to the query.
func (mweq *MessageWithEnumQuery) Limit(limit int) *MessageWithEnumQuery {
	mweq.limit = &limit
	return mweq
}

// Offset adds an offset step to the query.
func (mweq *MessageWithEnumQuery) Offset(offset int) *MessageWithEnumQuery {
	mweq.offset = &offset
	return mweq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mweq *MessageWithEnumQuery) Unique(unique bool) *MessageWithEnumQuery {
	mweq.unique = &unique
	return mweq
}

// Order adds an order step to the query.
func (mweq *MessageWithEnumQuery) Order(o ...OrderFunc) *MessageWithEnumQuery {
	mweq.order = append(mweq.order, o...)
	return mweq
}

// First returns the first MessageWithEnum entity from the query.
// Returns a *NotFoundError when no MessageWithEnum was found.
func (mweq *MessageWithEnumQuery) First(ctx context.Context) (*MessageWithEnum, error) {
	nodes, err := mweq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{messagewithenum.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) FirstX(ctx context.Context) *MessageWithEnum {
	node, err := mweq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first MessageWithEnum ID from the query.
// Returns a *NotFoundError when no MessageWithEnum ID was found.
func (mweq *MessageWithEnumQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mweq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{messagewithenum.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) FirstIDX(ctx context.Context) int {
	id, err := mweq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single MessageWithEnum entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one MessageWithEnum entity is found.
// Returns a *NotFoundError when no MessageWithEnum entities are found.
func (mweq *MessageWithEnumQuery) Only(ctx context.Context) (*MessageWithEnum, error) {
	nodes, err := mweq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{messagewithenum.Label}
	default:
		return nil, &NotSingularError{messagewithenum.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) OnlyX(ctx context.Context) *MessageWithEnum {
	node, err := mweq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only MessageWithEnum ID in the query.
// Returns a *NotSingularError when more than one MessageWithEnum ID is found.
// Returns a *NotFoundError when no entities are found.
func (mweq *MessageWithEnumQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mweq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{messagewithenum.Label}
	default:
		err = &NotSingularError{messagewithenum.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) OnlyIDX(ctx context.Context) int {
	id, err := mweq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of MessageWithEnums.
func (mweq *MessageWithEnumQuery) All(ctx context.Context) ([]*MessageWithEnum, error) {
	if err := mweq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return mweq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) AllX(ctx context.Context) []*MessageWithEnum {
	nodes, err := mweq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of MessageWithEnum IDs.
func (mweq *MessageWithEnumQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := mweq.Select(messagewithenum.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) IDsX(ctx context.Context) []int {
	ids, err := mweq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mweq *MessageWithEnumQuery) Count(ctx context.Context) (int, error) {
	if err := mweq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return mweq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) CountX(ctx context.Context) int {
	count, err := mweq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mweq *MessageWithEnumQuery) Exist(ctx context.Context) (bool, error) {
	if err := mweq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return mweq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (mweq *MessageWithEnumQuery) ExistX(ctx context.Context) bool {
	exist, err := mweq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MessageWithEnumQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mweq *MessageWithEnumQuery) Clone() *MessageWithEnumQuery {
	if mweq == nil {
		return nil
	}
	return &MessageWithEnumQuery{
		config:     mweq.config,
		limit:      mweq.limit,
		offset:     mweq.offset,
		order:      append([]OrderFunc{}, mweq.order...),
		predicates: append([]predicate.MessageWithEnum{}, mweq.predicates...),
		// clone intermediate query.
		sql:    mweq.sql.Clone(),
		path:   mweq.path,
		unique: mweq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		EnumType messagewithenum.EnumType `json:"enum_type,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.MessageWithEnum.Query().
//		GroupBy(messagewithenum.FieldEnumType).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (mweq *MessageWithEnumQuery) GroupBy(field string, fields ...string) *MessageWithEnumGroupBy {
	grbuild := &MessageWithEnumGroupBy{config: mweq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := mweq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return mweq.sqlQuery(ctx), nil
	}
	grbuild.label = messagewithenum.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		EnumType messagewithenum.EnumType `json:"enum_type,omitempty"`
//	}
//
//	client.MessageWithEnum.Query().
//		Select(messagewithenum.FieldEnumType).
//		Scan(ctx, &v)
//
func (mweq *MessageWithEnumQuery) Select(fields ...string) *MessageWithEnumSelect {
	mweq.fields = append(mweq.fields, fields...)
	selbuild := &MessageWithEnumSelect{MessageWithEnumQuery: mweq}
	selbuild.label = messagewithenum.Label
	selbuild.flds, selbuild.scan = &mweq.fields, selbuild.Scan
	return selbuild
}

// Aggregate returns a MessageWithEnumSelect configured with the given aggregations.
func (mweq *MessageWithEnumQuery) Aggregate(fns ...AggregateFunc) *MessageWithEnumSelect {
	return mweq.Select().Aggregate(fns...)
}

func (mweq *MessageWithEnumQuery) prepareQuery(ctx context.Context) error {
	for _, f := range mweq.fields {
		if !messagewithenum.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mweq.path != nil {
		prev, err := mweq.path(ctx)
		if err != nil {
			return err
		}
		mweq.sql = prev
	}
	return nil
}

func (mweq *MessageWithEnumQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*MessageWithEnum, error) {
	var (
		nodes = []*MessageWithEnum{}
		_spec = mweq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*MessageWithEnum).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &MessageWithEnum{config: mweq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mweq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (mweq *MessageWithEnumQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mweq.querySpec()
	_spec.Node.Columns = mweq.fields
	if len(mweq.fields) > 0 {
		_spec.Unique = mweq.unique != nil && *mweq.unique
	}
	return sqlgraph.CountNodes(ctx, mweq.driver, _spec)
}

func (mweq *MessageWithEnumQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := mweq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (mweq *MessageWithEnumQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   messagewithenum.Table,
			Columns: messagewithenum.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithenum.FieldID,
			},
		},
		From:   mweq.sql,
		Unique: true,
	}
	if unique := mweq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := mweq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, messagewithenum.FieldID)
		for i := range fields {
			if fields[i] != messagewithenum.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mweq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mweq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mweq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mweq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mweq *MessageWithEnumQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mweq.driver.Dialect())
	t1 := builder.Table(messagewithenum.Table)
	columns := mweq.fields
	if len(columns) == 0 {
		columns = messagewithenum.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mweq.sql != nil {
		selector = mweq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mweq.unique != nil && *mweq.unique {
		selector.Distinct()
	}
	for _, p := range mweq.predicates {
		p(selector)
	}
	for _, p := range mweq.order {
		p(selector)
	}
	if offset := mweq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mweq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MessageWithEnumGroupBy is the group-by builder for MessageWithEnum entities.
type MessageWithEnumGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mwegb *MessageWithEnumGroupBy) Aggregate(fns ...AggregateFunc) *MessageWithEnumGroupBy {
	mwegb.fns = append(mwegb.fns, fns...)
	return mwegb
}

// Scan applies the group-by query and scans the result into the given value.
func (mwegb *MessageWithEnumGroupBy) Scan(ctx context.Context, v any) error {
	query, err := mwegb.path(ctx)
	if err != nil {
		return err
	}
	mwegb.sql = query
	return mwegb.sqlScan(ctx, v)
}

func (mwegb *MessageWithEnumGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range mwegb.fields {
		if !messagewithenum.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := mwegb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mwegb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (mwegb *MessageWithEnumGroupBy) sqlQuery() *sql.Selector {
	selector := mwegb.sql.Select()
	aggregation := make([]string, 0, len(mwegb.fns))
	for _, fn := range mwegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(mwegb.fields)+len(mwegb.fns))
		for _, f := range mwegb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(mwegb.fields...)...)
}

// MessageWithEnumSelect is the builder for selecting fields of MessageWithEnum entities.
type MessageWithEnumSelect struct {
	*MessageWithEnumQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (mwes *MessageWithEnumSelect) Aggregate(fns ...AggregateFunc) *MessageWithEnumSelect {
	mwes.fns = append(mwes.fns, fns...)
	return mwes
}

// Scan applies the selector query and scans the result into the given value.
func (mwes *MessageWithEnumSelect) Scan(ctx context.Context, v any) error {
	if err := mwes.prepareQuery(ctx); err != nil {
		return err
	}
	mwes.sql = mwes.MessageWithEnumQuery.sqlQuery(ctx)
	return mwes.sqlScan(ctx, v)
}

func (mwes *MessageWithEnumSelect) sqlScan(ctx context.Context, v any) error {
	aggregation := make([]string, 0, len(mwes.fns))
	for _, fn := range mwes.fns {
		aggregation = append(aggregation, fn(mwes.sql))
	}
	switch n := len(*mwes.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		mwes.sql.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		mwes.sql.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := mwes.sql.Query()
	if err := mwes.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
