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

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
)

func TestModifyConfig_empty(t *testing.T) {
	e, err := newSchemaGenerator(&gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	require.NoError(t, err)
	e.relaySpec = false

	cfg, err := e.genModels()
	require.NoError(t, err)

	expected := map[string]string{}
	require.Equal(t, expected, cfg)
}

func createGraph(relayConnection bool) *gen.Graph {
	return &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
			IDType: &field.TypeInfo{
				Type: field.TypeInt,
			},
		},
		Nodes: []*gen.Type{
			{
				Name: "Todo",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"RelayConnection": relayConnection,
					},
				},
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
						"Skip": true,
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
						"RelayConnection": relayConnection,
					},
				},
			},
			{
				Name: "GroupWithSort",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
					Annotations: map[string]interface{}{
						annotationName: map[string]interface{}{
							"OrderField": "NAME",
						},
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"RelayConnection": relayConnection,
					},
				},
			},
		},
	}
}

func TestModifyConfig(t *testing.T) {
	e, err := newSchemaGenerator(createGraph(false))
	require.NoError(t, err)

	e.relaySpec = false
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Todo":          "example.com.Todo",
		"Group":         "example.com.Group",
		"GroupWithSort": "example.com.GroupWithSort",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_relay(t *testing.T) {
	g := createGraph(true)
	e, err := newSchemaGenerator(g)
	e.relaySpec = false
	require.NoError(t, err)
	_, err = e.genModels()
	require.Error(t, err, ErrRelaySpecDisabled)

	e.relaySpec = true
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Cursor":                  "example.com.Cursor",
		"Group":                   "example.com.Group",
		"GroupConnection":         "example.com.GroupConnection",
		"GroupEdge":               "example.com.GroupEdge",
		"GroupWithSort":           "example.com.GroupWithSort",
		"GroupWithSortConnection": "example.com.GroupWithSortConnection",
		"GroupWithSortEdge":       "example.com.GroupWithSortEdge",
		"GroupWithSortOrder":      "example.com.GroupWithSortOrder",
		"GroupWithSortOrderField": "example.com.GroupWithSortOrderField",
		"Node":                    "example.com.Noder",
		"OrderDirection":          "example.com.OrderDirection",
		"PageInfo":                "example.com.PageInfo",
		"Todo":                    "example.com.Todo",
		"TodoConnection":          "example.com.TodoConnection",
		"TodoEdge":                "example.com.TodoEdge",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	disableRelayConnection(graph)

	e, err := newSchemaGenerator(graph)
	require.NoError(t, err)
	e.relaySpec = false

	cfg, err := e.genModels()
	require.NoError(t, err)

	expected := map[string]string{
		"Category":         "entgo.io/contrib/entgql/internal/todoplugin/ent.Category",
		"CategoryStatus":   "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status",
		"CategoryConfig":   "entgo.io/contrib/entgql/internal/todo/ent/schema/schematype.CategoryConfig",
		"MasterUser":       "entgo.io/contrib/entgql/internal/todoplugin/ent.User",
		"Role":             "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role",
		"Status":           "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status",
		"Todo":             "entgo.io/contrib/entgql/internal/todoplugin/ent.Todo",
		"VisibilityStatus": "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)

	e, err := newSchemaGenerator(graph)
	require.NoError(t, err)
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Category":             "entgo.io/contrib/entgql/internal/todoplugin/ent.Category",
		"CategoryConfig":       "entgo.io/contrib/entgql/internal/todo/ent/schema/schematype.CategoryConfig",
		"CategoryConnection":   "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryConnection",
		"CategoryEdge":         "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryEdge",
		"CategoryOrder":        "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrder",
		"CategoryOrderField":   "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrderField",
		"CategoryStatus":       "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status",
		"Cursor":               "entgo.io/contrib/entgql/internal/todoplugin/ent.Cursor",
		"MasterUser":           "entgo.io/contrib/entgql/internal/todoplugin/ent.User",
		"MasterUserConnection": "entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserConnection",
		"MasterUserEdge":       "entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserEdge",
		"Node":                 "entgo.io/contrib/entgql/internal/todoplugin/ent.Noder",
		"OrderDirection":       "entgo.io/contrib/entgql/internal/todoplugin/ent.OrderDirection",
		"PageInfo":             "entgo.io/contrib/entgql/internal/todoplugin/ent.PageInfo",
		"Role":                 "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role",
		"Status":               "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status",
		"Todo":                 "entgo.io/contrib/entgql/internal/todoplugin/ent.Todo",
		"TodoConnection":       "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoConnection",
		"TodoEdge":             "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoEdge",
		"TodoOrder":            "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrder",
		"TodoOrderField":       "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrderField",
		"VisibilityStatus":     "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus",
	}
	require.Equal(t, expected, cfg)
}
