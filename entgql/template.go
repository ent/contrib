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
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"

	"entgo.io/ent/entc/gen"
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
		"filterNodes":  filterNodes,
		"filterEdges":  filterEdges,
		"filterFields": filterFields,
	}

	//go:embed template/*
	templates embed.FS
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Funcs(TemplateFuncs).
		ParseFS(templates, path))
}

func filterNodes(nodes []*gen.Type) ([]*gen.Type, error) {
	var filteredNodes []*gen.Type
	for _, n := range nodes {
		ant := &Annotation{}
		if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes, nil
}

func filterEdges(edges []*gen.Edge) ([]*gen.Edge, error) {
	var filteredEdges []*gen.Edge
	for _, e := range edges {
		ant := &Annotation{}
		if e.Annotations != nil && e.Annotations[ant.Name()] != nil {
			if err := ant.Decode(e.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		// Check if type is skipped
		if e.Type.Annotations != nil && e.Type.Annotations[ant.Name()] != nil {
			if err := ant.Decode(e.Type.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredEdges = append(filteredEdges, e)
	}
	return filteredEdges, nil
}

func filterFields(fields []*gen.Field) ([]*gen.Field, error) {
	var filteredFields []*gen.Field
	for _, f := range fields {
		ant := &Annotation{}
		if f.Annotations != nil && f.Annotations[ant.Name()] != nil {
			if err := ant.Decode(f.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}
		}
		filteredFields = append(filteredFields, f)
	}
	return filteredFields, nil
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
