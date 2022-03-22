// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/contrib/entproto/internal/entprototest/ent/twomethodservice"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TwoMethodServiceQuery is the builder for querying TwoMethodService entities.
type TwoMethodServiceQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.TwoMethodService
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the TwoMethodServiceQuery builder.
func (tmsq *TwoMethodServiceQuery) Where(ps ...predicate.TwoMethodService) *TwoMethodServiceQuery {
	tmsq.predicates = append(tmsq.predicates, ps...)
	return tmsq
}

// Limit adds a limit step to the query.
func (tmsq *TwoMethodServiceQuery) Limit(limit int) *TwoMethodServiceQuery {
	tmsq.limit = &limit
	return tmsq
}

// Offset adds an offset step to the query.
func (tmsq *TwoMethodServiceQuery) Offset(offset int) *TwoMethodServiceQuery {
	tmsq.offset = &offset
	return tmsq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (tmsq *TwoMethodServiceQuery) Unique(unique bool) *TwoMethodServiceQuery {
	tmsq.unique = &unique
	return tmsq
}

// Order adds an order step to the query.
func (tmsq *TwoMethodServiceQuery) Order(o ...OrderFunc) *TwoMethodServiceQuery {
	tmsq.order = append(tmsq.order, o...)
	return tmsq
}

// First returns the first TwoMethodService entity from the query.
// Returns a *NotFoundError when no TwoMethodService was found.
func (tmsq *TwoMethodServiceQuery) First(ctx context.Context) (*TwoMethodService, error) {
	nodes, err := tmsq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{twomethodservice.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) FirstX(ctx context.Context) *TwoMethodService {
	node, err := tmsq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first TwoMethodService ID from the query.
// Returns a *NotFoundError when no TwoMethodService ID was found.
func (tmsq *TwoMethodServiceQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tmsq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{twomethodservice.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) FirstIDX(ctx context.Context) int {
	id, err := tmsq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single TwoMethodService entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one TwoMethodService entity is found.
// Returns a *NotFoundError when no TwoMethodService entities are found.
func (tmsq *TwoMethodServiceQuery) Only(ctx context.Context) (*TwoMethodService, error) {
	nodes, err := tmsq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{twomethodservice.Label}
	default:
		return nil, &NotSingularError{twomethodservice.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) OnlyX(ctx context.Context) *TwoMethodService {
	node, err := tmsq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only TwoMethodService ID in the query.
// Returns a *NotSingularError when more than one TwoMethodService ID is found.
// Returns a *NotFoundError when no entities are found.
func (tmsq *TwoMethodServiceQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = tmsq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{twomethodservice.Label}
	default:
		err = &NotSingularError{twomethodservice.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) OnlyIDX(ctx context.Context) int {
	id, err := tmsq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of TwoMethodServices.
func (tmsq *TwoMethodServiceQuery) All(ctx context.Context) ([]*TwoMethodService, error) {
	if err := tmsq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return tmsq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) AllX(ctx context.Context) []*TwoMethodService {
	nodes, err := tmsq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of TwoMethodService IDs.
func (tmsq *TwoMethodServiceQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := tmsq.Select(twomethodservice.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) IDsX(ctx context.Context) []int {
	ids, err := tmsq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (tmsq *TwoMethodServiceQuery) Count(ctx context.Context) (int, error) {
	if err := tmsq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return tmsq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) CountX(ctx context.Context) int {
	count, err := tmsq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (tmsq *TwoMethodServiceQuery) Exist(ctx context.Context) (bool, error) {
	if err := tmsq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return tmsq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (tmsq *TwoMethodServiceQuery) ExistX(ctx context.Context) bool {
	exist, err := tmsq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the TwoMethodServiceQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (tmsq *TwoMethodServiceQuery) Clone() *TwoMethodServiceQuery {
	if tmsq == nil {
		return nil
	}
	return &TwoMethodServiceQuery{
		config:     tmsq.config,
		limit:      tmsq.limit,
		offset:     tmsq.offset,
		order:      append([]OrderFunc{}, tmsq.order...),
		predicates: append([]predicate.TwoMethodService{}, tmsq.predicates...),
		// clone intermediate query.
		sql:    tmsq.sql.Clone(),
		path:   tmsq.path,
		unique: tmsq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (tmsq *TwoMethodServiceQuery) GroupBy(field string, fields ...string) *TwoMethodServiceGroupBy {
	grbuild := &TwoMethodServiceGroupBy{config: tmsq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := tmsq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return tmsq.sqlQuery(ctx), nil
	}
	grbuild.label = twomethodservice.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (tmsq *TwoMethodServiceQuery) Select(fields ...string) *TwoMethodServiceSelect {
	tmsq.fields = append(tmsq.fields, fields...)
	selbuild := &TwoMethodServiceSelect{TwoMethodServiceQuery: tmsq}
	selbuild.label = twomethodservice.Label
	selbuild.flds, selbuild.scan = &tmsq.fields, selbuild.Scan
	return selbuild
}

func (tmsq *TwoMethodServiceQuery) prepareQuery(ctx context.Context) error {
	for _, f := range tmsq.fields {
		if !twomethodservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if tmsq.path != nil {
		prev, err := tmsq.path(ctx)
		if err != nil {
			return err
		}
		tmsq.sql = prev
	}
	return nil
}

func (tmsq *TwoMethodServiceQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*TwoMethodService, error) {
	var (
		nodes = []*TwoMethodService{}
		_spec = tmsq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		return (*TwoMethodService).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		node := &TwoMethodService{config: tmsq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, tmsq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (tmsq *TwoMethodServiceQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := tmsq.querySpec()
	_spec.Node.Columns = tmsq.fields
	if len(tmsq.fields) > 0 {
		_spec.Unique = tmsq.unique != nil && *tmsq.unique
	}
	return sqlgraph.CountNodes(ctx, tmsq.driver, _spec)
}

func (tmsq *TwoMethodServiceQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := tmsq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (tmsq *TwoMethodServiceQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   twomethodservice.Table,
			Columns: twomethodservice.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: twomethodservice.FieldID,
			},
		},
		From:   tmsq.sql,
		Unique: true,
	}
	if unique := tmsq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := tmsq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, twomethodservice.FieldID)
		for i := range fields {
			if fields[i] != twomethodservice.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := tmsq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := tmsq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := tmsq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := tmsq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (tmsq *TwoMethodServiceQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(tmsq.driver.Dialect())
	t1 := builder.Table(twomethodservice.Table)
	columns := tmsq.fields
	if len(columns) == 0 {
		columns = twomethodservice.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if tmsq.sql != nil {
		selector = tmsq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if tmsq.unique != nil && *tmsq.unique {
		selector.Distinct()
	}
	for _, p := range tmsq.predicates {
		p(selector)
	}
	for _, p := range tmsq.order {
		p(selector)
	}
	if offset := tmsq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := tmsq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// TwoMethodServiceGroupBy is the group-by builder for TwoMethodService entities.
type TwoMethodServiceGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (tmsgb *TwoMethodServiceGroupBy) Aggregate(fns ...AggregateFunc) *TwoMethodServiceGroupBy {
	tmsgb.fns = append(tmsgb.fns, fns...)
	return tmsgb
}

// Scan applies the group-by query and scans the result into the given value.
func (tmsgb *TwoMethodServiceGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := tmsgb.path(ctx)
	if err != nil {
		return err
	}
	tmsgb.sql = query
	return tmsgb.sqlScan(ctx, v)
}

func (tmsgb *TwoMethodServiceGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range tmsgb.fields {
		if !twomethodservice.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := tmsgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := tmsgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (tmsgb *TwoMethodServiceGroupBy) sqlQuery() *sql.Selector {
	selector := tmsgb.sql.Select()
	aggregation := make([]string, 0, len(tmsgb.fns))
	for _, fn := range tmsgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(tmsgb.fields)+len(tmsgb.fns))
		for _, f := range tmsgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(tmsgb.fields...)...)
}

// TwoMethodServiceSelect is the builder for selecting fields of TwoMethodService entities.
type TwoMethodServiceSelect struct {
	*TwoMethodServiceQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (tmss *TwoMethodServiceSelect) Scan(ctx context.Context, v interface{}) error {
	if err := tmss.prepareQuery(ctx); err != nil {
		return err
	}
	tmss.sql = tmss.TwoMethodServiceQuery.sqlQuery(ctx)
	return tmss.sqlScan(ctx, v)
}

func (tmss *TwoMethodServiceSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := tmss.sql.Query()
	if err := tmss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
