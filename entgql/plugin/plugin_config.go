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
	"fmt"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/codegen/config"
)

// MutateConfig implements the ConfigMutator interface
func (e *EntGQL) MutateConfig(cfg *config.Config) error {
	if e.relaySpec {
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
	}

	for _, node := range e.nodes {
		ant, err := entgql.DecodeAnnotation(node.Annotations)
		if err != nil {
			return err
		}
		if ant.Skip {
			continue
		}

		name := node.Name
		if ant.Type != "" {
			name = ant.Type
		}
		if !cfg.Models.Exists(name) {
			cfg.Models.Add(name, e.entGoType(node.Name))
		}
	}

	return nil
}

func (e *EntGQL) entGoType(name string) string {
	return fmt.Sprintf("%s.%s", e.graph.Package, name)
}
