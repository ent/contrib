// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/entprototest/ent/duplicatenumbermessage"
	"entgo.io/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// DuplicateNumberMessageQuery is the builder for querying DuplicateNumberMessage entities.
type DuplicateNumberMessageQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.DuplicateNumberMessage
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the DuplicateNumberMessageQuery builder.
func (dnmq *DuplicateNumberMessageQuery) Where(ps ...predicate.DuplicateNumberMessage) *DuplicateNumberMessageQuery {
	dnmq.predicates = append(dnmq.predicates, ps...)
	return dnmq
}

// Limit adds a limit step to the query.
func (dnmq *DuplicateNumberMessageQuery) Limit(limit int) *DuplicateNumberMessageQuery {
	dnmq.limit = &limit
	return dnmq
}

// Offset adds an offset step to the query.
func (dnmq *DuplicateNumberMessageQuery) Offset(offset int) *DuplicateNumberMessageQuery {
	dnmq.offset = &offset
	return dnmq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (dnmq *DuplicateNumberMessageQuery) Unique(unique bool) *DuplicateNumberMessageQuery {
	dnmq.unique = &unique
	return dnmq
}

// Order adds an order step to the query.
func (dnmq *DuplicateNumberMessageQuery) Order(o ...OrderFunc) *DuplicateNumberMessageQuery {
	dnmq.order = append(dnmq.order, o...)
	return dnmq
}

// First returns the first DuplicateNumberMessage entity from the query.
// Returns a *NotFoundError when no DuplicateNumberMessage was found.
func (dnmq *DuplicateNumberMessageQuery) First(ctx context.Context) (*DuplicateNumberMessage, error) {
	nodes, err := dnmq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{duplicatenumbermessage.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) FirstX(ctx context.Context) *DuplicateNumberMessage {
	node, err := dnmq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first DuplicateNumberMessage ID from the query.
// Returns a *NotFoundError when no DuplicateNumberMessage ID was found.
func (dnmq *DuplicateNumberMessageQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dnmq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{duplicatenumbermessage.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) FirstIDX(ctx context.Context) int {
	id, err := dnmq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single DuplicateNumberMessage entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one DuplicateNumberMessage entity is found.
// Returns a *NotFoundError when no DuplicateNumberMessage entities are found.
func (dnmq *DuplicateNumberMessageQuery) Only(ctx context.Context) (*DuplicateNumberMessage, error) {
	nodes, err := dnmq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{duplicatenumbermessage.Label}
	default:
		return nil, &NotSingularError{duplicatenumbermessage.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) OnlyX(ctx context.Context) *DuplicateNumberMessage {
	node, err := dnmq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only DuplicateNumberMessage ID in the query.
// Returns a *NotSingularError when more than one DuplicateNumberMessage ID is found.
// Returns a *NotFoundError when no entities are found.
func (dnmq *DuplicateNumberMessageQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dnmq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = &NotSingularError{duplicatenumbermessage.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) OnlyIDX(ctx context.Context) int {
	id, err := dnmq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of DuplicateNumberMessages.
func (dnmq *DuplicateNumberMessageQuery) All(ctx context.Context) ([]*DuplicateNumberMessage, error) {
	if err := dnmq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return dnmq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) AllX(ctx context.Context) []*DuplicateNumberMessage {
	nodes, err := dnmq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of DuplicateNumberMessage IDs.
func (dnmq *DuplicateNumberMessageQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := dnmq.Select(duplicatenumbermessage.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) IDsX(ctx context.Context) []int {
	ids, err := dnmq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (dnmq *DuplicateNumberMessageQuery) Count(ctx context.Context) (int, error) {
	if err := dnmq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return dnmq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) CountX(ctx context.Context) int {
	count, err := dnmq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (dnmq *DuplicateNumberMessageQuery) Exist(ctx context.Context) (bool, error) {
	if err := dnmq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return dnmq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (dnmq *DuplicateNumberMessageQuery) ExistX(ctx context.Context) bool {
	exist, err := dnmq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the DuplicateNumberMessageQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (dnmq *DuplicateNumberMessageQuery) Clone() *DuplicateNumberMessageQuery {
	if dnmq == nil {
		return nil
	}
	return &DuplicateNumberMessageQuery{
		config:     dnmq.config,
		limit:      dnmq.limit,
		offset:     dnmq.offset,
		order:      append([]OrderFunc{}, dnmq.order...),
		predicates: append([]predicate.DuplicateNumberMessage{}, dnmq.predicates...),
		// clone intermediate query.
		sql:    dnmq.sql.Clone(),
		path:   dnmq.path,
		unique: dnmq.unique,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Hello string `json:"hello,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.DuplicateNumberMessage.Query().
//		GroupBy(duplicatenumbermessage.FieldHello).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (dnmq *DuplicateNumberMessageQuery) GroupBy(field string, fields ...string) *DuplicateNumberMessageGroupBy {
	group := &DuplicateNumberMessageGroupBy{config: dnmq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := dnmq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return dnmq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Hello string `json:"hello,omitempty"`
//	}
//
//	client.DuplicateNumberMessage.Query().
//		Select(duplicatenumbermessage.FieldHello).
//		Scan(ctx, &v)
//
func (dnmq *DuplicateNumberMessageQuery) Select(fields ...string) *DuplicateNumberMessageSelect {
	dnmq.fields = append(dnmq.fields, fields...)
	return &DuplicateNumberMessageSelect{DuplicateNumberMessageQuery: dnmq}
}

func (dnmq *DuplicateNumberMessageQuery) prepareQuery(ctx context.Context) error {
	for _, f := range dnmq.fields {
		if !duplicatenumbermessage.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if dnmq.path != nil {
		prev, err := dnmq.path(ctx)
		if err != nil {
			return err
		}
		dnmq.sql = prev
	}
	return nil
}

func (dnmq *DuplicateNumberMessageQuery) sqlAll(ctx context.Context) ([]*DuplicateNumberMessage, error) {
	var (
		nodes = []*DuplicateNumberMessage{}
		_spec = dnmq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &DuplicateNumberMessage{config: dnmq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, dnmq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (dnmq *DuplicateNumberMessageQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := dnmq.querySpec()
	_spec.Node.Columns = dnmq.fields
	if len(dnmq.fields) > 0 {
		_spec.Unique = dnmq.unique != nil && *dnmq.unique
	}
	return sqlgraph.CountNodes(ctx, dnmq.driver, _spec)
}

func (dnmq *DuplicateNumberMessageQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := dnmq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (dnmq *DuplicateNumberMessageQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   duplicatenumbermessage.Table,
			Columns: duplicatenumbermessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: duplicatenumbermessage.FieldID,
			},
		},
		From:   dnmq.sql,
		Unique: true,
	}
	if unique := dnmq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := dnmq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, duplicatenumbermessage.FieldID)
		for i := range fields {
			if fields[i] != duplicatenumbermessage.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := dnmq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := dnmq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := dnmq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := dnmq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (dnmq *DuplicateNumberMessageQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(dnmq.driver.Dialect())
	t1 := builder.Table(duplicatenumbermessage.Table)
	columns := dnmq.fields
	if len(columns) == 0 {
		columns = duplicatenumbermessage.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if dnmq.sql != nil {
		selector = dnmq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if dnmq.unique != nil && *dnmq.unique {
		selector.Distinct()
	}
	for _, p := range dnmq.predicates {
		p(selector)
	}
	for _, p := range dnmq.order {
		p(selector)
	}
	if offset := dnmq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := dnmq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// DuplicateNumberMessageGroupBy is the group-by builder for DuplicateNumberMessage entities.
type DuplicateNumberMessageGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (dnmgb *DuplicateNumberMessageGroupBy) Aggregate(fns ...AggregateFunc) *DuplicateNumberMessageGroupBy {
	dnmgb.fns = append(dnmgb.fns, fns...)
	return dnmgb
}

// Scan applies the group-by query and scans the result into the given value.
func (dnmgb *DuplicateNumberMessageGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := dnmgb.path(ctx)
	if err != nil {
		return err
	}
	dnmgb.sql = query
	return dnmgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := dnmgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(dnmgb.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := dnmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) StringsX(ctx context.Context) []string {
	v, err := dnmgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = dnmgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) StringX(ctx context.Context) string {
	v, err := dnmgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(dnmgb.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := dnmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) IntsX(ctx context.Context) []int {
	v, err := dnmgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = dnmgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) IntX(ctx context.Context) int {
	v, err := dnmgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(dnmgb.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := dnmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := dnmgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = dnmgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) Float64X(ctx context.Context) float64 {
	v, err := dnmgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(dnmgb.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := dnmgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := dnmgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dnmgb *DuplicateNumberMessageGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = dnmgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (dnmgb *DuplicateNumberMessageGroupBy) BoolX(ctx context.Context) bool {
	v, err := dnmgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (dnmgb *DuplicateNumberMessageGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range dnmgb.fields {
		if !duplicatenumbermessage.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := dnmgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := dnmgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (dnmgb *DuplicateNumberMessageGroupBy) sqlQuery() *sql.Selector {
	selector := dnmgb.sql.Select()
	aggregation := make([]string, 0, len(dnmgb.fns))
	for _, fn := range dnmgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(dnmgb.fields)+len(dnmgb.fns))
		for _, f := range dnmgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(dnmgb.fields...)...)
}

// DuplicateNumberMessageSelect is the builder for selecting fields of DuplicateNumberMessage entities.
type DuplicateNumberMessageSelect struct {
	*DuplicateNumberMessageQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (dnms *DuplicateNumberMessageSelect) Scan(ctx context.Context, v interface{}) error {
	if err := dnms.prepareQuery(ctx); err != nil {
		return err
	}
	dnms.sql = dnms.DuplicateNumberMessageQuery.sqlQuery(ctx)
	return dnms.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) ScanX(ctx context.Context, v interface{}) {
	if err := dnms.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Strings(ctx context.Context) ([]string, error) {
	if len(dnms.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := dnms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) StringsX(ctx context.Context) []string {
	v, err := dnms.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = dnms.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) StringX(ctx context.Context) string {
	v, err := dnms.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Ints(ctx context.Context) ([]int, error) {
	if len(dnms.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := dnms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) IntsX(ctx context.Context) []int {
	v, err := dnms.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = dnms.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) IntX(ctx context.Context) int {
	v, err := dnms.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(dnms.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := dnms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) Float64sX(ctx context.Context) []float64 {
	v, err := dnms.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = dnms.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) Float64X(ctx context.Context) float64 {
	v, err := dnms.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(dnms.fields) > 1 {
		return nil, errors.New("ent: DuplicateNumberMessageSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := dnms.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) BoolsX(ctx context.Context) []bool {
	v, err := dnms.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (dnms *DuplicateNumberMessageSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = dnms.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{duplicatenumbermessage.Label}
	default:
		err = fmt.Errorf("ent: DuplicateNumberMessageSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (dnms *DuplicateNumberMessageSelect) BoolX(ctx context.Context) bool {
	v, err := dnms.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (dnms *DuplicateNumberMessageSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := dnms.sql.Query()
	if err := dnms.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
