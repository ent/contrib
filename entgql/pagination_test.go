// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package entgql

import (
	"testing"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/require"
)

func TestMultiCursorsPredicate(t *testing.T) {
	t.Run("joined order field", func(t *testing.T) {
		builder := sql.Dialect(dialect.Postgres)
		users, groups := builder.Table("users"), builder.Table("groups")
		group := builder.Select(groups.C("id")).
			From(groups).
			AppendSelectAs(groups.C("name"), "group_name")
		selector := builder.Select(users.C("id")).From(users)
		selector.LeftJoin(group).On(users.C("group_id"), group.C("id"))
		selector.AppendSelect(group.C("group_name"))

		predicates, err := MultiCursorsPredicate(
			&Cursor[int]{ID: 10, Value: []any{"engineering"}},
			nil,
			&MultiCursorsOptions{
				FieldID:     "id",
				DirectionID: OrderDirectionAsc,
				Fields:      []string{"group_name"},
				Directions:  []OrderDirection{OrderDirectionAsc},
			},
		)
		require.NoError(t, err)
		for _, predicate := range predicates {
			predicate(selector)
		}
		query, _ := selector.Query()
		require.NotContains(t, query, `"users"."group_name"`)
		require.Contains(t, query, `"t1"."group_name" > $1 OR ("t1"."group_name" = $2 AND "users"."id" > $3)`)
	})

	t.Run("does not mutate options", func(t *testing.T) {
		opts := &MultiCursorsOptions{
			FieldID:     "id",
			DirectionID: OrderDirectionAsc,
			Fields:      []string{"name"},
			Directions:  []OrderDirection{OrderDirectionAsc},
		}
		predicates, err := MultiCursorsPredicate(
			&Cursor[int]{ID: 10, Value: []any{"alice"}},
			&Cursor[int]{ID: 20, Value: []any{"bob"}},
			opts,
		)
		require.NoError(t, err)
		require.Len(t, predicates, 2)
		require.Equal(t, []string{"name"}, opts.Fields)
		require.Equal(t, []OrderDirection{OrderDirectionAsc}, opts.Directions)
	})
}
