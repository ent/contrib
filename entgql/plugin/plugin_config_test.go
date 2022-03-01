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

func TestModifyConfig(t *testing.T) {
	e, err := NewEntGQLPlugin(&gen.Graph{
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
	})
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
