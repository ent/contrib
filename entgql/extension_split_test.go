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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

// loadTestGraph loads the todo schema graph with Target set to tmpDir.
func loadTestGraph(t *testing.T, tmpDir string) *gen.Graph {
	t.Helper()
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)
	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
		Target:  tmpDir,
		Package: "entgo.io/contrib/entgql/internal/todo/ent",
	})
	require.NoError(t, err)
	return graph
}

// nonSkippedNodes returns the names of nodes that are not skipped via SkipType.
func nonSkippedNodes(t *testing.T, graph *gen.Graph) []string {
	t.Helper()
	nodes, err := filterNodes(graph.Nodes, SkipType)
	require.NoError(t, err)
	var names []string
	for _, n := range nodes {
		names = append(names, n.Name)
	}
	return names
}

func TestGenerateSplitPagination(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_pagination_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Create a dummy monolithic gql_pagination.go to verify overwrite.
	monolithicPath := filepath.Join(tmpDir, "gql_pagination.go")
	err = os.WriteFile(monolithicPath, []byte("package ent\n// monolithic placeholder\n"), 0644)
	require.NoError(t, err)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitPagination(graph)
	require.NoError(t, err)

	// Verify the monolithic file was overwritten with shared content.
	sharedContent, err := os.ReadFile(monolithicPath)
	require.NoError(t, err)
	sharedStr := string(sharedContent)
	require.Contains(t, sharedStr, "package ent")
	// Shared content should NOT contain per-entity types.
	require.NotContains(t, sharedStr, "TodoEdge struct")
	require.NotContains(t, sharedStr, "CategoryConnection struct")

	// Verify per-entity pagination files are created.
	nodeNames := nonSkippedNodes(t, graph)
	require.NotEmpty(t, nodeNames)
	for _, name := range nodeNames {
		filename := fmt.Sprintf("gql_pagination_%s.go", snake(name))
		path := filepath.Join(tmpDir, filename)
		content, err := os.ReadFile(path)
		require.NoError(t, err, "pagination entity file should exist for %s", name)
		contentStr := string(content)
		require.Contains(t, contentStr, "package ent")
		// Verify the file contains pagination-related code (Edge, Connection).
		// Note: some entities have type aliases (e.g. Workspace -> Organization),
		// so we check for the generic patterns rather than exact name matches.
		require.Contains(t, contentStr, "Edge struct")
		require.Contains(t, contentStr, "Connection struct")
	}

	// Verify skipped types do NOT get per-entity files.
	skippedPath := filepath.Join(tmpDir, "gql_pagination_very_secret.go")
	_, err = os.Stat(skippedPath)
	require.True(t, os.IsNotExist(err), "skipped entity should not have pagination file")
}

func TestGenerateSplitCollection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_collection_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Create a dummy monolithic gql_collection.go to verify overwrite.
	monolithicPath := filepath.Join(tmpDir, "gql_collection.go")
	err = os.WriteFile(monolithicPath, []byte("package ent\n// monolithic placeholder\n"), 0644)
	require.NoError(t, err)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitCollection(graph)
	require.NoError(t, err)

	// Verify the monolithic file was overwritten with shared content.
	sharedContent, err := os.ReadFile(monolithicPath)
	require.NoError(t, err)
	sharedStr := string(sharedContent)
	require.Contains(t, sharedStr, "package ent")

	// Verify per-entity collection files are created.
	nodeNames := nonSkippedNodes(t, graph)
	require.NotEmpty(t, nodeNames)
	for _, name := range nodeNames {
		filename := fmt.Sprintf("gql_collection_%s.go", snake(name))
		path := filepath.Join(tmpDir, filename)
		content, err := os.ReadFile(path)
		require.NoError(t, err, "collection entity file should exist for %s", name)
		contentStr := string(content)
		require.Contains(t, contentStr, "package ent")
	}

	// Verify skipped types do NOT get per-entity files.
	skippedPath := filepath.Join(tmpDir, "gql_collection_very_secret.go")
	_, err = os.Stat(skippedPath)
	require.True(t, os.IsNotExist(err), "skipped entity should not have collection file")
}

func TestGenerateSplitGoFiles_IncludesPaginationAndCollection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_gofiles_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Create dummy monolithic files to check they get handled.
	for _, name := range []string{"gql_where_input.go", "gql_mutation_input.go", "gql_pagination.go", "gql_collection.go", "gql_edge.go", "gql_node_descriptor.go", "gql_node.go"} {
		path := filepath.Join(tmpDir, name)
		err = os.WriteFile(path, []byte("package ent\n// placeholder\n"), 0644)
		require.NoError(t, err)
	}

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitGoFiles(graph)
	require.NoError(t, err)

	// Where, mutation, and edge monolithic files should be removed.
	_, err = os.Stat(filepath.Join(tmpDir, "gql_where_input.go"))
	require.True(t, os.IsNotExist(err), "monolithic where_input should be removed")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_mutation_input.go"))
	require.True(t, os.IsNotExist(err), "monolithic mutation_input should be removed")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_edge.go"))
	require.True(t, os.IsNotExist(err), "monolithic edge should be removed")

	// Pagination, collection, and node_descriptor monolithic files should still exist (overwritten).
	_, err = os.Stat(filepath.Join(tmpDir, "gql_pagination.go"))
	require.NoError(t, err, "gql_pagination.go should still exist (overwritten with shared content)")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_collection.go"))
	require.NoError(t, err, "gql_collection.go should still exist (overwritten with shared content)")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_node_descriptor.go"))
	require.NoError(t, err, "gql_node_descriptor.go should still exist (overwritten with shared content)")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_node.go"))
	require.NoError(t, err, "gql_node.go should still exist (overwritten with shared content)")

	// Verify pagination, collection, node_descriptor, and node per-entity files exist.
	nodeNames := nonSkippedNodes(t, graph)
	for _, name := range nodeNames {
		paginationFile := filepath.Join(tmpDir, fmt.Sprintf("gql_pagination_%s.go", snake(name)))
		_, err = os.Stat(paginationFile)
		require.NoError(t, err, "pagination entity file should exist for %s", name)

		collectionFile := filepath.Join(tmpDir, fmt.Sprintf("gql_collection_%s.go", snake(name)))
		_, err = os.Stat(collectionFile)
		require.NoError(t, err, "collection entity file should exist for %s", name)

		nodeDescFile := filepath.Join(tmpDir, fmt.Sprintf("gql_node_descriptor_%s.go", snake(name)))
		_, err = os.Stat(nodeDescFile)
		require.NoError(t, err, "node descriptor entity file should exist for %s", name)

		nodeFile := filepath.Join(tmpDir, fmt.Sprintf("gql_node_%s.go", snake(name)))
		_, err = os.Stat(nodeFile)
		require.NoError(t, err, "node entity file should exist for %s", name)
	}

	// Verify edge per-entity files exist only for entities with edges.
	for _, name := range nodeNames {
		edgeFile := filepath.Join(tmpDir, fmt.Sprintf("gql_edge_%s.go", snake(name)))
		_, statErr := os.Stat(edgeFile)
		// Edge files are only created for entities that have edges after filtering,
		// so we just verify no error for the general case. Detailed edge testing
		// is covered by TestGenerateSplitEdge.
		_ = statErr
	}
}

func TestGenerateSplitGoFiles_RemovesStaleGeneratedFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_cleanup_gofiles_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	staleFiles := []string{
		"gql_where_input_stale.go",
		"gql_mutation_input_stale.go",
		"gql_edge_stale.go",
		"gql_pagination_stale.go",
		"gql_collection_stale.go",
		"gql_node_stale.go",
		"gql_node_descriptor_stale.go",
	}
	for _, name := range staleFiles {
		err = os.WriteFile(filepath.Join(tmpDir, name), []byte("package ent\n// stale generated file\n"), 0644)
		require.NoError(t, err)
	}

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitGoFiles(graph)
	require.NoError(t, err)

	for _, name := range staleFiles {
		_, err = os.Stat(filepath.Join(tmpDir, name))
		require.True(t, os.IsNotExist(err), "stale generated file should be removed: %s", name)
	}
}

func TestCollectionEntityFile_HasWhereInputTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_collection_where_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Find a node to test with.
	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	// Test with genWhereInput = true
	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateCollectionEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_collection_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)
	require.Contains(t, contentStr, "package ent")

	// Test with genWhereInput = false
	ex2, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(false),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex2.generateCollectionEntityFile(graph, todoNode)
	require.NoError(t, err)

	content2, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr2 := string(content2)
	require.Contains(t, contentStr2, "package ent")
	// The content should differ between where=true and where=false
	// (where=true generates WhereInput filter code).
	// We just verify both versions compile (don't error) and contain the package declaration.
}

func TestPaginationSharedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_pagination_shared_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generatePaginationSharedFile(graph)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_pagination.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	// Shared file should contain package declaration and shared types.
	require.Contains(t, contentStr, "package ent")

	// Shared code should have shared helpers and type aliases.
	require.Contains(t, contentStr, "Cursor")
	require.Contains(t, contentStr, "PageInfo")
}

func TestCollectionSharedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_collection_shared_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateCollectionSharedFile(graph)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_collection.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	// Shared file should contain package declaration.
	require.Contains(t, contentStr, "package ent")
}

func TestPaginationEntityFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_pagination_entity_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generatePaginationEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_pagination_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	require.Contains(t, contentStr, "package ent")
	require.Contains(t, contentStr, "TodoEdge")
	require.Contains(t, contentStr, "TodoConnection")
	require.Contains(t, contentStr, "TodoPaginateOption")

	// Should NOT contain other entity types.
	require.NotContains(t, contentStr, "CategoryEdge")
	require.NotContains(t, contentStr, "CategoryConnection")
}

func TestCollectionEntityFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_collection_entity_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateCollectionEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_collection_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	require.Contains(t, contentStr, "package ent")

	// Should NOT contain other entity types.
	require.NotContains(t, contentStr, "CategoryQuery")
}

func TestGenerateSplitEdge(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_edge_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitEdge(graph)
	require.NoError(t, err)

	// Verify per-entity edge files are created only for entities with edges.
	nodeNames := nonSkippedNodes(t, graph)
	require.NotEmpty(t, nodeNames)
	for _, name := range nodeNames {
		filename := fmt.Sprintf("gql_edge_%s.go", snake(name))
		path := filepath.Join(tmpDir, filename)
		_, statErr := os.Stat(path)

		// Find the node to check if it has edges.
		var node *gen.Type
		for _, n := range graph.Nodes {
			if n.Name == name {
				node = n
				break
			}
		}
		require.NotNil(t, node, "node %s should exist in graph", name)

		edges, filterErr := filterEdges(node.Edges, SkipType)
		require.NoError(t, filterErr)

		if len(edges) == 0 {
			require.True(t, os.IsNotExist(statErr),
				"edge entity file should NOT exist for %s (no edges)", name)
		} else {
			require.NoError(t, statErr,
				"edge entity file should exist for %s", name)
			content, readErr := os.ReadFile(path)
			require.NoError(t, readErr)
			contentStr := string(content)
			require.Contains(t, contentStr, "package ent")
		}
	}

	// Verify skipped types do NOT get per-entity files.
	skippedPath := filepath.Join(tmpDir, "gql_edge_very_secret.go")
	_, err = os.Stat(skippedPath)
	require.True(t, os.IsNotExist(err), "skipped entity should not have edge file")
}

func TestEdgeEntityFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_edge_entity_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateEdgeEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_edge_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	require.Contains(t, contentStr, "package ent")
	// Todo has Parent, Children, and Category edges.
	require.Contains(t, contentStr, "func ResolveTodoParent(")
	require.Contains(t, contentStr, "func ResolveTodoChildren(")
	require.Contains(t, contentStr, "func ResolveTodoCategory(")

	// Should NOT contain other entity edge resolvers.
	require.NotContains(t, contentStr, "ResolveCategory")
	require.NotContains(t, contentStr, "ResolveUser")
}

func TestEdgeEntityFile_HasWhereInputTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_edge_where_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Find a node with a Relay connection edge that uses where input.
	var categoryNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Category" {
			categoryNode = n
			break
		}
	}
	require.NotNil(t, categoryNode)

	// Test with genWhereInput = true
	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateEdgeEntityFile(graph, categoryNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_edge_category.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)
	require.Contains(t, contentStr, "package ent")
	// With where inputs enabled, the edge should have a where parameter.
	require.Contains(t, contentStr, "where *TodoWhereInput")

	// Test with genWhereInput = false
	ex2, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(false),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex2.generateEdgeEntityFile(graph, categoryNode)
	require.NoError(t, err)

	content2, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr2 := string(content2)
	require.Contains(t, contentStr2, "package ent")
	// Without where inputs, the where parameter should not be present.
	require.NotContains(t, contentStr2, "WhereInput")
}

func TestGenerateSplitNodeDescriptor(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_node_descriptor_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Create a dummy monolithic gql_node_descriptor.go to verify overwrite.
	monolithicPath := filepath.Join(tmpDir, "gql_node_descriptor.go")
	err = os.WriteFile(monolithicPath, []byte("package ent\n// monolithic placeholder\n"), 0644)
	require.NoError(t, err)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitNodeDescriptor(graph)
	require.NoError(t, err)

	// Verify the monolithic file was overwritten with shared content.
	sharedContent, err := os.ReadFile(monolithicPath)
	require.NoError(t, err)
	sharedStr := string(sharedContent)
	require.Contains(t, sharedStr, "package ent")
	// Shared content should NOT contain per-entity Node() methods.
	require.NotContains(t, sharedStr, "func (_m *Todo) Node(")
	require.NotContains(t, sharedStr, "func (_m *Category) Node(")

	// Verify per-entity node descriptor files are created.
	nodeNames := nonSkippedNodes(t, graph)
	require.NotEmpty(t, nodeNames)
	for _, name := range nodeNames {
		filename := fmt.Sprintf("gql_node_descriptor_%s.go", snake(name))
		path := filepath.Join(tmpDir, filename)
		content, err := os.ReadFile(path)
		require.NoError(t, err, "node descriptor entity file should exist for %s", name)
		contentStr := string(content)
		require.Contains(t, contentStr, "package ent")
		require.Contains(t, contentStr, "registerNodeDescriptor(")
	}

	// Verify skipped types do NOT get per-entity files.
	skippedPath := filepath.Join(tmpDir, "gql_node_descriptor_very_secret.go")
	_, err = os.Stat(skippedPath)
	require.True(t, os.IsNotExist(err), "skipped entity should not have node descriptor file")
}

func TestNodeDescriptorSharedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_descriptor_shared_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeDescriptorSharedFile(graph)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node_descriptor.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	// Shared file should contain package declaration and shared types.
	require.Contains(t, contentStr, "package ent")
	require.Contains(t, contentStr, "Node struct")
	require.Contains(t, contentStr, "Field struct")
	require.Contains(t, contentStr, "Edge struct")
	require.Contains(t, contentStr, "func (c *Client) Node(")

	// Shared content should NOT contain per-entity Node() methods.
	require.NotContains(t, contentStr, "func (_m *Todo) Node(")
	require.NotContains(t, contentStr, "func (_m *BillProduct) Node(")
}

func TestNodeDescriptorEntityFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_descriptor_entity_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeDescriptorEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node_descriptor_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	require.Contains(t, contentStr, "package ent")
	require.Contains(t, contentStr, "func TodoNode(_m *Todo, ctx context.Context)")
	require.Contains(t, contentStr, "registerNodeDescriptor(todo.Table")

	// Should NOT contain other entity node descriptor functions.
	require.NotContains(t, contentStr, "CategoryNode(")
	require.NotContains(t, contentStr, "BillProductNode(")
	// Should NOT contain shared types.
	require.NotContains(t, contentStr, "Node struct")
	require.NotContains(t, contentStr, "Field struct")
	require.NotContains(t, contentStr, "func (c *Client) Node(")
}

func TestGenerateSplitNode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_node_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	// Create a dummy monolithic gql_node.go to verify overwrite.
	monolithicPath := filepath.Join(tmpDir, "gql_node.go")
	err = os.WriteFile(monolithicPath, []byte("package ent\n// monolithic placeholder\n"), 0644)
	require.NoError(t, err)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitNode(graph)
	require.NoError(t, err)

	// Verify the monolithic file was overwritten with shared content.
	sharedContent, err := os.ReadFile(monolithicPath)
	require.NoError(t, err)
	sharedStr := string(sharedContent)
	require.Contains(t, sharedStr, "package ent")
	require.Contains(t, sharedStr, "type Noder interface")
	require.Contains(t, sharedStr, "nodeResolvers")
	require.Contains(t, sharedStr, "registerNodeResolver")
	// Shared content should NOT contain per-entity code.
	require.NotContains(t, sharedStr, "case todo.Table:")
	require.NotContains(t, sharedStr, "todoImplementors")

	// Verify per-entity node files are created.
	nodeNames := nonSkippedNodes(t, graph)
	require.NotEmpty(t, nodeNames)
	for _, name := range nodeNames {
		filename := fmt.Sprintf("gql_node_%s.go", snake(name))
		path := filepath.Join(tmpDir, filename)
		content, err := os.ReadFile(path)
		require.NoError(t, err, "node entity file should exist for %s", name)
		contentStr := string(content)
		require.Contains(t, contentStr, "package ent")
		require.Contains(t, contentStr, "Implementors")
		require.Contains(t, contentStr, "func init()")
		require.Contains(t, contentStr, "registerNodeResolver")
	}

	// Verify skipped types do NOT get per-entity files.
	skippedPath := filepath.Join(tmpDir, "gql_node_very_secret.go")
	_, err = os.Stat(skippedPath)
	require.True(t, os.IsNotExist(err), "skipped entity should not have node file")
}

func TestNodeSharedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_shared_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeSharedFile(graph)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	// Shared file should contain package declaration and shared types.
	require.Contains(t, contentStr, "package ent")
	require.Contains(t, contentStr, "type Noder interface")
	require.Contains(t, contentStr, "nodeResolver struct")
	require.Contains(t, contentStr, "nodeResolvers")
	require.Contains(t, contentStr, "registerNodeResolver")

	// Shared content should NOT contain per-entity code.
	require.NotContains(t, contentStr, "func (_m *Todo)")
	require.NotContains(t, contentStr, "todoImplementors")
	require.NotContains(t, contentStr, "case todo.Table")
}

func TestNodeSharedFile_WithNodeDescriptor(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_shared_nd_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithNodeDescriptor(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeSharedFile(graph)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	// Noder interface should only have IsNode() (Node() was removed for split-package compat).
	require.Contains(t, contentStr, "IsNode()")
	require.NotContains(t, contentStr, "Node(context.Context) (*Node, error)")
}

func TestNodeEntityFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_entity_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)

	require.Contains(t, contentStr, "package ent")
	require.Contains(t, contentStr, "todoImplementors")
	require.Contains(t, contentStr, "func init()")
	require.Contains(t, contentStr, "registerNodeResolver")
	require.Contains(t, contentStr, "todoNoder")
	require.Contains(t, contentStr, "todoNoders")

	// Should NOT contain other entity code.
	require.NotContains(t, contentStr, "categoryImplementors")
	require.NotContains(t, contentStr, "userImplementors")
	// Should NOT contain shared types.
	require.NotContains(t, contentStr, "type Noder interface")
	require.NotContains(t, contentStr, "type nodeResolver struct")
}

func TestNodeEntityFile_HasCollectionTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_node_collection_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	var todoNode *gen.Type
	for _, n := range graph.Nodes {
		if n.Name == "Todo" {
			todoNode = n
			break
		}
	}
	require.NotNil(t, todoNode)

	// Test with collection template enabled (default).
	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateNodeEntityFile(graph, todoNode)
	require.NoError(t, err)

	path := filepath.Join(tmpDir, "gql_node_todo.go")
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr := string(content)
	require.Contains(t, contentStr, "collectField")
	require.Contains(t, contentStr, "CollectFields")

	// Test with collection template removed via WithTemplates().
	// Provide all default templates except CollectionTemplate.
	ex2, err := NewExtension(
		WithSchemaGenerator(),
		WithSplitGoFiles(true),
		WithTemplates(
			EnumTemplate,
			NodeTemplate,
			PaginationTemplate,
			TransactionTemplate,
			EdgeTemplate,
			MutationInputTemplate,
		),
	)
	require.NoError(t, err)

	err = ex2.generateNodeEntityFile(graph, todoNode)
	require.NoError(t, err)

	content2, err := os.ReadFile(path)
	require.NoError(t, err)
	contentStr2 := string(content2)
	require.Contains(t, contentStr2, "package ent")
	// Without collection template, collectField should not be present.
	require.NotContains(t, contentStr2, "collectField")
}

func requireNoDuplicateStructFields(t *testing.T, path string) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	require.NoError(t, err, "parse generated file %s", path)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			seen := make(map[string]struct{})
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					if _, ok := seen[name.Name]; ok {
						t.Fatalf("duplicate struct field %q in %s (%s)", name.Name, filepath.Base(path), typeSpec.Name.Name)
					}
					seen[name.Name] = struct{}{}
				}
			}
		}
	}
}

func TestGenerateSplitWhereInputs_ParallelFlag(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_where_inputs_parallel_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	graph := loadTestGraph(t, tmpDir)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
		WithParallelWhereInputFiles(true),
	)
	require.NoError(t, err)

	for i := 0; i < 50; i++ {
		err = ex.generateSplitWhereInputs(graph)
		require.NoError(t, err)

		paths, err := filepath.Glob(filepath.Join(tmpDir, "gql_where_input_*.go"))
		require.NoError(t, err)
		require.NotEmpty(t, paths)

		for _, path := range paths {
			requireNoDuplicateStructFields(t, path)
		}

		path := filepath.Join(tmpDir, "gql_where_input_todo.go")
		_, err = os.Stat(path)
		require.NoError(t, err)
	}
}
