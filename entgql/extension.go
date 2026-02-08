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

type (
	// Extension implements the entc.Extension for providing GraphQL integration.
	Extension struct {
		schemaGenerator
		entc.DefaultExtension
		outputWriter func(*ast.Schema) error
		hooks        []gen.Hook
		templates    []*gen.Template
		schemaDir    string // directory for split schema output
		splitGoFiles bool   // split generated Go files into per-entity files
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
	// Create schema directory if needed
	if err := os.MkdirAll(e.schemaDir, 0755); err != nil {
		return fmt.Errorf("entgql: failed to create schema directory: %w", err)
	}
	// Write shared types
	sharedPath := filepath.Join(e.schemaDir, "ent_shared.graphql")
	if err := os.WriteFile(sharedPath, []byte(printSchema(split.Shared)), 0644); err != nil {
		return fmt.Errorf("entgql: failed to write shared schema: %w", err)
	}
	// Write query type
	if split.Query != nil && len(split.Query.Types) > 0 {
		queryPath := filepath.Join(e.schemaDir, "ent_query.graphql")
		if err := os.WriteFile(queryPath, []byte(printSchema(split.Query)), 0644); err != nil {
			return fmt.Errorf("entgql: failed to write query schema: %w", err)
		}
	}
	// Write per-entity schemas in parallel.
	var fns []func() error
	for name, schema := range split.Entities {
		fns = append(fns, func() error {
			entityPath := filepath.Join(e.schemaDir, fmt.Sprintf("ent_%s.graphql", snake(name)))
			if err := os.WriteFile(entityPath, []byte(printSchema(schema)), 0644); err != nil {
				return fmt.Errorf("entgql: failed to write entity schema %s: %w", name, err)
			}
			return nil
		})
	}
	return parallelGenerate(fns)
}

// genSplitGoFilesHook returns a hook that generates per-entity Go files
// for WhereInput and MutationInput types.
func (e *Extension) genSplitGoFilesHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
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
	// Remove old monolithic files to avoid duplicate type definitions.
	if err := e.removeMonolithicGoFiles(g); err != nil {
		return err
	}

	// Collect all independent file-generation tasks into a single slice
	// so they can be executed concurrently.
	var fns []func() error

	// Shared files (one per type, no per-entity loop).
	fns = append(fns,
		func() error { return e.generatePaginationSharedFile(g) },
		func() error { return e.generateCollectionSharedFile(g) },
		func() error { return e.generateNodeSharedFile(g) },
	)
	if _, exists := e.hasTemplate(NodeDescriptorTemplate); exists {
		fns = append(fns, func() error { return e.generateNodeDescriptorSharedFile(g) })
	}

	// Per-entity where input files.
	if e.genWhereInput {
		nodes, err := filterNodes(g.Nodes, SkipWhereInput)
		if err != nil {
			return err
		}
		for _, n := range nodes {
			fns = append(fns, func() error { return e.generateWhereInputFile(g, n) })
		}
	}

	// Per-entity mutation input files.
	if e.genMutations {
		inputs, err := mutationInputs(g.Nodes)
		if err != nil {
			return err
		}
		byEntity := make(map[string][]*MutationDescriptor)
		for _, input := range inputs {
			byEntity[input.Type.Name] = append(byEntity[input.Type.Name], input)
		}
		for name, entityInputs := range byEntity {
			fns = append(fns, func() error { return e.generateMutationInputFile(g, name, entityInputs) })
		}
	}

	// Per-entity pagination, collection, edge, node descriptor, and node files.
	nodes, err := filterNodes(g.Nodes, SkipType)
	if err != nil {
		return err
	}
	_, hasNodeDescriptor := e.hasTemplate(NodeDescriptorTemplate)
	for _, n := range nodes {
		fns = append(fns,
			func() error { return e.generatePaginationEntityFile(g, n) },
			func() error { return e.generateCollectionEntityFile(g, n) },
			func() error { return e.generateNodeEntityFile(g, n) },
		)
		if hasNodeDescriptor {
			fns = append(fns, func() error { return e.generateNodeDescriptorEntityFile(g, n) })
		}
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

// removeMonolithicGoFiles removes the old monolithic Go files when split mode is enabled.
func (e *Extension) removeMonolithicGoFiles(g *gen.Graph) error {
	filesToRemove := []string{
		"gql_where_input.go",
		"gql_mutation_input.go",
		"gql_edge.go",
	}
	for _, filename := range filesToRemove {
		path := filepath.Join(g.Target, filename)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("entgql: failed to remove monolithic file %s: %w", filename, err)
		}
	}
	return nil
}


// generateSplitWhereInputs generates per-entity where input files.
func (e *Extension) generateSplitWhereInputs(g *gen.Graph) error {
	nodes, err := filterNodes(g.Nodes, SkipWhereInput)
	if err != nil {
		return err
	}
	var fns []func() error
	for _, n := range nodes {
		fns = append(fns, func() error { return e.generateWhereInputFile(g, n) })
	}
	return parallelGenerate(fns)
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
