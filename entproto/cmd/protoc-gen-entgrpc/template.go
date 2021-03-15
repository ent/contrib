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

import "fmt"

// printTemplate is a utility function to make working with protogen have a more declarative interface.
// It receives a protogenPrinter (in practice a *protogen.GeneratedFile), a template string with placeholder
// formatted like "%(variableName)" and a tmplValues map containing the values that should be replaced when
// the template is rendered.
//
// Instead of:
//	g.P("func ", svcName, "(", paramName, " string)")
// We can use
//	printTemplate(g, "func %(svcName)(%(paramName) string)
func printTemplate(printer protogenPrinter, template string, values tmplValues) error {
	var inToken bool
	var buf string
	var output []interface{}
	for _, c := range template {
		str := string(c)
		if inToken {
			if len(buf) == 1 && str != "(" {
				return fmt.Errorf("entproto: corrupt template, percent must be followed by left parenthesis")
			}
			buf += str
			if str == ")" {
				inToken = false
				val, err := values.retrieve(buf)
				if err != nil {
					return err
				}
				output = append(output, val)
				buf = ""
			}
		} else {
			if str == "%" {
				inToken = true
				output = append(output, buf)
				buf = str
			} else {
				buf += str
			}
		}
	}
	if inToken {
		return fmt.Errorf("entproto: corrupt template, must close parenthesis")
	}
	output = append(output, buf)
	printer.P(output...)
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

type protogenPrinter interface {
	P(...interface{})
}
