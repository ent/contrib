// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schemast

import "go/ast"

// kind stores configuration relevant to processing AST for different kinds of ent schema attributes
// (fields, edges, annotation, etc.).
type kind struct {
	// methodName is the name of the method on the schema that returns this kind.
	// For example, the "Fields" method returns the list of fields that a schema has.
	methodName string

	// ifaceSelector is the selector expression representing the type that is returned by the method.
	// For example, the Fields method returns a slice of "ent.Field".
	ifaceSelector *ast.SelectorExpr
}

var (
	kindEdge = kind{
		methodName:    "Edges",
		ifaceSelector: selectorLit("ent", "Edge"),
	}
	kindField = kind{
		methodName:    "Fields",
		ifaceSelector: selectorLit("ent", "Field"),
	}
	kindAnnot = kind{
		methodName:    "Annotations",
		ifaceSelector: selectorLit("schema", "Annotation"),
	}
	kindIndex = kind{
		methodName:    "Indexes",
		ifaceSelector: selectorLit("ent", "Index"),
	}
)
