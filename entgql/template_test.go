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

package entgql

import (
	"bytes"
	"strings"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

var annotationName = Annotation{}.Name()

func TestFilterNodes(t *testing.T) {
	nodes, err := filterNodes([]*gen.Type{
		{
			Name: "Type1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name:   "Type2",
			Config: &gen.Config{},
		},
		{
			Name: "SkippedType",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
	}, SkipType)
	require.NoError(t, err)
	require.Equal(t, []*gen.Type{
		{
			Name: "Type1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name:   "Type2",
			Config: &gen.Config{},
		},
	}, nodes)
}

func TestFilterEdges(t *testing.T) {
	edges, err := filterEdges([]*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{},
		},
		{
			Name: "SkippedEdge",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
		{
			Name: "SkippedEdgeType",
			Type: &gen.Type{
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{"Skip": SkipAll},
				},
			},
		},
	}, SkipType)
	require.NoError(t, err)
	require.Equal(t, []*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{},
		},
	}, edges)
}

func TestFieldCollections(t *testing.T) {
	edges := []*gen.Edge{
		{
			Name: "Edge1",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Edge2",
			Type: &gen.Type{
				Name: "Todo",
			},
		},
		{
			Name: "two_words",
			Type: &gen.Type{
				Name: "Todo",
			},
		},
		{
			Name: "Unbind",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"Unbind": true,
				},
			},
		},
		{
			Name: "EdgeMapping",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"Unbind":  true,
					"Mapping": []string{"field1", "field2"},
				},
			},
		},
	}
	collect, err := fieldCollections(edges)
	require.NoError(t, err)
	require.Equal(t, []*fieldCollection{
		{
			Edge:    edges[0],
			Mapping: []string{"edge1"},
		},
		{
			Edge:    edges[1],
			Mapping: []string{"edge2"},
		},
		{
			Edge:    edges[2],
			Mapping: []string{"twoWords"},
		},
		{
			Edge:    edges[4],
			Mapping: []string{"field1", "field2"},
		},
	}, collect)

	_, err = fieldCollections([]*gen.Edge{
		{
			Name: "EdgeInvalid",
			Type: &gen.Type{
				Name: "Todo",
			},
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{
					"Unbind":  false,
					"Mapping": []string{"field1", "field2"},
				},
			},
		},
	})
	require.Errorf(t, err, "bind and mapping annotations are mutually exclusive")
}

func TestPaginationSharedTemplateParsed(t *testing.T) {
	// Verify the PaginationSharedTemplate was parsed successfully during init().
	require.NotNil(t, PaginationSharedTemplate, "PaginationSharedTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_pagination_shared", PaginationSharedTemplate.Name())
	// Verify it has the expected define block.
	tmpl := PaginationSharedTemplate.Lookup("gql_pagination_shared")
	require.NotNil(t, tmpl, "template should contain 'gql_pagination_shared' define block")
}

func TestPaginationSharedTemplateContent(t *testing.T) {
	// Verify the template source contains the expected shared code elements.
	tmpl := PaginationSharedTemplate.Lookup("gql_pagination_shared")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify shared type aliases are present.
	require.Contains(t, src, "Cursor = entgql.Cursor")
	require.Contains(t, src, "PageInfo = entgql.PageInfo")
	require.Contains(t, src, "OrderDirection = entgql.OrderDirection")
	require.Contains(t, src, "NullsDirection = entgql.NullsDirection")

	// Verify shared functions are present.
	require.Contains(t, src, "func orderFunc")
	require.Contains(t, src, "func validateFirstLast")
	require.Contains(t, src, "func collectedField")
	require.Contains(t, src, "func hasCollectedField")
	require.Contains(t, src, "func paginateLimit")

	// Verify shared constants are present.
	// errInvalidPagination is a literal string constant in the template.
	require.Contains(t, src, "errInvalidPagination")
	// The field constants (edgesField, nodeField, etc.) are generated via a range
	// over list "edges" "node" "pageInfo" "totalCount", so check the list items.
	require.Contains(t, src, `"edges"`)
	require.Contains(t, src, `"node"`)
	require.Contains(t, src, `"pageInfo"`)
	require.Contains(t, src, `"totalCount"`)
	// Verify the Field suffix pattern is in the template.
	require.Contains(t, src, `Field = "`)

	// Verify NO per-entity code is present.
	require.NotContains(t, src, "Connection struct")
	require.NotContains(t, src, "Edge struct")
	require.NotContains(t, src, "Pager")
	require.NotContains(t, src, "Paginate")
	require.NotContains(t, src, "applyOrder")
	require.NotContains(t, src, "applyCursors")
	require.NotContains(t, src, "applyFilter")
	require.NotContains(t, src, "ToEdge")
}

func TestPaginationSharedTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph and verify
	// the generated output contains expected shared code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	tmpl := PaginationSharedTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_pagination_shared", struct {
		*gen.Graph
	}{graph})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify shared type aliases are generated.
	require.Contains(t, output, "Cursor = entgql.Cursor[")
	require.Contains(t, output, "PageInfo = entgql.PageInfo[")
	require.Contains(t, output, "OrderDirection = entgql.OrderDirection")
	require.Contains(t, output, "NullsDirection = entgql.NullsDirection")

	// Verify shared functions are generated.
	require.Contains(t, output, "func orderFunc(")
	require.Contains(t, output, "func validateFirstLast(")
	require.Contains(t, output, "func collectedField(")
	require.Contains(t, output, "func hasCollectedField(")
	require.Contains(t, output, "func paginateLimit(")

	// Verify shared constants are generated.
	require.Contains(t, output, `errInvalidPagination`)
	require.Contains(t, output, `edgesField = "edges"`)
	require.Contains(t, output, `nodeField = "node"`)
	require.Contains(t, output, `pageInfoField = "pageInfo"`)
	require.Contains(t, output, `totalCountField = "totalCount"`)

	// Verify NO per-entity code is present (e.g., no TodoConnection, TodoEdge, todoPager).
	require.False(t, strings.Contains(output, "TodoConnection"),
		"shared template should not generate per-entity TodoConnection type")
	require.False(t, strings.Contains(output, "TodoEdge"),
		"shared template should not generate per-entity TodoEdge type")
	require.False(t, strings.Contains(output, "todoPager"),
		"shared template should not generate per-entity todoPager type")
	require.False(t, strings.Contains(output, "func (_m *TodoQuery) Paginate"),
		"shared template should not generate per-entity Paginate method")
}

func TestCollectionSharedTemplateParsed(t *testing.T) {
	// Verify the CollectionSharedTemplate was parsed successfully during init().
	require.NotNil(t, CollectionSharedTemplate, "CollectionSharedTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_collection_shared", CollectionSharedTemplate.Name())
	// Verify it has the expected define block.
	tmpl := CollectionSharedTemplate.Lookup("gql_collection_shared")
	require.NotNil(t, tmpl, "template should contain 'gql_collection_shared' define block")
}

func TestCollectionSharedTemplateContent(t *testing.T) {
	// Verify the template contains the expected shared code elements
	// by checking the template tree for key identifiers.
	tmpl := CollectionSharedTemplate.Lookup("gql_collection_shared")
	require.NotNil(t, tmpl)

	// The template source should reference our shared constants and functions.
	// We verify this by checking the template's parse tree root text contains expected strings.
	src := tmpl.Tree.Root.String()

	// The constants are generated via a range over a list, so we verify
	// the list items are present in the template source rather than the
	// expanded constant names (which only appear after template execution).
	require.Contains(t, src, `"after"`)
	require.Contains(t, src, `"first"`)
	require.Contains(t, src, `"before"`)
	require.Contains(t, src, `"last"`)
	require.Contains(t, src, `"orderBy"`)
	require.Contains(t, src, `"direction"`)
	require.Contains(t, src, `"where"`)
	// Verify the Field suffix pattern is in the template.
	require.Contains(t, src, `Field = "`)

	// Verify shared functions are present.
	require.Contains(t, src, "func fieldArgs")
	require.Contains(t, src, "func unmarshalArgs")
	require.Contains(t, src, "func mayAddCondition")

	// Verify NO per-entity code is present.
	require.NotContains(t, src, "CollectFields")
	require.NotContains(t, src, "collectField")
	require.NotContains(t, src, "PaginateArgs")
	require.NotContains(t, src, "newPaginateArg")
}

func TestCollectionSharedTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph and verify
	// the generated output contains expected shared code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	tmpl := CollectionSharedTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_collection_shared", struct {
		*gen.Graph
	}{graph})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify shared constants are generated.
	require.Contains(t, output, `afterField = "after"`)
	require.Contains(t, output, `firstField = "first"`)
	require.Contains(t, output, `beforeField = "before"`)
	require.Contains(t, output, `lastField = "last"`)
	require.Contains(t, output, `orderByField = "orderBy"`)
	require.Contains(t, output, `directionField = "direction"`)
	require.Contains(t, output, `fieldField = "field"`)
	require.Contains(t, output, `whereField = "where"`)

	// Verify shared functions are generated.
	require.Contains(t, output, "func fieldArgs(")
	require.Contains(t, output, "func unmarshalArgs(")
	require.Contains(t, output, "func mayAddCondition(")

	// Verify NO per-entity code is present.
	require.False(t, strings.Contains(output, "CollectFields"),
		"shared template should not generate per-entity CollectFields method")
	require.False(t, strings.Contains(output, "collectField"),
		"shared template should not generate per-entity collectField method")
	require.False(t, strings.Contains(output, "PaginateArgs"),
		"shared template should not generate per-entity PaginateArgs type")
	require.False(t, strings.Contains(output, "newPaginateArg"),
		"shared template should not generate per-entity newPaginateArgs function")
}

func TestPaginationEntityTemplateParsed(t *testing.T) {
	// Verify the PaginationEntityTemplate was parsed successfully during init().
	require.NotNil(t, PaginationEntityTemplate, "PaginationEntityTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_pagination_entity", PaginationEntityTemplate.Name())
	// Verify it has the expected define block.
	tmpl := PaginationEntityTemplate.Lookup("gql_pagination_entity")
	require.NotNil(t, tmpl, "template should contain 'gql_pagination_entity' define block")
}

func TestPaginationEntityTemplateContent(t *testing.T) {
	// Verify the template source contains the expected per-entity code elements.
	tmpl := PaginationEntityTemplate.Lookup("gql_pagination_entity")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify per-entity types are present (template source uses {{$edge}}, {{$conn}} etc.).
	// In the parsed tree, template vars are rendered as {{$varname}}.
	require.Contains(t, src, "}} struct")       // Edge/Connection struct declarations
	require.Contains(t, src, "PaginateOption")  // PaginateOption type
	require.Contains(t, src, "OrderField")      // OrderField struct

	// Verify per-entity methods are present.
	require.Contains(t, src, "Paginate")
	require.Contains(t, src, "applyOrder")
	require.Contains(t, src, "applyCursors")
	require.Contains(t, src, "applyFilter")
	require.Contains(t, src, "toCursor")
	require.Contains(t, src, "orderExpr")
	require.Contains(t, src, "ToEdge")
	require.Contains(t, src, "MarshalGQL")
	require.Contains(t, src, "UnmarshalGQL")

	// Verify the paginate helper is inlined (no template call).
	require.NotContains(t, src, `template "gql_pagination/helper/paginate"`)
	require.Contains(t, src, "validateFirstLast")
	require.Contains(t, src, "paginateLimit")

	// Verify no $Scope references remain.
	require.NotContains(t, src, "$.Scope")

	// Verify the template uses $.Node instead of range loop.
	require.Contains(t, src, "$.Node")
	require.NotContains(t, src, "range $node := $gqlNodes")
}

func TestPaginationEntityTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph for one entity
	// and verify the generated output contains expected per-entity code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := PaginationEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_pagination_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, todoNode})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify per-entity types are generated.
	require.Contains(t, output, "TodoEdge struct")
	require.Contains(t, output, "TodoConnection struct")
	require.Contains(t, output, "TodoPaginateOption")
	require.Contains(t, output, "todoPager")
	require.Contains(t, output, "TodoOrderField struct")
	require.Contains(t, output, "TodoOrder struct")
	require.Contains(t, output, "DefaultTodoOrder")

	// Verify per-entity methods are generated.
	require.Contains(t, output, "func (_m *TodoQuery) Paginate(")
	require.Contains(t, output, "func TodoToEdge(_m *Todo,")

	// Verify the paginate helper is inlined (no template calls in output).
	require.Contains(t, output, "validateFirstLast(first, last)")
	require.Contains(t, output, "newTodoPager(opts, last != nil)")
	require.Contains(t, output, "pager.applyFilter(_m)")
	require.Contains(t, output, "pager.applyCursors(_m, after, before)")
	require.Contains(t, output, "pager.applyOrder(_m)")

	// Verify order fields are generated (Todo has OrderField annotations).
	// VarName uses the field's StructField name, not the GQL order field name.
	require.Contains(t, output, "TodoOrderFieldCreatedAt")
	require.Contains(t, output, "TodoOrderFieldStatus")
	require.Contains(t, output, "TodoOrderFieldText")
	require.Contains(t, output, "TodoOrderFieldPriority")

	// Verify NO shared code is present (shared types, shared funcs).
	require.NotContains(t, output, "Cursor = entgql.Cursor[")
	require.NotContains(t, output, "PageInfo = entgql.PageInfo[")
	require.NotContains(t, output, "func orderFunc(")
	require.NotContains(t, output, "func validateFirstLast(")
	require.NotContains(t, output, "func paginateLimit(")
	require.NotContains(t, output, "errInvalidPagination")
}

func TestPaginationEntityTemplateMultipleEntities(t *testing.T) {
	// Verify the template can be executed for multiple different entities
	// and produces distinct output for each.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find both Todo and Category nodes.
	var todoNode, categoryNode *gen.Type
	for _, n := range graph.Nodes {
		switch n.Name {
		case "Todo":
			todoNode = n
		case "Category":
			categoryNode = n
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist")
	require.NotNil(t, categoryNode, "Category node should exist")

	tmpl := PaginationEntityTemplate

	// Generate for Todo.
	var todoBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&todoBuf, "gql_pagination_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, todoNode})
	require.NoError(t, err)
	todoOutput := todoBuf.String()

	// Generate for Category.
	var catBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&catBuf, "gql_pagination_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, categoryNode})
	require.NoError(t, err)
	catOutput := catBuf.String()

	// Verify Todo output has Todo-specific types.
	require.Contains(t, todoOutput, "TodoEdge struct")
	require.Contains(t, todoOutput, "TodoConnection struct")
	require.NotContains(t, todoOutput, "CategoryEdge struct")
	require.NotContains(t, todoOutput, "CategoryConnection struct")

	// Verify Category output has Category-specific types.
	require.Contains(t, catOutput, "CategoryEdge struct")
	require.Contains(t, catOutput, "CategoryConnection struct")
	require.NotContains(t, catOutput, "TodoEdge struct")
	require.NotContains(t, catOutput, "TodoConnection struct")
}

func TestFilterFields(t *testing.T) {
	fields, err := filterFields([]*gen.Field{
		{
			Name: "Field1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Field2",
		},
		{
			Name: "SkippedField",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{"Skip": SkipAll},
			},
		},
	}, SkipType)
	require.NoError(t, err)
	require.Equal(t, []*gen.Field{
		{
			Name: "Field1",
			Annotations: map[string]interface{}{
				annotationName: map[string]interface{}{},
			},
		},
		{
			Name: "Field2",
		},
	}, fields)
}

func TestCollectionEntityTemplateParsed(t *testing.T) {
	// Verify the CollectionEntityTemplate was parsed successfully during init().
	require.NotNil(t, CollectionEntityTemplate, "CollectionEntityTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_collection_entity", CollectionEntityTemplate.Name())
	// Verify it has the expected define block.
	tmpl := CollectionEntityTemplate.Lookup("gql_collection_entity")
	require.NotNil(t, tmpl, "template should contain 'gql_collection_entity' define block")
}

func TestCollectionEntityTemplateContent(t *testing.T) {
	// Verify the template source contains the expected per-entity code elements
	// and does not contain monolithic template artifacts.
	tmpl := CollectionEntityTemplate.Lookup("gql_collection_entity")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify per-entity collection code is present.
	require.Contains(t, src, "CollectFields")
	require.Contains(t, src, "collectField")
	require.Contains(t, src, "PaginateArgs")

	// Verify load_total helper is inlined (no template call).
	require.NotContains(t, src, `template "gql_pagination/helper/load_total"`,
		"load_total helper should be inlined, not called as a sub-template")

	// Verify no $Scope references remain.
	require.NotContains(t, src, "Scope",
		"entity template should not contain $Scope references")

	// Verify hasTemplate is not used.
	require.NotContains(t, src, "hasTemplate",
		"entity template should use $.HasWhereInputTemplate instead of hasTemplate")

	// Verify $.HasWhereInputTemplate is used instead.
	require.Contains(t, src, "HasWhereInputTemplate",
		"entity template should reference HasWhereInputTemplate")

	// Verify it references $.Node (single entity, not range loop).
	require.Contains(t, src, "$.Node",
		"entity template should reference $.Node for the single entity")
}

func TestCollectionEntityTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph and verify
	// the generated output contains expected per-entity code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := CollectionEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_collection_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, todoNode, true})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify CollectFields method is generated for the Todo query.
	require.Contains(t, output, "CollectFields(ctx context.Context, satisfies ...string)")
	require.Contains(t, output, "collectField(ctx context.Context, oneNode bool")

	// Verify PaginateArgs struct is generated.
	require.Contains(t, output, "todoPaginateArgs")
	require.Contains(t, output, "newTodoPaginateArgs")

	// Verify no template call to load_total (it should be inlined).
	require.NotContains(t, output, "gql_pagination/helper/load_total")

	// Verify no Scope references in output.
	require.NotContains(t, output, "Scope")

	// Verify where input filter is present (since HasWhereInputTemplate=true).
	require.Contains(t, output, "whereField")

	// Verify per-entity imports.
	require.Contains(t, output, `"entgo.io/contrib/entgql/internal/todo/ent/todo"`)
}

func TestCollectionEntityTemplateNoWhereInput(t *testing.T) {
	// Execute the template with HasWhereInputTemplate=false and verify
	// the where input filter code is not generated.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the BillProduct node (simple entity with no edges).
	var node *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "BillProduct" {
			node = n
			break
		}
	}
	require.NotNil(t, node, "BillProduct node should exist in the schema")

	tmpl := CollectionEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_collection_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, node, false})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify CollectFields is present.
	require.Contains(t, output, "CollectFields(ctx context.Context, satisfies ...string)")

	// Verify PaginateArgs struct is generated for this entity.
	// Note: camel("BillProduct") produces "billproduct" for the struct type name.
	require.Contains(t, output, "billproductPaginateArgs")
	require.Contains(t, output, "newBillProductPaginateArgs")

	// Since HasWhereInputTemplate=false, the where filter block should not appear
	// in the newPaginateArgs function.
	// Check that the output does NOT contain the where input filter assignment.
	require.NotContains(t, output, "BillProductWhereInput",
		"where input filter should not be generated when HasWhereInputTemplate=false")
}

func TestEdgeEntityTemplateParsed(t *testing.T) {
	// Verify the EdgeEntityTemplate was parsed successfully during init().
	require.NotNil(t, EdgeEntityTemplate, "EdgeEntityTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_edge_entity", EdgeEntityTemplate.Name())
	// Verify it has the expected define block.
	tmpl := EdgeEntityTemplate.Lookup("gql_edge_entity")
	require.NotNil(t, tmpl, "template should contain 'gql_edge_entity' define block")
}

func TestEdgeEntityTemplateContent(t *testing.T) {
	// Verify the template source contains the expected per-entity code elements
	// and does not contain monolithic template artifacts.
	tmpl := EdgeEntityTemplate.Lookup("gql_edge_entity")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify the paginate helper is inlined (no template call).
	require.NotContains(t, src, `template "gql_edge/helper/paginate"`,
		"paginate helper should be inlined, not called as a sub-template")

	// Verify no $Scope references remain.
	require.NotContains(t, src, "$.Scope",
		"entity template should not contain $.Scope references")

	// Verify hasTemplate is not used.
	require.NotContains(t, src, "hasTemplate",
		"entity template should use $.HasWhereInputTemplate instead of hasTemplate")

	// Verify $.HasWhereInputTemplate is used instead.
	require.Contains(t, src, "HasWhereInputTemplate",
		"entity template should reference HasWhereInputTemplate")

	// Verify it references $.Node (single entity, not range loop).
	require.Contains(t, src, "$.Node",
		"entity template should reference $.Node for the single entity")

	// Verify all three edge types are handled.
	require.Contains(t, src, "isRelayConn",
		"template should handle Relay connection edges")
	require.Contains(t, src, "IsNotLoaded",
		"template should handle non-Relay edges with IsNotLoaded fallback")
	require.Contains(t, src, "MaskNotFound",
		"template should handle optional unique edges with MaskNotFound")

	// Verify Relay connection inlined code has expected elements.
	require.Contains(t, src, "nodePaginationNames",
		"template should use nodePaginationNames for Relay edges")
	require.Contains(t, src, "Paginate",
		"template should call Paginate as fallback for Relay edges")
}

func TestEdgeEntityTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph for Todo entity
	// and verify the generated output contains expected edge resolver code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := EdgeEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_edge_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, todoNode, true})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify edge resolver standalone functions are generated.
	// Todo has edges: parent, children, category, secret.
	// Children is a Relay connection edge, so it should have pagination params.
	require.Contains(t, output, "func ResolveTodo")

	// Verify import statements.
	require.Contains(t, output, `"context"`)
	require.Contains(t, output, `"github.com/99designs/gqlgen/graphql"`)

	// Verify no template call to gql_edge/helper/paginate (it should be inlined).
	require.NotContains(t, output, "gql_edge/helper/paginate")

	// Verify no Scope references in output.
	require.NotContains(t, output, "Scope")
}

func TestEdgeEntityTemplateNoWhereInput(t *testing.T) {
	// Execute the template with HasWhereInputTemplate=false and verify
	// the where input filter code is not generated.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := EdgeEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_edge_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, todoNode, false})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify edge resolver standalone functions are generated.
	require.Contains(t, output, "func ResolveTodo")

	// Since HasWhereInputTemplate=false, the where input filter should not appear.
	require.NotContains(t, output, "WhereInput",
		"where input should not be generated when HasWhereInputTemplate=false")
}

func TestEdgeEntityTemplateMultipleEntities(t *testing.T) {
	// Verify the template can be executed for multiple different entities
	// and produces distinct output for each.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find both Todo and Category nodes.
	var todoNode, categoryNode *gen.Type
	for _, n := range graph.Nodes {
		switch n.Name {
		case "Todo":
			todoNode = n
		case "Category":
			categoryNode = n
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist")
	require.NotNil(t, categoryNode, "Category node should exist")

	tmpl := EdgeEntityTemplate

	// Generate for Todo.
	var todoBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&todoBuf, "gql_edge_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, todoNode, true})
	require.NoError(t, err)
	todoOutput := todoBuf.String()

	// Generate for Category.
	var catBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&catBuf, "gql_edge_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{graph, categoryNode, true})
	require.NoError(t, err)
	catOutput := catBuf.String()

	// Verify Todo output has Todo-specific standalone functions.
	require.Contains(t, todoOutput, "func ResolveTodo")
	require.NotContains(t, todoOutput, "ResolveCategory")

	// Verify Category output has Category-specific standalone functions.
	require.Contains(t, catOutput, "func ResolveCategory")
	require.NotContains(t, catOutput, "ResolveTodo")
}

func TestNodeDescriptorSharedTemplateParsed(t *testing.T) {
	// Verify the NodeDescriptorSharedTemplate was parsed successfully during init().
	require.NotNil(t, NodeDescriptorSharedTemplate, "NodeDescriptorSharedTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_node_descriptor_shared", NodeDescriptorSharedTemplate.Name())
	// Verify it has the expected define block.
	tmpl := NodeDescriptorSharedTemplate.Lookup("gql_node_descriptor_shared")
	require.NotNil(t, tmpl, "template should contain 'gql_node_descriptor_shared' define block")
}

func TestNodeDescriptorSharedTemplateContent(t *testing.T) {
	// Verify the template source contains the expected shared code elements.
	tmpl := NodeDescriptorSharedTemplate.Lookup("gql_node_descriptor_shared")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify shared struct types are present.
	require.Contains(t, src, "Node struct")
	require.Contains(t, src, "Field struct")
	require.Contains(t, src, "Edge struct")

	// Verify Client.Node() method is present.
	require.Contains(t, src, "func (c *Client) Node(")
	require.Contains(t, src, "c.noder(ctx, table, id)")

	// Verify node descriptor registry is present.
	require.Contains(t, src, "nodeDescriptors")
	require.Contains(t, src, "registerNodeDescriptor")

	// Verify NO per-entity code is present.
	// Note: $.Nodes (plural) is expected in the shared template for filterNodes;
	// we check that $.Node (singular entity reference) is not used as a data field.
	require.NotContains(t, src, "$.Node }}")
	require.NotContains(t, src, "$.Node.Name")
	require.NotContains(t, src, "json.Marshal")
	require.NotContains(t, src, "QueryTodos")
	require.NotContains(t, src, "Noder interface")
}

func TestNodeDescriptorSharedTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph and verify
	// the generated output contains expected shared code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	tmpl := NodeDescriptorSharedTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_descriptor_shared", struct {
		*gen.Graph
	}{graph})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify shared struct types are generated.
	require.Contains(t, output, "type Node struct")
	require.Contains(t, output, "type Field struct")
	require.Contains(t, output, "type Edge struct")

	// Verify Client.Node() method is generated.
	require.Contains(t, output, "func (c *Client) Node(ctx context.Context, id int) (*Node, error)")
	require.Contains(t, output, "c.noder(ctx, table, id)")

	// Verify node descriptor registry is generated.
	require.Contains(t, output, "var nodeDescriptors")
	require.Contains(t, output, "func registerNodeDescriptor(")

	// Verify NO per-entity code is present.
	require.False(t, strings.Contains(output, "BillProductNode("),
		"shared template should not generate per-entity BillProductNode() function")
	require.False(t, strings.Contains(output, "TodoNode("),
		"shared template should not generate per-entity TodoNode() function")
	require.False(t, strings.Contains(output, "json.Marshal"),
		"shared template should not contain json.Marshal (per-entity code)")
}

func TestNodeDescriptorEntityTemplateParsed(t *testing.T) {
	// Verify the NodeDescriptorEntityTemplate was parsed successfully during init().
	require.NotNil(t, NodeDescriptorEntityTemplate, "NodeDescriptorEntityTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_node_descriptor_entity", NodeDescriptorEntityTemplate.Name())
	// Verify it has the expected define block.
	tmpl := NodeDescriptorEntityTemplate.Lookup("gql_node_descriptor_entity")
	require.NotNil(t, tmpl, "template should contain 'gql_node_descriptor_entity' define block")
}

func TestNodeDescriptorEntityTemplateContent(t *testing.T) {
	// Verify the template source contains the expected per-entity code elements.
	tmpl := NodeDescriptorEntityTemplate.Lookup("gql_node_descriptor_entity")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify per-entity node descriptor function is present.
	require.Contains(t, src, "registerNodeDescriptor")
	require.Contains(t, src, "json.Marshal")
	require.Contains(t, src, "$.Node")

	// Verify NO shared code is present.
	require.NotContains(t, src, "Node struct")
	require.NotContains(t, src, "Field struct")
	require.NotContains(t, src, "Edge struct")
	require.NotContains(t, src, "func (c *Client) Node(")
}

func TestNodeDescriptorEntityTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph for one entity
	// and verify the generated output contains expected per-entity code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := NodeDescriptorEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_descriptor_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, todoNode})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify per-entity standalone node descriptor function is generated for Todo.
	require.Contains(t, output, "func TodoNode(_m *Todo, ctx context.Context) (node *Node, err error)")

	// Verify node descriptor registration.
	require.Contains(t, output, "registerNodeDescriptor(todo.Table")

	// Verify field serialization is present.
	require.Contains(t, output, "json.Marshal(_m.CreatedAt)")
	require.Contains(t, output, "json.Marshal(_m.Status)")
	require.Contains(t, output, "json.Marshal(_m.Priority)")
	require.Contains(t, output, "json.Marshal(_m.Text)")

	// Verify edge queries use FromContext pattern.
	require.Contains(t, output, "FromContext(ctx).Todo.QueryParent(_m)")
	require.Contains(t, output, "FromContext(ctx).Todo.QueryChildren(_m)")
	require.Contains(t, output, "FromContext(ctx).Todo.QueryCategory(_m)")

	// Verify edge type imports.
	require.Contains(t, output, `"entgo.io/contrib/entgql/internal/todo/ent/todo"`)
	require.Contains(t, output, `"entgo.io/contrib/entgql/internal/todo/ent/category"`)

	// Verify NO shared code is present.
	require.NotContains(t, output, "type Node struct")
	require.NotContains(t, output, "type Field struct")
	require.NotContains(t, output, "type Edge struct")
	require.NotContains(t, output, "func (c *Client) Node(")
}

func TestNodeDescriptorEntityTemplateMultipleEntities(t *testing.T) {
	// Verify the template can be executed for multiple different entities
	// and produces distinct output for each.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find both Todo and BillProduct nodes.
	var todoNode, billProductNode *gen.Type
	for _, n := range graph.Nodes {
		switch n.Name {
		case "Todo":
			todoNode = n
		case "BillProduct":
			billProductNode = n
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist")
	require.NotNil(t, billProductNode, "BillProduct node should exist")

	tmpl := NodeDescriptorEntityTemplate

	// Generate for Todo.
	var todoBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&todoBuf, "gql_node_descriptor_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, todoNode})
	require.NoError(t, err)
	todoOutput := todoBuf.String()

	// Generate for BillProduct.
	var bpBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&bpBuf, "gql_node_descriptor_entity", struct {
		*gen.Graph
		Node *gen.Type
	}{graph, billProductNode})
	require.NoError(t, err)
	bpOutput := bpBuf.String()

	// Verify Todo output has Todo-specific standalone function.
	require.Contains(t, todoOutput, "func TodoNode(_m *Todo, ctx context.Context)")
	require.NotContains(t, todoOutput, "BillProductNode(")

	// Verify BillProduct output has BillProduct-specific standalone function.
	require.Contains(t, bpOutput, "func BillProductNode(_m *BillProduct, ctx context.Context)")
	require.NotContains(t, bpOutput, "TodoNode(")

	// BillProduct has no edges, so no edge queries should appear.
	require.NotContains(t, bpOutput, "QueryTodos")
	require.NotContains(t, bpOutput, "QueryParent")

	// BillProduct has fields: name, sku, quantity.
	require.Contains(t, bpOutput, "json.Marshal(_m.Name)")
	require.Contains(t, bpOutput, "json.Marshal(_m.Sku)")
	require.Contains(t, bpOutput, "json.Marshal(_m.Quantity)")
}

func TestNodeSharedTemplateParsed(t *testing.T) {
	// Verify the NodeSharedTemplate was parsed successfully during init().
	require.NotNil(t, NodeSharedTemplate, "NodeSharedTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_node_shared", NodeSharedTemplate.Name())
	// Verify it has the expected define block.
	tmpl := NodeSharedTemplate.Lookup("gql_node_shared")
	require.NotNil(t, tmpl, "template should contain 'gql_node_shared' define block")
}

func TestNodeSharedTemplateContent(t *testing.T) {
	// Verify the template source contains the expected shared code elements.
	tmpl := NodeSharedTemplate.Lookup("gql_node_shared")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify shared types and infrastructure are present.
	require.Contains(t, src, "Noder interface")
	require.Contains(t, src, "errNodeInvalidID")
	require.Contains(t, src, "NodeOption")
	require.Contains(t, src, "WithNodeType")
	require.Contains(t, src, "WithFixedNodeType")
	require.Contains(t, src, "nodeOptions struct")

	// Verify new map-based dispatch infrastructure.
	require.Contains(t, src, "nodeResolver struct")
	require.Contains(t, src, "nodeResolvers")
	require.Contains(t, src, "registerNodeResolver")

	// Verify Client methods are present.
	require.Contains(t, src, "func (c *Client) Noder(")
	require.Contains(t, src, "func (c *Client) noder(")
	require.Contains(t, src, "func (c *Client) Noders(")
	require.Contains(t, src, "func (c *Client) noders(")
	require.Contains(t, src, "func (c *Client) newNodeOpts(")

	// Verify map-based dispatch instead of switch.
	require.Contains(t, src, "nodeResolvers[table]")
	require.Contains(t, src, "resolver.byID(")
	require.Contains(t, src, "resolver.byIDs(")

	// Verify NO per-entity code is present.
	require.NotContains(t, src, "range $n")
	require.NotContains(t, src, "$.Node.Name")
	require.NotContains(t, src, "nodeImplementorsVar")

	// Verify no hasTemplate references.
	require.NotContains(t, src, "hasTemplate",
		"shared template should not use hasTemplate")
}

func TestNodeSharedTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph and verify
	// the generated output contains expected shared code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	tmpl := NodeSharedTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_shared", struct {
		*gen.Graph
		HasNodeDescriptorTemplate bool
	}{graph, true})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify Noder interface is generated with IsNode() method.
	require.Contains(t, output, "type Noder interface")
	require.Contains(t, output, "IsNode()")
	require.NotContains(t, output, "Node(context.Context) (*Node, error)")

	// Verify shared infrastructure types are generated.
	require.Contains(t, output, "var errNodeInvalidID")
	require.Contains(t, output, "type NodeOption func(*nodeOptions)")
	require.Contains(t, output, "func WithNodeType(")
	require.Contains(t, output, "func WithFixedNodeType(")

	// Verify new map-based dispatch types.
	require.Contains(t, output, "type nodeResolver struct")
	require.Contains(t, output, "var nodeResolvers = map[string]nodeResolver{}")
	require.Contains(t, output, "func registerNodeResolver(")

	// Verify Client methods.
	require.Contains(t, output, "func (c *Client) Noder(ctx context.Context, id int")
	require.Contains(t, output, "func (c *Client) Noders(ctx context.Context, ids []int")

	// Verify NO per-entity code is present (no entity package imports, no switch cases per entity).
	require.NotContains(t, output, `"/todo"`)
	require.NotContains(t, output, `"/category"`)
	require.NotContains(t, output, "case todo.Table")
	require.NotContains(t, output, "case category.Table")
}

func TestNodeSharedTemplateExecution_WithoutNodeDescriptor(t *testing.T) {
	// Execute with HasNodeDescriptorTemplate=false and verify Node() method is not in Noder interface.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	tmpl := NodeSharedTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_shared", struct {
		*gen.Graph
		HasNodeDescriptorTemplate bool
	}{graph, false})
	require.NoError(t, err)

	output := buf.String()

	// Noder interface should still exist but without Node() method.
	require.Contains(t, output, "type Noder interface")
	require.Contains(t, output, "IsNode()")
	require.NotContains(t, output, "Node(context.Context) (*Node, error)")
}

func TestNodeEntityTemplateParsed(t *testing.T) {
	// Verify the NodeEntityTemplate was parsed successfully during init().
	require.NotNil(t, NodeEntityTemplate, "NodeEntityTemplate should be parsed during init()")
	// Verify the template has the expected name.
	require.Equal(t, "gql_node_entity", NodeEntityTemplate.Name())
	// Verify it has the expected define block.
	tmpl := NodeEntityTemplate.Lookup("gql_node_entity")
	require.NotNil(t, tmpl, "template should contain 'gql_node_entity' define block")
}

func TestNodeEntityTemplateContent(t *testing.T) {
	// Verify the template source contains the expected per-entity code elements.
	tmpl := NodeEntityTemplate.Lookup("gql_node_entity")
	require.NotNil(t, tmpl)
	src := tmpl.Tree.Root.String()

	// Verify per-entity code is present.
	require.Contains(t, src, "nodeImplementorsVar")
	require.Contains(t, src, "nodeImplementors")
	require.Contains(t, src, "$.Node")
	require.Contains(t, src, "registerNodeResolver")
	require.Contains(t, src, "func init()")

	// Verify HasCollectionTemplate is used instead of hasTemplate.
	require.NotContains(t, src, "hasTemplate",
		"entity template should use $.HasCollectionTemplate instead of hasTemplate")
	require.Contains(t, src, "HasCollectionTemplate",
		"entity template should reference HasCollectionTemplate")

	// Verify NO shared code is present.
	require.NotContains(t, src, "Noder interface")
	require.NotContains(t, src, "errNodeInvalidID")
	require.NotContains(t, src, "NodeOption")
	require.NotContains(t, src, "func (c *Client) Noder(")
	require.NotContains(t, src, "func (c *Client) Noders(")
}

func TestNodeEntityTemplateExecution(t *testing.T) {
	// Execute the template against the real todo schema graph for one entity
	// and verify the generated output contains expected per-entity code.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find the Todo node.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist in the schema")

	tmpl := NodeEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasCollectionTemplate bool
	}{graph, todoNode, true})
	require.NoError(t, err)

	output := buf.String()

	// Verify the generated output contains the package declaration.
	require.Contains(t, output, "package ent")

	// Verify implementors var.
	require.Contains(t, output, "todoImplementors")

	// Verify init() registration.
	require.Contains(t, output, "func init()")
	require.Contains(t, output, "registerNodeResolver(todo.Table")
	require.Contains(t, output, "todoNoder")
	require.Contains(t, output, "todoNoders")

	// Verify noder functions.
	require.Contains(t, output, "func todoNoder(ctx context.Context, c *Client, id int) (Noder, error)")
	require.Contains(t, output, "func todoNoders(ctx context.Context, c *Client, ids []int")
	require.Contains(t, output, "c.Todo.Query()")
	require.Contains(t, output, "todo.ID(id)")
	require.Contains(t, output, "todo.IDIn(ids...)")

	// Verify collectField is present since HasCollectionTemplate=true.
	require.Contains(t, output, "collectField")

	// Verify entity package import.
	require.Contains(t, output, `"entgo.io/contrib/entgql/internal/todo/ent/todo"`)

	// Verify NO shared code is present.
	require.NotContains(t, output, "type Noder interface")
	require.NotContains(t, output, "var errNodeInvalidID")
	require.NotContains(t, output, "func (c *Client) Noder(")
	require.NotContains(t, output, "func (c *Client) Noders(")
}

func TestNodeEntityTemplateExecution_WithoutCollection(t *testing.T) {
	// Execute with HasCollectionTemplate=false and verify collectField is absent.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	tmpl := NodeEntityTemplate
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "gql_node_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasCollectionTemplate bool
	}{graph, todoNode, false})
	require.NoError(t, err)

	output := buf.String()

	// Verify collectField is NOT present since HasCollectionTemplate=false.
	require.NotContains(t, output, "collectField")
	require.NotContains(t, output, "CollectFields")

	// But the core noder functions should still be present.
	require.Contains(t, output, "func todoNoder(")
	require.Contains(t, output, "func todoNoders(")
	require.Contains(t, output, "func init()")
}

func TestNodeEntityTemplateMultipleEntities(t *testing.T) {
	// Verify the template can be executed for multiple different entities
	// and produces distinct output for each.
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	// Find both Todo and Category nodes.
	var todoNode, categoryNode *gen.Type
	for _, n := range graph.Nodes {
		switch n.Name {
		case "Todo":
			todoNode = n
		case "Category":
			categoryNode = n
		}
	}
	require.NotNil(t, todoNode, "Todo node should exist")
	require.NotNil(t, categoryNode, "Category node should exist")

	tmpl := NodeEntityTemplate

	// Generate for Todo.
	var todoBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&todoBuf, "gql_node_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasCollectionTemplate bool
	}{graph, todoNode, true})
	require.NoError(t, err)
	todoOutput := todoBuf.String()

	// Generate for Category.
	var catBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&catBuf, "gql_node_entity", struct {
		*gen.Graph
		Node                  *gen.Type
		HasCollectionTemplate bool
	}{graph, categoryNode, true})
	require.NoError(t, err)
	catOutput := catBuf.String()

	// Verify Todo output has Todo-specific code.
	require.Contains(t, todoOutput, "todoImplementors")
	require.Contains(t, todoOutput, "todoNoder")
	require.Contains(t, todoOutput, "todoNoders")
	require.Contains(t, todoOutput, "registerNodeResolver(todo.Table")
	require.NotContains(t, todoOutput, "categoryNoder")

	// Verify Category output has Category-specific code.
	require.Contains(t, catOutput, "categoryImplementors")
	require.Contains(t, catOutput, "categoryNoder")
	require.Contains(t, catOutput, "categoryNoders")
	require.Contains(t, catOutput, "registerNodeResolver(category.Table")
	require.NotContains(t, catOutput, "todoNoder")
}
