// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/explicitskippedmessage"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ExplicitSkippedMessageQuery is the builder for querying ExplicitSkippedMessage entities.
type ExplicitSkippedMessageQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.ExplicitSkippedMessage
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ExplicitSkippedMessageQuery builder.
func (esmq *ExplicitSkippedMessageQuery) Where(ps ...predicate.ExplicitSkippedMessage) *ExplicitSkippedMessageQuery {
	esmq.predicates = append(esmq.predicates, ps...)
	return esmq
}

// Limit adds a limit step to the query.
func (esmq *ExplicitSkippedMessageQuery) Limit(limit int) *ExplicitSkippedMessageQuery {
	esmq.limit = &limit
	return esmq
}

// Offset adds an offset step to the query.
func (esmq *ExplicitSkippedMessageQuery) Offset(offset int) *ExplicitSkippedMessageQuery {
	esmq.offset = &offset
	return esmq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (esmq *ExplicitSkippedMessageQuery) Unique(unique bool) *ExplicitSkippedMessageQuery {
	esmq.unique = &unique
	return esmq
}

// Order adds an order step to the query.
func (esmq *ExplicitSkippedMessageQuery) Order(o ...OrderFunc) *ExplicitSkippedMessageQuery {
	esmq.order = append(esmq.order, o...)
	return esmq
}

// First returns the first ExplicitSkippedMessage entity from the query.
// Returns a *NotFoundError when no ExplicitSkippedMessage was found.
func (esmq *ExplicitSkippedMessageQuery) First(ctx context.Context) (*ExplicitSkippedMessage, error) {
	nodes, err := esmq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{explicitskippedmessage.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) FirstX(ctx context.Context) *ExplicitSkippedMessage {
	node, err := esmq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first ExplicitSkippedMessage ID from the query.
// Returns a *NotFoundError when no ExplicitSkippedMessage ID was found.
func (esmq *ExplicitSkippedMessageQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = esmq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{explicitskippedmessage.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) FirstIDX(ctx context.Context) int {
	id, err := esmq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single ExplicitSkippedMessage entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ExplicitSkippedMessage entity is found.
// Returns a *NotFoundError when no ExplicitSkippedMessage entities are found.
func (esmq *ExplicitSkippedMessageQuery) Only(ctx context.Context) (*ExplicitSkippedMessage, error) {
	nodes, err := esmq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{explicitskippedmessage.Label}
	default:
		return nil, &NotSingularError{explicitskippedmessage.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) OnlyX(ctx context.Context) *ExplicitSkippedMessage {
	node, err := esmq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only ExplicitSkippedMessage ID in the query.
// Returns a *NotSingularError when more than one ExplicitSkippedMessage ID is found.
// Returns a *NotFoundError when no entities are found.
func (esmq *ExplicitSkippedMessageQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = esmq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{explicitskippedmessage.Label}
	default:
		err = &NotSingularError{explicitskippedmessage.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) OnlyIDX(ctx context.Context) int {
	id, err := esmq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of ExplicitSkippedMessages.
func (esmq *ExplicitSkippedMessageQuery) All(ctx context.Context) ([]*ExplicitSkippedMessage, error) {
	if err := esmq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return esmq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) AllX(ctx context.Context) []*ExplicitSkippedMessage {
	nodes, err := esmq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of ExplicitSkippedMessage IDs.
func (esmq *ExplicitSkippedMessageQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := esmq.Select(explicitskippedmessage.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) IDsX(ctx context.Context) []int {
	ids, err := esmq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (esmq *ExplicitSkippedMessageQuery) Count(ctx context.Context) (int, error) {
	if err := esmq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return esmq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) CountX(ctx context.Context) int {
	count, err := esmq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (esmq *ExplicitSkippedMessageQuery) Exist(ctx context.Context) (bool, error) {
	if err := esmq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return esmq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (esmq *ExplicitSkippedMessageQuery) ExistX(ctx context.Context) bool {
	exist, err := esmq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ExplicitSkippedMessageQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (esmq *ExplicitSkippedMessageQuery) Clone() *ExplicitSkippedMessageQuery {
	if esmq == nil {
		return nil
	}
	return &ExplicitSkippedMessageQuery{
		config:     esmq.config,
		limit:      esmq.limit,
		offset:     esmq.offset,
		order:      append([]OrderFunc{}, esmq.order...),
		predicates: append([]predicate.ExplicitSkippedMessage{}, esmq.predicates...),
		// clone intermediate query.
		sql:    esmq.sql.Clone(),
		path:   esmq.path,
		unique: esmq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (esmq *ExplicitSkippedMessageQuery) GroupBy(field string, fields ...string) *ExplicitSkippedMessageGroupBy {
	grbuild := &ExplicitSkippedMessageGroupBy{config: esmq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := esmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return esmq.sqlQuery(ctx), nil
	}
	grbuild.label = explicitskippedmessage.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (esmq *ExplicitSkippedMessageQuery) Select(fields ...string) *ExplicitSkippedMessageSelect {
	esmq.fields = append(esmq.fields, fields...)
	selbuild := &ExplicitSkippedMessageSelect{ExplicitSkippedMessageQuery: esmq}
	selbuild.label = explicitskippedmessage.Label
	selbuild.flds, selbuild.scan = &esmq.fields, selbuild.Scan
	return selbuild
}

// Aggregate returns a ExplicitSkippedMessageSelect configured with the given aggregations.
func (esmq *ExplicitSkippedMessageQuery) Aggregate(fns ...AggregateFunc) *ExplicitSkippedMessageSelect {
	return esmq.Select().Aggregate(fns...)
}

func (esmq *ExplicitSkippedMessageQuery) prepareQuery(ctx context.Context) error {
	for _, f := range esmq.fields {
		if !explicitskippedmessage.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if esmq.path != nil {
		prev, err := esmq.path(ctx)
		if err != nil {
			return err
		}
		esmq.sql = prev
	}
	return nil
}

func (esmq *ExplicitSkippedMessageQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ExplicitSkippedMessage, error) {
	var (
		nodes = []*ExplicitSkippedMessage{}
		_spec = esmq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ExplicitSkippedMessage).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ExplicitSkippedMessage{config: esmq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, esmq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (esmq *ExplicitSkippedMessageQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := esmq.querySpec()
	_spec.Node.Columns = esmq.fields
	if len(esmq.fields) > 0 {
		_spec.Unique = esmq.unique != nil && *esmq.unique
	}
	return sqlgraph.CountNodes(ctx, esmq.driver, _spec)
}

func (esmq *ExplicitSkippedMessageQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := esmq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (esmq *ExplicitSkippedMessageQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   explicitskippedmessage.Table,
			Columns: explicitskippedmessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: explicitskippedmessage.FieldID,
			},
		},
		From:   esmq.sql,
		Unique: true,
	}
	if unique := esmq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := esmq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, explicitskippedmessage.FieldID)
		for i := range fields {
			if fields[i] != explicitskippedmessage.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := esmq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := esmq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := esmq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := esmq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (esmq *ExplicitSkippedMessageQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(esmq.driver.Dialect())
	t1 := builder.Table(explicitskippedmessage.Table)
	columns := esmq.fields
	if len(columns) == 0 {
		columns = explicitskippedmessage.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if esmq.sql != nil {
		selector = esmq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if esmq.unique != nil && *esmq.unique {
		selector.Distinct()
	}
	for _, p := range esmq.predicates {
		p(selector)
	}
	for _, p := range esmq.order {
		p(selector)
	}
	if offset := esmq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := esmq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ExplicitSkippedMessageGroupBy is the group-by builder for ExplicitSkippedMessage entities.
type ExplicitSkippedMessageGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (esmgb *ExplicitSkippedMessageGroupBy) Aggregate(fns ...AggregateFunc) *ExplicitSkippedMessageGroupBy {
	esmgb.fns = append(esmgb.fns, fns...)
	return esmgb
}

// Scan applies the group-by query and scans the result into the given value.
func (esmgb *ExplicitSkippedMessageGroupBy) Scan(ctx context.Context, v any) error {
	query, err := esmgb.path(ctx)
	if err != nil {
		return err
	}
	esmgb.sql = query
	return esmgb.sqlScan(ctx, v)
}

func (esmgb *ExplicitSkippedMessageGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range esmgb.fields {
		if !explicitskippedmessage.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := esmgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := esmgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (esmgb *ExplicitSkippedMessageGroupBy) sqlQuery() *sql.Selector {
	selector := esmgb.sql.Select()
	aggregation := make([]string, 0, len(esmgb.fns))
	for _, fn := range esmgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(esmgb.fields)+len(esmgb.fns))
		for _, f := range esmgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(esmgb.fields...)...)
}

// ExplicitSkippedMessageSelect is the builder for selecting fields of ExplicitSkippedMessage entities.
type ExplicitSkippedMessageSelect struct {
	*ExplicitSkippedMessageQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (esms *ExplicitSkippedMessageSelect) Aggregate(fns ...AggregateFunc) *ExplicitSkippedMessageSelect {
	esms.fns = append(esms.fns, fns...)
	return esms
}

// Scan applies the selector query and scans the result into the given value.
func (esms *ExplicitSkippedMessageSelect) Scan(ctx context.Context, v any) error {
	if err := esms.prepareQuery(ctx); err != nil {
		return err
	}
	esms.sql = esms.ExplicitSkippedMessageQuery.sqlQuery(ctx)
	return esms.sqlScan(ctx, v)
}

func (esms *ExplicitSkippedMessageSelect) sqlScan(ctx context.Context, v any) error {
	aggregation := make([]string, 0, len(esms.fns))
	for _, fn := range esms.fns {
		aggregation = append(aggregation, fn(esms.sql))
	}
	switch n := len(*esms.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		esms.sql.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		esms.sql.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := esms.sql.Query()
	if err := esms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
