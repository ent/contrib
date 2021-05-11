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
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEntBuiltins(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	e.entBuiltins()
	require.Equal(t, `enum OrderDirection {
	ASC
	DESC
}
`, e.print())
}

func TestEntOrderBy(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
		Nodes: []*gen.Type{
			{
				Name: "Excluded",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
			},
			{
				Name: "Included",
				Fields: []*gen.Field{
					{
						Name: "name",
						Type: &field.TypeInfo{
							Type: field.TypeString,
						},
						Annotations: map[string]interface{}{
							annotationName: map[string]interface{}{
								"order_field": "NAME",
							},
						},
					},
					{
						Name: "active",
						Type: &field.TypeInfo{
							Type: field.TypeString,
						},
						Annotations: map[string]interface{}{
							annotationName: map[string]interface{}{
								"order_field": "ACTIVE",
							},
						},
					},
					{
						Name: "gender",
						Type: &field.TypeInfo{
							Type: field.TypeString,
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	err = e.entOrderBy()
	require.NoError(t, err)
	require.Equal(t, `input IncludedOrder {
	direction: OrderDirection!
	field: IncludedOrderField!
}
enum IncludedOrderField {
	NAME
	ACTIVE
}
`, e.print())
}
