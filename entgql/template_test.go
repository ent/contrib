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
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
	}, SkipType)
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
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
		{
			Name: "SkippedEdgeType",
			Type: &gen.Type{
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{"Skip": SkipAll},
				},
			},
		},
	}, SkipType)
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
	edges := []*gen.Edge{
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
			Name: "two_words",
			Type: &gen.Type{
				Name: "Todo",
			},
		},
		{
			Name: "Unbind",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"Unbind": true,
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
					"Unbind":  true,
					"Mapping": []string{"field1", "field2"},
				},
			},
		},
	}
	collect, err := fieldCollections(edges)
	require.NoError(t, err)
	require.Equal(t, []*fieldCollection{
		{
			Edge:    edges[0],
			Mapping: []string{"edge1"},
		},
		{
			Edge:    edges[1],
			Mapping: []string{"edge2"},
		},
		{
			Edge:    edges[2],
			Mapping: []string{"twoWords"},
		},
		{
			Edge:    edges[4],
			Mapping: []string{"field1", "field2"},
		},
	}, collect)

	_, err = fieldCollections([]*gen.Edge{
		{
			Name: "EdgeInvalid",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"Unbind":  false,
					"Mapping": []string{"field1", "field2"},
				},
			},
		},
	})
	require.Errorf(t, err, "bind and mapping annotations are mutually exclusive")
}

func TestFieldCollectorCases(t *testing.T) {
	t.Parallel()
	nameField := &gen.Field{
		Name: "name",
	}
	firstNameField := &gen.Field{
		Name: "first_name",
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"CollectedFor": []string{"fullName"},
			},
		},
	}
	lastNameField := &gen.Field{
		Name: "last_name",
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"CollectedFor": []string{"fullName"},
			},
		},
	}
	nameWithCollectedFor := &gen.Field{
		Name: "name",
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"CollectedFor": []string{"uppercaseName"},
			},
		},
	}
	categoryField := &gen.Field{
		Name: "category_id",
		Annotations: map[string]interface{}{
			annotationName: map[string]interface{}{
				"Unbind":  true,
				"Mapping": []string{"categoryID", "category_id", "categoryX"},
			},
		},
	}

	cases, err := fieldCollectorCases([]*gen.Field{
		nameWithCollectedFor,
		categoryField,
	})
	require.NoError(t, err)
	require.Equal(t, []*fieldCollectorCase{
		{
			Mapping: []string{"categoryID", "categoryX", "category_id"},
			Fields:  []*gen.Field{categoryField},
		},
		{
			Mapping: []string{"name", "uppercaseName"},
			Fields:  []*gen.Field{nameWithCollectedFor},
		},
	}, cases)

	multiCases, err := fieldCollectorCases([]*gen.Field{firstNameField, lastNameField})
	require.NoError(t, err)
	require.Equal(t, []*fieldCollectorCase{
		{
			Mapping: []string{"firstName"},
			Fields:  []*gen.Field{firstNameField},
		},
		{
			Mapping: []string{"fullName"},
			Fields:  []*gen.Field{firstNameField, lastNameField},
		},
		{
			Mapping: []string{"lastName"},
			Fields:  []*gen.Field{lastNameField},
		},
	}, multiCases)

	_, err = fieldCollectorCases([]*gen.Field{nameField})
	require.NoError(t, err)
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
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
	}, SkipType)
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
