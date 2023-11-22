// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package entgql

import (
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
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

// OrderTermOption returns the OrderTermOption for setting the order direction.
func (o OrderDirection) OrderTermOption() sql.OrderTermOption {
	if o == OrderDirectionAsc {
		return sql.OrderAsc()
	}
	return sql.OrderDesc()
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

// Reverse the direction.
func (o OrderDirection) Reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

// PageInfo of a connection type.
type PageInfo[T any] struct {
	HasNextPage     bool       `json:"hasNextPage"`
	HasPreviousPage bool       `json:"hasPreviousPage"`
	StartCursor     *Cursor[T] `json:"startCursor"`
	EndCursor       *Cursor[T] `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor[T any] struct {
	ID    T         `msgpack:"i"`
	Value ent.Value `msgpack:"v,omitempty"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor[T]) MarshalGQL(w io.Writer) {
	quote := []byte{'"'}
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor[T]) UnmarshalGQL(v interface{}) error {
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

// CursorsPredicate converts the given cursors to predicates.
func CursorsPredicate[T any](after, before *Cursor[T], idField, field string, direction OrderDirection) []func(s *sql.Selector) {
	var predicates []func(s *sql.Selector)
	for _, cursor := range []*Cursor[T]{after, before} {
		if cursor == nil {
			continue
		}
		if cursor.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeGT
			} else {
				predicate = sql.CompositeLT
			}
			// Scope the cursor of the current iteration
			// because it will be used in the closure.
			cursor := cursor
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(sql.P(func(b *sql.Builder) {
					// The predicate function is executed on query generation time.
					column := s.C(field)
					// If there is a non-ambiguous match, we use it. That is because
					// some order terms may append joined information to query selection.
					if matches := s.FindSelection(field); len(matches) == 1 {
						column = matches[0]
					}
					b.Join(predicate([]string{column, s.C(idField)}, cursor.Value, cursor.ID))
				}))
			})
		} else {
			if direction == OrderDirectionAsc {
				predicates = append(predicates, sql.FieldGT(idField, cursor.ID))
			} else {
				predicates = append(predicates, sql.FieldLT(idField, cursor.ID))
			}
		}
	}
	return predicates
}

// MultiCursorOptions are the options for building the cursor predicates.
type MultiCursorsOptions struct {
	FieldID     string           // ID field name.
	DirectionID OrderDirection   // ID field direction.
	Fields      []string         // OrderBy fields used by the cursor.
	Directions  []OrderDirection // OrderBy directions used by the cursor.
}

// MultiCursorsPredicate returns a predicate that filters records by the given cursors.
func MultiCursorsPredicate[T any](after, before *Cursor[T], opts *MultiCursorsOptions) ([]func(s *sql.Selector), error) {
	var predicates []func(s *sql.Selector)
	for _, cursor := range []*Cursor[T]{after, before} {
		if cursor == nil {
			continue
		}
		if cursor.Value != nil {
			predicate, err := multiPredicate(cursor, opts)
			if err != nil {
				return nil, err
			}
			predicates = append(predicates, predicate)
		} else {
			if opts.DirectionID == OrderDirectionAsc {
				predicates = append(predicates, sql.FieldGT(opts.FieldID, cursor.ID))
			} else {
				predicates = append(predicates, sql.FieldLT(opts.FieldID, cursor.ID))
			}
		}
	}
	return predicates, nil
}

func multiPredicate[T any](cursor *Cursor[T], opts *MultiCursorsOptions) (func(*sql.Selector), error) {
	values, ok := cursor.Value.([]any)
	if !ok {
		return nil, fmt.Errorf("cursor %T is not a slice", cursor.Value)
	}
	if len(values) != len(opts.Fields) {
		return nil, fmt.Errorf("cursor values length %d do not match orderBy fields length %d", len(values), len(opts.Fields))
	}
	if len(opts.Directions) != len(opts.Fields) {
		return nil, fmt.Errorf("orderBy directions length %d do not match orderBy fields length %d", len(opts.Directions), len(opts.Fields))
	}
	// Ensure the row value is unique by adding
	// the ID field, if not already present.
	if slices.Index(opts.Fields, opts.FieldID) == -1 {
		values = append(values, cursor.ID)
		opts.Fields = append(opts.Fields, opts.FieldID)
		opts.Directions = append(opts.Directions, opts.DirectionID)
	}
	return func(s *sql.Selector) {
		// Given the following terms: x DESC, y ASC, etc. The following predicate will be
		// generated: (x < x1 OR (x = x1 AND y > y1) OR (x = x1 AND y = y1 AND id > last)).

		// getColumnNameForField gets the name for the term and considers non-ambigous matching of
		// terms that may be joined instead of a column on the table.
		getColumnNameForField := func(field string) string  {
			// The predicate function is executed on query generation time.
			column := s.C(field)
			// If there is a non-ambiguous match, we use it. That is because
			// some order terms may append joined information to query selection.
			if matches := s.FindSelection(field); len(matches) == 1 {
				column = matches[0]
			}
			return column
		}

		var or []*sql.Predicate
		for i := range opts.Fields {
			var ands []*sql.Predicate
			for j := 0; j < i; j++ {
				c := getColumnNameForField(opts.Fields[j])
				ands = append(ands, sql.EQ(c, values[j]))
			}
			c := getColumnNameForField(opts.Fields[i])
			if opts.Directions[i] == OrderDirectionAsc {
				ands = append(ands, sql.GT(c, values[i]))
			} else {
				ands = append(ands, sql.LT(c, values[i]))
			}
			or = append(or, sql.And(ands...))
		}
		s.Where(sql.Or(or...))
	}, nil
}
