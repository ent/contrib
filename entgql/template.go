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

package entgql

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

var (
	// CollectionTemplate adds fields collection support using auto eager-load ent edges.
	// More info can be found here: https://spec.graphql.org/June2018/#sec-Field-Collection.
	CollectionTemplate = parseT("template/collection.tmpl")

	// EnumTemplate adds a template implementing MarshalGQL/UnmarshalGQL methods for enums.
	EnumTemplate = parseT("template/enum.tmpl")

	// NodeTemplate implements the Relay Node interface for all types.
	NodeTemplate = parseT("template/node.tmpl")

	// PaginationTemplate adds pagination support according to the GraphQL Cursor Connections Spec.
	// More info can be found in the following link: https://relay.dev/graphql/connections.htm.
	PaginationTemplate = parseT("template/pagination.tmpl")

	// TransactionTemplate adds support for ent.Client for opening transactions for the transaction
	// middleware. See transaction.go for for information.
	TransactionTemplate = parseT("template/transaction.tmpl")

	// EdgeTemplate adds edge resolution using eager-loading with a query fallback.
	EdgeTemplate = parseT("template/edge.tmpl")

	// WhereTemplate adds a template for generating <T>WhereInput filters for each schema type.
	WhereTemplate = parseT("template/where_input.tmpl")

	// AllTemplates holds all templates for extending ent to support GraphQL.
	AllTemplates = []*gen.Template{
		CollectionTemplate,
		EnumTemplate,
		NodeTemplate,
		PaginationTemplate,
		TransactionTemplate,
		EdgeTemplate,
	}

	// TemplateFuncs contains the extra template functions used by entgql.
	TemplateFuncs = template.FuncMap{
		"fieldCollections":    fieldCollections,
		"filterEdges":         filterEdges,
		"filterFields":        filterFields,
		"orderFields":         orderFields,
		"filterNodes":         filterNodes,
		"findIDType":          findIDType,
		"nodePaginationNames": nodePaginationNames,
		"skipMode":            skipModeFromString,
		"isSkipMode":          isSkipMode,
	}

	//go:embed template/*
	templates embed.FS
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(TemplateFuncs).
		ParseFS(templates, path))
}

// findIDType returns the type of the ID field of the given type.
func findIDType(nodes []*gen.Type, defaultType *field.TypeInfo) (*field.TypeInfo, error) {
	t := defaultType
	if len(nodes) > 0 {
		t = nodes[0].ID.Type
		// Ensure all id types have the same type.
		for _, n := range nodes[1:] {
			if n.ID.Type.Type != t.Type {
				return nil, errors.New("node does not support multiple id types")
			}
		}
	}
	return t, nil
}

type fieldCollection struct {
	Name    string
	Mapping []string
}

func fieldCollections(edges []*gen.Edge) (map[string]fieldCollection, error) {
	result := make(map[string]fieldCollection)
	for _, e := range edges {
		result[e.Name] = fieldCollection{
			Name:    e.Type.Name,
			Mapping: []string{e.Name},
		}
		ant, err := annotation(e.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Unbind {
			delete(result, e.Name)
		}
		if len(ant.Mapping) > 0 {
			if _, bind := result[e.Name]; bind {
				return nil, errors.New("bind and mapping annotations are mutually exclusive")
			}
			result[e.Name] = fieldCollection{
				Name:    e.Type.Name,
				Mapping: ant.Mapping,
			}
		}
	}
	return result, nil
}

// filterNodes filters out nodes that should not be included in the GraphQL schema.
func filterNodes(nodes []*gen.Type, skip SkipMode) ([]*gen.Type, error) {
	filteredNodes := make([]*gen.Type, 0, len(nodes))
	for _, n := range nodes {
		ant, err := annotation(n.Annotations)
		if err != nil {
			return nil, err
		}
		if !ant.Skip.Is(skip) {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes, nil
}

// filterEdges filters out edges that should not be included in the GraphQL schema.
func filterEdges(edges []*gen.Edge, skip SkipMode) ([]*gen.Edge, error) {
	filteredEdges := make([]*gen.Edge, 0, len(edges))
	for _, e := range edges {
		antE, err := annotation(e.Annotations)
		if err != nil {
			return nil, err
		}
		antT, err := annotation(e.Type.Annotations)
		if err != nil {
			return nil, err
		}
		if !antE.Skip.Is(skip) && !antT.Skip.Is(skip) {
			filteredEdges = append(filteredEdges, e)
		}
	}
	return filteredEdges, nil
}

// filterFields filters out fields that should not be included in the GraphQL schema.
func filterFields(fields []*gen.Field, skip SkipMode) ([]*gen.Field, error) {
	filteredFields := make([]*gen.Field, 0, len(fields))
	for _, f := range fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if !ant.Skip.Is(skip) {
			filteredFields = append(filteredFields, f)
		}
	}
	return filteredFields, nil
}

// orderFields returns the fields of the given node with the `OrderField` annotation.
func orderFields(n *gen.Type) ([]*gen.Field, error) {
	var ordered []*gen.Field
	for _, f := range n.Fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.OrderField == "" {
			continue
		}
		if !f.Type.Comparable() {
			return nil, fmt.Errorf("entgql: ordered field %s.%s must be comparable", n.Name, f.Name)
		}
		if ant.Skip.Is(SkipOrderField) {
			return nil, fmt.Errorf("entgql: ordered field %s.%s cannot be skipped", n.Name, f.Name)
		}
		ordered = append(ordered, f)
	}
	return ordered, nil
}

// skipModeFromString returns SkipFlag from a string
func skipModeFromString(s string) (SkipMode, error) {
	switch s {
	case "type":
		return SkipType, nil
	case "enum_field":
		return SkipEnumField, nil
	case "order_field":
		return SkipOrderField, nil
	case "where_input":
		return SkipWhereInput, nil
	}
	return 0, fmt.Errorf("invalid skip mode: %s", s)
}

func isSkipMode(antSkip interface{}, m string) (bool, error) {
	skip, err := skipModeFromString(m)
	if err != nil || antSkip == nil {
		return false, err
	}
	if raw, ok := antSkip.(float64); ok {
		return SkipMode(raw).Is(skip), nil
	}
	return false, fmt.Errorf("invalid annotation skip: %v", antSkip)
}

// PaginationNames holds the names of the pagination fields.
type PaginationNames struct {
	Connection string
	Edge       string
	Node       string
	Order      string
	OrderField string
}

// nodePaginationNames returns the names of the pagination types for the node.
func nodePaginationNames(t *gen.Type) (*PaginationNames, error) {
	node := t.Name
	ant, err := annotation(t.Annotations)
	if err != nil {
		return nil, err
	}
	if ant.Type != "" {
		node = ant.Type
	}
	return &PaginationNames{
		Connection: fmt.Sprintf("%sConnection", node),
		Edge:       fmt.Sprintf("%sEdge", node),
		Node:       node,
		Order:      fmt.Sprintf("%sOrder", node),
		OrderField: fmt.Sprintf("%sOrderField", node),
	}, nil
}

// removeOldAssets removes files that were generated before v0.1.0.
func removeOldAssets(next gen.Generator) gen.Generator {
	const prefix = "gql_"
	templates := []*gen.Template{WhereTemplate}
	templates = append(templates, AllTemplates...)
	return gen.GenerateFunc(func(g *gen.Graph) error {
		for _, rootT := range templates {
			for _, t := range rootT.Templates() {
				if parse.IsEmptyTree(t.Root) {
					continue
				}
				if !strings.HasPrefix(t.Name(), prefix) {
					continue
				}
				name := strings.TrimPrefix(t.Name(), prefix)
				if err := removeOldTemplate(g, name); err != nil {
					return err
				}
			}
		}
		return next.Generate(g)
	})
}

func removeOldTemplate(g *gen.Graph, name string) error {
	// Check if name already taken by existing schema field.
	for _, n := range g.Nodes {
		if n.Package() == name {
			return nil
		}
	}
	err := os.Remove(filepath.Join(g.Target, name+".go"))
	if !os.IsNotExist(err) {
		return err
	}
	return nil
}
