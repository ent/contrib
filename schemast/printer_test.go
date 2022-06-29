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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)

	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	require.Len(t, tt.graph.Nodes, 2)
}

func TestPrintWithAdded(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)

	typeName := "Dog"
	require.NoError(t, tt.ctx.AddType(typeName))
	require.NoError(t, tt.ctx.AppendField(typeName, field.String("name").Descriptor()))

	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	require.Len(t, tt.graph.Nodes, 3)
	g := tt.getType(typeName)
	require.NotNil(t, g, "expected to find a type named Dog")
	require.Len(t, g.Fields, 1)
	require.EqualValues(t, "name", g.Fields[0].Name)
}

func TestPrintAppended(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)

	require.NoError(t, tt.ctx.AppendField("Message", field.String("title").Optional().Descriptor()))
	require.NoError(t, tt.ctx.AppendField("Message", field.String("author").Optional().Descriptor()))
	require.NoError(t, tt.ctx.AppendIndex("Message", index.Fields("title", "author")))
	require.NoError(t, tt.print())
	require.NoError(t, tt.load())
	g := tt.getType("Message")
	require.NotNil(t, g, "expected to find a type named Message")
	require.Len(t, g.Fields, 2)
	require.Len(t, g.Indexes, 1)
	title := g.Fields[0]
	require.EqualValues(t, "title", title.Name)
	require.EqualValues(t, true, title.Optional)
}

func TestPrintHeaderComment(t *testing.T) {
	tt, err := newPrintTest(t)
	commentRegexp := regexp.MustCompile("(?m)^// File updated by test.$")

	require.NoError(t, err)
	require.NoError(t, tt.print(Header("File updated by test.")))
	contents, err := os.ReadFile(filepath.Join(tt.schemaDir(), "message.go")) // A file that didn't have the header.
	require.NoError(t, err)
	require.Regexp(t, commentRegexp, string(contents))

	contents, err = os.ReadFile(filepath.Join(tt.schemaDir(), "user.go")) // A file that had the header, but not on the first line.
	require.NoError(t, err)
	matches := commentRegexp.FindAllString(string(contents), -1)
	require.Len(t, matches, 1)
}

func TestPrintAddImport(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)

	require.NoError(t, tt.ctx.AppendField("Message", field.UUID("uuid", uuid.UUID{}).Descriptor()))
	require.NoError(t, tt.ctx.AppendField("Message", field.UUID("hash", uuid.UUID{}).Descriptor()))
	require.NoError(t, tt.print())

	contents, err := os.ReadFile(filepath.Join(tt.schemaDir(), "message.go"))
	require.NoError(t, err)
	matches := strings.Count(string(contents), "github.com/google/uuid")
	require.Equal(t, matches, 1)
}

func TestPrintStructOnSameLine(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)

	require.NoError(t, tt.ctx.AppendField("Message", field.String("name").Descriptor()))
	require.NoError(t, tt.ctx.AppendField("Message", field.JSON("json", struct{}{}).Descriptor()))
	require.NoError(t, tt.print())

	contents, err := os.ReadFile(filepath.Join(tt.schemaDir(), "message.go"))
	require.NoError(t, err)
	require.Contains(t, string(contents), "struct{}{}")
}

func newPrintTest(t *testing.T) (*printTest, error) {
	dir, err := os.MkdirTemp(".", "printtest-")
	if err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	p := &printTest{dir: abs, t: t}
	if err := os.MkdirAll(p.schemaDir(), 0700); err != nil {
		return p, err
	}
	p.ctx, err = Load("./internal/printtest/ent/schema")
	if err != nil {
		return p, err
	}
	return p, nil
}

type printTest struct {
	dir   string
	t     *testing.T
	ctx   *Context
	graph *gen.Graph
}

func (p *printTest) schemaDir() string {
	return filepath.Join(p.dir, "ent", "schema")
}

func (p *printTest) print(opts ...PrintOption) error {
	return p.ctx.Print(p.schemaDir(), opts...)
}

func (p *printTest) load() error {
	graph, err := entc.LoadGraph(p.schemaDir(), &gen.Config{})
	if err != nil {
		return err
	}
	p.graph = graph
	return nil
}

func (p *printTest) getType(name string) *gen.Type {
	for _, t := range p.graph.Nodes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (p *printTest) contents(fname string) string {
	file, err := os.ReadFile(filepath.Join(p.schemaDir(), fname))
	require.NoError(p.t, err)
	return string(file)
}
