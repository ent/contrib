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
	"google.golang.org/protobuf/types/pluginpb"
)

func TestPrintTemplate(t *testing.T) {
	values := tmplValues{
		"world": "world",
		"ctx":   protogen.GoImportPath("context").Ident("Context"),
	}
	tests := []struct {
		tmpl             string
		expectedErr      string
		expectedContents string
	}{
		{
			tmpl:        "// %(missing key)",
			expectedErr: "could not find token \"%(missing key)\" in map",
		},
		{
			tmpl:             "// hello %(world)",
			expectedContents: "// hello world",
		},
		{
			tmpl:             "func c(ctx %(ctx)) {}",
			expectedContents: "func c(ctx context.Context)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tmpl, func(t *testing.T) {
			g, err := initGeneratedFile()
			require.NoError(t, err)
			err = printTemplate(g, tt.tmpl, values)
			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			}
			if tt.expectedContents != "" {
				bytes, err := g.Content()
				require.NoError(t, err)
				require.Contains(t, string(bytes), tt.expectedContents)
			}
		})
	}
}

func initGeneratedFile() (*protogen.GeneratedFile, error) {
	gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{})
	if err != nil {
		return nil, err
	}
	g := gen.NewGeneratedFile("foo.go", "golang.org/x/foo")
	g.P("package foo")
	return g, err
}
