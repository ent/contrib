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

package main

import (
	"embed"
	"flag"
	"fmt"
	"path"
	"strconv"
	"strings"
	"text/template"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	entSchemaPath *string
	snake         = gen.Funcs["snake"].(func(string) string)
	status        = protogen.GoImportPath("google.golang.org/grpc/status")
	codes         = protogen.GoImportPath("google.golang.org/grpc/codes")
)

func main() {
	var flags flag.FlagSet
	entSchemaPath = flags.String("schema_path", "", "ent schema path")
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(plg *protogen.Plugin) error {
		g, err := entc.LoadGraph(*entSchemaPath, &gen.Config{})
		if err != nil {
			return err
		}
		for _, f := range plg.Files {
			if !f.Generate {
				continue
			}
			if err := processFile(plg, f, g); err != nil {
				return err
			}
		}
		return nil
	})
}

// processFile generates service implementations from all services defined in the file.
func processFile(gen *protogen.Plugin, file *protogen.File, graph *gen.Graph) error {
	if len(file.Services) == 0 {
		return nil
	}
	adapter, err := entproto.LoadAdapter(graph)
	if err != nil {
		return err
	}
	for _, s := range file.Services {
		if name := string(s.Desc.Name()); !containsSvc(adapter, name) {
			continue
		}
		sg, err := newServiceGenerator(gen, file, graph, adapter, s)
		if err != nil {
			return err
		}
		if err := sg.generate(); err != nil {
			return err
		}
	}
	return nil
}

// containsSvc reports if the service definition for svc is created by the adapter.
func containsSvc(adapter *entproto.Adapter, svc string) bool {
	for _, d := range adapter.AllFileDescriptors() {
		for _, s := range d.GetServices() {
			if s.GetName() == svc {
				return true
			}
		}
	}
	return false
}

func newServiceGenerator(plugin *protogen.Plugin, file *protogen.File, graph *gen.Graph, adapter *entproto.Adapter, service *protogen.Service) (*serviceGenerator, error) {
	typ, err := extractEntTypeName(service, graph)
	if err != nil {
		return nil, err
	}
	filename := file.GeneratedFilenamePrefix + "_" + snake(service.GoName) + ".go"
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)
	fieldMap, err := adapter.FieldMap(typ.Name)
	if err != nil {
		return nil, err
	}
	return &serviceGenerator{
		GeneratedFile: g,
		EntPackage:    protogen.GoImportPath(graph.Config.Package),
		File:          file,
		Service:       service,
		EntType:       typ,
		FieldMap:      fieldMap,
	}, nil
}

func (g *serviceGenerator) generate() error {
	tmpl, err := gen.NewTemplate("service").
		Funcs(template.FuncMap{
			"ident":        g.QualifiedGoIdent,
			"entIdent":     g.entIdent,
			"newConverter": g.newConverter,
			"unquote":      strconv.Unquote,
			"qualify": func(pkg, ident string) string {
				return g.QualifiedGoIdent(protogen.GoImportPath(pkg).Ident(ident))
			},
			"protoIdentNormalize": entproto.NormalizeEnumIdentifier,
			"statusErr": func(code, msg string) string {
				return fmt.Sprintf("%s(%s, %q)",
					g.QualifiedGoIdent(status.Ident("Error")),
					g.QualifiedGoIdent(codes.Ident(code)),
					msg,
				)
			},
			"statusErrf": func(code, format string, args ...string) string {
				return fmt.Sprintf("%s(%s, %s, %s)",
					g.QualifiedGoIdent(status.Ident("Errorf")),
					g.QualifiedGoIdent(codes.Ident(code)),
					strconv.Quote(format),
					strings.Join(args, ","),
				)
			},
			"method": func(m *protogen.Method) *methodInput {
				return &methodInput{
					G:      g,
					Method: m,
				}
			},
		}).
		ParseFS(templates, "template/*.tmpl")
	if err != nil {
		return err
	}
	if err := tmpl.ExecuteTemplate(g, "service", g); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}
	return nil
}

type (
	serviceGenerator struct {
		*protogen.GeneratedFile
		EntPackage protogen.GoImportPath
		File       *protogen.File
		Service    *protogen.Service
		EntType    *gen.Type
		FieldMap   entproto.FieldMap
	}
	methodInput struct {
		G      *serviceGenerator
		Method *protogen.Method
	}
)

//go:embed template/*
var templates embed.FS

func extractEntTypeName(s *protogen.Service, g *gen.Graph) (*gen.Type, error) {
	typeName := strings.TrimSuffix(s.GoName, "Service")
	for _, gt := range g.Nodes {
		if gt.Name == typeName {
			return gt, nil
		}
	}
	return nil, fmt.Errorf("entproto: type %q of service %q not found in graph", typeName, s.GoName)
}

func (g *serviceGenerator) entIdent(subpath string, ident string) protogen.GoIdent {
	ip := path.Join(string(g.EntPackage), subpath)
	return protogen.GoImportPath(ip).Ident(ident)
}
