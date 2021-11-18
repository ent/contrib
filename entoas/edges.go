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
	"entgo.io/contrib/entoas/serialization"
	"entgo.io/ent/entc/gen"
)

type (
	// Edge wraps a gen.Edge and denotes an edge to be returned in an operation response. It recursively defines
	// edges to load on the gen.Type the wrapped edge is pointing at.
	Edge struct {
		*gen.Edge
		Edges Edges
	}
	Edges []*Edge
	// A step when traversing the schema graph.
	step struct {
		from *gen.Type
		over *gen.Edge
	}
	// walk is a sequence of steps.
	walk []step
)

// EdgeTree returns the Edges to include on a type for the given serialization groups.
func EdgeTree(n *gen.Type, gs serialization.Groups) (Edges, error) { return edgeTree(n, walk{}, gs) }

// Flatten returns a list of all gen.Edge present in the tree.
func (es Edges) Flatten() []*gen.Edge {
	var r []*gen.Edge
	for _, t := range edges(es) {
		r = append(r, t)
	}
	return r
}

// edges recursively adds all gen.Edge present on the given Edges to the given map.
func edges(es Edges) map[string]*gen.Edge {
	m := make(map[string]*gen.Edge)
	for _, e := range es {
		m[e.Name] = e.Edge
		for k, v := range edges(e.Edges) {
			m[k] = v
		}
	}
	return m
}

// edgeTree recursively collects the edges to load on this type for the requested groups.
func edgeTree(n *gen.Type, w walk, gs serialization.Groups) (Edges, error) {
	// Iterate over the edges of the given type.
	// If the type has an edge we need to eager load, do so.
	// Recursively go down the current types edges and, if requested, eager load those too.
	var es Edges
	for _, e := range n.Edges {
		a, err := EdgeAnnotation(e)
		if err != nil {
			return nil, err
		}
		// If the edge has at least one of the groups requested, load the edge.
		if a.Groups.Match(gs) {
			s := step{n, e}
			// If we already visited this edge before don't do it again to prevent an endless cycle.
			if w.visited(s) {
				continue
			}
			w.push(s)
			// Recursively collect the eager loads of edge-types edges.
			es1, err := edgeTree(e.Type, w, gs)
			if err != nil {
				return nil, err
			}
			// Done visiting this edge.
			w.pop()
			es = append(es, &Edge{Edge: e, Edges: es1})
		}
	}
	return es, nil
}

// visited returns if the given step has been done before.
func (w walk) visited(s step) bool {
	if len(w) == 0 {
		return false
	}
	for i := len(w) - 1; i >= 0; i-- {
		if w[i].equal(s) {
			return true
		}
	}
	return false
}

// push adds a new step to the walk.
func (w *walk) push(s step) { *w = append(*w, s) }

// pop removes the last step of the walk.
func (w *walk) pop() {
	if len(*w) > 0 {
		*w = (*w)[:len(*w)-1]
	}
}

// equal returns if the given step o is equal to the current step s.
func (s step) equal(o step) bool { return s.from.Name == o.from.Name && s.over.Name == o.over.Name }
