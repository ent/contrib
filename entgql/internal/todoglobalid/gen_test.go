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

package todo_test

import (
	"os"
	"path/filepath"
	"testing"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestGeneratedSchema(t *testing.T) {
	tempDir := t.TempDir()
	gqlcfg, err := os.ReadFile("./gqlgen.yml")
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, "gqlgen.yml"), gqlcfg, 0644)
	require.NoError(t, err)
	ex, err := entgql.NewExtension(
		entgql.WithConfigPath(filepath.Join(tempDir, "gqlgen.yml")),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath(filepath.Join(tempDir, "ent.graphql")),
		entgql.WithWhereInputs(true),
		entgql.WithNodeDescriptor(true),
	)
	require.NoError(t, err)
	err = entc.Generate("./ent/schema", &gen.Config{
		Target: tempDir,
		Features: []gen.Feature{
			gen.FeatureModifier,
		},
	}, entc.Extensions(ex))
	require.NoError(t, err)
	expected, err := os.ReadFile("./ent.graphql")
	require.NoError(t, err)
	actual, err := os.ReadFile(filepath.Join(tempDir, "ent.graphql"))
	require.NoError(t, err)
	require.Equal(t, string(expected), string(actual))
}
