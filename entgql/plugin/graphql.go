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
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

func (e *Entgqlgen) scalars() {
	for scalar := range e.scalarMappings {
		switch scalar {
		case "Int", "Float", "String", "Boolean", "ID":
			// Ignore builtins
		default:
			e.schema.Types[scalar] = &ast.Definition{
				Kind: ast.Scalar,
				Name: scalar,
			}
		}
	}
}

func (e *Entgqlgen) enums() error {
	enums := make(map[string][]string)
	for _, t := range e.genTypes {
		for _, f := range t.Fields {
			ann := &entgql.Annotation{}
			err := ann.Decode(f.Annotations[ann.Name()])
			if err != nil {
				return err
			}
			if ann.Skip {
				continue
			}
			if f.IsEnum() {
				enumName := strings.Title(f.Name)
				if ann.GqlType != "" {
					enumName = ann.GqlType
				}
				if values, ok := enums[enumName]; ok {
					if !unorderedEqual(values, f.EnumValues()) {
						panic(fmt.Errorf("enums are not equal. Name: %s, Values1: %s, Values: %s", enumName, values, f.EnumValues()))
					}
				} else {
					enums[enumName] = f.EnumValues()
				}
			}
		}
	}
	for name, values := range enums {
		var valueDefinitions ast.EnumValueList
		for _, v := range values {
			valueDefinitions = append(valueDefinitions, &ast.EnumValueDefinition{
				Name: v,
			})
		}
		e.insertDefinition(&ast.Definition{
			Name:       name,
			Kind:       ast.Enum,
			EnumValues: valueDefinitions,
		})
	}
	return nil
}

func unorderedEqual(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value] = true
	}
	for _, value := range second {
		if !exists[value] {
			return false
		}
	}
	return true
}

func (e *Entgqlgen) types() error {
	for _, t := range e.genTypes {
		// TODO: make relay config opt in
		interfaces := []string{"Node"}
		ann := &entgql.Annotation{}
		err := ann.Decode(t.Annotations[ann.Name()])
		if err != nil {
			return err
		}
		interfaces = append(interfaces, ann.GqlImplements...)
		fields, err := e.typeFields(t)
		if err != nil {
			return fmt.Errorf("type(%s): %w", t.Name, err)
		}
		e.insertDefinition(&ast.Definition{
			Name:       t.Name,
			Kind:       ast.Object,
			Fields:     fields,
			Interfaces: interfaces,
			Directives: e.directives(ann.GqlDirectives),
		})
		if ann.RelayConnection {
			e.relayConnection(t)
		}
	}
	return nil
}

func (e *Entgqlgen) directives(directives []entgql.Directive) ast.DirectiveList {
	var list ast.DirectiveList
	for _, d := range directives {
		var args ast.ArgumentList
		for _, arg := range d.Arguments {
			args = append(args, &ast.Argument{
				Name: arg.Name,
				Value: &ast.Value{
					Raw:  arg.Value,
					Kind: arg.Kind,
				},
			})
		}
		list = append(list, &ast.Directive{
			Name:      d.Name,
			Arguments: args,
		})
	}
	return list
}

func (e *Entgqlgen) insertDefinitions(defs []*ast.Definition) {
	for _, d := range defs {
		e.schema.Types[d.Name] = d
	}
}

func (e *Entgqlgen) insertDefinition(d *ast.Definition) {
	e.schema.Types[d.Name] = d
}
