// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/allmethodsservice"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// AllMethodsServiceQuery is the builder for querying AllMethodsService entities.
type AllMethodsServiceQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.AllMethodsService
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AllMethodsServiceQuery builder.
func (amsq *AllMethodsServiceQuery) Where(ps ...predicate.AllMethodsService) *AllMethodsServiceQuery {
	amsq.predicates = append(amsq.predicates, ps...)
	return amsq
}

// Limit adds a limit step to the query.
func (amsq *AllMethodsServiceQuery) Limit(limit int) *AllMethodsServiceQuery {
	amsq.limit = &limit
	return amsq
}

// Offset adds an offset step to the query.
func (amsq *AllMethodsServiceQuery) Offset(offset int) *AllMethodsServiceQuery {
	amsq.offset = &offset
	return amsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (amsq *AllMethodsServiceQuery) Unique(unique bool) *AllMethodsServiceQuery {
	amsq.unique = &unique
	return amsq
}

// Order adds an order step to the query.
func (amsq *AllMethodsServiceQuery) Order(o ...OrderFunc) *AllMethodsServiceQuery {
	amsq.order = append(amsq.order, o...)
	return amsq
}

// First returns the first AllMethodsService entity from the query.
// Returns a *NotFoundError when no AllMethodsService was found.
func (amsq *AllMethodsServiceQuery) First(ctx context.Context) (*AllMethodsService, error) {
	nodes, err := amsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{allmethodsservice.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) FirstX(ctx context.Context) *AllMethodsService {
	node, err := amsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AllMethodsService ID from the query.
// Returns a *NotFoundError when no AllMethodsService ID was found.
func (amsq *AllMethodsServiceQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = amsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{allmethodsservice.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) FirstIDX(ctx context.Context) int {
	id, err := amsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AllMethodsService entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AllMethodsService entity is found.
// Returns a *NotFoundError when no AllMethodsService entities are found.
func (amsq *AllMethodsServiceQuery) Only(ctx context.Context) (*AllMethodsService, error) {
	nodes, err := amsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{allmethodsservice.Label}
	default:
		return nil, &NotSingularError{allmethodsservice.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) OnlyX(ctx context.Context) *AllMethodsService {
	node, err := amsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AllMethodsService ID in the query.
// Returns a *NotSingularError when more than one AllMethodsService ID is found.
// Returns a *NotFoundError when no entities are found.
func (amsq *AllMethodsServiceQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = amsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{allmethodsservice.Label}
	default:
		err = &NotSingularError{allmethodsservice.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) OnlyIDX(ctx context.Context) int {
	id, err := amsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AllMethodsServices.
func (amsq *AllMethodsServiceQuery) All(ctx context.Context) ([]*AllMethodsService, error) {
	if err := amsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return amsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) AllX(ctx context.Context) []*AllMethodsService {
	nodes, err := amsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AllMethodsService IDs.
func (amsq *AllMethodsServiceQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := amsq.Select(allmethodsservice.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) IDsX(ctx context.Context) []int {
	ids, err := amsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (amsq *AllMethodsServiceQuery) Count(ctx context.Context) (int, error) {
	if err := amsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return amsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) CountX(ctx context.Context) int {
	count, err := amsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (amsq *AllMethodsServiceQuery) Exist(ctx context.Context) (bool, error) {
	if err := amsq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return amsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (amsq *AllMethodsServiceQuery) ExistX(ctx context.Context) bool {
	exist, err := amsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AllMethodsServiceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (amsq *AllMethodsServiceQuery) Clone() *AllMethodsServiceQuery {
	if amsq == nil {
		return nil
	}
	return &AllMethodsServiceQuery{
		config:     amsq.config,
		limit:      amsq.limit,
		offset:     amsq.offset,
		order:      append([]OrderFunc{}, amsq.order...),
		predicates: append([]predicate.AllMethodsService{}, amsq.predicates...),
		// clone intermediate query.
		sql:    amsq.sql.Clone(),
		path:   amsq.path,
		unique: amsq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (amsq *AllMethodsServiceQuery) GroupBy(field string, fields ...string) *AllMethodsServiceGroupBy {
	grbuild := &AllMethodsServiceGroupBy{config: amsq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := amsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return amsq.sqlQuery(ctx), nil
	}
	grbuild.label = allmethodsservice.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (amsq *AllMethodsServiceQuery) Select(fields ...string) *AllMethodsServiceSelect {
	amsq.fields = append(amsq.fields, fields...)
	selbuild := &AllMethodsServiceSelect{AllMethodsServiceQuery: amsq}
	selbuild.label = allmethodsservice.Label
	selbuild.flds, selbuild.scan = &amsq.fields, selbuild.Scan
	return selbuild
}

func (amsq *AllMethodsServiceQuery) prepareQuery(ctx context.Context) error {
	for _, f := range amsq.fields {
		if !allmethodsservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if amsq.path != nil {
		prev, err := amsq.path(ctx)
		if err != nil {
			return err
		}
		amsq.sql = prev
	}
	return nil
}

func (amsq *AllMethodsServiceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AllMethodsService, error) {
	var (
		nodes = []*AllMethodsService{}
		_spec = amsq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*AllMethodsService).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &AllMethodsService{config: amsq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, amsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (amsq *AllMethodsServiceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := amsq.querySpec()
	_spec.Node.Columns = amsq.fields
	if len(amsq.fields) > 0 {
		_spec.Unique = amsq.unique != nil && *amsq.unique
	}
	return sqlgraph.CountNodes(ctx, amsq.driver, _spec)
}

func (amsq *AllMethodsServiceQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := amsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (amsq *AllMethodsServiceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   allmethodsservice.Table,
			Columns: allmethodsservice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: allmethodsservice.FieldID,
			},
		},
		From:   amsq.sql,
		Unique: true,
	}
	if unique := amsq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := amsq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, allmethodsservice.FieldID)
		for i := range fields {
			if fields[i] != allmethodsservice.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := amsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := amsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := amsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := amsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (amsq *AllMethodsServiceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(amsq.driver.Dialect())
	t1 := builder.Table(allmethodsservice.Table)
	columns := amsq.fields
	if len(columns) == 0 {
		columns = allmethodsservice.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if amsq.sql != nil {
		selector = amsq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if amsq.unique != nil && *amsq.unique {
		selector.Distinct()
	}
	for _, p := range amsq.predicates {
		p(selector)
	}
	for _, p := range amsq.order {
		p(selector)
	}
	if offset := amsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := amsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AllMethodsServiceGroupBy is the group-by builder for AllMethodsService entities.
type AllMethodsServiceGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (amsgb *AllMethodsServiceGroupBy) Aggregate(fns ...AggregateFunc) *AllMethodsServiceGroupBy {
	amsgb.fns = append(amsgb.fns, fns...)
	return amsgb
}

// Scan applies the group-by query and scans the result into the given value.
func (amsgb *AllMethodsServiceGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := amsgb.path(ctx)
	if err != nil {
		return err
	}
	amsgb.sql = query
	return amsgb.sqlScan(ctx, v)
}

func (amsgb *AllMethodsServiceGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range amsgb.fields {
		if !allmethodsservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := amsgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := amsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (amsgb *AllMethodsServiceGroupBy) sqlQuery() *sql.Selector {
	selector := amsgb.sql.Select()
	aggregation := make([]string, 0, len(amsgb.fns))
	for _, fn := range amsgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(amsgb.fields)+len(amsgb.fns))
		for _, f := range amsgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(amsgb.fields...)...)
}

// AllMethodsServiceSelect is the builder for selecting fields of AllMethodsService entities.
type AllMethodsServiceSelect struct {
	*AllMethodsServiceQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (amss *AllMethodsServiceSelect) Scan(ctx context.Context, v interface{}) error {
	if err := amss.prepareQuery(ctx); err != nil {
		return err
	}
	amss.sql = amss.AllMethodsServiceQuery.sqlQuery(ctx)
	return amss.sqlScan(ctx, v)
}

func (amss *AllMethodsServiceSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := amss.sql.Query()
	if err := amss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
