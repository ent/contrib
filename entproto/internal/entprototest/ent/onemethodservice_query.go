// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/onemethodservice"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// OneMethodServiceQuery is the builder for querying OneMethodService entities.
type OneMethodServiceQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.OneMethodService
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the OneMethodServiceQuery builder.
func (omsq *OneMethodServiceQuery) Where(ps ...predicate.OneMethodService) *OneMethodServiceQuery {
	omsq.predicates = append(omsq.predicates, ps...)
	return omsq
}

// Limit adds a limit step to the query.
func (omsq *OneMethodServiceQuery) Limit(limit int) *OneMethodServiceQuery {
	omsq.limit = &limit
	return omsq
}

// Offset adds an offset step to the query.
func (omsq *OneMethodServiceQuery) Offset(offset int) *OneMethodServiceQuery {
	omsq.offset = &offset
	return omsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (omsq *OneMethodServiceQuery) Unique(unique bool) *OneMethodServiceQuery {
	omsq.unique = &unique
	return omsq
}

// Order adds an order step to the query.
func (omsq *OneMethodServiceQuery) Order(o ...OrderFunc) *OneMethodServiceQuery {
	omsq.order = append(omsq.order, o...)
	return omsq
}

// First returns the first OneMethodService entity from the query.
// Returns a *NotFoundError when no OneMethodService was found.
func (omsq *OneMethodServiceQuery) First(ctx context.Context) (*OneMethodService, error) {
	nodes, err := omsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{onemethodservice.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) FirstX(ctx context.Context) *OneMethodService {
	node, err := omsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first OneMethodService ID from the query.
// Returns a *NotFoundError when no OneMethodService ID was found.
func (omsq *OneMethodServiceQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = omsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{onemethodservice.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) FirstIDX(ctx context.Context) int {
	id, err := omsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single OneMethodService entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one OneMethodService entity is found.
// Returns a *NotFoundError when no OneMethodService entities are found.
func (omsq *OneMethodServiceQuery) Only(ctx context.Context) (*OneMethodService, error) {
	nodes, err := omsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{onemethodservice.Label}
	default:
		return nil, &NotSingularError{onemethodservice.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) OnlyX(ctx context.Context) *OneMethodService {
	node, err := omsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only OneMethodService ID in the query.
// Returns a *NotSingularError when more than one OneMethodService ID is found.
// Returns a *NotFoundError when no entities are found.
func (omsq *OneMethodServiceQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = omsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{onemethodservice.Label}
	default:
		err = &NotSingularError{onemethodservice.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) OnlyIDX(ctx context.Context) int {
	id, err := omsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of OneMethodServices.
func (omsq *OneMethodServiceQuery) All(ctx context.Context) ([]*OneMethodService, error) {
	if err := omsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return omsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) AllX(ctx context.Context) []*OneMethodService {
	nodes, err := omsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of OneMethodService IDs.
func (omsq *OneMethodServiceQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := omsq.Select(onemethodservice.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) IDsX(ctx context.Context) []int {
	ids, err := omsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (omsq *OneMethodServiceQuery) Count(ctx context.Context) (int, error) {
	if err := omsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return omsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) CountX(ctx context.Context) int {
	count, err := omsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (omsq *OneMethodServiceQuery) Exist(ctx context.Context) (bool, error) {
	if err := omsq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return omsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (omsq *OneMethodServiceQuery) ExistX(ctx context.Context) bool {
	exist, err := omsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the OneMethodServiceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (omsq *OneMethodServiceQuery) Clone() *OneMethodServiceQuery {
	if omsq == nil {
		return nil
	}
	return &OneMethodServiceQuery{
		config:     omsq.config,
		limit:      omsq.limit,
		offset:     omsq.offset,
		order:      append([]OrderFunc{}, omsq.order...),
		predicates: append([]predicate.OneMethodService{}, omsq.predicates...),
		// clone intermediate query.
		sql:    omsq.sql.Clone(),
		path:   omsq.path,
		unique: omsq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (omsq *OneMethodServiceQuery) GroupBy(field string, fields ...string) *OneMethodServiceGroupBy {
	grbuild := &OneMethodServiceGroupBy{config: omsq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := omsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return omsq.sqlQuery(ctx), nil
	}
	grbuild.label = onemethodservice.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (omsq *OneMethodServiceQuery) Select(fields ...string) *OneMethodServiceSelect {
	omsq.fields = append(omsq.fields, fields...)
	selbuild := &OneMethodServiceSelect{OneMethodServiceQuery: omsq}
	selbuild.label = onemethodservice.Label
	selbuild.flds, selbuild.scan = &omsq.fields, selbuild.Scan
	return selbuild
}

func (omsq *OneMethodServiceQuery) prepareQuery(ctx context.Context) error {
	for _, f := range omsq.fields {
		if !onemethodservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if omsq.path != nil {
		prev, err := omsq.path(ctx)
		if err != nil {
			return err
		}
		omsq.sql = prev
	}
	return nil
}

func (omsq *OneMethodServiceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*OneMethodService, error) {
	var (
		nodes = []*OneMethodService{}
		_spec = omsq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*OneMethodService).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &OneMethodService{config: omsq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, omsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (omsq *OneMethodServiceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := omsq.querySpec()
	_spec.Node.Columns = omsq.fields
	if len(omsq.fields) > 0 {
		_spec.Unique = omsq.unique != nil && *omsq.unique
	}
	return sqlgraph.CountNodes(ctx, omsq.driver, _spec)
}

func (omsq *OneMethodServiceQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := omsq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (omsq *OneMethodServiceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   onemethodservice.Table,
			Columns: onemethodservice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: onemethodservice.FieldID,
			},
		},
		From:   omsq.sql,
		Unique: true,
	}
	if unique := omsq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := omsq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, onemethodservice.FieldID)
		for i := range fields {
			if fields[i] != onemethodservice.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := omsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := omsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := omsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := omsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (omsq *OneMethodServiceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(omsq.driver.Dialect())
	t1 := builder.Table(onemethodservice.Table)
	columns := omsq.fields
	if len(columns) == 0 {
		columns = onemethodservice.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if omsq.sql != nil {
		selector = omsq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if omsq.unique != nil && *omsq.unique {
		selector.Distinct()
	}
	for _, p := range omsq.predicates {
		p(selector)
	}
	for _, p := range omsq.order {
		p(selector)
	}
	if offset := omsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := omsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// OneMethodServiceGroupBy is the group-by builder for OneMethodService entities.
type OneMethodServiceGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (omsgb *OneMethodServiceGroupBy) Aggregate(fns ...AggregateFunc) *OneMethodServiceGroupBy {
	omsgb.fns = append(omsgb.fns, fns...)
	return omsgb
}

// Scan applies the group-by query and scans the result into the given value.
func (omsgb *OneMethodServiceGroupBy) Scan(ctx context.Context, v any) error {
	query, err := omsgb.path(ctx)
	if err != nil {
		return err
	}
	omsgb.sql = query
	return omsgb.sqlScan(ctx, v)
}

func (omsgb *OneMethodServiceGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range omsgb.fields {
		if !onemethodservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := omsgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := omsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (omsgb *OneMethodServiceGroupBy) sqlQuery() *sql.Selector {
	selector := omsgb.sql.Select()
	aggregation := make([]string, 0, len(omsgb.fns))
	for _, fn := range omsgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(omsgb.fields)+len(omsgb.fns))
		for _, f := range omsgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(omsgb.fields...)...)
}

// OneMethodServiceSelect is the builder for selecting fields of OneMethodService entities.
type OneMethodServiceSelect struct {
	*OneMethodServiceQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (omss *OneMethodServiceSelect) Scan(ctx context.Context, v any) error {
	if err := omss.prepareQuery(ctx); err != nil {
		return err
	}
	omss.sql = omss.OneMethodServiceQuery.sqlQuery(ctx)
	return omss.sqlScan(ctx, v)
}

func (omss *OneMethodServiceSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := omss.sql.Query()
	if err := omss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
