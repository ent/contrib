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

package entoas

import (
	"fmt"

	"entgo.io/contrib/entoas/serialization"
	"entgo.io/ent/entc/gen"
)

const maxDepth = 25

type (
	// Edge wraps a gen.Edge and denotes an edge to be returned in an operation response. It recursively defines
	// edges to load on the gen.Type the wrapped edge is pointing at.
	Edge struct {
		*gen.Edge
		Edges Edges
	}
	Edges []*Edge
	// walk is a node sequence in the schema graph. Used to keep track when creating an Edge.
	walk []string
)

// EdgeTree returns the Edges to include on a type for the given serialization groups.
func EdgeTree(n *gen.Type, gs serialization.Groups) (Edges, error) { return edgeTree(n, walk{}, gs) }

// Flatten returns a list of all gen.Types present in the tree.
func (es Edges) Flatten() []*gen.Edge {
	em := make(map[string]*gen.Edge)
	types(em, es)
 	r := make([]*gen.Edge, 0, len(em))
	for _, t := range em {
		r = append(r, t)
	}
	return r
}

// types recursively adds all gen.Types present on the given Edges to the given map.
func types(dest map[string]*gen.Edge, es Edges) {
	for _, e := range es {
		dest[e.Type.Name] = e.Edge
		types(dest, e.Edges)
	}
}

// edgeTree recursively collects the edges to load on this type for the requested groups.
func edgeTree(n *gen.Type, w walk, gs serialization.Groups) (Edges, error) {
	// If we have reached maxDepth there most possibly is an unwanted circular reference.
	if w.reachedMaxDepth() {
		return nil, fmt.Errorf("entoas: max depth of %d reached: ", maxDepth)
	}
	// Iterate over the edges of the given type.
	// If the type has an edge we need to eager load, do so.
	// Recursively go down the current types edges and, if requested, eager load those too.
	var es Edges
	for _, e := range n.Edges {
		a, err := EdgeAnnotation(e)
		if err != nil {
			return nil, err
		}
		if a.MaxDepth == 0 {
			a.MaxDepth = 1
		}
		// If the edge has at least one of the groups requested, load the edge.
		if a.Groups.Match(gs) {
			// Add the current step to our walk, since we will add this edge.
			w.push(n.Name + "." + e.Name)
			// If we have reached the max depth on this field for the given type stop the recursion. Backtrack!
			if w.cycleDepth() > a.MaxDepth {
				w.pop()
				continue
			}
			// Recursively collect the eager loads of edge-types edges.
			es1, err := edgeTree(e.Type, w, gs)
			if err != nil {
				return nil, err
			}
			// Done visiting this node. Remove this node from our walk.
			w.pop()
			es = append(es, &Edge{Edge: e, Edges: es1})
		}
	}
	return es, nil
}

// cycleDepth determines the length of a cycle on the last visited node.
//   <nil>: 0 -> no visits at all
// a->b->c: 1 -> 1st visit on c
// a->b->b: 2 -> 2nd visit on b
// a->a->a: 3 -> 3rd visit on a
// a->b->a: 2 -> 2nd visit on a
func (w walk) cycleDepth() uint {
	if len(w) == 0 {
		return 0
	}
	n := w[len(w)-1]
	c := uint(1)
	for i := len(w) - 2; i >= 0; i-- {
		if n == w[i] {
			c++
		}
	}
	return c
}

// reachedMaxDepth returns if the walk has reached a depth greater than maxDepth.
func (w walk) reachedMaxDepth() bool { return len(w) > maxDepth }

// push adds a new step to the walk.
func (w *walk) push(s string) { *w = append(*w, s) }

// pop removed the last step of the walk.
func (w *walk) pop() {
	if len(*w) > 0 {
		*w = (*w)[:len(*w)-1]
	}
}
