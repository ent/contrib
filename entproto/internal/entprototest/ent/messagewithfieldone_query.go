// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/bionicstork/contrib/entproto/internal/entprototest/ent/messagewithfieldone"
	"github.com/bionicstork/contrib/entproto/internal/entprototest/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// MessageWithFieldOneQuery is the builder for querying MessageWithFieldOne entities.
type MessageWithFieldOneQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.MessageWithFieldOne
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MessageWithFieldOneQuery builder.
func (mwfoq *MessageWithFieldOneQuery) Where(ps ...predicate.MessageWithFieldOne) *MessageWithFieldOneQuery {
	mwfoq.predicates = append(mwfoq.predicates, ps...)
	return mwfoq
}

// Limit adds a limit step to the query.
func (mwfoq *MessageWithFieldOneQuery) Limit(limit int) *MessageWithFieldOneQuery {
	mwfoq.limit = &limit
	return mwfoq
}

// Offset adds an offset step to the query.
func (mwfoq *MessageWithFieldOneQuery) Offset(offset int) *MessageWithFieldOneQuery {
	mwfoq.offset = &offset
	return mwfoq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mwfoq *MessageWithFieldOneQuery) Unique(unique bool) *MessageWithFieldOneQuery {
	mwfoq.unique = &unique
	return mwfoq
}

// Order adds an order step to the query.
func (mwfoq *MessageWithFieldOneQuery) Order(o ...OrderFunc) *MessageWithFieldOneQuery {
	mwfoq.order = append(mwfoq.order, o...)
	return mwfoq
}

// First returns the first MessageWithFieldOne entity from the query.
// Returns a *NotFoundError when no MessageWithFieldOne was found.
func (mwfoq *MessageWithFieldOneQuery) First(ctx context.Context) (*MessageWithFieldOne, error) {
	nodes, err := mwfoq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{messagewithfieldone.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) FirstX(ctx context.Context) *MessageWithFieldOne {
	node, err := mwfoq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first MessageWithFieldOne ID from the query.
// Returns a *NotFoundError when no MessageWithFieldOne ID was found.
func (mwfoq *MessageWithFieldOneQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mwfoq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{messagewithfieldone.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) FirstIDX(ctx context.Context) int {
	id, err := mwfoq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single MessageWithFieldOne entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one MessageWithFieldOne entity is not found.
// Returns a *NotFoundError when no MessageWithFieldOne entities are found.
func (mwfoq *MessageWithFieldOneQuery) Only(ctx context.Context) (*MessageWithFieldOne, error) {
	nodes, err := mwfoq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{messagewithfieldone.Label}
	default:
		return nil, &NotSingularError{messagewithfieldone.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) OnlyX(ctx context.Context) *MessageWithFieldOne {
	node, err := mwfoq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only MessageWithFieldOne ID in the query.
// Returns a *NotSingularError when exactly one MessageWithFieldOne ID is not found.
// Returns a *NotFoundError when no entities are found.
func (mwfoq *MessageWithFieldOneQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = mwfoq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = &NotSingularError{messagewithfieldone.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) OnlyIDX(ctx context.Context) int {
	id, err := mwfoq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of MessageWithFieldOnes.
func (mwfoq *MessageWithFieldOneQuery) All(ctx context.Context) ([]*MessageWithFieldOne, error) {
	if err := mwfoq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return mwfoq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) AllX(ctx context.Context) []*MessageWithFieldOne {
	nodes, err := mwfoq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of MessageWithFieldOne IDs.
func (mwfoq *MessageWithFieldOneQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := mwfoq.Select(messagewithfieldone.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) IDsX(ctx context.Context) []int {
	ids, err := mwfoq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mwfoq *MessageWithFieldOneQuery) Count(ctx context.Context) (int, error) {
	if err := mwfoq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return mwfoq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) CountX(ctx context.Context) int {
	count, err := mwfoq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mwfoq *MessageWithFieldOneQuery) Exist(ctx context.Context) (bool, error) {
	if err := mwfoq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return mwfoq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (mwfoq *MessageWithFieldOneQuery) ExistX(ctx context.Context) bool {
	exist, err := mwfoq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MessageWithFieldOneQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mwfoq *MessageWithFieldOneQuery) Clone() *MessageWithFieldOneQuery {
	if mwfoq == nil {
		return nil
	}
	return &MessageWithFieldOneQuery{
		config:     mwfoq.config,
		limit:      mwfoq.limit,
		offset:     mwfoq.offset,
		order:      append([]OrderFunc{}, mwfoq.order...),
		predicates: append([]predicate.MessageWithFieldOne{}, mwfoq.predicates...),
		// clone intermediate query.
		sql:  mwfoq.sql.Clone(),
		path: mwfoq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		FieldOne int32 `json:"field_one,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.MessageWithFieldOne.Query().
//		GroupBy(messagewithfieldone.FieldFieldOne).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (mwfoq *MessageWithFieldOneQuery) GroupBy(field string, fields ...string) *MessageWithFieldOneGroupBy {
	group := &MessageWithFieldOneGroupBy{config: mwfoq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := mwfoq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return mwfoq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		FieldOne int32 `json:"field_one,omitempty"`
//	}
//
//	client.MessageWithFieldOne.Query().
//		Select(messagewithfieldone.FieldFieldOne).
//		Scan(ctx, &v)
//
func (mwfoq *MessageWithFieldOneQuery) Select(fields ...string) *MessageWithFieldOneSelect {
	mwfoq.fields = append(mwfoq.fields, fields...)
	return &MessageWithFieldOneSelect{MessageWithFieldOneQuery: mwfoq}
}

func (mwfoq *MessageWithFieldOneQuery) prepareQuery(ctx context.Context) error {
	for _, f := range mwfoq.fields {
		if !messagewithfieldone.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mwfoq.path != nil {
		prev, err := mwfoq.path(ctx)
		if err != nil {
			return err
		}
		mwfoq.sql = prev
	}
	return nil
}

func (mwfoq *MessageWithFieldOneQuery) sqlAll(ctx context.Context) ([]*MessageWithFieldOne, error) {
	var (
		nodes = []*MessageWithFieldOne{}
		_spec = mwfoq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &MessageWithFieldOne{config: mwfoq.config}
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
	if err := sqlgraph.QueryNodes(ctx, mwfoq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (mwfoq *MessageWithFieldOneQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mwfoq.querySpec()
	_spec.Node.Columns = mwfoq.fields
	if len(mwfoq.fields) > 0 {
		_spec.Unique = mwfoq.unique != nil && *mwfoq.unique
	}
	return sqlgraph.CountNodes(ctx, mwfoq.driver, _spec)
}

func (mwfoq *MessageWithFieldOneQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := mwfoq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (mwfoq *MessageWithFieldOneQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   messagewithfieldone.Table,
			Columns: messagewithfieldone.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: messagewithfieldone.FieldID,
			},
		},
		From:   mwfoq.sql,
		Unique: true,
	}
	if unique := mwfoq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := mwfoq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, messagewithfieldone.FieldID)
		for i := range fields {
			if fields[i] != messagewithfieldone.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mwfoq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mwfoq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mwfoq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mwfoq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mwfoq *MessageWithFieldOneQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mwfoq.driver.Dialect())
	t1 := builder.Table(messagewithfieldone.Table)
	columns := mwfoq.fields
	if len(columns) == 0 {
		columns = messagewithfieldone.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mwfoq.sql != nil {
		selector = mwfoq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mwfoq.unique != nil && *mwfoq.unique {
		selector.Distinct()
	}
	for _, p := range mwfoq.predicates {
		p(selector)
	}
	for _, p := range mwfoq.order {
		p(selector)
	}
	if offset := mwfoq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mwfoq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MessageWithFieldOneGroupBy is the group-by builder for MessageWithFieldOne entities.
type MessageWithFieldOneGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mwfogb *MessageWithFieldOneGroupBy) Aggregate(fns ...AggregateFunc) *MessageWithFieldOneGroupBy {
	mwfogb.fns = append(mwfogb.fns, fns...)
	return mwfogb
}

// Scan applies the group-by query and scans the result into the given value.
func (mwfogb *MessageWithFieldOneGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := mwfogb.path(ctx)
	if err != nil {
		return err
	}
	mwfogb.sql = query
	return mwfogb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := mwfogb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(mwfogb.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := mwfogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) StringsX(ctx context.Context) []string {
	v, err := mwfogb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = mwfogb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) StringX(ctx context.Context) string {
	v, err := mwfogb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(mwfogb.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := mwfogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) IntsX(ctx context.Context) []int {
	v, err := mwfogb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = mwfogb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) IntX(ctx context.Context) int {
	v, err := mwfogb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(mwfogb.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := mwfogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := mwfogb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = mwfogb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) Float64X(ctx context.Context) float64 {
	v, err := mwfogb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(mwfogb.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := mwfogb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := mwfogb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (mwfogb *MessageWithFieldOneGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = mwfogb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (mwfogb *MessageWithFieldOneGroupBy) BoolX(ctx context.Context) bool {
	v, err := mwfogb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (mwfogb *MessageWithFieldOneGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range mwfogb.fields {
		if !messagewithfieldone.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := mwfogb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mwfogb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (mwfogb *MessageWithFieldOneGroupBy) sqlQuery() *sql.Selector {
	selector := mwfogb.sql.Select()
	aggregation := make([]string, 0, len(mwfogb.fns))
	for _, fn := range mwfogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(mwfogb.fields)+len(mwfogb.fns))
		for _, f := range mwfogb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(mwfogb.fields...)...)
}

// MessageWithFieldOneSelect is the builder for selecting fields of MessageWithFieldOne entities.
type MessageWithFieldOneSelect struct {
	*MessageWithFieldOneQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (mwfos *MessageWithFieldOneSelect) Scan(ctx context.Context, v interface{}) error {
	if err := mwfos.prepareQuery(ctx); err != nil {
		return err
	}
	mwfos.sql = mwfos.MessageWithFieldOneQuery.sqlQuery(ctx)
	return mwfos.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) ScanX(ctx context.Context, v interface{}) {
	if err := mwfos.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Strings(ctx context.Context) ([]string, error) {
	if len(mwfos.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := mwfos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) StringsX(ctx context.Context) []string {
	v, err := mwfos.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = mwfos.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) StringX(ctx context.Context) string {
	v, err := mwfos.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Ints(ctx context.Context) ([]int, error) {
	if len(mwfos.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := mwfos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) IntsX(ctx context.Context) []int {
	v, err := mwfos.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = mwfos.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) IntX(ctx context.Context) int {
	v, err := mwfos.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(mwfos.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := mwfos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) Float64sX(ctx context.Context) []float64 {
	v, err := mwfos.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = mwfos.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) Float64X(ctx context.Context) float64 {
	v, err := mwfos.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(mwfos.fields) > 1 {
		return nil, errors.New("ent: MessageWithFieldOneSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := mwfos.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) BoolsX(ctx context.Context) []bool {
	v, err := mwfos.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (mwfos *MessageWithFieldOneSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = mwfos.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{messagewithfieldone.Label}
	default:
		err = fmt.Errorf("ent: MessageWithFieldOneSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (mwfos *MessageWithFieldOneSelect) BoolX(ctx context.Context) bool {
	v, err := mwfos.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (mwfos *MessageWithFieldOneSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := mwfos.sql.Query()
	if err := mwfos.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
