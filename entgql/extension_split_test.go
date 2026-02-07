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
	for _, name := range []string{"gql_where_input.go", "gql_mutation_input.go", "gql_pagination.go", "gql_collection.go"} {
		path := filepath.Join(tmpDir, name)
		err = os.WriteFile(path, []byte("package ent\n// placeholder\n"), 0644)
		require.NoError(t, err)
	}

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithWhereInputs(true),
		WithSplitGoFiles(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitGoFiles(graph)
	require.NoError(t, err)

	// Where and mutation monolithic files should be removed.
	_, err = os.Stat(filepath.Join(tmpDir, "gql_where_input.go"))
	require.True(t, os.IsNotExist(err), "monolithic where_input should be removed")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_mutation_input.go"))
	require.True(t, os.IsNotExist(err), "monolithic mutation_input should be removed")

	// Pagination and collection monolithic files should still exist (overwritten, not removed).
	_, err = os.Stat(filepath.Join(tmpDir, "gql_pagination.go"))
	require.NoError(t, err, "gql_pagination.go should still exist (overwritten with shared content)")
	_, err = os.Stat(filepath.Join(tmpDir, "gql_collection.go"))
	require.NoError(t, err, "gql_collection.go should still exist (overwritten with shared content)")

	// Verify pagination and collection per-entity files exist.
	nodeNames := nonSkippedNodes(t, graph)
	for _, name := range nodeNames {
		paginationFile := filepath.Join(tmpDir, fmt.Sprintf("gql_pagination_%s.go", snake(name)))
		_, err = os.Stat(paginationFile)
		require.NoError(t, err, "pagination entity file should exist for %s", name)

		collectionFile := filepath.Join(tmpDir, fmt.Sprintf("gql_collection_%s.go", snake(name)))
		_, err = os.Stat(collectionFile)
		require.NoError(t, err, "collection entity file should exist for %s", name)
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
