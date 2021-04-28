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
	"go/printer"
	"go/token"
	"testing"

	"entgo.io/contrib/schemast/internal/mutatetest/ent/schema"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"github.com/stretchr/testify/require"
)

func TestFromEdgeDescriptor(t *testing.T) {
	tests := []struct {
		name           string
		edge           ent.Edge
		expected       string
		expectedErrMsg string
	}{
		{
			name:     "basic",
			edge:     edge.To("entity", Entity.Type),
			expected: `edge.To("entity", Entity.Type)`,
		},
		{
			name:     "inverse",
			edge:     edge.From("entity", Entity.Type).Ref("related"),
			expected: `edge.From("entity", Entity.Type).Ref("related")`,
		},
		{
			name:           "annotations",
			edge:           edge.To("entity", Entity.Type).Annotations(annotation("x")),
			expectedErrMsg: "schemast: unsupported feature: Annotations",
		},
		{
			name:     "required",
			edge:     edge.To("entity", Entity.Type).Required(),
			expected: `edge.To("entity", Entity.Type).Required()`,
		},
		{
			name:     "unique",
			edge:     edge.To("entity", Entity.Type).Unique(),
			expected: `edge.To("entity", Entity.Type).Unique()`,
		},
		{
			name:     "field",
			edge:     edge.To("entity", Entity.Type).Field("field"),
			expected: `edge.To("entity", Entity.Type).Field("field")`,
		},
		{
			name:     "struct_tag",
			edge:     edge.To("entity", Entity.Type).StructTag("tag"),
			expected: `edge.To("entity", Entity.Type).StructTag("tag")`,
		},
		{
			name:     "storage_key_one_col",
			edge:     edge.To("entity", Entity.Type).StorageKey(edge.Table("table"), edge.Column("to")),
			expected: `edge.To("entity", Entity.Type).StorageKey(edge.Table("table"), edge.Column("to"))`,
		},
		{
			name:     "storage_key_two_col",
			edge:     edge.To("entity", Entity.Type).StorageKey(edge.Table("table"), edge.Columns("to", "from")),
			expected: `edge.To("entity", Entity.Type).StorageKey(edge.Table("table"), edge.Columns("to", "from"))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Edge(tt.edge.Descriptor())
			if tt.expectedErrMsg != "" {
				require.EqualError(t, err, tt.expectedErrMsg)
				return
			}
			require.NoError(t, err)
			var buf bytes.Buffer
			fst := token.NewFileSet()
			err = printer.Fprint(&buf, fst, r)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, buf.String())
		})
	}
}

//nolint:golint,unused
type Entity struct {
	ent.Schema
}

func TestAppendEdge(t *testing.T) {
	tests := []struct {
		typeName     string
		expectedBody string
		expectedErr  string
	}{
		{
			typeName: "WithFields",
			expectedBody: `// Edges of the WithFields.
func (WithFields) Edges() []ent.Edge {
	return []ent.Edge{edge.To("owner", User.Type).Unique()}
}`,
		},
		{
			typeName:    "WithoutFields",
			expectedErr: `schemast: could not find method "Edges" for type "WithoutFields"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			ctx, err := Load("./internal/mutatetest/ent/schema")
			require.NoError(t, err)
			err = ctx.AppendEdge(tt.typeName, edge.To("owner", schema.User.Type).Unique().Descriptor())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			var buf bytes.Buffer
			method, _ := ctx.lookupMethod(tt.typeName, "Edges")
			err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedBody, buf.String())
		})
	}
}

func TestRemoveEdge(t *testing.T) {
	ctx, err := Load("./internal/mutatetest/ent/schema")
	require.NoError(t, err)
	err = ctx.RemoveEdge("WithModifiedField", "non_existent")
	require.EqualError(t, err, `schemast: could not find edge "non_existent" in type "WithModifiedField"`)
	err = ctx.RemoveEdge("WithModifiedField", "owner")
	require.NoError(t, err)

	var buf bytes.Buffer
	method, _ := ctx.lookupMethod("WithModifiedField", "Edges")
	err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
	require.NoError(t, err)
	require.EqualValues(t, `func (WithModifiedField) Edges() []ent.Edge {
	return []ent.Edge{}
}`, buf.String())
}
