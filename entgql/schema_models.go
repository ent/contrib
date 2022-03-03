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
	"fmt"
)

func (e *schemaGenerator) genModels() (map[string]string, error) {
	models := map[string]string{}

	if e.relaySpec {
		models[RelayPageInfo] = e.entGoType(RelayPageInfo)
		models[RelayNode] = e.entGoType("Noder")
		models[RelayCursor] = e.entGoType(RelayCursor)
	}
	for _, node := range e.nodes {
		ant, err := DecodeAnnotation(node.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip {
			continue
		}

		name := node.Name
		if ant.Type != "" {
			name = ant.Type
		}
		models[name] = e.entGoType(node.Name)

		hasOrderBy := false
		for _, field := range node.Fields {
			ant, err := DecodeAnnotation(field.Annotations)
			if err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}

			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				hasOrderBy = true
			}

			goType := ""
			switch {
			case field.IsOther() || (field.IsEnum() && field.HasGoType()):
				goType = fmt.Sprintf("%s.%s", field.Type.RType.PkgPath, field.Type.RType.Name)
			case field.IsEnum():
				goType = fmt.Sprintf("%s/%s", e.graph.Package, field.Type.Ident)
			default:
				continue
			}

			// NOTE(giautm): I'm not sure this is
			// the right approach, but it passed the test
			defs, err := e.typeFromField(field, false, ant.Type)
			if err != nil {
				return nil, err
			}
			name := defs.Name()

			models[name] = goType
		}

		// TODO(giautm): Added RelayConnection annotation check
		if e.relaySpec {
			pagination, err := NodePaginationNames(node)
			if err != nil {
				return nil, err
			}

			models[pagination.Connection] = e.entGoType(pagination.Connection)
			models[pagination.Edge] = e.entGoType(pagination.Edge)

			if hasOrderBy {
				models["OrderDirection"] = e.entGoType("OrderDirection")
				models[pagination.Order] = e.entGoType(pagination.Order)
				models[pagination.OrderField] = e.entGoType(pagination.OrderField)
			}
		}
	}

	return models, nil
}

func (e *schemaGenerator) entGoType(name string) string {
	return fmt.Sprintf("%s.%s", e.graph.Package, name)
}
