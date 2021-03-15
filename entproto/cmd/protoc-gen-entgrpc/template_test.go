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
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
)

func TestTemplate(t *testing.T) {
	printer := &mockPrint{}
	err := printTemplate(printer, "func hello(ctx %(ctx) bool {}", tmplValues{
		"ctx": protogen.GoImportPath("context").Ident("Context"),
	})
	require.NoError(t, err)
	require.Len(t, printer.memory, 3)
	ctx := printer.memory[1]
	require.IsType(t, ctx, protogen.GoIdent{})

	err = printTemplate(printer, "hello %(c)", tmplValues{})
	require.EqualError(t, err, "entproto: could not find token \"%(c)\" in map")

	err = printTemplate(printer, "hello world %(not closing this", tmplValues{})
	require.EqualError(t, err, "entproto: corrupt template, must close parenthesis")
}

type mockPrint struct {
	memory []interface{}
}

func (m *mockPrint) P(i ...interface{}) {
	m.memory = append(m.memory, i...)
}
