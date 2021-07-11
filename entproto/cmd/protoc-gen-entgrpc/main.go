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

package main

import (
	"flag"
	"fmt"
	"path"
	"strings"
	"text/template"

	"entgo.io/contrib/entproto"
	"entgo.io/contrib/entproto/cmd/protoc-gen-entgrpc/internal"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	entSchemaPath *string
	snake         = gen.Funcs["snake"].(func(string) string)
	contextImp    = protogen.GoImportPath("context")
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
	for _, s := range file.Services {
		sg, err := newServiceGenerator(gen, file, graph, s)
		if err != nil {
			return err
		}
		if err := sg.generate(); err != nil {
			return err
		}
	}
	return nil
}

func newServiceGenerator(plugin *protogen.Plugin, file *protogen.File, graph *gen.Graph, service *protogen.Service) (*serviceGenerator, error) {
	adapter, err := entproto.LoadAdapter(graph)
	if err != nil {
		return nil, err
	}
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
	tmpl := template.New("").
		Funcs(gen.Funcs).
		Funcs(template.FuncMap{
			"ident":    g.QualifiedGoIdent,
			"entIdent": g.entIdent,
			"qualify": func(pkg, ident string) string {
				return g.QualifiedGoIdent(protogen.GoImportPath(pkg).Ident(ident))
			},
		})
	for _, t := range templates {
		if _, err := tmpl.Parse(string(t)); err != nil {
			return err
		}
	}
	if err := tmpl.Execute(g, g); err != nil {
		return err
	}
	if err := g.generateToProtoFunc(); err != nil {
		return err
	}
	g.P()
	for _, method := range g.Service.Methods {
		if err := g.generateMethod(method); err != nil {
			return err
		}
	}
	return nil
}

func (g *serviceGenerator) generateToProtoFunc() error {
	// Mapper from the ent type to the proto type.
	g.Tmpl(`
	// toProto%(typeName) transforms the ent type to the pb type
	func toProto%(typeName)(e *%(entTypeIdent)) (*%(typeName), error) {
		v := &%(typeName){}`, tmplValues{
		"typeName":     g.EntType.Name,
		"entTypeIdent": g.EntPackage.Ident(g.EntType.Name),
	})
	for _, fld := range g.FieldMap.Fields() {
		conv, err := g.newConverter(fld)
		if err != nil {
			return err
		}
		varName := camel(fld.EntField.StructField())
		g.Tmpl(`
		%(toProto)
		v.%(pbStructField) = %(varName)`, g.withGlobals(tmplValues{
			"pbStructField": fld.PbStructField(),
			"varName":       varName,
			"toProto":       g.renderToProto(conv, varName, "e."+fld.EntField.StructField()),
			"conv":          conv,
		}))
	}

	for _, edg := range g.FieldMap.Edges() {
		conv, err := g.newConverter(edg)
		if err != nil {
			return err
		}
		tmpl := `for _, edg := range e.Edges.%(edgeName) {
					%(converted)
					v.%(edgeName) = append(v.%(edgeName), &%(refType){
						%(pbIdField): %(varName),					
					})
				 }`
		if edg.EntEdge.Unique {
			tmpl = `if edg := e.Edges.%(edgeName); edg != nil {
						%(converted)
						v.%(edgeName) = &%(refType){
							%(pbIdField): %(varName),
						}
					}`
		}
		varName := camel(edg.EntEdge.Type.ID.StructField())
		g.Tmpl(tmpl, g.withGlobals(tmplValues{
			"edgeName":   edg.EntEdge.StructField(),
			"refType":    edg.EntEdge.Type.Name,
			"pbIdField":  edg.EdgeIDPbStructField(),
			"entIdField": edg.EntEdge.Type.ID.StructField(),
			"varName":    varName,
			"converted":  g.renderToProto(conv, varName, "edg."+edg.EntEdge.Type.ID.StructField()),
		}))
	}
	g.P("  return v, nil")
	g.P("}")
	return nil
}

type serviceGenerator struct {
	*protogen.GeneratedFile
	EntPackage protogen.GoImportPath
	File       *protogen.File
	Service    *protogen.Service
	EntType    *gen.Type
	FieldMap   entproto.FieldMap
}

func (g *serviceGenerator) Tmpl(s string, values tmplValues) {
	if err := printTemplate(g.GeneratedFile, s, values); err != nil {
		panic(err)
	}
}

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template
var templates = [][]byte{
	internal.MustAsset("template/service.tmpl"),
	internal.MustAsset("template/enums.tmpl"),
}

func (g *serviceGenerator) generateMethod(me *protogen.Method) error {
	g.Tmpl(`
	// %(name) implements %(svcName)Server.%(name)
	func (svc *%(svcName)) %(name)(ctx %(ctx), req *%(inputIdent)) (*%(outputIdent), error) {`, tmplValues{
		"name":        me.GoName,
		"svcName":     g.Service.GoName,
		"ctx":         contextImp.Ident("Context"),
		"inputIdent":  me.Input.GoIdent,
		"outputIdent": me.Output.GoIdent,
	})

	switch me.GoName {
	case "Create":
		if err := g.generateCreateMethod(); err != nil {
			return err
		}
	case "Get":
		if err := g.generateGetMethod(me.Input.GoIdent.GoName); err != nil {
			return err
		}
	case "Delete":
		if err := g.generateDeleteMethod(); err != nil {
			return err
		}
	case "Update":
		if err := g.generateUpdateMethod(); err != nil {
			return err
		}
	default:
		g.Tmpl(`return nil, %(grpcStatusError)(%(notImplemented), "error")`, tmplValues{
			"grpcStatusError": status.Ident("Error"),
			"notImplemented":  codes.Ident("Unimplemented"),
		})
	}
	g.P("}")
	return nil
}

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
