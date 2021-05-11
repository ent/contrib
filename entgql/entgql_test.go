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
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConnections_empty(t *testing.T) {
	connections := getConnections(&gen.Graph{
		Nodes: []*gen.Type{},
	})
	require.Equal(t, connections, []string(nil))
}

func TestGetConnections(t *testing.T) {
	connections := getConnections(&gen.Graph{
		Nodes: []*gen.Type{
			{
				Name: "User",
				Annotations: map[string]interface{}{
					"EntGQL": map[string]interface{}{
						"RelayConnection": true,
					},
				},
			},
			{
				Name: "Todo",
				Annotations: map[string]interface{}{
					"EntGQL": map[string]interface{}{},
				},
			},
			{
				Name: "User_Todo",
			},
		},
	})
	require.Equal(t, connections, []string{"User"})
}

func TestModifyConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.Noder"},
	}
	require.Equal(t, cfg, expected)
}

func TestModifyConfig_autoBindPresent(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.AutoBind = append(cfg.AutoBind, "example.com")
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.Noder"},
	}
	require.Equal(t, cfg, expected)
}

func TestModifyConfig_noderPresent(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.CustomNoder"},
	}
	modifyConfig(cfg, &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	expected := config.DefaultConfig()
	expected.AutoBind = append(expected.AutoBind, "example.com")
	expected.Models["Node"] = config.TypeMapEntry{
		Model: []string{"example.com.CustomNoder"},
	}
	require.Equal(t, cfg, expected)
}
