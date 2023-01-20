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

	"entgo.io/ent"
	"entgo.io/ent/schema/index"
	"github.com/stretchr/testify/require"
)

func TestFromIndexDescriptor(t *testing.T) {
	tests := []struct {
		name           string
		index          ent.Index
		expected       string
		expectedErrMsg string
	}{
		{
			name:     "basic",
			index:    index.Fields("cat_id"),
			expected: `index.Fields("cat_id")`,
		},
		{
			name:     "multi",
			index:    index.Fields("cat_id", "dog_id"),
			expected: `index.Fields("cat_id", "dog_id")`,
		},
		{
			name:     "unique",
			index:    index.Fields("cat_id").Unique(),
			expected: `index.Fields("cat_id").Unique()`,
		},
		{
			name:     "storage key",
			index:    index.Fields("cat_id").StorageKey("skey"),
			expected: `index.Fields("cat_id").StorageKey("skey")`,
		},
		{
			name:     "single edge",
			index:    index.Fields("cat_id").Edges("edge_id"),
			expected: `index.Fields("cat_id").Edges("edge_id")`,
		},
		{
			name:     "multi edge",
			index:    index.Fields("cat_id").Edges("edge", "other_edge"),
			expected: `index.Fields("cat_id").Edges("edge", "other_edge")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Index(tt.index.Descriptor())
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

func TestAppendIndex(t *testing.T) {
	tests := []struct {
		typeName     string
		expectedBody string
		expectedErr  string
	}{
		{
			typeName: "WithFields",
			expectedBody: `// Indexes of the WithFields.
func (WithFields) Indexes() []ent.Index {
	return []ent.Index{index.Fields("a", "b").Unique()}
}`,
		},
		{
			typeName: "WithoutFields",
			expectedBody: `func (WithoutFields) Indexes() []ent.Index {
	return []ent.Index{index.Fields("a", "b").Unique()}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			ctx, err := Load("./internal/mutatetest/ent/schema")
			require.NoError(t, err)
			err = ctx.AppendIndex(tt.typeName, index.Fields("a", "b").Unique())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			var buf bytes.Buffer
			method, _ := ctx.lookupMethod(tt.typeName, "Indexes")
			err = printer.Fprint(&buf, ctx.SchemaPackage.Fset, method)
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedBody, buf.String())
		})
	}
}
