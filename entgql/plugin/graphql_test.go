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
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
func TestEnums(t *testing.T) {
	e := New(&gen.Graph{
		Config: &gen.Config{},
		Nodes: []*gen.Type{
			{
				Name: "User",
				Annotations: map[string]interface{}{
					"EntGQL": map[string]interface{}{
						"GenType": true,
					},
				},
				Fields: []*gen.Field{
					{
						Name: "Status",
						Type: &field.TypeInfo{
							Type: field.TypeEnum,
						},
						Enums: []gen.Enum{
							{
								Name:  "Pending",
								Value: "PENDING",
							}, {
								Name:  "Approved",
								Value: "APPROVED",
							},
						},
					},
				},
			},
		},
	})
	// TODO: .EnumValues is unavailable - maybe switch to Enums
	e.enums()
	require.Equal(t, e.print(), ``)
}
*/

func TestScalars(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	e.scalars()
	require.Equal(t, "scalar Time\n", e.print())
	e, err = New(&gen.Graph{
		Config: &gen.Config{
			Annotations: map[string]interface{}{
				"EntGQL": entgql.Annotation{
					GqlScalarMappings: map[string]string{
						"Date": "Date",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	e.scalars()
	require.Equal(t, "scalar Date\n", e.print())
	e, err = New(&gen.Graph{
		Config: &gen.Config{
			Annotations: map[string]interface{}{
				"EntGQL": entgql.Annotation{
					GqlScalarMappings: map[string]string{
						"Time":    "Time",
						"Int":     "Int",
						"Float":   "Float",
						"Boolean": "Boolean",
						"String":  "String",
						"ID":      "ID",
					},
				},
			},
		},
	})
	require.NoError(t, err)
	e.scalars()
	require.Equal(t, "scalar Time\n", e.print())
}

func TestTypes(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
		Nodes: []*gen.Type{
			{
				Name: "Todo",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
			},
			{
				Name: "User",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"skip": true,
					},
				},
			},
			{
				Name: "Group",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"relay_connection": true,
						"gql_implements":   []string{"SomeInterface"},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	err = e.types()
	require.NoError(t, err)
	require.Equal(t, `type Group implements Node & SomeInterface {
	name: String!
}
type GroupConnection {
	edges: [GroupEdge]
	pageInfo: PageInfo!
	totalCount: Int!
}
type GroupEdge {
	node: Group
	cursor: Cursor
}
type Todo implements Node {
	name: String!
}
`, e.print())
}
