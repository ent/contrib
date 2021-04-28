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
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestBasic(t *testing.T) {
	tt := newGenTest(t, "testdata/basic.proto")
	contents, err := tt.fileContents("basic.go")
	require.NoError(t, err)
	require.Contains(t, contents, "type Basic struct")
	require.Contains(t, contents, `field.String("name")`)
	_, err = tt.fileContents("skipped.go")
	require.EqualError(t, err, `file "skipped.go" not generated`)
	require.Len(t, tt.output, 1)
}

func TestCustomName(t *testing.T) {
	tt := newGenTest(t, "testdata/custom_name.proto")
	contents, err := tt.fileContents("rotemtam.go")
	require.NoError(t, err)
	require.Contains(t, contents, "type Rotemtam struct")
}

func TestFieldModifier(t *testing.T) {
	tt := newGenTest(t, "testdata/fields.proto")
	contents, err := tt.fileContents("pet.go")
	require.NoError(t, err)
	require.Contains(t, contents, "type Pet struct")
	require.Contains(t, contents, `field.String("name").Optional().StorageKey("shem")`)
}

func TestEdges_O2M(t *testing.T) {
	tt := newGenTest(t, "testdata/edges.proto")
	catContents, err := tt.fileContents("cat.go")
	require.NoError(t, err)
	require.Contains(t, catContents, `edge.To("owner", Human.Type)`)
	humanContents, err := tt.fileContents("human.go")
	require.NoError(t, err)
	require.Contains(t, humanContents, `edge.From("cats", Cat.Type)`)
}

func TestEdges_M2M(t *testing.T) {
	tt := newGenTest(t, "testdata/edges.proto")
	articleContents, err := tt.fileContents("article.go")
	require.NoError(t, err)
	require.Contains(t, articleContents, `edge.To("categories", Category.Type).StorageKey(edge.Table("table"), edge.Columns("a", "b"))`)
	categoryContents, err := tt.fileContents("category.go")
	require.NoError(t, err)
	require.Contains(t, categoryContents, `edge.From("articles", Article.Type)`)
}

func TestEnum(t *testing.T) {
	tt := newGenTest(t, "testdata/enums.proto")
	contents, err := tt.fileContents("job.go")
	require.NoError(t, err)
	require.Contains(t, contents, `field.Enum("priority").Values("PRIORITY_UNSPECIFIED", "LOW", "HIGH")`)
	require.Contains(t, contents, `field.Enum("status").Values("STATUS_UNSPECIFIED", "PENDING", "ACTIVE", "COMPLETE", "FAILED")`)
}

type genTest struct {
	output map[string]string
}

func newGenTest(t *testing.T, files ...string) *genTest {
	tmp, err := ioutil.TempDir("", "protoc-gen-ent-")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmp)
	})
	var parser protoparse.Parser
	var descs []*descriptorpb.FileDescriptorProto
	tgts := []string{"google/protobuf/descriptor.proto", "options/ent/opts.proto"}
	tgts = append(tgts, files...)
	parsed, err := parser.ParseFiles(tgts...)
	require.NoError(t, err)
	for _, p := range parsed {
		descs = append(descs, p.AsFileDescriptorProto())
	}
	gen, err := protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate:  files,
		Parameter:       nil,
		ProtoFile:       descs,
		CompilerVersion: nil,
	})
	require.NoError(t, err)
	err = printSchemas(tmp, gen)
	require.NoError(t, err)
	output := make(map[string]string)
	filepath.Walk(tmp, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		contents, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		output[filepath.Base(path)] = string(contents)
		return nil
	})
	return &genTest{output: output}
}

func (g *genTest) fileContents(name string) (string, error) {
	contents, ok := g.output[name]
	if !ok {
		return "", fmt.Errorf("file %q not generated", name)
	}
	return contents, nil
}
