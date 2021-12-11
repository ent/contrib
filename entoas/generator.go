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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/go-openapi/inflect"
	"github.com/ogen-go/ogen"
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

func generate(g *gen.Graph, spec *ogen.Spec) error {
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
func schemas(g *gen.Graph, spec *ogen.Spec) error {
	// Loop over every defined node and add it to the spec.
	for _, n := range g.Nodes {
		s := ogen.NewSchema()
		if err := addSchemaFields(s, append([]*gen.Field{n.ID}, n.Fields...)); err != nil {
			return err
		}
		spec.AddSchema(n.Name, s)
	}
	// Loop over every node once more to add the edges.
	for _, n := range g.Nodes {
		for _, e := range n.Edges {
			es, ok := spec.Components.Schemas[e.Type.Name]
			if !ok {
				return fmt.Errorf("schema %q not found for edge %q on %q", e.Type.Name, e.Name, n.Name)
			}
			es = es.ToNamed(e.Type.Name).AsLocalRef()
			if !e.Unique {
				es = es.AsArray()
			}
			addProperty(
				spec.Components.Schemas[n.Name],
				ogen.NewProperty().SetName(e.Name).SetSchema(es),
				!e.Optional,
			)
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
			s := ogen.NewSchema()
			if err := addSchemaFields(s, v.Fields); err != nil {
				return err
			}
			spec.AddSchema(n, s)
		}
		// Loop over every view once more to add the edges.
		for n, v := range vs {
			for _, e := range v.Edges {
				vn, err := viewNameEdge(strings.Split(n, "_")[0], e)
				if err != nil {
					return err
				}
				es, ok := spec.Components.Schemas[vn]
				if !ok {
					return fmt.Errorf("schema %q not found for edge %q on %q", vn, e.Name, n)
				}
				es = es.ToNamed(e.Type.Name).AsLocalRef()
				if !e.Unique {
					es = es.AsArray()
				}
				addProperty(
					spec.Components.Schemas[n],
					ogen.NewProperty().SetName(e.Name).SetSchema(es),
					!e.Optional,
				)
			}
		}
	}
	return nil
}

// addSchemaFields adds the given gen.Field slice to the ogen.Schema.
func addSchemaFields(s *ogen.Schema, fs []*gen.Field) error {
	for _, f := range fs {
		p, err := property(f)
		if err != nil {
			return err
		}
		addProperty(s, p, !f.Optional)
	}
	return nil
}

// addProperty adds the ogen.Property to the ogen.Schema and marks it as required if needed.
func addProperty(s *ogen.Schema, p *ogen.Property, req bool) {
	if req {
		s.AddRequiredProperties(p)
	} else {
		s.AddOptionalProperties(p)
	}
}

// errResponses adds all responses to the spec responses.
func errorResponses(s *ogen.Spec) {
	for c, d := range map[int]string{
		http.StatusBadRequest:          "invalid input, data invalid",
		http.StatusConflict:            "conflicting resources",
		http.StatusForbidden:           "insufficient permissions",
		http.StatusInternalServerError: "unexpected error",
		http.StatusNotFound:            "resource not found",
	} {
		s.AddResponse(
			strconv.Itoa(c),
			ogen.NewResponse().
				SetDescription(d).
				SetContent(nil), // TODO: Implement when Content builders are present ...
		)
		// s.Components.Responses[strconv.Itoa(c)] = &spec.Response{
		// 	Name:        strconv.Itoa(c),
		// 	Description: d,
		// 	Headers:     nil, // TODO
		// 	Content: spec.Content{
		// 		spec.JSON: &spec.MediaTypeObject{
		// 			Unique: true,
		// 			Schema: spec.Schema{
		// 				Fields: map[string]*spec.Field{
		// 					"code": {
		// 						Type:    _int32,
		// 						Unique:  true,
		// 						Example: c,
		// 					},
		// 					"status": {
		// 						Type:    _string,
		// 						Unique:  true,
		// 						Example: http.StatusText(c),
		// 					},
		// 				},
		// 				Edges: map[string]*spec.Edge{
		// 					"errors": {
		// 						Schema: new(spec.Schema),
		// 						Unique: true,
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}
}

var rules = inflect.NewDefaultRuleset()

// paths adds all operations to the spec paths.
func paths(g *gen.Graph, spec *ogen.Spec) error {
	for _, n := range g.Nodes {
		// Add schema operations.
		ops, err := NodeOperations(n)
		if err != nil {
			return err
		}
		// root for all operations on this node.
		root := "/" + rules.Pluralize(strcase.KebabCase(n.Name))
		// // Create operation.
		// if contains(ops, OpCreate) {
		// 	path(spec, root).Post, err = createOp(spec, n)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// Read operation.
		if contains(ops, OpRead) {
			path(spec, root+"/{id}").Get, err = readOp(n)
			if err != nil {
				return err
			}
		}
		// // Update operation.
		// if contains(ops, OpUpdate) {
		// 	path(spec, root+"/{id}").Patch, err = updateOp(spec, n)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// // Delete operation.
		// if contains(ops, OpDelete) {
		// 	path(spec, root+"/{id}").Delete, err = deleteOp(spec, n)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// // List operation.
		// if contains(ops, OpList) {
		// 	path(spec, root).Get, err = listOp(spec, n)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// // Sub-Resource operations.
		// for _, e := range n.Edges {
		// 	subRoot := root + "/{id}/" + strcase.KebabCase(e.Name)
		// 	ops, err := EdgeOperations(e)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// Create operation.
		// 	if contains(ops, OpCreate) {
		// 		path(spec, subRoot).Post, err = createEdgeOp(spec, n, e)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	// Read operation.
		// 	if contains(ops, OpRead) {
		// 		path(spec, subRoot).Get, err = readEdgeOp(spec, n, e)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	// Delete operation.
		// 	if contains(ops, OpDelete) {
		// 		path(spec, subRoot).Delete, err = deleteEdgeOp(spec, n, e)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	// List operation.
		// 	if contains(ops, OpList) {
		// 		path(spec, subRoot).Get, err = listEdgeOp(spec, n, e)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// }
	}
	return nil
}

// path returns the correct spec.Path for the given root. Creates and sets a fresh instance if non does yet exist.
func path(s *ogen.Spec, root string) *ogen.PathItem {
	if s.Paths == nil {
		s.Paths = make(ogen.Paths)
	}
	if _, ok := s.Paths[root]; !ok {
		s.Paths[root] = ogen.NewPathItem()
	}
	return s.Paths[root]
}

// // createOp returns the spec description for a create operation on the given node.
// func createOp(s *ogen.Spec, n *gen.Type) (*ogen.Operation, error) {
// 	req, err := reqBody(n, OpCreate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	vn, err := viewName(n, OpCreate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	op := ogen.NewOperation().
// 		// SetSummary(fmt.Sprintf("Create a new %s", n.Name)) // TODO: re-enable after https://github.com/ogen-go/ogen/pull/73
// 		SetDescription(fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name))
// 	return op, nil
// 	return &spec.Operation{
// 		Summary:     fmt.Sprintf("Create a new %s", n.Name),
// 		Description: fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name),
// 		Tags:        []string{n.Name},
// 		OperationID: string(OpCreate) + n.Name,
// 		RequestBody: req,
// 		Responses: map[string]*spec.OperationResponse{
// 			strconv.Itoa(http.StatusOK): {
// 				Response: &spec.Response{
// 					Description: fmt.Sprintf("%s created", n.Name),
// 					Headers:     nil, // TODO
// 					Content: spec.Content{
// 						spec.JSON: &spec.MediaTypeObject{
// 							Unique: true,
// 							Ref:    s.Components.Schemas[vn],
// 						},
// 					},
// 				},
// 			},
// 			strconv.Itoa(http.StatusBadRequest): {
// 				Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)],
// 			},
// 			strconv.Itoa(http.StatusInternalServerError): {
// 				Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)],
// 			},
// 		},
// 		// Security: ant.CreateSecurity,
// 	}, nil
// }
//
// // createEdgeOp returns the spec description for a create operation on a subresource.
// func createEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
// 	// Create a basic create operation as if this was a first level operation.
// 	op, err := createOp(s, e.Type)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// But now alter the fields required to make this a second level operation.
// 	op.Summary = fmt.Sprintf("Create a new %s and attach it to the %s", e.Type.Name, n.Name)
// 	op.Description = fmt.Sprintf("Creates a new %s and attaches it to the %s", e.Type.Name, n.Name)
// 	op.Tags = []string{n.Name}
// 	op.OperationID = string(OpCreate) + n.Name + strcase.UpperCamelCase(e.Name)
// 	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
// 	rp.Description = fmt.Sprintf("%s created and attached to the %s", e.Type.Name, n.Name)
// 	vn, err := edgeViewName(n, e, OpCreate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	op.Parameters = []*spec.Parameter{id}
// 	return op, nil
// }

// readOp returns a spec.OperationConfig for a read operation on the given node.
func readOp(n *gen.Type) (*ogen.Operation, error) {
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	// vn, err := viewName(n, OpRead)
	// if err != nil {
	// 	return nil, err
	// }
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Find a %s by ID", n.Name)).
		SetDescription(fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpRead)+n.Name).
		AddParameters(id).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s with requested ID was found", n.Name)),
		)
	return op, nil
	// return &spec.Operation{
	// 	Summary:     fmt.Sprintf("Find a %s by ID", n.Name),
	// 	Description: fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name),
	// 	Tags:        []string{n.Name},
	// 	OperationID: string(OpRead) + n.Name,
	// 	Parameters:  []*spec.Parameter{id},
	// 	Responses: map[string]*spec.OperationResponse{
	// 		strconv.Itoa(http.StatusOK): {
	// 			Response: &spec.Response{
	// 				Description: fmt.Sprintf("%s with requested ID was found", n.Name),
	// 				Headers:     nil, // TODO
	// 				Content: spec.Content{
	// 					spec.JSON: &spec.MediaTypeObject{
	// 						Unique: true,
	// 						Ref:    s.Components.Schemas[vn],
	// 					},
	// 				},
	// 			},
	// 		},
	// 		strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
	// 		strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
	// 		strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
	// 	},
	// 	// Security: ant.ReadSecurity,
	// }, nil
}

// // readEdgeOp returns the spec description for a read operation on a subresource.
// func readEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
// 	if !e.Unique {
// 		return nil, errors.New("read operations are not allowed on non unique edges")
// 	}
// 	// Create a basic read operation as if this was a first level operation.
// 	op, err := readOp(s, e.Type)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// But now alter the fields required to make this a second level operation.
// 	op.Summary = fmt.Sprintf("Find the attached %s", e.Type.Name)
// 	op.Description = fmt.Sprintf("Find the attached %s of the %s with the given ID", e.Type.Name, n.Name)
// 	op.Tags = []string{n.Name}
// 	op.OperationID = string(OpRead) + n.Name + strcase.UpperCamelCase(e.Name)
// 	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
// 	rp.Description = fmt.Sprintf("%s attached to %s with requested ID was found", e.Type.Name, n.Name)
// 	vn, err := edgeViewName(n, e, OpRead)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	op.Parameters = []*spec.Parameter{id}
// 	return op, nil
// }
//
// // updateOp returns a spec.OperationConfig for an update operation on the given node.
// func updateOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
// 	req, err := reqBody(n, OpUpdate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	vn, err := viewName(n, OpUpdate)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &spec.Operation{
// 		Summary:     fmt.Sprintf("Updates a %s", n.Name),
// 		Description: fmt.Sprintf("Updates a %s and persists changes to storage.", n.Name),
// 		Tags:        []string{n.Name},
// 		OperationID: string(OpUpdate) + n.Name,
// 		Parameters:  []*spec.Parameter{id},
// 		RequestBody: req,
// 		Responses: map[string]*spec.OperationResponse{
// 			strconv.Itoa(http.StatusOK): {
// 				Response: &spec.Response{
// 					Description: fmt.Sprintf("%s updated", n.Name),
// 					Headers:     nil, // TODO
// 					Content: spec.Content{
// 						spec.JSON: &spec.MediaTypeObject{
// 							Unique: true,
// 							Ref:    s.Components.Schemas[vn],
// 						},
// 					},
// 				},
// 			},
// 			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
// 			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
// 			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
// 		},
// 	}, nil
// }
//
// // deleteOp returns a spec.OperationConfig for a delete operation on the given node.
// func deleteOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &spec.Operation{
// 		Summary:     fmt.Sprintf("Deletes a %s by ID", n.Name),
// 		Description: fmt.Sprintf("Deletes the %s with the requested ID.", n.Name),
// 		Tags:        []string{n.Name},
// 		OperationID: string(OpDelete) + n.Name,
// 		Parameters:  []*spec.Parameter{id},
// 		Responses: map[string]*spec.OperationResponse{
// 			strconv.Itoa(http.StatusNoContent): {
// 				Response: &spec.Response{
// 					Description: fmt.Sprintf("%s with requested ID was deleted", n.Name),
// 					Headers:     nil, // TODO
// 				},
// 			},
// 			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
// 			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
// 			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
// 		},
// 	}, nil
// }
//
// // deleteEdgeOp returns the spec description for a delete operation on a subresource.
// func deleteEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
// 	if !e.Unique {
// 		return nil, errors.New("delete operations are not allowed on non unique edges")
// 	}
// 	// Create a basic delete operation as if this was a first level operation.
// 	op, err := deleteOp(s, e.Type)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// But now alter the fields required to make this a second level operation.
// 	op.Summary = fmt.Sprintf("Delete the attached %s", strcase.UpperCamelCase(e.Name))
// 	op.Description = fmt.Sprintf(
// 		"Delete the attached %s of the %s with the given ID", strcase.UpperCamelCase(e.Name), n.Name,
// 	)
// 	op.Tags = []string{n.Name}
// 	op.OperationID = string(OpDelete) + n.Name + strcase.UpperCamelCase(e.Name)
// 	op.Responses[strconv.Itoa(http.StatusNoContent)].Response.Description = fmt.Sprintf(
// 		"%s with requested ID was deleted", strcase.UpperCamelCase(e.Name),
// 	)
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	op.Parameters = []*spec.Parameter{id}
// 	return op, nil
// }
//
// // listOp returns a spec.OperationConfig for a list operation on the given node.
// func listOp(s *spec.Spec, n *gen.Type) (*spec.Operation, error) {
// 	vn, err := viewName(n, OpList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &spec.Operation{
// 		Summary:     fmt.Sprintf("List %s", rules.Pluralize(n.Name)),
// 		Description: fmt.Sprintf("List %s.", rules.Pluralize(n.Name)),
// 		Tags:        []string{n.Name},
// 		OperationID: string(OpList) + n.Name,
// 		Parameters: []*spec.Parameter{{
// 			Name:        "page",
// 			In:          spec.InQuery,
// 			Description: "what page to render",
// 			Schema:      _int32,
// 		}, {
// 			Name:        "itemsPerPage",
// 			In:          spec.InQuery,
// 			Description: "item count to render per page",
// 			Schema:      _int32,
// 		}},
// 		Responses: map[string]*spec.OperationResponse{
// 			strconv.Itoa(http.StatusOK): {
// 				Response: &spec.Response{
// 					Description: fmt.Sprintf("result %s list", n.Name),
// 					Headers:     nil, // TODO
// 					Content: spec.Content{
// 						spec.JSON: &spec.MediaTypeObject{
// 							Ref: s.Components.Schemas[vn],
// 						},
// 					},
// 				},
// 			},
// 			strconv.Itoa(http.StatusBadRequest):          {Ref: s.Components.Responses[strconv.Itoa(http.StatusBadRequest)]},
// 			strconv.Itoa(http.StatusNotFound):            {Ref: s.Components.Responses[strconv.Itoa(http.StatusNotFound)]},
// 			strconv.Itoa(http.StatusInternalServerError): {Ref: s.Components.Responses[strconv.Itoa(http.StatusInternalServerError)]},
// 		},
// 	}, nil
// }
//
// // listEdgeOp returns the spec description for a read operation on a subresource.
// func listEdgeOp(s *spec.Spec, n *gen.Type, e *gen.Edge) (*spec.Operation, error) {
// 	if e.Unique {
// 		return nil, errors.New("list operations are not allowed on unique edges")
// 	}
// 	// Create a basic read operation as if this was a first level operation.
// 	op, err := listOp(s, e.Type)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// But now alter the fields required to make this a second level operation.
// 	op.Summary = fmt.Sprintf("List attached %s", rules.Pluralize(strcase.UpperCamelCase(e.Name)))
// 	op.Description = fmt.Sprintf("List attached %s.", rules.Pluralize(strcase.UpperCamelCase(e.Name)))
// 	op.Tags = []string{n.Name}
// 	op.OperationID = string(OpList) + n.Name + strcase.UpperCamelCase(e.Name)
// 	rp := op.Responses[strconv.Itoa(http.StatusOK)].Response
// 	rp.Description = fmt.Sprintf("result %s list", rules.Pluralize(strcase.UpperCamelCase(n.Name)))
// 	vn, err := edgeViewName(n, e, OpList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rp.Content[spec.JSON].Ref = s.Components.Schemas[vn]
// 	id, err := pathParam(n)
// 	if err != nil {
// 		return nil, err
// 	}
// 	op.Parameters = []*spec.Parameter{id}
// 	return op, nil
// }
//

// property creates an ogen.Property out of an ent schema field.
func property(f *gen.Field) (*ogen.Property, error) {
	s, err := ogenSchema(f)
	if err != nil {
		return nil, err
	}
	return ogen.NewProperty().SetName(f.Name).SetSchema(s), nil
}

var _types = map[string]*ogen.Schema{
	"bool":      ogen.Bool(),
	"time.Time": ogen.DateTime(),
	"string":    ogen.String(),
	"[]byte":    ogen.Bytes(),
	"uuid.UUID": ogen.String(),
	"int":       ogen.Int32(),
	"int8":      ogen.Int32(),
	"int16":     ogen.Int32(),
	"int32":     ogen.Int32(),
	"uint":      ogen.Int32(),
	"uint8":     ogen.Int32(),
	"uint16":    ogen.Int32(),
	"uint32":    ogen.Int32(),
	"int64":     ogen.Int64(),
	"uint64":    ogen.Int64(),
	"float32":   ogen.Float(),
	"float64":   ogen.Double(),
}

// ogenSchema returns the ogen.Schema to use for the given gen.Field.
func ogenSchema(f *gen.Field) (*ogen.Schema, error) {
	// If there is a custom property given on the field use it.
	ant, err := FieldAnnotation(f)
	if err != nil {
		return nil, err
	}
	if ant.Schema != nil {
		return ant.Schema, nil
	}
	// Enum values need special case.
	if f.IsEnum() {
		var d json.RawMessage
		if f.Default {
			d = json.RawMessage((f.DefaultValue()).(string))
		}
		vs := make([]json.RawMessage, len(f.EnumValues()))
		for i, v := range f.EnumValues() {
			vs[i] = json.RawMessage(v)
		}
		return ogen.String().AsEnum(d, vs...), nil
	}
	s := f.Type.String()
	// Handle slice types.
	if strings.HasPrefix(s, "[]") {
		if t, ok := _types[s[2:]]; ok {
			return t.AsArray(), nil
		}
	}
	t, ok := _types[s]
	if !ok {
		return nil, fmt.Errorf("no OAS-type exists for type %q of field %s", s, f.StructField())
	}
	return t, nil
}

// // field created a spec schema field from an ent schema field.
// func field(f *gen.Field) (*spec.Field, error) {
// 	t, err := oasType(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ex, err := exampleValue(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &spec.Field{Unique: true, Required: !f.Optional, Type: t, Example: ex}, nil
// }

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

// // reqBody returns the request body for the given node and operation.
// func reqBody(n *gen.Type, op Operation) (*spec.RequestBody, error) {
// 	req := &spec.RequestBody{}
// 	switch op {
// 	case OpCreate:
// 		req.Description = fmt.Sprintf("%s to create", n.Name)
// 	case OpUpdate:
// 		req.Description = fmt.Sprintf("%s properties to update", n.Name)
// 	default:
// 		return nil, fmt.Errorf("requestBody: unsupported operation %q", op)
// 	}
// 	fs := make(spec.Fields)
// 	for _, f := range n.Fields {
// 		if op == OpCreate || !f.Immutable {
// 			sf, err := field(f)
// 			if err != nil {
// 				return nil, err
// 			}
// 			fs[f.Name] = sf
// 		}
// 	}
// 	for _, e := range n.Edges {
// 		t, err := oasType(e.Type.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fs[e.Name] = &spec.Field{
// 			Unique:   e.Unique,
// 			Required: !e.Optional,
// 			Type:     t,
// 			Example:  nil, // TODO: Example for a unique / non-unique edge
// 		}
// 	}
// 	req.Content = spec.Content{
// 		spec.JSON: &spec.MediaTypeObject{
// 			Unique: true,
// 			Schema: spec.Schema{
// 				Name:   n.Name + op.Title() + "Request",
// 				Fields: fs,
// 			},
// 		},
// 	}
// 	return req, nil
// }
//
// // exampleValue returns the user defined example value for the ent schema field.
// func exampleValue(f *gen.Field) (interface{}, error) {
// 	a, err := FieldAnnotation(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if a != nil && a.Example != nil {
// 		return a.Example, err
// 	}
// 	if f.IsEnum() {
// 		return f.EnumValues()[0], nil
// 	}
// 	return nil, nil
// }

// contains checks if a string slice contains the given value.
func contains(xs []Operation, s Operation) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

// pathParam creates a new Parameter in path for the ID of gen.Type.
func pathParam(n *gen.Type) (*ogen.Parameter, error) {
	t, err := ogenSchema(n.ID)
	if err != nil {
		return nil, err
	}
	return ogen.NewParameter().
		InPath().
		SetName("id").
		SetDescription(fmt.Sprintf("ID of the %s", n.Name)).
		SetRequired(true).
		SetSchema(t), nil
}
