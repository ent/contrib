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
	require.False(t, strings.Contains(output, "func (t *TodoQuery) Paginate"),
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
