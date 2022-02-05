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

package entgql

import (
	"testing"

	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

var annotationName = Annotation{}.Name()

func TestFilterNodes(t *testing.T) {
	nodes, err := filterNodes([]*gen.Type{
		{
			Name: "Type1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name:   "Type2",
			Config: &gen.Config{},
		},
		{
			Name: "SkippedType",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": true},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []*gen.Type{
		{
			Name: "Type1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name:   "Type2",
			Config: &gen.Config{},
		},
	}, nodes)
}

func TestFilterEdges(t *testing.T) {
	edges, err := filterEdges([]*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{},
		},
		{
			Name: "SkippedEdge",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": true},
			},
		},
		{
			Name: "SkippedEdgeType",
			Type: &gen.Type{
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{"Skip": true},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{},
		},
	}, edges)
}

func TestFieldCollections(t *testing.T) {
	edges, err := fieldCollections([]*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{
				Name: "Todo",
			},
		},
		{
			Name: "BindDisabled",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"BindDisabled": true,
				},
			},
		},
		{
			Name: "EdgeMapping",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"BindDisabled": true,
					"Mapping":      []string{"field1", "field2"},
				},
			},
		},
	})
	require.NoError(t, err)

	require.Equal(t, map[string]fieldCollection{
		"Edge1": {
			Name:    "Todo",
			Mapping: []string{"Edge1"},
		},
		"Edge2": {
			Name:    "Todo",
			Mapping: []string{"Edge2"},
		},
		"EdgeMapping": {
			Name:    "Todo",
			Mapping: []string{"field1", "field2"},
		},
	}, edges)

	_, err = fieldCollections([]*gen.Edge{
		{
			Name: "EdgeInvalid",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"BindDisabled": false,
					"Mapping":      []string{"field1", "field2"},
				},
			},
		},
	})
	require.Errorf(t, err, "bind and mapping annotations are mutually exclusive")
}

func TestFilterFields(t *testing.T) {
	fields, err := filterFields([]*gen.Field{
		{
			Name: "Field1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Field2",
		},
		{
			Name: "SkippedField",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": true},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []*gen.Field{
		{
			Name: "Field1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Field2",
		},
	}, fields)
}
