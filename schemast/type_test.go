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

package schemast

import (
	"bytes"
	"go/ast"
	"go/printer"
	"os"
	"path"
	"testing"

	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
)

func TestContext_AddType(t *testing.T) {
	ctx, err := Load("./internal/mutatetest/ent/schema")
	require.NoError(t, err)
	err = ctx.AddType("Cat")
	require.NoError(t, err)
	err = ctx.AppendField("Cat", field.String("name").Descriptor())
	require.NoError(t, err)

	var buf bytes.Buffer
	method, _ := ctx.lookupMethod("Cat", "Fields")
	err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
	require.NoError(t, err)
	require.EqualValues(t, `func (Cat) Fields() []ent.Field {
	return []ent.Field{field.String("name")}
}`, buf.String())

	file, _, ok := ctx.lookupTypeDecl("Cat")
	require.True(t, ok)
	require.Len(t, file.Imports, 2)
}

func TestContext_AddTypeWithBase(t *testing.T) {
	ctx, err := Load("./internal/printtest/ent/schema")
	require.NoError(t, err)
	err = ctx.AddTypeWithBase("Kitten", "mutatetest.TestBaseSchema", "entgo.io/contrib/schemast/internal/mutatetest")
	require.NoError(t, err)

	file, _, ok := ctx.lookupTypeDecl("Kitten")
	require.True(t, ok)
	require.Len(t, file.Decls, 5)
	require.Len(t, file.Imports, 3)

	structDef, ok := ctx.lookupBaseStruct("Kitten")
	require.True(t, ok)
	require.Len(t, structDef.Fields.List, 1)

	fieldDef, ok := structDef.Fields.List[0].Type.(*ast.SelectorExpr)
	require.True(t, ok)
	require.Equal(t, "mutatetest", fieldDef.X.(*ast.Ident).Name)
	require.Equal(t, "TestBaseSchema", fieldDef.Sel.Name)
}

func TestContext_RemoveType(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)
	err = tt.ctx.AddType("NewType")
	require.NoError(t, err)
	err = tt.ctx.RemoveType("Message")
	require.NoError(t, err)
	err = tt.ctx.RemoveType("NewType")
	require.NoError(t, err)
	err = tt.ctx.RemoveType("Nothing")
	require.EqualError(t, err, `schemast: type "Nothing" not found`)
	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	removed := tt.getType("Message")
	require.Nil(t, removed)
	removed = tt.getType("NewType")
	require.Nil(t, removed)

	file, err := os.ReadFile(path.Join(tt.schemaDir(), "user.go"))
	require.NoError(t, err)
	require.NotContains(t, string(file), "// Message holds the schema definition for the Message entity.")
}
