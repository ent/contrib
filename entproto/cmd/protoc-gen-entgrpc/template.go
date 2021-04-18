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

package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	// match tokens like %(hello) in a string
	templateTokens = regexp.MustCompile(`(%\(.+?\))`)
)

// printTemplate is a utility function to make working with protogen have a more declarative interface.
// It receives a *protogen.GeneratedFile, a template string with placeholder formatted like "%(variableName)"
// and a tmplValues map containing the values that should be replaced when the template is rendered.
//
// Instead of:
//	g.P("func New", svcName, "(p string) *", protogen.GoImportPath("context").Ident("Context"))
// We can use
//	printTemplate(g, "func New%(svcName)(p string) *%(ctx)", tmplValues{
//		"svcName": "UserService",
//		"ctx": protogen.GoImportPath("context").Ident("Context"),
//	})
func printTemplate(g *protogen.GeneratedFile, tmpl string, values tmplValues) error {
	var errors error
	out := templateTokens.ReplaceAllStringFunc(tmpl, func(s string) string {
		val, err := values.retrieve(s)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		switch p := val.(type) {
		case protogen.GoIdent:
			return g.QualifiedGoIdent(p) // The P(..) magic for Go identifiers.
		case wrappedCalls:
			return renderWrappedCalls(g, p)
		default:
			return fmt.Sprint(p)
		}
	})
	if errors != nil {
		return errors
	}
	g.P(out)
	return nil
}

type tmplValues map[string]interface{}

func (t tmplValues) retrieve(token string) (interface{}, error) {
	k, ok := t[token[2:len(token)-1]] // from %(hello) => hello
	if !ok {
		return nil, fmt.Errorf("entproto: could not find token %q in map", token)
	}
	return k, nil
}

type wrappedCalls struct {
	invocations []interface{}
	arguments   []string
}

func renderWrappedCalls(g *protogen.GeneratedFile, chain wrappedCalls) string {
	var out string
	for _, inv := range chain.invocations {
		switch t := inv.(type) {
		case protogen.GoIdent:
			out += g.QualifiedGoIdent(t)
		default:
			out += fmt.Sprint(t)
		}
		out += "("
	}
	out += strings.Join(chain.arguments, ",")
	return out + strings.Repeat(")", len(chain.invocations))
}
