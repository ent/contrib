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

package plugin

import (
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEntGQL_buildTypes(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin, err := NewEntGQLPlugin(graph)

	require.NoError(t, err)
	types, err := plugin.buildTypes()
	require.NoError(t, err)

	require.Equal(t, `type Category {
	id: ID!
	text: String!
	uuidA: UUID
	status: CategoryStatus!
	config: CategoryConfig!
	duration: Duration!
	count: Uint64!
}
type Todo {
	id: ID!
	createdAt: Time!
	visibilityStatus: VisibilityStatus!
	status: Status!
	priority: Int!
	text: String!
}
type User {
	id: ID!
	username: String!
	age: Float!
	amount: Float!
	role: Role!
}
`, printSchema(&ast.Schema{
		Types: types,
	}))
}
