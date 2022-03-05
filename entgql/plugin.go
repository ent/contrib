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
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

// CreatePlugin create the plugin for GQLGen
func (e *Extension) CreatePlugin() plugin.Plugin {
	return &gqlgenPlugin{
		schema: e.schema,
		models: e.models,
	}
}

type gqlgenPlugin struct {
	schema *ast.Schema
	models map[string]string
}

var (
	_ plugin.Plugin              = (*gqlgenPlugin)(nil)
	_ plugin.EarlySourceInjector = (*gqlgenPlugin)(nil)
	_ plugin.ConfigMutator       = (*gqlgenPlugin)(nil)
)

func (gqlgenPlugin) Name() string {
	return "entgql"
}

func (e *gqlgenPlugin) InjectSourceEarly() *ast.Source {
	if e.schema == nil {
		return nil
	}

	sb := &strings.Builder{}
	formatter.NewFormatter(sb).FormatSchema(e.schema)

	return &ast.Source{
		Name:    "entgql.graphql",
		Input:   sb.String(),
		BuiltIn: false,
	}
}

func (e *gqlgenPlugin) MutateConfig(cfg *config.Config) error {
	if e.models != nil {
		for name, goType := range e.models {
			if !cfg.Models.Exists(name) {
				cfg.Models.Add(name, goType)
			}
		}
	}
	return nil
}
