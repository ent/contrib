// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package entgql_test

import (
	"testing"

	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebookincubator/ent-contrib/entgql"
	"github.com/stretchr/testify/require"
)

func TestAnnotation(t *testing.T) {
	require.Implements(t, (*field.Annotation)(nil), entgql.Annotation{})
	require.Implements(t, (*edge.Annotation)(nil), entgql.Annotation{})

	annotation := entgql.OrderField("foo")
	require.Equal(t, "foo", annotation.OrderField)

	annotation = entgql.Bind()
	require.True(t, annotation.Bind)
	require.Empty(t, annotation.Mapping)

	names := []string{"foo", "bar", "baz"}
	annotation = entgql.MapsTo(names...)
	require.False(t, annotation.Bind)
	require.ElementsMatch(t, names, annotation.Mapping)
}
