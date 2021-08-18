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
	"os"
	"path/filepath"
	"strings"
	"text/template"
	tparse "text/template/parse"

	"entgo.io/contrib/entgql/internal"

	"entgo.io/ent/entc/gen"
	_ "github.com/go-bindata/go-bindata"
)

var (
	// CollectionTemplate adds fields collection support using auto eager-load ent edges.
	// More info can be found here: https://spec.graphql.org/June2018/#sec-Field-Collection.
	CollectionTemplate = parse("template/collection.tmpl")

	// EnumTemplate adds a template implementing MarshalGQL/UnmarshalGQL methods for enums.
	EnumTemplate = parse("template/enum.tmpl")

	// NodeTemplate implements the Relay Node interface for all types.
	NodeTemplate = parse("template/node.tmpl")

	// PaginationTemplate adds pagination support according to the GraphQL Cursor Connections Spec.
	// More info can be found in the following link: https://relay.dev/graphql/connections.htm.
	PaginationTemplate = parse("template/pagination.tmpl")

	// TransactionTemplate adds support for ent.Client for opening transactions for the transaction
	// middleware. See transaction.go for for information.
	TransactionTemplate = parse("template/transaction.tmpl")

	// EdgeTemplate adds edge resolution using eager-loading with a query fallback.
	EdgeTemplate = parse("template/edge.tmpl")

	// WhereTemplate adds a template for generating <T>WhereInput filters for each schema type.
	WhereTemplate = parse("template/where_input.tmpl")

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
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template

func parse(path string) *gen.Template {
	text := string(internal.MustAsset(path))
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Funcs(TemplateFuncs).
		Parse(text))
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

func removeOldGeneratedFiles() gen.Hook {
	prefix := "gql_"
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			for _, rootTemplate := range AllTemplates {
				for _, t := range rootTemplate.Templates() {
					// Check if template is empty
					if tparse.IsEmptyTree(t.Root) {
						continue
					}
					// Check if template has correct prefix
					if !strings.HasPrefix(t.Name(), prefix) {
						continue
					}
					name := strings.TrimPrefix(t.Name(), prefix)
					err := deleteOldTemplateByName(g, name)
					if err != nil {
						return err
					}
				}
			}
			return next.Generate(g)
		})
	}
}

func deleteOldTemplateByName(g *gen.Graph, name string) error {
	// Check if name already taken by existing schema field
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
