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

package entoas

import (
	"os"
	"path/filepath"
	"testing"

	"entgo.io/contrib/entoas/serialization"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestView(t *testing.T) {
	t.Parallel()
	// Load a graph.
	wd, err := os.Getwd()
	require.NoError(t, err)
	g, err := entc.LoadGraph(filepath.Join(wd, "internal", "simple", "schema"), &gen.Config{})
	require.NoError(t, err)
	// Extract the Edges for a read operation on the Pet entity.
	var p *gen.Type
	var o *gen.Edge
	for _, n := range g.Nodes {
		if n.Name == "Pet" {
			p = n
			for _, e := range n.Edges {
				if e.Name == "owner" {
					o = e
					break
				}
			}
			break
		}
	}
	v, err := view(p, serialization.Groups{"test:view"})
	require.NoError(t, err)
	require.Equal(t, &View{
		Type:   p,
		Fields: []*gen.Field{p.ID},
		Edges:  []*gen.Edge{o},
	}, v)
}
