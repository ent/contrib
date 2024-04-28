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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
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
				vn, err := ViewNameEdge(strings.Split(n, "_")[0], e)
				if err != nil {
					return err
				}
				es, ok := spec.Components.Schemas[vn]
				if !ok {
					return fmt.Errorf("schema %q not found for edge %q on %q", vn, e.Name, n)
				}
				es = es.ToNamed(vn).AsLocalRef()
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
		ant, err := FieldAnnotation(f)
		if err != nil {
			return err
		}
		if ant.Skip {
			continue
		}
		p, err := property(f)
		if err != nil {
			return err
		}
		addProperty(s, p, !(f.Optional || f.Nillable))
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
				SetJSONContent(ogen.NewSchema().
					AddRequiredProperties(
						ogen.Int().ToProperty("code"),
						ogen.String().ToProperty("status"),
					).
					AddOptionalProperties(
						ogen.NewSchema().ToProperty("errors"),
					),
				), // TODO(masseelch): Add examples once present https://github.com/ogen-go/ogen/issues/70
		)
	}
}

var rules = inflect.NewDefaultRuleset()

// paths adds all operations to the spec paths.
func paths(g *gen.Graph, spec *ogen.Spec) error {
	cfg, err := GetConfig(g.Config)
	if err != nil {
		return err
	}

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
			path(spec, root).Post, err = createOp(spec, n, cfg.AllowClientUUIDs)
			if err != nil {
				return err
			}
		}
		// Read operation.
		if contains(ops, OpRead) {
			path(spec, root+"/{id}").Get, err = readOp(spec, n)
			if err != nil {
				return err
			}
		}
		// Update operation.
		if contains(ops, OpUpdate) {
			path(spec, root+"/{id}").Patch, err = updateOp(spec, n)
			if err != nil {
				return err
			}
		}
		// Delete operation.
		if contains(ops, OpDelete) {
			path(spec, root+"/{id}").Delete, err = deleteOp(spec, n)
			if err != nil {
				return err
			}
		}
		// List operation.
		if contains(ops, OpList) {
			path(spec, root).Get, err = listOp(spec, n)
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
			// Read operation.
			if contains(ops, OpRead) {
				path(spec, subRoot).Get, err = readEdgeOp(spec, n, e)
				if err != nil {
					return err
				}
			}
			// List operation.
			if contains(ops, OpList) {
				path(spec, subRoot).Get, err = listEdgeOp(spec, n, e)
				if err != nil {
					return err
				}
			}
		}
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

// createOp returns an ogen.Operation for a create operation on the given node.
func createOp(spec *ogen.Spec, n *gen.Type, allowClientUUIDs bool) (*ogen.Operation, error) {
	req, err := reqBody(n, OpCreate, allowClientUUIDs)
	if err != nil {
		return nil, err
	}
	vn, err := ViewName(n, OpCreate)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Create a new %s", n.Name)).
		SetDescription(fmt.Sprintf("Creates a new %s and persists it to storage.", n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpCreate)+n.Name).
		SetRequestBody(req).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s created", n.Name)).
				SetJSONContent(spec.RefSchema(vn).Schema),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// readOp returns an ogen.Operation for a read operation on the given node.
func readOp(spec *ogen.Spec, n *gen.Type) (*ogen.Operation, error) {
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := ViewName(n, OpRead)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Find a %s by ID", n.Name)).
		SetDescription(fmt.Sprintf("Finds the %s with the requested ID and returns it.", n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpRead)+n.Name).
		AddParameters(id).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s with requested ID was found", n.Name)).
				SetJSONContent(spec.RefSchema(vn).Schema),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// readEdgeOp returns the spec description for a read operation on a subresource.
func readEdgeOp(spec *ogen.Spec, n *gen.Type, e *gen.Edge) (*ogen.Operation, error) {
	if !e.Unique {
		return nil, errors.New("read operations are not allowed on non unique edges")
	}
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := EdgeViewName(n, e, OpRead)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Find the attached %s", e.Type.Name)).
		SetDescription(fmt.Sprintf("Find the attached %s of the %s with the given ID", e.Type.Name, n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpRead)+n.Name+strcase.UpperCamelCase(e.Name)).
		AddParameters(id).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s attached to %s with requested ID was found", e.Type.Name, n.Name)).
				SetJSONContent(spec.RefSchema(vn).Schema),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// updateOp returns a spec.OperationConfig for an update operation on the given node.
func updateOp(spec *ogen.Spec, n *gen.Type) (*ogen.Operation, error) {
	req, err := reqBody(n, OpUpdate, false)
	if err != nil {
		return nil, err
	}
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := ViewName(n, OpUpdate)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Updates a %s", n.Name)).
		SetDescription(fmt.Sprintf("Updates a %s and persists changes to storage.", n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpUpdate)+n.Name).
		AddParameters(id).
		SetRequestBody(req).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s updated", n.Name)).
				SetJSONContent(spec.RefSchema(vn).Schema),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// deleteOp returns a spec.Operation for a delete operation on the given node.
func deleteOp(spec *ogen.Spec, n *gen.Type) (*ogen.Operation, error) {
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("Deletes a %s by ID", n.Name)).
		SetDescription(fmt.Sprintf("Deletes the %s with the requested ID.", n.Name)).
		AddTags(n.Name).
		SetOperationID(string(OpDelete)+n.Name).
		AddParameters(id).
		AddResponse(
			strconv.Itoa(http.StatusNoContent),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("%s with requested ID was deleted", n.Name)),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// listOp returns a spec.OperationConfig for a list operation on the given node.
func listOp(spec *ogen.Spec, n *gen.Type) (*ogen.Operation, error) {
	cfg, err := GetConfig(n.Config)
	if err != nil {
		return nil, err
	}
	vn, err := ViewName(n, OpList)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("List %s", rules.Pluralize(n.Name))).
		SetDescription(fmt.Sprintf("List %s.", rules.Pluralize(n.Name))).
		AddTags(n.Name).
		SetOperationID(string(OpList)+n.Name).
		AddParameters( // TODO(masseelch): Add cursor based pagination to entoas and ogent.
			ogen.NewParameter().
				InQuery().
				SetName("page").
				SetDescription("what page to render").
				SetSchema(ogen.Int().SetMinimum(ptr(int64(1)))),
			ogen.NewParameter().
				InQuery().
				SetName("itemsPerPage").
				SetDescription("item count to render per page").
				SetSchema(ogen.Int().
					SetMinimum(&cfg.MinItemsPerPage).
					SetMaximum(&cfg.MaxItemsPerPage),
				),
		).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("result %s list", n.Name)).
				SetJSONContent(spec.RefSchema(vn).Schema.AsArray()),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// listEdgeOp returns the spec description for a read operation on a subresource.
func listEdgeOp(spec *ogen.Spec, n *gen.Type, e *gen.Edge) (*ogen.Operation, error) {
	if e.Unique {
		return nil, errors.New("list operations are not allowed on unique edges")
	}
	id, err := pathParam(n)
	if err != nil {
		return nil, err
	}
	vn, err := EdgeViewName(n, e, OpList)
	if err != nil {
		return nil, err
	}
	op := ogen.NewOperation().
		SetSummary(fmt.Sprintf("List attached %s", rules.Pluralize(strcase.UpperCamelCase(e.Name)))).
		SetDescription(fmt.Sprintf("List attached %s.", rules.Pluralize(strcase.UpperCamelCase(e.Name)))).
		AddTags(n.Name).
		SetOperationID(string(OpList)+n.Name+strcase.UpperCamelCase(e.Name)).
		AddParameters( // TODO(masseelch): Add cursor based pagination to entoas and ogent.
			id,
			ogen.NewParameter().
				InQuery().
				SetName("page").
				SetDescription("what page to render").
				SetSchema(ogen.Int()),
			ogen.NewParameter().
				InQuery().
				SetName("itemsPerPage").
				SetDescription("item count to render per page").
				SetSchema(ogen.Int()),
		).
		AddResponse(
			strconv.Itoa(http.StatusOK),
			ogen.NewResponse().
				SetDescription(fmt.Sprintf("result %s list", rules.Pluralize(strcase.UpperCamelCase(n.Name)))).
				SetJSONContent(spec.RefSchema(vn).Schema.AsArray()),
		).
		AddNamedResponses(
			spec.RefResponse(strconv.Itoa(http.StatusBadRequest)),
			spec.RefResponse(strconv.Itoa(http.StatusConflict)),
			spec.RefResponse(strconv.Itoa(http.StatusNotFound)),
			spec.RefResponse(strconv.Itoa(http.StatusInternalServerError)),
		)
	return op, nil
}

// property creates an ogen.Property out of an ent schema field.
func property(f *gen.Field) (*ogen.Property, error) {
	s, err := OgenSchema(f)
	if err != nil {
		return nil, err
	}
	return ogen.NewProperty().SetName(f.Name).SetSchema(s), nil
}

// ptr returns a pointer to v, for the purposes of base types (int, string, etc).
func ptr[T any](v T) *T {
	return &v
}

// mapTypeToSchema returns an ogen.Schema for the given gen.Field, if it exists.
// returns nil if the type is not supported.
func mapTypeToSchema(baseType string) *ogen.Schema {
	switch baseType {
	case "bool":
		return ogen.Bool()
	case "time.Time":
		return ogen.DateTime()
	case "string":
		return ogen.String()
	case "[]byte":
		return ogen.Bytes()
	case "uuid.UUID":
		return ogen.UUID()
	case "int":
		return ogen.Int()
	case "int8":
		return ogen.Int32().SetMinimum(ptr(int64(math.MinInt8))).SetMaximum(ptr(int64(math.MaxInt8)))
	case "int16":
		return ogen.Int32().SetMinimum(ptr(int64(math.MinInt16))).SetMaximum(ptr(int64(math.MaxInt16)))
	case "int32":
		return ogen.Int32().SetMinimum(ptr(int64(math.MinInt32))).SetMaximum(ptr(int64(math.MaxInt32)))
	case "int64":
		return ogen.Int64().SetMinimum(ptr(int64(math.MinInt64))).SetMaximum(ptr(int64(math.MaxInt64)))
	case "uint":
		return ogen.Int64().SetMinimum(ptr(int64(0))).SetMaximum(ptr(int64(math.MaxUint32)))
	case "uint8":
		return ogen.Int32().SetMinimum(ptr(int64(0))).SetMaximum(ptr(int64(math.MaxUint8)))
	case "uint16":
		return ogen.Int32().SetMinimum(ptr(int64(0))).SetMaximum(ptr(int64(math.MaxUint16)))
	case "uint32":
		return ogen.Int64().SetMinimum(ptr(int64(0))).SetMaximum(ptr(int64(math.MaxUint32)))
	case "uint64":
		return ogen.Int64().SetMinimum(ptr(int64(0)))
	case "float32":
		return ogen.Float()
	case "float64":
		return ogen.Double()
	default:
		return nil
	}
}

// OgenSchema returns the ogen.Schema to use for the given gen.Field.
func OgenSchema(f *gen.Field) (*ogen.Schema, error) {
	// If there is a custom property/schema given on the field, use it.
	ant, err := FieldAnnotation(f)
	if err != nil {
		return nil, err
	}
	if ant.Schema != nil {
		return ant.Schema, nil
	}

	var schema *ogen.Schema
	baseType := f.Type.String()

	// Enum values need special case.
	if f.IsEnum() {
		var d json.RawMessage
		if f.Default {
			d, err = json.Marshal(f.DefaultValue().(string))
			if err != nil {
				return nil, err
			}
		}
		vs := make([]json.RawMessage, len(f.EnumValues()))
		for i, e := range f.EnumValues() {
			vs[i], err = json.Marshal(e)
			if err != nil {
				return nil, err
			}
		}
		schema = ogen.String().AsEnum(d, vs...)
	}

	if schema == nil {
		if strings.HasPrefix(baseType, "[]") { // Handle slice types.
			schema = mapTypeToSchema(baseType[2:])
			if schema != nil {
				schema = schema.AsArray()
			}
		}

		if schema == nil {
			schema = mapTypeToSchema(baseType)
		}
	}

	if schema == nil {
		return nil, fmt.Errorf("no OAS-type exists for type %q of field %s", baseType, f.StructField())
	}

	if schema.Description == "" {
		schema.Description = f.Comment()
	}

	if ant.Example != nil {
		schema.Example, err = json.Marshal(ant.Example)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal example for field %s: %w", f.StructField(), err)
		}
	}

	return schema, nil
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
		sort.Slice(ops, func(i, j int) bool {
			return ops[i] < ops[j]
		})
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
				return []Operation{OpRead}, nil
			} else {
				return []Operation{OpList}, nil
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
			m[OpRead] = ant.Read
		} else {
			m[OpList] = ant.List
		}
		for op, opn := range m {
			if opn.Policy == PolicyExpose || (opn.Policy == PolicyNone && c.DefaultPolicy == PolicyExpose) {
				ops = append(ops, op)
				continue
			}
		}
		sort.Slice(ops, func(i, j int) bool {
			return ops[i] < ops[j]
		})
		return ops, nil
	}
}

// reqBody returns the request body for the given node and operation.
func reqBody(n *gen.Type, op Operation, allowClientUUIDs bool) (*ogen.RequestBody, error) {
	req := ogen.NewRequestBody().SetRequired(true)
	c := ogen.NewSchema()
	switch op {
	case OpCreate:
		// add the ID field as client setable if it is a UUID.
		if allowClientUUIDs && n.ID.Type.Type == field.TypeUUID {
			p, err := property(n.ID)
			if err != nil {
				return nil, err
			}
			addProperty(c, p, !n.ID.Default)
		}

		req.SetDescription(fmt.Sprintf("%s to create", n.Name))
	case OpUpdate:
		req.SetDescription(fmt.Sprintf("%s properties to update", n.Name))
	default:
		return nil, fmt.Errorf("requestBody: unsupported operation %q", op)
	}
	for _, f := range n.Fields {
		a, err := FieldAnnotation(f)
		if err != nil {
			return nil, err
		}
		if a.ReadOnly || a.Skip {
			continue
		}
		if op == OpCreate || !f.Immutable {
			p, err := property(f)
			if err != nil {
				return nil, err
			}
			addProperty(c, p, op == OpCreate && !f.Optional && !f.Default)
		}
	}
	for _, e := range n.Edges {
		if op == OpUpdate && (e.Immutable || (e.Field() != nil && e.Field().Immutable)) {
			continue
		}

		s, err := OgenSchema(e.Type.ID)
		if err != nil {
			return nil, err
		}
		if !e.Unique {
			s = s.AsArray()
		}

		if e.Field() != nil {
			f := e.Field()
			a, err := FieldAnnotation(f)
			if err != nil {
				return nil, err
			}

			if a.ReadOnly {
				continue
			}

			if !a.Skip {
				// If the edge has a field, and the field isn't skipped, then there is no
				// point in having two fields that can be used during create (especially
				// if both are required).
				continue
			}

			addProperty(c, s.ToProperty(e.Name), (op == OpCreate && !e.Optional && !f.Default) || (op == OpUpdate && !f.UpdateDefault))
		} else {
			addProperty(c, s.ToProperty(e.Name), (op == OpCreate && !e.Optional))
		}
	}
	req.SetJSONContent(c)
	return req, nil
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

// pathParam creates a new Parameter in path for the ID of gen.Type.
func pathParam(n *gen.Type) (*ogen.Parameter, error) {
	t, err := OgenSchema(n.ID)
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
