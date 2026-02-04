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

package entgql_test

import (
	"testing"

	"entgo.io/contrib/entgql"
	"github.com/stretchr/testify/require"
)

func TestWithSplitGoFiles(t *testing.T) {
	t.Parallel()

	// Test that WithSplitGoFiles can be used with NewExtension
	ext, err := entgql.NewExtension(
		entgql.WithSplitGoFiles(true),
	)
	require.NoError(t, err)
	require.NotNil(t, ext)

	// Test that WithSplitGoFiles(false) doesn't cause issues
	ext, err = entgql.NewExtension(
		entgql.WithSplitGoFiles(false),
	)
	require.NoError(t, err)
	require.NotNil(t, ext)

	// Test combination with other options
	ext, err = entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithWhereInputs(true),
		entgql.WithSplitGoFiles(true),
	)
	require.NoError(t, err)
	require.NotNil(t, ext)
}

func TestExtensionOptionsMutualExclusivity(t *testing.T) {
	t.Parallel()

	// WithSchemaPath and WithSchemaDir are mutually exclusive
	_, err := entgql.NewExtension(
		entgql.WithSchemaPath("./ent.graphql"),
		entgql.WithSchemaDir("./gql"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")

	// But the reverse order should also fail
	_, err = entgql.NewExtension(
		entgql.WithSchemaDir("./gql"),
		entgql.WithSchemaPath("./ent.graphql"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")
}
