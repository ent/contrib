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
	"encoding/json"
	"go/printer"
	"go/token"
	"net/http"
	"net/url"
	"testing"
	"time"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFromFieldDescriptor(t *testing.T) {
	tests := []struct {
		name           string
		field          ent.Field
		expected       string
		expectedErrMsg string
	}{
		{
			name:     "basic",
			field:    field.String("x"),
			expected: `field.String("x")`,
		},
		{
			name:     "optional",
			field:    field.String("x").Optional(),
			expected: `field.String("x").Optional()`,
		},
		{
			name:     "int64",
			field:    field.Int64("x"),
			expected: `field.Int64("x")`,
		},
		{
			name:           "unsupported type",
			field:          field.Bytes("unsupported"),
			expectedErrMsg: "schemast: unsupported type TypeBytes",
		},
		{
			name:     "json_ptr",
			field:    field.JSON("x", &SomeJSON{}).Optional(),
			expected: `field.JSON("x", &schemast.SomeJSON{}).Optional()`,
		},
		{
			name:     "json_primitive_slice",
			field:    field.JSON("x", []string{}),
			expected: `field.JSON("x", []string{})`,
		},
		{
			name:     "json_type_slice",
			field:    field.JSON("x", []SomeJSON{}),
			expected: `field.JSON("x", []schemast.SomeJSON{})`,
		},
		{
			name:     "json_pointer_type_slice",
			field:    field.JSON("x", []*url.Values{}),
			expected: `field.JSON("x", []*url.Values{})`,
		},
		{
			name:     "json_map_primitive_types",
			field:    field.JSON("x", map[string]string{}),
			expected: `field.JSON("x", map[string]string{})`,
		},
		{
			name:     "json_map_primitive_to_type",
			field:    field.JSON("x", map[string]url.Values{}),
			expected: `field.JSON("x", map[string]url.Values{})`,
		},
		{
			name:     "json_map_primitive_to_pointer_type",
			field:    field.JSON("x", map[string]*url.Values{}),
			expected: `field.JSON("x", map[string]*url.Values{})`,
		},
		{
			name:     "json_map_primitive_to_interface",
			field:    field.JSON("x", map[string]interface{}{}),
			expected: `field.JSON("x", map[string]interface {}{})`,
		},
		{
			name:     "json_ptr_imported_type",
			field:    field.JSON("x", &url.URL{}),
			expected: `field.JSON("x", &url.URL{})`,
		},
		{
			name:     "json_slice_imported_type",
			field:    field.JSON("x", []http.Dir{}).Comment("some_comment"),
			expected: `field.JSON("x", []http.Dir{}).Comment("some_comment")`,
		},
		{
			name:     "json_type_alias",
			field:    field.JSON("x", json.RawMessage{}),
			expected: `field.JSON("x", json.RawMessage{})`,
		},
		{
			name:     "time",
			field:    field.Time("time").Default(time.Now),
			expected: `field.Time("time").Default(time.Now)`,
		},
		{
			name: "time anonymous",
			field: field.Time("time").Default(func() time.Time {
				return time.Time{}
			}),
			expectedErrMsg: "schemast: only selector exprs are supported for default func",
		},
		{
			name:     "struct tag",
			field:    field.String("x").StructTag(`j:"hi"`),
			expected: `field.String("x").StructTag("j:\"hi\"")`,
		},
		{
			name:           "enums:values",
			field:          field.Enum("x").Values("a", "b"),
			expected:       `field.Enum("x").Values("a", "b")`,
			expectedErrMsg: "",
		},
		{
			name:     "enums:named values",
			field:    field.Enum("x").NamedValues("a", "b"),
			expected: `field.Enum("x").NamedValues("a", "b")`,
		},
		{
			name:     "storage key",
			field:    field.String("x").StorageKey("s"),
			expected: `field.String("x").StorageKey("s")`,
		},
		{
			name: "schema type",
			field: field.String("x").SchemaType(map[string]string{
				dialect.SQLite: "VARCHAR",
			}),
			expected: `field.String("x").SchemaType(map[string]string{"sqlite3": "VARCHAR"})`,
		},
		{
			name:     "annotations",
			field:    field.String("x").Annotations(entproto.Message()),
			expected: `field.String("x").Annotations(entproto.Message())`,
		},
		{
			name:     "default:string",
			field:    field.String("x").Default("x"),
			expected: `field.String("x").Default("x")`,
		},
		{
			name:     "default:int",
			field:    field.Int("x").Default(1),
			expected: `field.Int("x").Default(1)`,
		},
		{
			name:     "default:uint64",
			field:    field.Uint64("x").Default(1),
			expected: `field.Uint64("x").Default(1)`,
		},
		{
			name:     "default:float32",
			field:    field.Float32("x").Default(3.14),
			expected: `field.Float32("x").Default(3.14)`,
		},
		{
			name:     "default:bool",
			field:    field.Bool("x").Default(true),
			expected: `field.Bool("x").Default(true)`,
		},
		{
			name: "unsupported validator",
			field: field.String("x").Validate(func(s string) error {
				return nil
			}),
			expectedErrMsg: "schemast: unsupported feature Descriptor.Validators",
		},
		{
			name:     "bytes",
			field:    field.Bytes("x"),
			expected: `field.Bytes("x")`,
		},
		{
			name:     "uuid",
			field:    field.UUID("x", uuid.UUID{}),
			expected: `field.UUID("x", uuid.UUID{})`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Field(tt.field.Descriptor())
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

type (
	SomeJSON   struct{}
	annotation string
)

func (a annotation) Name() string { return string(a) }

func TestAppendField(t *testing.T) {
	tests := []struct {
		typeName     string
		expectedBody string
		expectedErr  string
	}{
		{
			typeName: "WithFields",
			expectedBody: `// Fields of the WithFields.
func (WithFields) Fields() []ent.Field {
	return []ent.Field{
		field.String("existing"), field.String("newField"),
	}
}`,
		},
		{
			typeName: "WithNilFields",
			expectedBody: `// Fields of the WithNilFields.
func (WithNilFields) Fields() []ent.Field {
	return []ent.Field{field.String("newField")}
}`,
		},
		{
			typeName: "WithoutFields",
			expectedBody: `func (WithoutFields) Fields() []ent.Field {
	return []ent.Field{field.String("newField")}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			ctx, err := Load("./internal/mutatetest/ent/schema")
			require.NoError(t, err)
			err = ctx.AppendField(tt.typeName, field.String("newField").Descriptor())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			var buf bytes.Buffer
			method, _ := ctx.lookupMethod(tt.typeName, "Fields")
			err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedBody, buf.String())
		})
	}
}

func TestRemoveField(t *testing.T) {
	ctx, err := Load("./internal/mutatetest/ent/schema")
	require.NoError(t, err)
	err = ctx.RemoveField("WithModifiedField", "non_existent")
	require.EqualError(t, err, `schemast: could not find field "non_existent" in type "WithModifiedField"`)
	err = ctx.RemoveField("WithModifiedField", "name")
	require.NoError(t, err)

	var buf bytes.Buffer
	method, _ := ctx.lookupMethod("WithModifiedField", "Fields")
	err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
	require.NoError(t, err)
	require.EqualValues(t, `func (WithModifiedField) Fields() []ent.Field {
	return []ent.Field{}
}`, buf.String())
}
