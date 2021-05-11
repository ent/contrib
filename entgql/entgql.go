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
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

func Generate(cfg *config.Config, graph *gen.Graph) error {
	modifyConfig(cfg, graph)
	return api.Generate(cfg,
		api.AddPlugin(NewPlugin(getConnections(graph))),
	)
}

func getConnections(graph *gen.Graph) []string {
	var connections []string
	for _, n := range graph.Nodes {
		if ann, ok := n.Annotations["EntGQL"]; ok {
			entgqlAnn := ann.(map[string]interface{})
			if entgqlAnn["RelayConnection"] == true {
				connections = append(connections, n.Name)
			}
		}
	}
	return connections
}

func modifyConfig(cfg *config.Config, graph *gen.Graph) {
	autobindPresent := false
	for _, ab := range cfg.AutoBind {
		if ab == graph.Package {
			autobindPresent = true
		}
	}
	if !autobindPresent {
		cfg.AutoBind = append(cfg.AutoBind, graph.Package)
	}
	if !cfg.Models.Exists("Node") {
		cfg.Models.Add("Node", fmt.Sprintf("%s.Noder", graph.Package))
	}
}
