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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"entgo.io/contrib/entoas/spec"
	"entgo.io/ent/entc/gen"
	"github.com/go-openapi/inflect"
	"github.com/stoewer/go-strcase"
)

type Operation string

const (
	OpCreate Operation = "create"
	OpRead   Operation = "read"
	OpUpdate Operation = "update"
	OpDelete Operation = "delete"
	OpList   Operation = "list"
)

func generate(g *gen.Graph, spec *spec.Spec) error {
	// Add all schemas.
	if err := schemas(g, spec); err != nil {
		return err
	}
	// Add error responses.
	errorResponses(spec)
	// Add all paths.
	return paths(g, spec)
}

// schemas adds schemas for every node to the spec.
func schemas(g *gen.Graph, s *spec.Spec) error {
	// Loop over every defined node and add it to the spec.
	for _, n := range g.Nodes {
		// Create the schema.
		s.Components.Schemas[n.Name] = &spec.Schema{
			Name:   n.Name,
			Fields: make(spec.Fields, len(n.Fields)+1),
			Edges:  make(spec.Edges, len(n.Edges)),
		}
		// Add all fields and the ID.
		for _, f := range append([]*gen.Field{n.ID}, n.Fields...) {
			sf, err := field(f)
			if err != nil {
				return err
			}
			s.Components.Schemas[n.Name].Fields[f.Name] = sf
		}
	}
	// Loop over every node once more to add the edges.
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			s.Components.Schemas[n.Name].Edges[e.Name] = &spec.Edge{
				Ref:    s.Components.Schemas[e.Type.Name],
				Unique: e.Unique,
			}
		}
	}
	// If the SimpleModels feature is enabled to not generate a schema per response.
	cfg, err := GetConfig(g.Config)
	if err != nil {
		return err
	}
	if !cfg.SimpleModels {
		// Add all the views for the paths to the schemas.
		vs, err := Views(g)
		if err != nil {
			return err
		}
		for n, v := range vs {
			// Create the schema.
			s.Components.Schemas[n] = &spec.Schema{
				Name:   n,
				Fields: make(spec.Fields, len(v.Fields)),
				Edges:  make(spec.Edges, len(v.Edges)),
			}
			// Add all fields (ID is already part of them).
			for _, f := range v.Fields {
				sf, err := field(f)
				if err != nil {
					return err
				}
				s.Components.Schemas[n].Fields[f.Name] = sf
			}
		}
		// Loop over every view once more to add the edges.
		for n, v := range vs {
			for _, e := range v.Edges {
				vn, err := viewNameEdge(strings.Split(n, "_")[0], e)
				if err != nil {
					return err
				}
				if _, ok := s.Components.Schemas[vn]; !ok {
					return fmt.Errorf("entoas: view %q does not exist", vn)
				}
				s.Components.Schemas[n].Edges[e.Name] = &spec.Edge{
					Ref:    s.Components.Schemas[vn],
					Unique: e.Unique,
				}
			}
		}
	}
	return nil
}

// errResponses adds all responses to the spec responses.
func errorResponses(s *spec.Spec) {
	for c, d := range map[int]string{
		http.StatusBadRequest:          "invalid input, data invalid",
		http.StatusConflict:            "conflicting resources",
		http.StatusForbidden:           "insufficient permissions",
		http.StatusInternalServerError: "unexpected error",
		http.StatusNotFound:            "resource not found",
	} {
		s.Components.Responses[strconv.Itoa(c)] = &spec.Response{
			Name:        strconv.Itoa(c),
			Description: d,
			Headers:     nil, // TODO
			Content: spec.Content{
				spec.JSON: &spec.MediaTypeObject{
					Unique: true,
					Schema: spec.Schema{
						Fields: map[string]*spec.Field{
							"code": {
								Type:    _int32,
								Unique:  true,
								Example: c,
							},
							"status": {
								Type:    _string,
								Unique:  true,
								Example: http.StatusText(c),
							},
						},
						Edges: map[string]*spec.Edge{
							"errors": {
								Schema: new(spec.Schema),
								Unique: true,
							},
						},
					},
				},
			},
		}
	}
}

var rules = inflect.NewDefaultRuleset()

// paths adds all operations to the spec paths.
func paths(g *gen.Graph, s *spec.Spec) error {
	for _, n := range g.Nodes {
		// Add schema operations.
		ops, err := NodeOperations(n)
		if err != nil {
			return err
		}
		// root for all operations on this node.
		root := "/" + rules.Pluralize(strcase.KebabCase(n.Name))
		// Create operation.
		if contains(ops, OpCreate) {
			path(s, root).Post, err = createOp(s, n)
			if err != nil {
				return err
			}
		}
		// Read operation.
		if contains(ops, OpRead) {
			path(s, root+"/{id}").Get, err = readOp(s, n)
			if err != nil {
				return err
			}
		}
		// Update operation.
		if contains(ops, OpUpdate) {
			path(s, root+"/{id}").Patch, err = updateOp(s, n)
			if err != nil {
				return err
			}
		}
		// Delete operation.
		if contains(ops, OpDelete) {
			path(s, root+"/{id}").Delete, err = deleteOp(s, n)
			if err != nil {
				return err
			}
		}
		// List operation.
		if contains(ops, OpList) {
			path(s, root).Get, err = listOp(s, n)
			if err != nil {
				return err
			}
		}
		// Sub-Resource operations.
		for _, e := range n.Edges {
			subRoot := root + "/{id}/" + strcase.KebabCase(e.Name)
			ops, err := EdgeOperations(e)
			if err != nil {
				return err
			}
			// Create operation.
			if contains(ops, OpCreate) {
				path(s, subRoot).Post, err = createEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			}
			// Read operation.
			if contains(ops, OpRead) {
				path(s, subRoot).Get, err = readEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			}
			// Delete operation.
			if contains(ops, OpDelete) {
				path(s, subRoot).Delete, err = deleteEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			}
			// List operation.
			if contains(ops, OpList) {
				path(s, subRoot).Get, err = listEdgeOp(s, n, e)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// path returns the correct spec.Path for the given root. Creates and sets a fresh instance if non does yet exist.
func path(s *spec.Spec, root string) *spec.Path {
	if s.Paths == nil {
		s.Paths = make(map[string]*spec.Path)
	}
	if _, ok := s.Paths[root]; !ok {
		s.Paths[root] = new(spec.Path)
	}
	return s.Paths[root]
}

// createOp returns the spec description for a create operation on the given node.
func createOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	// ant, err := schemaAnnotation(n)
	// if err != nil {
	// 	return nil, err
	// }
	req, err := reqBody(n, OpCreate)
	if err != nil {
		return nil, err
	}
	vn, err := viewName(n, OpCreate)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Create a new %s", n.Name),
		Description: fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: string(OpCreate) + n.Name,
		RequestBody: req,
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: &spec.Response{
					Description: fmt.Sprintf("%s created", n.Name),
					Headers:     nil, // TODO
					Content: spec.Content{
						spec.JSON: &spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[vn],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest): {
				Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)],
			},
			strconv.Itoa(http.StatusInternalServerError): {
				Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)],
			},
		},
		// Security: ant.CreateSecurity,
	}, nil
}

// createEdgeOp returns the spec description for a create operation on a subresource.
func createEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	// Create a basic create operation as if this was a first level operation.
	op, err := createOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	// But now alter the fields required to make this a second level operation.
	op.Summary = fmt.Sprintf("Create a new %s and attach it to the %s", e.Type.Name, n.Name)
	op.Description = fmt.Sprintf("Creates a new %s and attaches it to the %s", e.Type.Name, n.Name)
	op.Tags = []string{n.Name}
	op.OperationID = string(OpCreate) + n.Name + strcase.UpperCamelCase(e.Name)
	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
	rp.Description = fmt.Sprintf("%s created and attached to the %s", e.Type.Name, n.Name)
	vn, err := edgeViewName(n, e, OpCreate)
	if err != nil {
		return nil, err
	}
	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	op.Parameters = []*spec.Parameter{id}
	return op, nil
}

// readOp returns a spec.OperationConfig for a read operation on the given node.
func readOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := viewName(n, OpRead)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Find a %s by ID", n.Name),
		Description: fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name),
		Tags:        []string{n.Name},
		OperationID: string(OpRead) + n.Name,
		Parameters:  []*spec.Parameter{id},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: &spec.Response{
					Description: fmt.Sprintf("%s with requested ID was found", n.Name),
					Headers:     nil, // TODO
					Content: spec.Content{
						spec.JSON: &spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[vn],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
		// Security: ant.ReadSecurity,
	}, nil
}

// readEdgeOp returns the spec description for a read operation on a subresource.
func readEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	if !e.Unique {
		return nil, errors.New("read operations are not allowed on non unique edges")
	}
	// Create a basic read operation as if this was a first level operation.
	op, err := readOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	// But now alter the fields required to make this a second level operation.
	op.Summary = fmt.Sprintf("Find the attached %s", e.Type.Name)
	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", e.Type.Name, n.Name)
	op.Tags = []string{n.Name}
	op.OperationID = string(OpRead) + n.Name + strcase.UpperCamelCase(e.Name)
	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
	rp.Description = fmt.Sprintf("%s attached to %s with requested ID was found", e.Type.Name, n.Name)
	vn, err := edgeViewName(n, e, OpRead)
	if err != nil {
		return nil, err
	}
	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	op.Parameters = []*spec.Parameter{id}
	return op, nil
}

// updateOp returns a spec.OperationConfig for an update operation on the given node.
func updateOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	req, err := reqBody(n, OpUpdate)
	if err != nil {
		return nil, err
	}
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := viewName(n, OpUpdate)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Updates a %s", n.Name),
		Description: fmt.Sprintf("Updates a %s and persists changes to storage.", n.Name),
		Tags:        []string{n.Name},
		OperationID: string(OpUpdate) + n.Name,
		Parameters:  []*spec.Parameter{id},
		RequestBody: req,
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: &spec.Response{
					Description: fmt.Sprintf("%s updated", n.Name),
					Headers:     nil, // TODO
					Content: spec.Content{
						spec.JSON: &spec.MediaTypeObject{
							Unique: true,
							Ref:    s.Components.Schemas[vn],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
	}, nil
}

// deleteOp returns a spec.OperationConfig for a delete operation on the given node.
func deleteOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("Deletes a %s by ID", n.Name),
		Description: fmt.Sprintf("Deletes the %s with the requested ID.", n.Name),
		Tags:        []string{n.Name},
		OperationID: string(OpDelete) + n.Name,
		Parameters:  []*spec.Parameter{id},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusNoContent): {
				Response: &spec.Response{
					Description: fmt.Sprintf("%s with requested ID was deleted", n.Name),
					Headers:     nil, // TODO
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
	}, nil
}

// deleteEdgeOp returns the spec description for a delete operation on a subresource.
func deleteEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	if !e.Unique {
		return nil, errors.New("delete operations are not allowed on non unique edges")
	}
	// Create a basic delete operation as if this was a first level operation.
	op, err := deleteOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	// But now alter the fields required to make this a second level operation.
	op.Summary = fmt.Sprintf("Delete the attached %s", strcase.UpperCamelCase(e.Name))
	op.Description = fmt.Sprintf(
		"Delete the attached %s of the %s with the given ID", strcase.UpperCamelCase(e.Name), n.Name,
	)
	op.Tags = []string{n.Name}
	op.OperationID = string(OpDelete) + n.Name + strcase.UpperCamelCase(e.Name)
	op.Responses[strconv.Itoa(http.StatusNoContent)].Response.Description = fmt.Sprintf(
		"%s with requested ID was deleted", strcase.UpperCamelCase(e.Name),
	)
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	op.Parameters = []*spec.Parameter{id}
	return op, nil
}

// listOp returns a spec.OperationConfig for a list operation on the given node.
func listOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
	vn, err := viewName(n, OpList)
	if err != nil {
		return nil, err
	}
	return &spec.Operation{
		Summary:     fmt.Sprintf("List %s", rules.Pluralize(n.Name)),
		Description: fmt.Sprintf("List %s.", rules.Pluralize(n.Name)),
		Tags:        []string{n.Name},
		OperationID: string(OpList) + n.Name,
		Parameters: []*spec.Parameter{{
			Name:        "page",
			In:          spec.InQuery,
			Description: "what page to render",
			Schema:      _int32,
		}, {
			Name:        "itemsPerPage",
			In:          spec.InQuery,
			Description: "item count to render per page",
			Schema:      _int32,
		}},
		Responses: map[string]*spec.OperationResponse{
			strconv.Itoa(http.StatusOK): {
				Response: &spec.Response{
					Description: fmt.Sprintf("result %s list", n.Name),
					Headers:     nil, // TODO
					Content: spec.Content{
						spec.JSON: &spec.MediaTypeObject{
							Ref: s.Components.Schemas[vn],
						},
					},
				},
			},
			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
		},
	}, nil
}

// listEdgeOp returns the spec description for a read operation on a subresource.
func listEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
	if e.Unique {
		return nil, errors.New("list operations are not allowed on unique edges")
	}
	// Create a basic read operation as if this was a first level operation.
	op, err := listOp(s, e.Type)
	if err != nil {
		return nil, err
	}
	// But now alter the fields required to make this a second level operation.
	op.Summary = fmt.Sprintf("List attached %s", rules.Pluralize(strcase.UpperCamelCase(e.Name)))
	op.Description = fmt.Sprintf("List attached %s.", rules.Pluralize(strcase.UpperCamelCase(e.Name)))
	op.Tags = []string{n.Name}
	op.OperationID = string(OpList) + n.Name + strcase.UpperCamelCase(e.Name)
	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
	rp.Description = fmt.Sprintf("result %s list", rules.Pluralize(strcase.UpperCamelCase(n.Name)))
	vn, err := edgeViewName(n, e, OpList)
	if err != nil {
		return nil, err
	}
	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	op.Parameters = []*spec.Parameter{id}
	return op, nil
}

// field created a spec schema field from an ent schema field.
func field(f *gen.Field) (*spec.Field, error) {
	t, err := oasType(f)
	if err != nil {
		return nil, err
	}
	ex, err := exampleValue(f)
	if err != nil {
		return nil, err
	}
	return &spec.Field{Unique: true, Required: !f.Optional, Type: t, Example: ex}, nil
}

// NodeOperations returns the list of operations to expose for this node.
func NodeOperations(n *gen.Type) ([]Operation, error) {
	c, err := GetConfig(n.Config)
	if err != nil {
		return nil, err
	}
	ant := &Annotation{}
	// If no policies are given follow the global policy.
	if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
		if c.DefaultPolicy == PolicyExpose {
			return []Operation{OpCreate, OpRead, OpUpdate, OpDelete, OpList}, nil
		}
		return nil, nil
	} else {
		// An operation gets exposed if it is either
		// - annotated with PolicyExpose or
		// - not annotated with PolicyExclude and the DefaultPolicy is PolicyExpose.
		if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
			return nil, err
		}
		var ops []Operation
		for op, opn := range map[Operation]OperationConfig{
			OpCreate: ant.Create,
			OpRead:   ant.Read,
			OpUpdate: ant.Update,
			OpDelete: ant.Delete,
			OpList:   ant.List,
		} {
			// If the operation is explicitly annotated to be exposed do so.
			if opn.Policy == PolicyExpose || (opn.Policy == PolicyNone && c.DefaultPolicy == PolicyExpose) {
				ops = append(ops, op)
				continue
			}
		}
		return ops, nil
	}
}

// EdgeOperations returns the list of operations to expose for this edge.
func EdgeOperations(e *gen.Edge) ([]Operation, error) {
	c, err := GetConfig(e.Type.Config)
	if err != nil {
		return nil, err
	}
	ant := &Annotation{}
	// If no policies are given follow the global policy.
	if e.Annotations == nil || e.Annotations[ant.Name()] == nil {
		if c.DefaultPolicy == PolicyExpose {
			if e.Unique {
				return []Operation{OpCreate, OpRead, OpDelete}, nil
			} else {
				return []Operation{OpCreate, OpList}, nil
			}
		}
		return nil, nil
	} else {
		// An edge-operation gets exposed if it is either
		// - annotated with PolicyExpose or
		// - not annotated with PolicyExclude and the DefaultPolicy is PolicyExpose.
		if err := ant.Decode(e.Annotations[ant.Name()]); err != nil {
			return nil, err
		}
		var ops []Operation
		m := make(map[Operation]OperationConfig)
		if e.Unique {
			m[OpCreate] = ant.Create
			m[OpRead] = ant.Read
			m[OpDelete] = ant.Delete
		} else {
			m[OpCreate] = ant.Create
			m[OpList] = ant.List
		}
		for op, opn := range m {
			if opn.Policy == PolicyExpose || (opn.Policy == PolicyNone && c.DefaultPolicy == PolicyExpose) {
				ops = append(ops, op)
				continue
			}
		}
		return ops, nil
	}
}

// reqBody returns the request body for the given node and operation.
func reqBody(n *gen.Type, op Operation) (*spec.RequestBody, error) {
	req := &spec.RequestBody{}
	switch op {
	case OpCreate:
		req.Description = fmt.Sprintf("%s to create", n.Name)
	case OpUpdate:
		req.Description = fmt.Sprintf("%s properties to update", n.Name)
	default:
		return nil, fmt.Errorf("requestBody: unsupported operation %q", op)
	}
	fs := make(spec.Fields)
	for _, f := range n.Fields {
		if op == OpCreate || !f.Immutable {
			sf, err := field(f)
			if err != nil {
				return nil, err
			}
			fs[f.Name] = sf
		}
	}
	for _, e := range n.Edges {
		t, err := oasType(e.Type.ID)
		if err != nil {
			return nil, err
		}
		fs[e.Name] = &spec.Field{
			Unique:   e.Unique,
			Required: !e.Optional,
			Type:     t,
			Example:  nil, // TODO: Example for a unique / non-unique edge
		}
	}
	req.Content = spec.Content{
		spec.JSON: &spec.MediaTypeObject{
			Unique: true,
			Schema: spec.Schema{
				Name:   n.Name + op.Title() + "Request",
				Fields: fs,
			},
		},
	}
	return req, nil
}

var (
	_empty    = &spec.Type{}
	_int32    = &spec.Type{Type: "integer", Format: "int32"}
	_int64    = &spec.Type{Type: "integer", Format: "int64"}
	_float    = &spec.Type{Type: "number", Format: "float"}
	_double   = &spec.Type{Type: "number", Format: "double"}
	_string   = &spec.Type{Type: "string"}
	_bytes    = &spec.Type{Type: "string", Format: "byte"}
	_bool     = &spec.Type{Type: "boolean"}
	_dateTime = &spec.Type{Type: "string", Format: "date-time"}
	_types    = map[string]*spec.Type{
		"bool":      _bool,
		"time.Time": _dateTime,
		"enum":      _string,
		"string":    _string,
		"[]byte":    _bytes,
		"uuid.UUID": _string,
		"int":       _int32,
		"int8":      _int32,
		"int16":     _int32,
		"int32":     _int32,
		"uint":      _int32,
		"uint8":     _int32,
		"uint16":    _int32,
		"uint32":    _int32,
		"int64":     _int64,
		"uint64":    _int64,
		"float32":   _float,
		"float64":   _double,
	}
)

// oasType returns the OAS primitive tye (if any) for the given ent schema field.
func oasType(f *gen.Field) (*spec.Type, error) {
	// If there is a custom type given on the field use it.
	ant, err := FieldAnnotation(f)
	if err != nil {
		return nil, err
	}
	if ant.OASType != nil {
		return ant.OASType, nil
	}
	if f.IsEnum() {
		return _string, nil
	}
	s := f.Type.String()
	// Handle slice types.
	if strings.HasPrefix(s, "[]") {
		if t, ok := _types[s[2:]]; ok {
			return &spec.Type{Type: "array", Items: t}, nil
		}
	}
	t, ok := _types[s]
	if !ok {
		return nil, fmt.Errorf("no OAS-type exists for type %q of field %s", s, f.StructField())
	}
	return t, nil
}

// exampleValue returns the user defined example value for the ent schema field.
func exampleValue(f *gen.Field) (interface{}, error) {
	a, err := FieldAnnotation(f)
	if err != nil {
		return nil, err
	}
	if a != nil && a.Example != nil {
		return a.Example, err
	}
	if f.IsEnum() {
		return f.EnumValues()[0], nil
	}
	return nil, nil
}

// contains checks if a string slice contains the given value.
func contains(xs []Operation, s Operation) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// pathParam created a new path parameter for the ID of gen.Type.
func pathParam(n *gen.Type) (*spec.Parameter, error) {
	t, err := oasType(n.ID)
	if err != nil {
		return nil, err
	}
	return &spec.Parameter{
		Name:        "id",
		In:          spec.InPath,
		Description: fmt.Sprintf("ID of the %s", n.Name),
		Required:    true,
		Schema:      t,
	}, nil
}
