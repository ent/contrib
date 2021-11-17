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
	"strings"

	"entgo.io/contrib/entoas/serialization"
	"entgo.io/ent/entc/gen"
)

// A View is a subset of a gen.Type. It may hold fewer Fields and Edges than the gen.Type it is derived from.
type View struct {
	Type   *gen.Type
	Fields []*gen.Field
	Edges  []*gen.Edge
}

// Views returns all views that are needed to fill the OAS.
func Views(g *gen.Graph) (map[string]*View, error) {
	m := make(map[string]*View)
	for _, n := range g.Nodes {
		ops, err := NodeOperations(n)
		if err != nil {
			return nil, err
		}
		// For every operation add a schema to use.
		for _, op := range ops {
			// Skip the delete operation (of course).
			if op == OpDelete {
				continue
			}
			gs, err := groupsForOperation(n.Annotations, op)
			if err != nil {
				return nil, err
			}
			v, err := view(n, gs)
			if err != nil {
				return nil, err
			}
			vn, err := viewName(n, op)
			if err != nil {
				return nil, err
			}
			m[vn] = v
			// Collect the "tree" of edges to load on this node and operation.
			es, err := EdgeTree(n, gs)
			if err != nil {
				return nil, err
			}
			// Flatten the tree. Create a view for every type involved.
			for _, e := range es.Flatten() {
				v, err := view(e.Type, gs)
				if err != nil {
					return nil, err
				}
				evn, err := viewNameEdge(vn, e)
				if err != nil {
					return nil, err
				}
				m[evn] = v
			}
		}
		// Look at the edges and the operations exposed on them. Do create the views too.
		for _, e := range n.Edges {
			ops, err := EdgeOperations(e)
			if err != nil {
				return nil, err
			}
			// For every operation add a schema to use.
			for _, op := range ops {
				// Skip the delete operation (of course).
				if op == OpDelete {
					continue
				}
				gs, err := groupsForOperation(e.Annotations, op)
				if err != nil {
					return nil, err
				}
				v, err := view(e.Type, gs)
				if err != nil {
					return nil, err
				}
				vn, err := edgeViewName(n, e, op)
				if err != nil {
					return nil, err
				}
				m[vn] = v
				// Collect the "tree" of edges to load on this edge and operation.
				es, err := EdgeTree(n, gs)
				if err != nil {
					return nil, err
				}
				// Flatten the tree. Create a view for every type involved.
				for _, t := range es.Flatten() {
					v, err := view(t.Type, gs)
					if err != nil {
						return nil, err
					}
					evn, err := viewNameEdge(vn, t)
					if err != nil {
						return nil, err
					}
					m[evn] = v
				}
			}
		}
	}
	return m, nil
}

// view creates a new view of the given type when serialized with the given groups.
func view(n *gen.Type, gs serialization.Groups) (*View, error) {
	v := &View{Type: n}
	for _, f := range append([]*gen.Field{n.ID}, n.Fields...) {
		ok, err := serializeField(f, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			v.Fields = append(v.Fields, f)
		}
	}
	for _, e := range n.Edges {
		ok, err := serializeEdge(e, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			v.Edges = append(v.Edges, e)
		}
	}
	return v, nil
}

// serializeField checks if a gen.Field is to be serialized for the requested groups.
func serializeField(f *gen.Field, g serialization.Groups) (bool, error) {
	// If the field is sensitive, don't serialize it.
	if f.Sensitive() {
		return false, nil
	}
	// If no groups are requested or the field has no groups defined render the field.
	if f.Annotations == nil || len(g) == 0 {
		return true, nil
	}
	// Extract the Groups defined on the edge.
	ant, err := FieldAnnotation(f)
	if err != nil {
		return false, err
	}
	// If no groups are given on the field default is to include it in the output.
	if len(ant.Groups) == 0 {
		return true, nil
	}
	// If there are groups given check if the groups match the requested ones.
	return g.Match(ant.Groups), nil
}

// serializeEdge checks if an edge is to be serialized according to its annotations and the requested groups.
func serializeEdge(e *gen.Edge, g serialization.Groups) (bool, error) {
	// If no groups are requested or the edge has no groups defined do not render the edge.
	if e.Annotations == nil || len(g) == 0 {
		return false, nil
	}
	// Extract the Groups defined on the edge.
	ant, err := EdgeAnnotation(e)
	if err != nil {
		return false, err
	}
	// If no groups are given on the edge default is to exclude it.
	if len(ant.Groups) == 0 {
		return false, nil
	}
	// If there are groups given check if the groups match the requested ones.
	return g.Match(ant.Groups), nil
}

// groupsForOperation returns the requested groups as defined on the given Annotations for the Operation.
func groupsForOperation(a gen.Annotations, op Operation) (serialization.Groups, error) {
	// If there are no annotations given do not load any groups.
	ant := &Annotation{}
	if a == nil || a[ant.Name()] == nil {
		return nil, nil
	}
	// Decode the types annotation and extract the groups requested for the given operation.
	if err := ant.Decode(a[ant.Name()]); err != nil {
		return nil, err
	}
	switch op {
	case OpCreate:
		return ant.Create.Groups, nil
	case OpRead:
		return ant.Read.Groups, nil
	case OpUpdate:
		return ant.Update.Groups, nil
	case OpList:
		return ant.List.Groups, nil
	}
	return nil, fmt.Errorf("unknown operation %q", op)
}

// viewName returns the name for a view for a given operation on a gen.Type.
func viewName(n *gen.Type, op Operation) (string, error) {
	cfg, err := GetConfig(n.Config)
	if err != nil {
		return "", err
	}
	if cfg.SimpleModels {
		return n.Name, nil
	}
	return n.Name + op.Title(), nil
}

// edgeViewName returns the name for a view for a given 2nd leve operation on a gen.Edge.
func edgeViewName(n *gen.Type, e *gen.Edge, op Operation) (string, error) {
	cfg, err := GetConfig(n.Config)
	if err != nil {
		return "", err
	}
	if cfg.SimpleModels {
		return e.Type.Name, nil
	}
	return fmt.Sprintf("%s_%s%s", n.Name, e.StructField(), op.Title()), nil
}

// viewNameEdge returns the name for a view that is an edge on the given node.
func viewNameEdge(vn string, e *gen.Edge) (string, error) {
	cfg, err := GetConfig(e.Type.Config)
	if err != nil {
		return "", err
	}
	if cfg.SimpleModels {
		return e.Type.Name, nil
	}
	return fmt.Sprintf("%s_%s", vn, e.StructField()), nil
}

// Title returns the title cases variant of the operation.
func (op Operation) Title() string { return strings.Title(string(op)) }
