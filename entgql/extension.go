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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/imports"
)

// ExtensionAnnotation carries graph-level entgql configuration for code generation.
// It is injected into gen.Graph.Annotations so templates can access the values.
type ExtensionAnnotation struct {
	MaxPageSize int
}

// Name implements the ent.Annotation interface.
func (ExtensionAnnotation) Name() string { return "EntGQLExtension" }

type (
	// Extension implements the entc.Extension for providing GraphQL integration.
	Extension struct {
		schemaGenerator
		entc.DefaultExtension
		outputWriter            func(*ast.Schema) error
		hooks                   []gen.Hook
		templates               []*gen.Template
		schemaDir               string // directory for split schema output
		splitGoFiles            bool   // split generated Go files into per-entity files
		parallelWhereInputFiles bool   // opt-in parallel generation for where-input split files
		importsDir              string // real target dir for imports.Process path resolution during split generation
		maxPageSize             int    // server-side default and cap for connection pagination
	}

	// ExtensionOption allows for managing the Extension configuration
	// using functional options.
	ExtensionOption func(*Extension) error

	// SchemaHook is the hook that run after the GQL schema generation.
	SchemaHook func(*gen.Graph, *ast.Schema) error
)

// WithSchemaPath sets the filepath to the GraphQL schema to write the
// generated Ent types. If the file does not exist, it will generate a
// new schema. Please note that your gqlgen.yml config file should be
// updated as follows to support multiple schema files:
//
//	schema:
//	 - schema.graphql // existing schema.
//	 - ent.graphql	  // generated schema.
//
// WithSchemaPath and WithSchemaDir are mutually exclusive.
func WithSchemaPath(path string) ExtensionOption {
	return func(ex *Extension) error {
		if ex.schemaDir != "" {
			return fmt.Errorf("entgql: WithSchemaPath and WithSchemaDir are mutually exclusive")
		}
		ex.path = path
		return nil
	}
}

// WithSchemaDir sets the directory for split schema output, generating
// per-entity GraphQL files instead of a single file. This can help gqlgen
// generate smaller resolver files, reducing memory usage during compilation.
//
// Output structure:
//
//	<schema_dir>/
//	  ent_shared.graphql      # Shared types: directives, Node, Cursor, PageInfo, etc.
//	  ent_query.graphql       # Query type with all entity fields
//	  ent_<entity>.graphql    # Per-entity types (Connection, Edge, Order, WhereInput, etc.)
//
// Please note that your gqlgen.yml config file should be updated to use a glob pattern:
//
//	schema:
//	 - schema.graphql            // existing schema.
//	 - <schema_dir>/*.graphql    // generated split schemas.
//
// WithSchemaDir and WithSchemaPath are mutually exclusive.
func WithSchemaDir(dir string) ExtensionOption {
	return func(ex *Extension) error {
		if ex.path != "" {
			return fmt.Errorf("entgql: WithSchemaDir and WithSchemaPath are mutually exclusive")
		}
		ex.schemaDir = dir
		return nil
	}
}

// WithSplitGoFiles enables splitting generated Go files (gql_where_input.go and
// gql_mutation_input.go) into per-entity files. This can reduce build times and
// memory usage during compilation by allowing the Go compiler to process smaller
// chunks in parallel.
//
// Output structure when enabled:
//
//	ent/
//	  gql_where_input_<entity>.go     # <Entity>WhereInput + methods
//	  gql_mutation_input_<entity>.go  # Create<Entity>Input, Update<Entity>Input + methods
//
// The original monolithic files (gql_where_input.go, gql_mutation_input.go) will
// not be generated when this option is enabled.
func WithSplitGoFiles(enabled bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.splitGoFiles = enabled
		return nil
	}
}

// WithParallelWhereInputFiles configures split where-input generation mode.
//
// By default, split where-input files are generated sequentially because some
// ent/op combinations may not be concurrency-safe during template execution.
// Enable this option only when using an ent version/path verified as race-free.
func WithParallelWhereInputFiles(enabled bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.parallelWhereInputFiles = enabled
		return nil
	}
}

// WithOutputWriter sets the function to write the generated schema.
func WithOutputWriter(w func(*ast.Schema) error) ExtensionOption {
	return func(ex *Extension) error {
		ex.outputWriter = w
		return nil
	}
}

// WithSchemaHook allows users to provide a list of hooks
// to run after the GQL schema generation.
func WithSchemaHook(hooks ...SchemaHook) ExtensionOption {
	return func(ex *Extension) error {
		ex.schemaHooks = append(ex.schemaHooks, hooks...)
		return nil
	}
}

// WithConfigPath sets the filepath to gqlgen.yml configuration file
// and injects its parsed version to the global annotations.
//
// Note that, enabling this option is recommended as it improves the
// GraphQL integration,
func WithConfigPath(path string, gqlgenOptions ...api.Option) ExtensionOption {
	return func(ex *Extension) (err error) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get working directory: %w", err)
		}
		if err := os.Chdir(filepath.Dir(path)); err != nil {
			return fmt.Errorf("unable to enter config dir: %w", err)
		}
		defer func() {
			if cerr := os.Chdir(cwd); cerr != nil {
				err = fmt.Errorf("unable to restore working directory: %w", cerr)
			}
		}()
		cfg, err := config.LoadConfig(filepath.Base(path))
		if err != nil {
			return err
		}
		ex.cfg = cfg
		return nil
	}
}

// WithTemplates overrides the default templates (entgql.AllTemplates)
// with specific templates.
func WithTemplates(templates ...*gen.Template) ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = templates
		return nil
	}
}

var (
	// WithWhereFilters configures the extension to either add or
	// remove the WhereTemplate from the code generation templates.
	//
	// Deprecated: use WithWhereInputs instead. This option is planned
	// to be removed in future versions.
	WithWhereFilters = WithWhereInputs
)

// WithWhereInputs configures the extension to either add or
// remove the WhereTemplate from the code generation templates.
//
// The WhereTemplate generates GraphQL filters to all types in the ent/schema.
func WithWhereInputs(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.genWhereInput = b
		i, exists := ex.hasTemplate(WhereTemplate)
		if b && !exists {
			ex.templates = append(ex.templates, WhereTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithNodeDescriptor configures the extension to either add or
// remove the NodeDescriptorTemplate from the code generation templates.
//
// In case this option is enabled, EntGQL generates a `Node()` method for each
// type that returns its representation in one standard way. A common use case for
// this option is to develop an administrator tool on top of Ent as implemented in:
// https://github.com/ent/ent/issues/1000#issuecomment-735663175.
func WithNodeDescriptor(b bool) ExtensionOption {
	return func(ex *Extension) error {
		i, exists := ex.hasTemplate(NodeDescriptorTemplate)
		if b && !exists {
			ex.templates = append(ex.templates, NodeDescriptorTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithRelaySpec enables or disables generating the Relay Node interface.
func WithRelaySpec(enabled bool) ExtensionOption {
	return func(e *Extension) error {
		e.relaySpec = enabled
		return nil
	}
}

// WithSchemaGenerator add a hook for generate GQL schema
func WithSchemaGenerator() ExtensionOption {
	return func(e *Extension) error {
		e.genSchema = true
		return nil
	}
}

// WithMax sets a global default and cap for connection pagination.
// When a query omits `first`/`last`, this value is used as the default limit.
// When a query provides a value exceeding this, it is silently capped.
// A value of 0 (default) disables enforcement, preserving the original behavior.
func WithMax(n int) ExtensionOption {
	return func(ex *Extension) error {
		if n < 0 {
			return fmt.Errorf("entgql: WithMax value must be non-negative, got %d", n)
		}
		ex.maxPageSize = n
		return nil
	}
}

// WithMapScalarFunc allows users to provide a custom function that
// maps an ent.Field (*gen.Field) into its GraphQL scalar type. If the
// function returns an empty string, the extension fallbacks to its
// default mapping.
//
//	ex, err := entgql.NewExtension(
//		entgql.WithMapScalarFunc(func(f *gen.Field, op gen.Op) string {
//			if t, ok := knowType(f, op); ok {
//				return t
//			}
//			// Fallback to the default mapping.
//			return ""
//		}),
//	)
func WithMapScalarFunc(scalarFunc func(*gen.Field, gen.Op) string) ExtensionOption {
	return func(ex *Extension) error {
		ex.scalarFunc = scalarFunc
		return nil
	}
}

// NewExtension creates a new extension with the given configuration.
//
//	ex, err := entgql.NewExtension(
//		entgql.WithSchemaGenerator(),
//		entgql.WithSchemaPath("../ent.graphql"),
//		entgql.WithWhereInputs(true),
//	)
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates: AllTemplates,
		schemaGenerator: schemaGenerator{
			relaySpec:    true,
			genMutations: true,
		},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	// When split Go files mode is enabled, add a hook to generate per-entity files
	// after the standard templates run (which generate the monolithic files).
	// The hook will replace the monolithic files with split files.
	if ex.splitGoFiles {
		ex.hooks = append(ex.hooks, ex.genSplitGoFilesHook())
	}
	ex.hooks = append(ex.hooks, ex.genSchemaHook(), removeOldAssets)
	return ex, nil
}

// Templates of the extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks of the extension.
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

// Options of the extension.
func (e *Extension) Options() []entc.Option {
	return []entc.Option{
		entc.FeatureNames(gen.FeatureNamedEdges.Name),
	}
}

// genSchema returns a new hook for generating
// the GraphQL schema from the graph.
func (e *Extension) genSchemaHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) (err error) {
			resetSafeOpsCache()
			defer resetSafeOpsCache()
			// Inject max page size into graph annotations before template rendering
			// so that pagination templates can read it via $.Annotations.
			if e.maxPageSize > 0 {
				if g.Annotations == nil {
					g.Annotations = make(gen.Annotations)
				}
				g.Annotations[ExtensionAnnotation{}.Name()] = ExtensionAnnotation{MaxPageSize: e.maxPageSize}
			}
			if err = next.Generate(g); err != nil {
				return err
			}
			for _, t := range g.Nodes {
				for _, f := range t.DeprecatedFields() {
					ant, err := annotation(f.Annotations)
					if err != nil {
						return err
					}
					if !slices.ContainsFunc(ant.Directives, func(d Directive) bool {
						return d.Name == "deprecated"
					}) {
						ant.Directives = append(ant.Directives, Deprecated(f.DeprecationReason()))
						if f.Annotations == nil {
							f.Annotations = make(map[string]interface{})
						}
						f.Annotations[ant.Name()] = ant
					}
				}
			}
			if !(e.genSchema || e.genWhereInput || e.genMutations) {
				return nil
			}
			// Handle split schema mode
			if e.schemaDir != "" {
				return e.generateSplitSchema(g)
			}
			schema, err := e.BuildSchema(g)
			if err != nil {
				return err
			}
			if e.outputWriter == nil {
				if e.path == "" {
					return nil
				}
				return os.WriteFile(e.path, []byte(printSchema(schema)), 0644)
			}
			return e.outputWriter(schema)
		})
	}
}

// generateSplitSchema generates per-entity GraphQL schema files.
func (e *Extension) generateSplitSchema(g *gen.Graph) error {
	split, err := e.BuildSplitSchema(g)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(e.schemaDir, 0755); err != nil {
		return fmt.Errorf("entgql: failed to create schema directory: %w", err)
	}

	expected := map[string]struct{}{"ent_shared.graphql": {}}
	if split.Query != nil && len(split.Query.Types) > 0 {
		expected["ent_query.graphql"] = struct{}{}
	}
	for name, schema := range split.Entities {
		if len(schema.Types) > 0 {
			expected[fmt.Sprintf("ent_%s.graphql", snake(name))] = struct{}{}
		}
	}
	if err := cleanupSplitSchemaFiles(e.schemaDir, expected); err != nil {
		return err
	}

	sharedPath := filepath.Join(e.schemaDir, "ent_shared.graphql")
	if err := writeFileAtomic(sharedPath, []byte(printSchema(split.Shared)), 0644); err != nil {
		return fmt.Errorf("entgql: failed to write shared schema: %w", err)
	}

	if split.Query != nil && len(split.Query.Types) > 0 {
		queryPath := filepath.Join(e.schemaDir, "ent_query.graphql")
		if err := writeFileAtomic(queryPath, []byte(printSchema(split.Query)), 0644); err != nil {
			return fmt.Errorf("entgql: failed to write query schema: %w", err)
		}
	}

	var fns []func() error
	for name, schema := range split.Entities {
		if len(schema.Types) == 0 {
			continue
		}
		name := name
		schema := schema
		fns = append(fns, func() error {
			entityPath := filepath.Join(e.schemaDir, fmt.Sprintf("ent_%s.graphql", snake(name)))
			if err := writeFileAtomic(entityPath, []byte(printSchema(schema)), 0644); err != nil {
				return fmt.Errorf("entgql: failed to write entity schema %s: %w", name, err)
			}
			return nil
		})
	}
	return parallelGenerate(fns)
}

func cleanupSplitSchemaFiles(schemaDir string, keep map[string]struct{}) error {
	matches, err := filepath.Glob(filepath.Join(schemaDir, "ent_*.graphql"))
	if err != nil {
		return fmt.Errorf("entgql: list split schema files: %w", err)
	}
	for _, path := range matches {
		name := filepath.Base(path)
		if _, ok := keep[name]; ok {
			continue
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("entgql: remove stale split schema file %s: %w", name, err)
		}
	}
	return nil
}

func writeFileAtomic(path string, content []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".entgql-write-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
			return fmt.Errorf("replace destination %s: %w", path, removeErr)
		}
		if err = os.Rename(tmpName, path); err != nil {
			return fmt.Errorf("rename temp file: %w", err)
		}
	}
	return nil
}

// genSplitGoFilesHook returns a hook that generates per-entity Go files
// for WhereInput and MutationInput types.
func (e *Extension) genSplitGoFilesHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			resetSafeOpsCache()
			defer resetSafeOpsCache()
			if err := next.Generate(g); err != nil {
				return err
			}
			return e.generateSplitGoFiles(g)
		})
	}
}

// parallelGenerate runs a slice of file-generation functions concurrently,
// limiting concurrency to the number of available CPUs.
func parallelGenerate(fns []func() error) error {
	g := new(errgroup.Group)
	g.SetLimit(runtime.GOMAXPROCS(0))
	for _, fn := range fns {
		g.Go(fn)
	}
	return g.Wait()
}

// generateSplitGoFiles generates per-entity Go files for WhereInput, MutationInput,
// Pagination, Collection, Edge, NodeDescriptor, and Node types.
func (e *Extension) generateSplitGoFiles(g *gen.Graph) error {
	resetSafeOpsCache()
	defer resetSafeOpsCache()

	tmpDir, err := os.MkdirTemp("", ".entgql-split-")
	if err != nil {
		return fmt.Errorf("entgql: create split go temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	staged := *g
	staged.Target = tmpDir

	var fns []func() error

	// Shared files (one per type, no per-entity loop).
	fns = append(fns,
		func() error { return e.generatePaginationSharedFile(&staged) },
		func() error { return e.generateCollectionSharedFile(&staged) },
		func() error { return e.generateNodeSharedFile(&staged) },
	)
	if _, exists := e.hasTemplate(NodeDescriptorTemplate); exists {
		fns = append(fns, func() error { return e.generateNodeDescriptorSharedFile(&staged) })
	}

	// Per-entity where input files.
	if e.genWhereInput {
		nodes, err := filterNodes(staged.Nodes, SkipWhereInput)
		if err != nil {
			return err
		}
		if e.parallelWhereInputFiles {
			for _, n := range nodes {
				n := n
				fns = append(fns, func() error { return e.generateWhereInputFile(&staged, n) })
			}
		} else {
			for _, n := range nodes {
				if err := e.generateWhereInputFile(&staged, n); err != nil {
					return err
				}
			}
		}
	}

	// Per-entity mutation input files.
	if e.genMutations {
		inputs, err := mutationInputs(staged.Nodes)
		if err != nil {
			return err
		}
		byEntity := make(map[string][]*MutationDescriptor)
		for _, input := range inputs {
			byEntity[input.Type.Name] = append(byEntity[input.Type.Name], input)
		}
		for name, entityInputs := range byEntity {
			name := name
			entityInputs := entityInputs
			fns = append(fns, func() error { return e.generateMutationInputFile(&staged, name, entityInputs) })
		}
	}

	// Per-entity pagination, collection, edge, node descriptor, and node files.
	nodes, err := filterNodes(staged.Nodes, SkipType)
	if err != nil {
		return err
	}
	_, hasNodeDescriptor := e.hasTemplate(NodeDescriptorTemplate)
	for _, n := range nodes {
		n := n
		fns = append(fns,
			func() error { return e.generatePaginationEntityFile(&staged, n) },
			func() error { return e.generateCollectionEntityFile(&staged, n) },
			func() error { return e.generateNodeEntityFile(&staged, n) },
		)
		if hasNodeDescriptor {
			fns = append(fns, func() error { return e.generateNodeDescriptorEntityFile(&staged, n) })
		}
		edges, err := filterEdges(n.Edges, SkipType)
		if err != nil {
			return err
		}
		if len(edges) > 0 {
			fns = append(fns, func() error { return e.generateEdgeEntityFile(&staged, n) })
		}
	}

	if err := parallelGenerate(fns); err != nil {
		return err
	}

	keep, names, err := listGeneratedGoFiles(tmpDir)
	if err != nil {
		return err
	}
	if err := promoteGeneratedGoFiles(tmpDir, g.Target, names); err != nil {
		return err
	}
	return cleanupSplitGeneratedGoFiles(g.Target, keep)
}

func listGeneratedGoFiles(dir string) (map[string]struct{}, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("entgql: read split go temp dir: %w", err)
	}
	keep := make(map[string]struct{})
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".go" {
			continue
		}
		keep[name] = struct{}{}
		names = append(names, name)
	}
	slices.Sort(names)
	return keep, names, nil
}

func promoteGeneratedGoFiles(fromDir, targetDir string, names []string) error {
	for _, name := range names {
		from := filepath.Join(fromDir, name)
		to := filepath.Join(targetDir, name)
		data, err := os.ReadFile(from)
		if err != nil {
			return fmt.Errorf("entgql: read generated go file %s: %w", name, err)
		}
		if err := os.WriteFile(to, data, 0644); err != nil {
			return fmt.Errorf("entgql: write generated go file %s: %w", name, err)
		}
	}
	return nil
}

func cleanupSplitGeneratedGoFiles(targetDir string, keep map[string]struct{}) error {
	patterns := []string{
		"gql_where_input.go",
		"gql_where_input_*.go",
		"gql_mutation_input.go",
		"gql_mutation_input_*.go",
		"gql_edge.go",
		"gql_edge_*.go",
		"gql_pagination.go",
		"gql_pagination_*.go",
		"gql_collection.go",
		"gql_collection_*.go",
		"gql_node.go",
		"gql_node_*.go",
		"gql_node_descriptor.go",
		"gql_node_descriptor_*.go",
	}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(targetDir, pattern))
		if err != nil {
			return fmt.Errorf("entgql: list generated go files (%s): %w", pattern, err)
		}
		for _, path := range matches {
			name := filepath.Base(path)
			if _, ok := keep[name]; ok {
				continue
			}
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("entgql: remove stale generated go file %s: %w", name, err)
			}
		}
	}
	return nil
}

// generateSplitWhereInputs generates per-entity where input files.
func (e *Extension) generateSplitWhereInputs(g *gen.Graph) error {
	resetSafeOpsCache()
	defer resetSafeOpsCache()

	nodes, err := filterNodes(g.Nodes, SkipWhereInput)
	if err != nil {
		return err
	}
	if e.parallelWhereInputFiles {
		var fns []func() error
		for _, n := range nodes {
			n := n
			fns = append(fns, func() error { return e.generateWhereInputFile(g, n) })
		}
		return parallelGenerate(fns)
	}
	for _, n := range nodes {
		if err := e.generateWhereInputFile(g, n); err != nil {
			return err
		}
	}
	return nil
}

// generateSplitMutationInputs generates per-entity mutation input files.
func (e *Extension) generateSplitMutationInputs(g *gen.Graph) error {
	inputs, err := mutationInputs(g.Nodes)
	if err != nil {
		return err
	}
	byEntity := make(map[string][]*MutationDescriptor)
	for _, input := range inputs {
		byEntity[input.Type.Name] = append(byEntity[input.Type.Name], input)
	}
	var fns []func() error
	for name, entityInputs := range byEntity {
		name := name
		entityInputs := entityInputs
		fns = append(fns, func() error { return e.generateMutationInputFile(g, name, entityInputs) })
	}
	return parallelGenerate(fns)
}

// generateSplitPagination overwrites the monolithic pagination file with shared-only
// content, then generates per-entity pagination files.
func (e *Extension) generateSplitPagination(g *gen.Graph) error {
	if err := e.generatePaginationSharedFile(g); err != nil {
		return err
	}
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		n := n
		fns = append(fns, func() error { return e.generatePaginationEntityFile(g, n) })
	}
	return parallelGenerate(fns)
}

// generateSplitCollection overwrites the monolithic collection file with shared-only
// content, then generates per-entity collection files.
func (e *Extension) generateSplitCollection(g *gen.Graph) error {
	if err := e.generateCollectionSharedFile(g); err != nil {
		return err
	}
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		n := n
		fns = append(fns, func() error { return e.generateCollectionEntityFile(g, n) })
	}
	return parallelGenerate(fns)
}

// generateSplitEdge generates per-entity edge files.
func (e *Extension) generateSplitEdge(g *gen.Graph) error {
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		n := n
		edges, err := filterEdges(n.Edges, SkipType)
		if err != nil {
			return err
		}
		if len(edges) > 0 {
			fns = append(fns, func() error { return e.generateEdgeEntityFile(g, n) })
		}
	}
	return parallelGenerate(fns)
}

// generateSplitNodeDescriptor overwrites the monolithic node descriptor file with shared-only
// content, then generates per-entity node descriptor files.
func (e *Extension) generateSplitNodeDescriptor(g *gen.Graph) error {
	if err := e.generateNodeDescriptorSharedFile(g); err != nil {
		return err
	}
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		n := n
		fns = append(fns, func() error { return e.generateNodeDescriptorEntityFile(g, n) })
	}
	return parallelGenerate(fns)
}

// generateSplitNode overwrites the monolithic node file with shared-only
// content, then generates per-entity node files.
func (e *Extension) generateSplitNode(g *gen.Graph) error {
	if err := e.generateNodeSharedFile(g); err != nil {
		return err
	}
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		n := n
		fns = append(fns, func() error { return e.generateNodeEntityFile(g, n) })
	}
	return parallelGenerate(fns)
}

// generateWhereInputFile generates a where input file for a single entity.
func (e *Extension) generateWhereInputFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_where_input_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)

	tmpl := WhereInputEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node *gen.Type
	}{g, n}); err != nil {
		return fmt.Errorf("entgql: execute where_input_entity template for %s: %w", n.Name, err)
	}

	// Format and add missing imports using goimports
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format where_input for %s: %w", n.Name, err)
	}

	return os.WriteFile(path, content, 0644)
}

// generateMutationInputFile generates a mutation input file for a single entity.
func (e *Extension) generateMutationInputFile(g *gen.Graph, name string, inputs []*MutationDescriptor) error {
	filename := fmt.Sprintf("gql_mutation_input_%s.go", snake(name))
	path := filepath.Join(g.Target, filename)

	tmpl := MutationInputEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		EntityName string
		Inputs     []*MutationDescriptor
	}{g, name, inputs}); err != nil {
		return fmt.Errorf("entgql: execute mutation_input_entity template for %s: %w", name, err)
	}

	// Format and add missing imports using goimports
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format mutation_input for %s: %w", name, err)
	}

	return os.WriteFile(path, content, 0644)
}

// generatePaginationSharedFile overwrites gql_pagination.go with shared-only content.
func (e *Extension) generatePaginationSharedFile(g *gen.Graph) error {
	path := filepath.Join(g.Target, "gql_pagination.go")
	tmpl := PaginationSharedTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
	}{g}); err != nil {
		return fmt.Errorf("entgql: execute pagination_shared template: %w", err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format pagination_shared: %w", err)
	}
	return os.WriteFile(path, content, 0644)
}

// generatePaginationEntityFile generates a pagination file for a single entity.
func (e *Extension) generatePaginationEntityFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_pagination_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)
	tmpl := PaginationEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node *gen.Type
	}{g, n}); err != nil {
		return fmt.Errorf("entgql: execute pagination_entity template for %s: %w", n.Name, err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format pagination for %s: %w", n.Name, err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateCollectionSharedFile overwrites gql_collection.go with shared-only content.
func (e *Extension) generateCollectionSharedFile(g *gen.Graph) error {
	path := filepath.Join(g.Target, "gql_collection.go")
	tmpl := CollectionSharedTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
	}{g}); err != nil {
		return fmt.Errorf("entgql: execute collection_shared template: %w", err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format collection_shared: %w", err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateCollectionEntityFile generates a collection file for a single entity.
func (e *Extension) generateCollectionEntityFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_collection_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)
	tmpl := CollectionEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{g, n, e.genWhereInput}); err != nil {
		return fmt.Errorf("entgql: execute collection_entity template for %s: %w", n.Name, err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format collection for %s: %w", n.Name, err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateEdgeEntityFile generates an edge file for a single entity.
func (e *Extension) generateEdgeEntityFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_edge_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)
	tmpl := EdgeEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node                  *gen.Type
		HasWhereInputTemplate bool
	}{g, n, e.genWhereInput}); err != nil {
		return fmt.Errorf("entgql: execute edge_entity template for %s: %w", n.Name, err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format edge for %s: %w", n.Name, err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateNodeDescriptorSharedFile overwrites gql_node_descriptor.go with shared-only content.
func (e *Extension) generateNodeDescriptorSharedFile(g *gen.Graph) error {
	path := filepath.Join(g.Target, "gql_node_descriptor.go")
	tmpl := NodeDescriptorSharedTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
	}{g}); err != nil {
		return fmt.Errorf("entgql: execute node_descriptor_shared template: %w", err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format node_descriptor_shared: %w", err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateNodeDescriptorEntityFile generates a node descriptor file for a single entity.
func (e *Extension) generateNodeDescriptorEntityFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_node_descriptor_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)
	tmpl := NodeDescriptorEntityTemplate
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node *gen.Type
	}{g, n}); err != nil {
		return fmt.Errorf("entgql: execute node_descriptor_entity template for %s: %w", n.Name, err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format node_descriptor for %s: %w", n.Name, err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateNodeSharedFile overwrites gql_node.go with shared-only content.
func (e *Extension) generateNodeSharedFile(g *gen.Graph) error {
	path := filepath.Join(g.Target, "gql_node.go")
	tmpl := NodeSharedTemplate
	_, hasNodeDescriptor := e.hasTemplate(NodeDescriptorTemplate)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		HasNodeDescriptorTemplate bool
	}{g, hasNodeDescriptor}); err != nil {
		return fmt.Errorf("entgql: execute node_shared template: %w", err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format node_shared: %w", err)
	}
	return os.WriteFile(path, content, 0644)
}

// generateNodeEntityFile generates a node file for a single entity.
func (e *Extension) generateNodeEntityFile(g *gen.Graph, n *gen.Type) error {
	filename := fmt.Sprintf("gql_node_%s.go", snake(n.Name))
	path := filepath.Join(g.Target, filename)
	tmpl := NodeEntityTemplate
	_, hasCollection := e.hasTemplate(CollectionTemplate)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		*gen.Graph
		Node                  *gen.Type
		HasCollectionTemplate bool
	}{g, n, hasCollection}); err != nil {
		return fmt.Errorf("entgql: execute node_entity template for %s: %w", n.Name, err)
	}
	content, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("entgql: format node for %s: %w", n.Name, err)
	}
	return os.WriteFile(path, content, 0644)
}

// hasTemplate reports if the template exists
// in the template list and returns its index.
func (e *Extension) hasTemplate(tem *gen.Template) (int, bool) {
	name := tem.Name()
	for i := range e.templates {
		if e.templates[i].Name() == name {
			return i, true
		}
	}
	return -1, false
}

var (
	_ entc.Extension = (*Extension)(nil)

	camel    = gen.Funcs["camel"].(func(string) string)
	pascal   = gen.Funcs["pascal"].(func(string) string)
	plural   = gen.Funcs["plural"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	snake    = gen.Funcs["snake"].(func(string) string)
)
