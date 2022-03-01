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
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	// EntGQL is a plugin that generates GQL schema from the Ent's Graph
	EntGQL struct {
		graph *gen.Graph
		nodes []*gen.Type
	}

	// EntGQLPluginOption is a option for the EntGQL plugin
	EntGQLPluginOption func(*EntGQL) error
)

var (
	_ plugin.Plugin              = (*EntGQL)(nil)
	_ plugin.EarlySourceInjector = (*EntGQL)(nil)
)

// NewEntGQLPlugin creates a new EntGQL plugin
func NewEntGQLPlugin(graph *gen.Graph, opts ...EntGQLPluginOption) (*EntGQL, error) {
	nodes, err := entgql.FilterNodes(graph.Nodes)
	if err != nil {
		return nil, err
	}

	e := &EntGQL{
		graph: graph,
		nodes: nodes,
	}
	for _, opt := range opts {
		if err = opt(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

// Name implements the Plugin interface.
func (*EntGQL) Name() string {
	return "entgql"
}

// InjectSourceEarly implements the EarlySourceInjector interface.
func (e *EntGQL) InjectSourceEarly() *ast.Source {
	return nil
}
