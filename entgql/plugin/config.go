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
	"fmt"
	"github.com/99designs/gqlgen/codegen/config"
)

func (e *Entgqlgen) MutateConfig(cfg *config.Config) error {
	if !cfg.Models.Exists(RelayPageInfo) {
		cfg.Models.Add(RelayPageInfo, e.entGoType(RelayPageInfo))
	}
	if !cfg.Models.Exists(RelayNode) {
		// Bind to Noder interface
		cfg.Models.Add(RelayNode, e.entGoType("Noder"))
	}
	if !cfg.Models.Exists(RelayCursor) {
		cfg.Models.Add(RelayCursor, e.entGoType(RelayCursor))
	}
	if !cfg.Models.Exists("OrderDirection") {
		cfg.Models.Add("OrderDirection", e.entGoType("OrderDirection"))
	}
	// Insert types
	for _, obj := range e.genTypes {
		ann := &entgql.Annotation{}
		err := ann.Decode(obj.Annotations[ann.Name()])
		if err != nil {
			return err
		}
		if ann.Skip {
			continue
		}
		if !cfg.Models.Exists(obj.Name) {
			cfg.Models.Add(obj.Name, e.entGoType(obj.Name))
		}
		if ann.RelayConnection {
			connection := fmt.Sprintf("%sConnection", obj.Name)
			edge := fmt.Sprintf("%sEdge", obj.Name)
			if !cfg.Models.Exists(connection) {
				cfg.Models.Add(connection, e.entGoType(connection))
			}
			if !cfg.Models.Exists(edge) {
				cfg.Models.Add(edge, e.entGoType(edge))
			}
			orderBy, err := hasOrderBy(obj)
			if err != nil {
				return err
			}
			if orderBy {
				order := fmt.Sprintf("%sOrder", obj.Name)
				cfg.Models[order] = config.TypeMapEntry{Model: []string{fmt.Sprintf("%s.%s", e.graph.Package, order)}}
				cfg.Models[order+"Field"] = config.TypeMapEntry{Model: []string{fmt.Sprintf("%s.%s", e.graph.Package, order+"Field")}}
			}
		}
	}
	return nil
}

func (e *Entgqlgen) entGoType(name string) string {
	return fmt.Sprintf("%s.%s", e.graph.Package, name)
}
