// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bionicstork/contrib/entgql/internal/todopulid/ent/category"
	"github.com/bionicstork/contrib/entgql/internal/todopulid/ent/schema/pulid"
	"github.com/bionicstork/contrib/entgql/internal/todopulid/ent/todo"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vmihailenco/msgpack/v5"
)

// OrderDirection defines the directions in which to order a list of items.
type OrderDirection string

const (
	// OrderDirectionAsc specifies an ascending order.
	OrderDirectionAsc OrderDirection = "ASC"
	// OrderDirectionDesc specifies a descending order.
	OrderDirectionDesc OrderDirection = "DESC"
)

// Validate the order direction value.
func (o OrderDirection) Validate() error {
	if o != OrderDirectionAsc && o != OrderDirectionDesc {
		return fmt.Errorf("%s is not a valid OrderDirection", o)
	}
	return nil
}

// String implements fmt.Stringer interface.
func (o OrderDirection) String() string {
	return string(o)
}

// MarshalGQL implements graphql.Marshaler interface.
func (o OrderDirection) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(o.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (o *OrderDirection) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("order direction %T must be a string", val)
	}
	*o = OrderDirection(str)
	return o.Validate()
}

func (o OrderDirection) reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

func (o OrderDirection) orderFunc(field string) OrderFunc {
	if o == OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

func cursorsToPredicates(direction OrderDirection, after, before *Cursor, field, idField string) []func(s *sql.Selector) {
	var predicates []func(s *sql.Selector)
	if after != nil {
		if after.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeGT
			} else {
				predicate = sql.CompositeLT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					after.Value, after.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.GT
			} else {
				predicate = sql.LT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					after.ID,
				))
			})
		}
	}
	if before != nil {
		if before.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeLT
			} else {
				predicate = sql.CompositeGT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					before.Value, before.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.LT
			} else {
				predicate = sql.GT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					before.ID,
				))
			})
		}
	}
	return predicates
}

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID    pulid.ID `msgpack:"i"`
	Value Value    `msgpack:"v,omitempty"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	quote := []byte{'"'}
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := msgpack.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(s),
		),
	).Decode(c); err != nil {
		return fmt.Errorf("cannot decode cursor: %w", err)
	}
	return nil
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func getCollectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	field := fc.Field

walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Name == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return getCollectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

// CategoryEdge is the edge representation of Category.
type CategoryEdge struct {
	Node   *Category `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// CategoryConnection is the connection containing edges to Category.
type CategoryConnection struct {
	Edges      []*CategoryEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

// CategoryPaginateOption enables pagination customization.
type CategoryPaginateOption func(*categoryPager) error

// WithCategoryOrder configures pagination ordering.
func WithCategoryOrder(order *CategoryOrder) CategoryPaginateOption {
	if order == nil {
		order = DefaultCategoryOrder
	}
	o := *order
	return func(pager *categoryPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultCategoryOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithCategoryFilter configures pagination filter.
func WithCategoryFilter(filter func(*CategoryQuery) (*CategoryQuery, error)) CategoryPaginateOption {
	return func(pager *categoryPager) error {
		if filter == nil {
			return errors.New("CategoryQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type categoryPager struct {
	order  *CategoryOrder
	filter func(*CategoryQuery) (*CategoryQuery, error)
}

func newCategoryPager(opts []CategoryPaginateOption) (*categoryPager, error) {
	pager := &categoryPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultCategoryOrder
	}
	return pager, nil
}

func (p *categoryPager) applyFilter(query *CategoryQuery) (*CategoryQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *categoryPager) toCursor(c *Category) Cursor {
	return p.order.Field.toCursor(c)
}

func (p *categoryPager) applyCursors(query *CategoryQuery, after, before *Cursor) *CategoryQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultCategoryOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *categoryPager) applyOrder(query *CategoryQuery, reverse bool) *CategoryQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultCategoryOrder.Field {
		query = query.Order(direction.orderFunc(DefaultCategoryOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Category.
func (c *CategoryQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...CategoryPaginateOption,
) (*CategoryConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newCategoryPager(opts)
	if err != nil {
		return nil, err
	}

	if c, err = pager.applyFilter(c); err != nil {
		return nil, err
	}

	conn := &CategoryConnection{Edges: []*CategoryEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := c.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := c.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	c = pager.applyCursors(c, after, before)
	c = pager.applyOrder(c, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		c = c.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := c.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Category
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Category {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Category {
			return nodes[i]
		}
	}

	conn.Edges = make([]*CategoryEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &CategoryEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

var (
	// CategoryOrderFieldText orders Category by text.
	CategoryOrderFieldText = &CategoryOrderField{
		field: category.FieldText,
		toCursor: func(c *Category) Cursor {
			return Cursor{
				ID:    c.ID,
				Value: c.Text,
			}
		},
	}
	// CategoryOrderFieldDuration orders Category by duration.
	CategoryOrderFieldDuration = &CategoryOrderField{
		field: category.FieldDuration,
		toCursor: func(c *Category) Cursor {
			return Cursor{
				ID:    c.ID,
				Value: c.Duration,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f CategoryOrderField) String() string {
	var str string
	switch f.field {
	case category.FieldText:
		str = "TEXT"
	case category.FieldDuration:
		str = "DURATION"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f CategoryOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *CategoryOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("CategoryOrderField %T must be a string", v)
	}
	switch str {
	case "TEXT":
		*f = *CategoryOrderFieldText
	case "DURATION":
		*f = *CategoryOrderFieldDuration
	default:
		return fmt.Errorf("%s is not a valid CategoryOrderField", str)
	}
	return nil
}

// CategoryOrderField defines the ordering field of Category.
type CategoryOrderField struct {
	field    string
	toCursor func(*Category) Cursor
}

// CategoryOrder defines the ordering of Category.
type CategoryOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *CategoryOrderField `json:"field"`
}

// DefaultCategoryOrder is the default ordering of Category.
var DefaultCategoryOrder = &CategoryOrder{
	Direction: OrderDirectionAsc,
	Field: &CategoryOrderField{
		field: category.FieldID,
		toCursor: func(c *Category) Cursor {
			return Cursor{ID: c.ID}
		},
	},
}

// ToEdge converts Category into CategoryEdge.
func (c *Category) ToEdge(order *CategoryOrder) *CategoryEdge {
	if order == nil {
		order = DefaultCategoryOrder
	}
	return &CategoryEdge{
		Node:   c,
		Cursor: order.Field.toCursor(c),
	}
}

// TodoEdge is the edge representation of Todo.
type TodoEdge struct {
	Node   *Todo  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// TodoConnection is the connection containing edges to Todo.
type TodoConnection struct {
	Edges      []*TodoEdge `json:"edges"`
	PageInfo   PageInfo    `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// TodoPaginateOption enables pagination customization.
type TodoPaginateOption func(*todoPager) error

// WithTodoOrder configures pagination ordering.
func WithTodoOrder(order *TodoOrder) TodoPaginateOption {
	if order == nil {
		order = DefaultTodoOrder
	}
	o := *order
	return func(pager *todoPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultTodoOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithTodoFilter configures pagination filter.
func WithTodoFilter(filter func(*TodoQuery) (*TodoQuery, error)) TodoPaginateOption {
	return func(pager *todoPager) error {
		if filter == nil {
			return errors.New("TodoQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type todoPager struct {
	order  *TodoOrder
	filter func(*TodoQuery) (*TodoQuery, error)
}

func newTodoPager(opts []TodoPaginateOption) (*todoPager, error) {
	pager := &todoPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultTodoOrder
	}
	return pager, nil
}

func (p *todoPager) applyFilter(query *TodoQuery) (*TodoQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *todoPager) toCursor(t *Todo) Cursor {
	return p.order.Field.toCursor(t)
}

func (p *todoPager) applyCursors(query *TodoQuery, after, before *Cursor) *TodoQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultTodoOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *todoPager) applyOrder(query *TodoQuery, reverse bool) *TodoQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultTodoOrder.Field {
		query = query.Order(direction.orderFunc(DefaultTodoOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Todo.
func (t *TodoQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...TodoPaginateOption,
) (*TodoConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newTodoPager(opts)
	if err != nil {
		return nil, err
	}

	if t, err = pager.applyFilter(t); err != nil {
		return nil, err
	}

	conn := &TodoConnection{Edges: []*TodoEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := t.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := t.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	t = pager.applyCursors(t, after, before)
	t = pager.applyOrder(t, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		t = t.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		t = t.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := t.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Todo
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Todo {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Todo {
			return nodes[i]
		}
	}

	conn.Edges = make([]*TodoEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &TodoEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

var (
	// TodoOrderFieldCreatedAt orders Todo by created_at.
	TodoOrderFieldCreatedAt = &TodoOrderField{
		field: todo.FieldCreatedAt,
		toCursor: func(t *Todo) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.CreatedAt,
			}
		},
	}
	// TodoOrderFieldStatus orders Todo by status.
	TodoOrderFieldStatus = &TodoOrderField{
		field: todo.FieldStatus,
		toCursor: func(t *Todo) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.Status,
			}
		},
	}
	// TodoOrderFieldPriority orders Todo by priority.
	TodoOrderFieldPriority = &TodoOrderField{
		field: todo.FieldPriority,
		toCursor: func(t *Todo) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.Priority,
			}
		},
	}
	// TodoOrderFieldText orders Todo by text.
	TodoOrderFieldText = &TodoOrderField{
		field: todo.FieldText,
		toCursor: func(t *Todo) Cursor {
			return Cursor{
				ID:    t.ID,
				Value: t.Text,
			}
		},
	}
)

// String implement fmt.Stringer interface.
func (f TodoOrderField) String() string {
	var str string
	switch f.field {
	case todo.FieldCreatedAt:
		str = "CREATED_AT"
	case todo.FieldStatus:
		str = "STATUS"
	case todo.FieldPriority:
		str = "PRIORITY"
	case todo.FieldText:
		str = "TEXT"
	}
	return str
}

// MarshalGQL implements graphql.Marshaler interface.
func (f TodoOrderField) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (f *TodoOrderField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("TodoOrderField %T must be a string", v)
	}
	switch str {
	case "CREATED_AT":
		*f = *TodoOrderFieldCreatedAt
	case "STATUS":
		*f = *TodoOrderFieldStatus
	case "PRIORITY":
		*f = *TodoOrderFieldPriority
	case "TEXT":
		*f = *TodoOrderFieldText
	default:
		return fmt.Errorf("%s is not a valid TodoOrderField", str)
	}
	return nil
}

// TodoOrderField defines the ordering field of Todo.
type TodoOrderField struct {
	field    string
	toCursor func(*Todo) Cursor
}

// TodoOrder defines the ordering of Todo.
type TodoOrder struct {
	Direction OrderDirection  `json:"direction"`
	Field     *TodoOrderField `json:"field"`
}

// DefaultTodoOrder is the default ordering of Todo.
var DefaultTodoOrder = &TodoOrder{
	Direction: OrderDirectionAsc,
	Field: &TodoOrderField{
		field: todo.FieldID,
		toCursor: func(t *Todo) Cursor {
			return Cursor{ID: t.ID}
		},
	},
}

// ToEdge converts Todo into TodoEdge.
func (t *Todo) ToEdge(order *TodoOrder) *TodoEdge {
	if order == nil {
		order = DefaultTodoOrder
	}
	return &TodoEdge{
		Node:   t,
		Cursor: order.Field.toCursor(t),
	}
}
