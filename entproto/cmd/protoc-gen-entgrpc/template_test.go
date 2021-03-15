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
