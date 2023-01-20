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

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestAnnotation(t *testing.T) {
	tests := []struct {
		name           string
		annot          schema.Annotation
		expected       string
		expectedOk     bool
		expectedErrMsg string
	}{
		{
			name:       "proto basic",
			annot:      entproto.Message(),
			expectedOk: true,
			expected:   `entproto.Message()`,
		},
		{
			name:       "proto skip",
			annot:      entproto.SkipGen(),
			expectedOk: true,
			expected:   `entproto.SkipGen()`,
		},
		{
			name:       "proto custom package",
			annot:      entproto.Message(entproto.PackageName("pkg")),
			expectedOk: true,
			expected:   `entproto.Message(entproto.PackageName("pkg"))`,
		},
		{
			name:       "proto service",
			annot:      entproto.Service(),
			expectedOk: true,
			expected:   `entproto.Service()`,
		},
		{
			name:       "proto field",
			annot:      entproto.Field(2),
			expectedOk: true,
			expected:   `entproto.Field(2)`,
		},
		{
			name:       "proto field with type",
			annot:      entproto.Field(2, entproto.Type(descriptorpb.FieldDescriptorProto_TYPE_UINT64)),
			expectedOk: true,
			expected:   `entproto.Field(2, entproto.Type(descriptorpb.FieldDescriptorProto_TYPE_UINT64))`,
		},
		{
			name:       "proto field with type name",
			annot:      entproto.Field(2, entproto.TypeName("TypeName")),
			expectedOk: true,
			expected:   `entproto.Field(2, entproto.TypeName("TypeName"))`,
		},
		{
			name: "proto enum",
			annot: entproto.Enum(map[string]int32{
				"unspecified": 0,
				"active":      1,
			}),
			expectedOk: true,
			expected:   `entproto.Enum(map[string]int32{"unspecified": 0, "active": 1})`,
		},
		{
			name: "entsql annotation table",
			annot: entsql.Annotation{
				Table: "table",
			},
			expectedOk: true,
			expected:   `entsql.Annotation{Table: "table"}`,
		},
		{
			name: "entsql annotation charset",
			annot: entsql.Annotation{
				Charset: "utf8mb4",
			},
			expectedOk: true,
			expected:   `entsql.Annotation{Charset: "utf8mb4"}`,
		},
		{
			name: "entsql annotation default",
			annot: entsql.Annotation{
				Default: "uuid_generate_v4()",
			},
			expectedOk: true,
			expected:   `entsql.Annotation{Default: "uuid_generate_v4()"}`,
		},
		{
			name: "entsql annotation collation",
			annot: entsql.Annotation{
				Collation: "utf8mb4_bin",
			},
			expectedOk: true,
			expected:   `entsql.Annotation{Collation: "utf8mb4_bin"}`,
		},
		{
			name: "entsql annotation size",
			annot: entsql.Annotation{
				Size: 128,
			},
			expectedOk: true,
			expected:   `entsql.Annotation{Size: 128}`,
		},
		{
			name: "entsql annotation on delete",
			annot: entsql.Annotation{
				OnDelete: entsql.NoAction,
			},
			expectedOk: true,
			expected:   `entsql.Annotation{OnDelete: entsql.NoAction}`,
		},
		{
			name: "entsql annotation unknown on delete",
			annot: entsql.Annotation{
				OnDelete: entsql.ReferenceOption("UNSUPPORTED"),
			},
			expectedOk:     false,
			expectedErrMsg: `schemast: unknown entsql ReferenceOption: "UNSUPPORTED"`,
		},
		{
			name:           "unsupported annotation",
			annot:          annotation("unsupported"),
			expectedErrMsg: `schemast: no Annotator configured for annotation "unsupported"`,
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, ok, err := Annotation(tt.annot)
			if tt.expectedErrMsg != "" {
				require.EqualError(t, err, tt.expectedErrMsg)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedOk, ok)
			var buf bytes.Buffer
			fst := token.NewFileSet()
			err = printer.Fprint(&buf, fst, r)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, buf.String())
		})
	}
}

func TestContext_AnnotateType(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)
	err = tt.ctx.AddType("NewType")
	require.NoError(t, err)
	err = tt.ctx.AppendTypeAnnotation("NewType", entproto.Message())
	require.NoError(t, err)

	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	nt := tt.getType("NewType")
	require.Len(t, nt.Annotations, 1)
	contents := tt.contents("new_type.go")
	require.Contains(t, contents, `func (NewType) Annotations() []schema.Annotation {
	return []schema.Annotation{entproto.Message()}
}`)
}

func TestContext_AnnotateTypeExisting(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)
	err = tt.ctx.AppendTypeAnnotation("Message", entproto.Message())
	require.NoError(t, err)
	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	nt := tt.getType("Message")
	require.Len(t, nt.Annotations, 1)
	contents := tt.contents("message.go")
	require.Contains(t, contents, `func (Message) Annotations() []schema.Annotation {
	return []schema.Annotation{entproto.Message()}
}`)
}
