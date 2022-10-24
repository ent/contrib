// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/contrib/entproto/internal/todo/ent/predicate"
	"entgo.io/contrib/entproto/internal/todo/ent/skipedgeexample"
	"entgo.io/contrib/entproto/internal/todo/ent/user"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SkipEdgeExampleQuery is the builder for querying SkipEdgeExample entities.
type SkipEdgeExampleQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.SkipEdgeExample
	withUser   *UserQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SkipEdgeExampleQuery builder.
func (seeq *SkipEdgeExampleQuery) Where(ps ...predicate.SkipEdgeExample) *SkipEdgeExampleQuery {
	seeq.predicates = append(seeq.predicates, ps...)
	return seeq
}

// Limit adds a limit step to the query.
func (seeq *SkipEdgeExampleQuery) Limit(limit int) *SkipEdgeExampleQuery {
	seeq.limit = &limit
	return seeq
}

// Offset adds an offset step to the query.
func (seeq *SkipEdgeExampleQuery) Offset(offset int) *SkipEdgeExampleQuery {
	seeq.offset = &offset
	return seeq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (seeq *SkipEdgeExampleQuery) Unique(unique bool) *SkipEdgeExampleQuery {
	seeq.unique = &unique
	return seeq
}

// Order adds an order step to the query.
func (seeq *SkipEdgeExampleQuery) Order(o ...OrderFunc) *SkipEdgeExampleQuery {
	seeq.order = append(seeq.order, o...)
	return seeq
}

// QueryUser chains the current query on the "user" edge.
func (seeq *SkipEdgeExampleQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: seeq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := seeq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := seeq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(skipedgeexample.Table, skipedgeexample.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, skipedgeexample.UserTable, skipedgeexample.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(seeq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first SkipEdgeExample entity from the query.
// Returns a *NotFoundError when no SkipEdgeExample was found.
func (seeq *SkipEdgeExampleQuery) First(ctx context.Context) (*SkipEdgeExample, error) {
	nodes, err := seeq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{skipedgeexample.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) FirstX(ctx context.Context) *SkipEdgeExample {
	node, err := seeq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first SkipEdgeExample ID from the query.
// Returns a *NotFoundError when no SkipEdgeExample ID was found.
func (seeq *SkipEdgeExampleQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = seeq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{skipedgeexample.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) FirstIDX(ctx context.Context) int {
	id, err := seeq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single SkipEdgeExample entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one SkipEdgeExample entity is found.
// Returns a *NotFoundError when no SkipEdgeExample entities are found.
func (seeq *SkipEdgeExampleQuery) Only(ctx context.Context) (*SkipEdgeExample, error) {
	nodes, err := seeq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{skipedgeexample.Label}
	default:
		return nil, &NotSingularError{skipedgeexample.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) OnlyX(ctx context.Context) *SkipEdgeExample {
	node, err := seeq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only SkipEdgeExample ID in the query.
// Returns a *NotSingularError when more than one SkipEdgeExample ID is found.
// Returns a *NotFoundError when no entities are found.
func (seeq *SkipEdgeExampleQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = seeq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{skipedgeexample.Label}
	default:
		err = &NotSingularError{skipedgeexample.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) OnlyIDX(ctx context.Context) int {
	id, err := seeq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SkipEdgeExamples.
func (seeq *SkipEdgeExampleQuery) All(ctx context.Context) ([]*SkipEdgeExample, error) {
	if err := seeq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return seeq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) AllX(ctx context.Context) []*SkipEdgeExample {
	nodes, err := seeq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of SkipEdgeExample IDs.
func (seeq *SkipEdgeExampleQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := seeq.Select(skipedgeexample.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) IDsX(ctx context.Context) []int {
	ids, err := seeq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (seeq *SkipEdgeExampleQuery) Count(ctx context.Context) (int, error) {
	if err := seeq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return seeq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) CountX(ctx context.Context) int {
	count, err := seeq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (seeq *SkipEdgeExampleQuery) Exist(ctx context.Context) (bool, error) {
	if err := seeq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return seeq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (seeq *SkipEdgeExampleQuery) ExistX(ctx context.Context) bool {
	exist, err := seeq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SkipEdgeExampleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (seeq *SkipEdgeExampleQuery) Clone() *SkipEdgeExampleQuery {
	if seeq == nil {
		return nil
	}
	return &SkipEdgeExampleQuery{
		config:     seeq.config,
		limit:      seeq.limit,
		offset:     seeq.offset,
		order:      append([]OrderFunc{}, seeq.order...),
		predicates: append([]predicate.SkipEdgeExample{}, seeq.predicates...),
		withUser:   seeq.withUser.Clone(),
		// clone intermediate query.
		sql:    seeq.sql.Clone(),
		path:   seeq.path,
		unique: seeq.unique,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (seeq *SkipEdgeExampleQuery) WithUser(opts ...func(*UserQuery)) *SkipEdgeExampleQuery {
	query := &UserQuery{config: seeq.config}
	for _, opt := range opts {
		opt(query)
	}
	seeq.withUser = query
	return seeq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (seeq *SkipEdgeExampleQuery) GroupBy(field string, fields ...string) *SkipEdgeExampleGroupBy {
	grbuild := &SkipEdgeExampleGroupBy{config: seeq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := seeq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return seeq.sqlQuery(ctx), nil
	}
	grbuild.label = skipedgeexample.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (seeq *SkipEdgeExampleQuery) Select(fields ...string) *SkipEdgeExampleSelect {
	seeq.fields = append(seeq.fields, fields...)
	selbuild := &SkipEdgeExampleSelect{SkipEdgeExampleQuery: seeq}
	selbuild.label = skipedgeexample.Label
	selbuild.flds, selbuild.scan = &seeq.fields, selbuild.Scan
	return selbuild
}

// Aggregate returns a SkipEdgeExampleSelect configured with the given aggregations.
func (seeq *SkipEdgeExampleQuery) Aggregate(fns ...AggregateFunc) *SkipEdgeExampleSelect {
	return seeq.Select().Aggregate(fns...)
}

func (seeq *SkipEdgeExampleQuery) prepareQuery(ctx context.Context) error {
	for _, f := range seeq.fields {
		if !skipedgeexample.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if seeq.path != nil {
		prev, err := seeq.path(ctx)
		if err != nil {
			return err
		}
		seeq.sql = prev
	}
	return nil
}

func (seeq *SkipEdgeExampleQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*SkipEdgeExample, error) {
	var (
		nodes       = []*SkipEdgeExample{}
		withFKs     = seeq.withFKs
		_spec       = seeq.querySpec()
		loadedTypes = [1]bool{
			seeq.withUser != nil,
		}
	)
	if seeq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, skipedgeexample.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*SkipEdgeExample).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &SkipEdgeExample{config: seeq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, seeq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := seeq.withUser; query != nil {
		if err := seeq.loadUser(ctx, query, nodes, nil,
			func(n *SkipEdgeExample, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (seeq *SkipEdgeExampleQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*SkipEdgeExample, init func(*SkipEdgeExample), assign func(*SkipEdgeExample, *User)) error {
	ids := make([]uint32, 0, len(nodes))
	nodeids := make(map[uint32][]*SkipEdgeExample)
	for i := range nodes {
		if nodes[i].user_skip_edge == nil {
			continue
		}
		fk := *nodes[i].user_skip_edge
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_skip_edge" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (seeq *SkipEdgeExampleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := seeq.querySpec()
	_spec.Node.Columns = seeq.fields
	if len(seeq.fields) > 0 {
		_spec.Unique = seeq.unique != nil && *seeq.unique
	}
	return sqlgraph.CountNodes(ctx, seeq.driver, _spec)
}

func (seeq *SkipEdgeExampleQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := seeq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (seeq *SkipEdgeExampleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   skipedgeexample.Table,
			Columns: skipedgeexample.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: skipedgeexample.FieldID,
			},
		},
		From:   seeq.sql,
		Unique: true,
	}
	if unique := seeq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := seeq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, skipedgeexample.FieldID)
		for i := range fields {
			if fields[i] != skipedgeexample.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := seeq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := seeq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := seeq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := seeq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (seeq *SkipEdgeExampleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(seeq.driver.Dialect())
	t1 := builder.Table(skipedgeexample.Table)
	columns := seeq.fields
	if len(columns) == 0 {
		columns = skipedgeexample.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if seeq.sql != nil {
		selector = seeq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if seeq.unique != nil && *seeq.unique {
		selector.Distinct()
	}
	for _, p := range seeq.predicates {
		p(selector)
	}
	for _, p := range seeq.order {
		p(selector)
	}
	if offset := seeq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := seeq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SkipEdgeExampleGroupBy is the group-by builder for SkipEdgeExample entities.
type SkipEdgeExampleGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (seegb *SkipEdgeExampleGroupBy) Aggregate(fns ...AggregateFunc) *SkipEdgeExampleGroupBy {
	seegb.fns = append(seegb.fns, fns...)
	return seegb
}

// Scan applies the group-by query and scans the result into the given value.
func (seegb *SkipEdgeExampleGroupBy) Scan(ctx context.Context, v any) error {
	query, err := seegb.path(ctx)
	if err != nil {
		return err
	}
	seegb.sql = query
	return seegb.sqlScan(ctx, v)
}

func (seegb *SkipEdgeExampleGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range seegb.fields {
		if !skipedgeexample.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := seegb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := seegb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (seegb *SkipEdgeExampleGroupBy) sqlQuery() *sql.Selector {
	selector := seegb.sql.Select()
	aggregation := make([]string, 0, len(seegb.fns))
	for _, fn := range seegb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(seegb.fields)+len(seegb.fns))
		for _, f := range seegb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(seegb.fields...)...)
}

// SkipEdgeExampleSelect is the builder for selecting fields of SkipEdgeExample entities.
type SkipEdgeExampleSelect struct {
	*SkipEdgeExampleQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (sees *SkipEdgeExampleSelect) Aggregate(fns ...AggregateFunc) *SkipEdgeExampleSelect {
	sees.fns = append(sees.fns, fns...)
	return sees
}

// Scan applies the selector query and scans the result into the given value.
func (sees *SkipEdgeExampleSelect) Scan(ctx context.Context, v any) error {
	if err := sees.prepareQuery(ctx); err != nil {
		return err
	}
	sees.sql = sees.SkipEdgeExampleQuery.sqlQuery(ctx)
	return sees.sqlScan(ctx, v)
}

func (sees *SkipEdgeExampleSelect) sqlScan(ctx context.Context, v any) error {
	aggregation := make([]string, 0, len(sees.fns))
	for _, fn := range sees.fns {
		aggregation = append(aggregation, fn(sees.sql))
	}
	switch n := len(*sees.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		sees.sql.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		sees.sql.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := sees.sql.Query()
	if err := sees.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
