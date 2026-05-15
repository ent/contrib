// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package entgql_test

import (
	"testing"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/require"
)

// applyPredicates runs predicates against a fresh Postgres-dialect selector
// and returns the SQL text + bound args. Using Postgres dialect to match
// our target deployment and make assertions stable.
func applyPredicates(t *testing.T, table string, preds []func(s *sql.Selector)) (string, []interface{}) {
	t.Helper()
	s := sql.Dialect(dialect.Postgres).Select("id").From(sql.Table(table))
	for _, p := range preds {
		p(s)
	}
	sqlText, args := s.Query()
	return sqlText, args
}

func TestCursorsPredicateExpr_AfterValueAscEmitsExpression(t *testing.T) {
	t.Parallel()
	after := &entgql.Cursor[string]{ID: "id-2", Value: "alpha"}
	preds := entgql.CursorsPredicateExpr(after, nil, "id", `left("name", 256)`, entgql.OrderDirectionAsc)
	require.Len(t, preds, 1)

	sqlText, args := applyPredicates(t, "contacts", preds)
	require.Contains(t, sqlText, `left("name", 256) > $1`)
	require.Contains(t, sqlText, `left("name", 256) = $2`)
	require.Contains(t, sqlText, `"contacts"."id" > $3`)
	require.Equal(t, []interface{}{"alpha", "alpha", "id-2"}, args)
}

func TestCursorsPredicateExpr_BeforeValueDescEmitsLessThan(t *testing.T) {
	t.Parallel()
	before := &entgql.Cursor[string]{ID: "id-5", Value: "zulu"}
	preds := entgql.CursorsPredicateExpr(nil, before, "id", `left("name", 256)`, entgql.OrderDirectionDesc)
	require.Len(t, preds, 1)

	sqlText, _ := applyPredicates(t, "contacts", preds)
	require.Contains(t, sqlText, `left("name", 256) < $1`)
	require.Contains(t, sqlText, `"contacts"."id" < $3`)
}

func TestCursorsPredicateExpr_NilValueFallsBackToIDComparison(t *testing.T) {
	t.Parallel()
	// Cursors with nil Value (e.g. when ordering by ID itself) should
	// behave identically to CursorsPredicate: just an ID comparison, no
	// expression substitution.
	after := &entgql.Cursor[string]{ID: "id-7"}
	preds := entgql.CursorsPredicateExpr(after, nil, "id", `left("name", 256)`, entgql.OrderDirectionAsc)
	require.Len(t, preds, 1)

	sqlText, args := applyPredicates(t, "contacts", preds)
	require.NotContains(t, sqlText, "left(", "no expression substitution when cursor has nil Value")
	require.Contains(t, sqlText, `"id"`)
	require.Equal(t, []interface{}{"id-7"}, args)
}

func TestCursorsPredicateExpr_NilCursorsProduceNoPredicates(t *testing.T) {
	t.Parallel()
	preds := entgql.CursorsPredicateExpr[string](nil, nil, "id", `left("name", 256)`, entgql.OrderDirectionAsc)
	require.Empty(t, preds)
}

func TestCursorsPredicateExpr_BothBoundsAppendInOrder(t *testing.T) {
	t.Parallel()
	after := &entgql.Cursor[string]{ID: "id-2", Value: "alpha"}
	before := &entgql.Cursor[string]{ID: "id-9", Value: "omega"}
	preds := entgql.CursorsPredicateExpr(after, before, "id", `left("name", 256)`, entgql.OrderDirectionAsc)
	require.Len(t, preds, 2)
}
