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
	"github.com/vektah/gqlparser/v2/ast"
)

// Directive to apply on the field/type
type Directive struct {
	Name      string              `json:"name,omitempty"`
	Arguments []DirectiveArgument `json:"arguments,omitempty"`
}

// DirectiveArgument return a GraphQL directive argument
type DirectiveArgument struct {
	Name  string        `json:"name,omitempty"`
	Value string        `json:"value,omitempty"`
	Kind  ast.ValueKind `json:"kind,omitempty"`
}

// NewDirective return a GraphQL directive
func NewDirective(name string, args ...DirectiveArgument) Directive {
	return Directive{
		Name:      name,
		Arguments: args,
	}
}

// DeprecatedDirective create `@deprecated` directive to apply on the field/type
func DeprecatedDirective(reason string) Directive {
	args := []DirectiveArgument{}
	if reason != "" {
		args = append(args, DirectiveArgument{
			Name:  "reason",
			Kind:  ast.StringValue,
			Value: reason,
		})
	}

	return NewDirective("deprecated", args...)
}
