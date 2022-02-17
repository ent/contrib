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
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/codegen/config"
)

func (e *EntGQL) MutateConfig(cfg *config.Config) error {
	idType, err := entgql.FindIDType(e.genTypes, e.graph.IDType)
	if err != nil {
		return err
	}
	// TODO: Add a warning for failure guess?
	if idType := guessTypeID(idType); idType != "" {
		cfg.Models["ID"] = config.TypeMapEntry{
			Model: []string{idType},
		}
	}

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

		for _, f := range obj.Fields {
			ann := &entgql.Annotation{}
			err := ann.Decode(f.Annotations[ann.Name()])
			if err != nil {
				return err
			}
			if ann.Skip {
				continue
			}

			goType := ""
			switch {
			case f.IsEnum():
				goType = fmt.Sprintf("%s/%s", e.graph.Package, f.Type.Ident)
			case f.IsOther():
				goType = fmt.Sprintf("%s.%s", f.Type.RType.PkgPath, f.Type.RType.Name)
			default:
				continue
			}

			name := strings.Title(f.Name)
			if ann.Type != "" {
				name = ann.Type
			}

			if !cfg.Models.Exists(name) {
				cfg.Models.Add(name, goType)
			}
		}

		if ann.RelayConnection {
			connection := fmt.Sprintf("%sConnection", obj.Name)
			if !cfg.Models.Exists(connection) {
				cfg.Models.Add(connection, e.entGoType(connection))
			}
			edge := fmt.Sprintf("%sEdge", obj.Name)
			if !cfg.Models.Exists(edge) {
				cfg.Models.Add(edge, e.entGoType(edge))
			}

			orderBy, err := hasOrderBy(obj)
			if err != nil {
				return err
			}
			if orderBy {
				order := fmt.Sprintf("%sOrder", obj.Name)
				cfg.Models[order] = config.TypeMapEntry{
					Model: []string{
						e.entGoType(order),
					},
				}
				cfg.Models[order+"Field"] = config.TypeMapEntry{
					Model: []string{
						e.entGoType(order + "Field"),
					},
				}
			}
		}
	}

	return nil
}

func (e *EntGQL) entGoType(name string) string {
	return fmt.Sprintf("%s.%s", e.graph.Package, name)
}

func guessTypeID(idType *field.TypeInfo) string {
	if idType == nil {
		return ""
	}

	t := idType.Type
	switch {
	case idType.RType != nil:
		return fmt.Sprintf("%s.%s", idType.RType.PkgPath, idType.RType.Name)
	case t == field.TypeInt, t == field.TypeInt32, t == field.TypeInt64:
		name := strings.TrimPrefix(t.ConstName(), "Type")
		return fmt.Sprintf("github.com/99designs/gqlgen/graphql.%sID", name)
	case t == field.TypeString:
		return "github.com/99designs/gqlgen/graphql.ID"
	}

	return ""
}
