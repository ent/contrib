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
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
)

func TestModifyConfig_empty(t *testing.T) {
	e, err := NewEntGQLPlugin(&gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	require.NoError(t, err)
	cfg := config.DefaultConfig()
	err = e.MutateConfig(cfg)
	require.NoError(t, err)
	expected := config.DefaultConfig()
	expected.Models = map[string]config.TypeMapEntry{}
	require.Equal(t, expected, cfg)
}

var g = &gen.Graph{
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
					"RelayConnection": true,
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
					"RelayConnection": true,
				},
			},
		},
	},
}

func TestModifyConfig(t *testing.T) {
	e, err := NewEntGQLPlugin(g)
	require.NoError(t, err)
	cfg := config.DefaultConfig()
	err = e.MutateConfig(cfg)
	require.NoError(t, err)
	expected := config.DefaultConfig()
	expected.Models = map[string]config.TypeMapEntry{
		"Todo":          {Model: []string{"example.com.Todo"}},
		"Group":         {Model: []string{"example.com.Group"}},
		"GroupWithSort": {Model: []string{"example.com.GroupWithSort"}},
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_relay(t *testing.T) {
	e, err := NewEntGQLPlugin(g, WithRelaySpecification(true))
	require.NoError(t, err)
	cfg := config.DefaultConfig()
	err = e.MutateConfig(cfg)
	require.NoError(t, err)
	expected := config.DefaultConfig()
	expected.Models = map[string]config.TypeMapEntry{
		"Cursor":                  {Model: []string{"example.com.Cursor"}},
		"Group":                   {Model: []string{"example.com.Group"}},
		"GroupConnection":         {Model: []string{"example.com.GroupConnection"}},
		"GroupEdge":               {Model: []string{"example.com.GroupEdge"}},
		"GroupWithSort":           {Model: []string{"example.com.GroupWithSort"}},
		"GroupWithSortConnection": {Model: []string{"example.com.GroupWithSortConnection"}},
		"GroupWithSortEdge":       {Model: []string{"example.com.GroupWithSortEdge"}},
		"GroupWithSortOrder":      {Model: []string{"example.com.GroupWithSortOrder"}},
		"GroupWithSortOrderField": {Model: []string{"example.com.GroupWithSortOrderField"}},
		"Node":                    {Model: []string{"example.com.Noder"}},
		"OrderDirection":          {Model: []string{"example.com.OrderDirection"}},
		"PageInfo":                {Model: []string{"example.com.PageInfo"}},
		"Todo":                    {Model: []string{"example.com.Todo"}},
		"TodoConnection":          {Model: []string{"example.com.TodoConnection"}},
		"TodoEdge":                {Model: []string{"example.com.TodoEdge"}},
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)

	e, err := NewEntGQLPlugin(graph)
	require.NoError(t, err)
	cfg := config.DefaultConfig()
	err = e.MutateConfig(cfg)
	require.NoError(t, err)
	expected := config.DefaultConfig()
	expected.Models = map[string]config.TypeMapEntry{
		"Category":         {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Category"}},
		"CategoryStatus":   {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status"}},
		"CategoryConfig":   {Model: []string{"entgo.io/contrib/entgql/internal/todo/ent/schema/schematype.CategoryConfig"}},
		"MasterUser":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.User"}},
		"Role":             {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role"}},
		"Status":           {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status"}},
		"Todo":             {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Todo"}},
		"VisibilityStatus": {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus"}},
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)

	e, err := NewEntGQLPlugin(graph, WithRelaySpecification(true))
	require.NoError(t, err)
	cfg := config.DefaultConfig()
	err = e.MutateConfig(cfg)
	require.NoError(t, err)
	expected := config.DefaultConfig()
	expected.Models = map[string]config.TypeMapEntry{
		"Category":             {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Category"}},
		"CategoryConfig":       {Model: []string{"entgo.io/contrib/entgql/internal/todo/ent/schema/schematype.CategoryConfig"}},
		"CategoryConnection":   {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryConnection"}},
		"CategoryEdge":         {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryEdge"}},
		"CategoryOrder":        {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrder"}},
		"CategoryOrderField":   {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrderField"}},
		"CategoryStatus":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status"}},
		"Cursor":               {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Cursor"}},
		"MasterUser":           {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.User"}},
		"MasterUserConnection": {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserConnection"}},
		"MasterUserEdge":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserEdge"}},
		"Node":                 {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Noder"}},
		"OrderDirection":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.OrderDirection"}},
		"PageInfo":             {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.PageInfo"}},
		"Role":                 {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role"}},
		"Status":               {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status"}},
		"Todo":                 {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.Todo"}},
		"TodoConnection":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.TodoConnection"}},
		"TodoEdge":             {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.TodoEdge"}},
		"TodoOrder":            {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrder"}},
		"TodoOrderField":       {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrderField"}},
		"VisibilityStatus":     {Model: []string{"entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus"}},
	}
	require.Equal(t, expected, cfg)
}
