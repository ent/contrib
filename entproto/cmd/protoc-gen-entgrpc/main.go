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

	"entgo.io/contrib/entproto"
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
	typeName, err := extractEntTypeName(service, graph)
	if err != nil {
		return nil, err
	}
	filename := file.GeneratedFilenamePrefix + "_" + snake(service.GoName) + ".go"
	g := plugin.NewGeneratedFile(filename, file.GoImportPath)
	fieldMap, err := adapter.FieldMap(typeName)
	if err != nil {
		return nil, err
	}
	return &serviceGenerator{
		GeneratedFile: g,
		entPackage:    protogen.GoImportPath(graph.Config.Package),
		file:          file,
		service:       service,
		typeName:      typeName,
		fieldMap:      fieldMap,
	}, nil
}

func (g *serviceGenerator) generate() error {
	g.P("// Code generated by protoc-gen-entgrpc. DO NOT EDIT.")
	g.P("package ", g.file.GoPackageName)
	g.P()
	g.generateConstructor()
	g.P()
	if err := g.generateEnumConvertFuncs(); err != nil {
		return err
	}
	if err := g.generateToProtoFunc(); err != nil {
		return err
	}
	if typeNeedsValidator(g.fieldMap) {
		g.generateValidator()
	}
	g.P()

	for _, method := range g.service.Methods {
		if err := g.generateMethod(method); err != nil {
			return err
		}
	}
	return nil
}

func (g *serviceGenerator) generateConstructor() {
	g.Tmpl(`
	// %(svcName) implements %(svcName)Server
	type %(svcName) struct {
		client *%(entClient)
		Unimplemented%(svcName)Server
	}

	// New%(svcName) returns a new %(svcName)
	func New%(svcName)(client *%(entClient)) *%(svcName) {
		return &%(svcName){
			client: client,
		}
	}`, tmplValues{
		"svcName":   g.service.GoName,
		"entClient": g.entPackage.Ident("Client"),
	})
}

func (g *serviceGenerator) generateEnumConvertFuncs() error {
	for _, ef := range g.fieldMap.Enums() {
		pbEnumIdent := g.pbEnumIdent(ef)
		g.Tmpl(`
		func toProto%(enumTypeName) (e %(entEnumIdent)) %(pbEnumIdent) {
			if v, ok := %(enumTypeName)_value[%(toUpper)(string(e))]; ok {
				return %(pbEnumIdent)(v)
			}
			return %(pbEnumIdent)(0)
		}

		func toEnt%(enumTypeName)(e %(pbEnumIdent)) %(entEnumIdent) {
			if v, ok := %(enumTypeName)_name[int32(e)]; ok {
				return %(entEnumIdent)(%(toLower)(v))
			}
			return ""
		}
`, tmplValues{
			"typeName":     g.typeName,
			"enumTypeName": pbEnumIdent.GoName,
			"entEnumIdent": g.entIdent(snake(g.typeName), ef.EntField.StructField()),
			"pbEnumIdent":  g.pbEnumIdent(ef),
			"toUpper":      protogen.GoImportPath("strings").Ident("ToUpper"),
			"toLower":      protogen.GoImportPath("strings").Ident("ToLower"),
		})
	}
	return nil
}

func (g *serviceGenerator) pbEnumIdent(fld *entproto.FieldMappingDescriptor) protogen.GoIdent {
	enumTypeName := fld.PbFieldDescriptor.GetEnumType().GetName()
	return g.file.GoImportPath.Ident(g.typeName + "_" + enumTypeName)
}

func (g *serviceGenerator) generateToProtoFunc() error {
	// Mapper from the ent type to the proto type.
	g.Tmpl(`
	// toProto%(typeName) transforms the ent type to the pb type
	func toProto%(typeName)(e *%(entTypeIdent)) *%(typeName){
		return &%(typeName) {`, tmplValues{
		"typeName":     g.typeName,
		"entTypeIdent": g.entPackage.Ident(g.typeName),
	})

	for _, fld := range g.fieldMap.Fields() {
		protoFunc, err := g.castToProtoFunc(fld)
		if err != nil {
			return err
		}
		g.Tmpl("%(pbStructField): %(castFunc)(e.%(entStructField)),", tmplValues{
			"pbStructField":  fld.PbStructField(),
			"entStructField": fld.EntField.StructField(),
			"castFunc":       protoFunc,
		})
	}
	g.P("	}")
	g.P("}")
	return nil
}

type serviceGenerator struct {
	*protogen.GeneratedFile
	entPackage protogen.GoImportPath
	file       *protogen.File
	service    *protogen.Service
	typeName   string
	fieldMap   entproto.FieldMap
}

func (g *serviceGenerator) Tmpl(s string, values tmplValues) {
	if err := printTemplate(g.GeneratedFile, s, values); err != nil {
		panic(err)
	}
}

func (g *serviceGenerator) generateMethod(me *protogen.Method) error {
	g.Tmpl(`
	// %(name) implements %(svcName)Server.%(name)
	func (svc *%(svcName)) %(name)(ctx %(ctx), req *%(inputIdent)) (*%(outputIdent), error) {`, tmplValues{
		"name":        me.GoName,
		"svcName":     g.service.GoName,
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
		if err := g.generateGetMethod(); err != nil {
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

func extractEntTypeName(s *protogen.Service, g *gen.Graph) (string, error) {
	typeName := strings.TrimSuffix(s.GoName, "Service")
	for _, gt := range g.Nodes {
		if gt.Name == typeName {
			return typeName, nil
		}
	}
	return "", fmt.Errorf("entproto: type %q of service %q not found in graph", typeName, s.GoName)
}

func (g *serviceGenerator) entIdent(subpath string, ident string) protogen.GoIdent {
	ip := path.Join(string(g.entPackage), subpath)
	return protogen.GoImportPath(ip).Ident(ident)
}
